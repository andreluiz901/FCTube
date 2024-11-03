//go:build e2e

package converter_test

import (
	"context"
	"encoding/json"
	"fmt"
	"imersaofc/internal/converter"
	"imersaofc/pkg/rabbitmq"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	videoID               = 1
	videoPath             = "../../mediatest/media/uploads/1"
	expectedPath          = "../../mediatest/media/uploads/1/mpeg-dash"
	outputFile            = "output.mpd"
	conversionKey         = "conversion"
	finishConversionKey   = "finish-conversion"
	conversionQueue       = "video_conversion_queue"
	finishConversionQueue = "video_confirmation_queue"
	exchangeName          = "conversion_exchange"
)

// TestTaskProcessing tests the full task processing pipeline, including message publishing, processing, and result verification.
func TestTaskProcessing(t *testing.T) {
	ctx := context.Background()

	// Setup PostgreSQL container
	postgresContainer, db, err := setupPostgresContainer(ctx)
	assert.NoError(t, err)
	defer postgresContainer.Terminate(ctx)
	defer db.Close()

	// Set up RabbitMQ container
	rabbitmqC, rabbitMQURL, err := startRabbitMQContainer(ctx)
	assert.NoError(t, err, "Failed to start RabbitMQ container")
	defer rabbitmqC.Terminate(ctx)

	// Set up RabbitMQ client
	rabbitClient, err := rabbitmq.NewRabbitClient(ctx, rabbitMQURL)
	assert.NoError(t, err, "Failed to connect to RabbitMQ")
	defer rabbitClient.Close()

	// Set up the video converter
	rootPath := "../../mediatest/media/uploads"
	videoConverter := converter.NewVideoConverter(rabbitClient, db, rootPath)

	// Prepare the message for conversion
	videoTask := converter.VideoTask{
		VideoID: videoID,
		Path:    videoPath,
	}

	// Marshal the task to JSON
	taskMessage, err := json.Marshal(videoTask)
	assert.NoError(t, err)

	msgs, err := rabbitClient.ConsumeMessages(exchangeName, conversionKey, conversionQueue)
	assert.NoError(t, err, "Failed to consume messages from video conversion queue")

	// Publish the message to the conversion queue
	err = rabbitClient.PublishMessage(exchangeName, conversionKey, conversionQueue, taskMessage)
	assert.NoError(t, err, "Failed to publish message to conversion queue")

	videoConverter.HandleMessage(ctx, <-msgs)

	// Publicar a mensagem de confirmação
	confirmationTask := converter.VideoTask{
		VideoID: videoID,
		Path:    expectedPath,
	}
	confirmationMessage, _ := json.Marshal(confirmationTask)
	err = rabbitClient.PublishMessage(exchangeName, finishConversionKey, finishConversionQueue, confirmationMessage)
	assert.NoError(t, err, "Failed to publish message to confirmation queue")

	// Check if the output file was created
	strVideoID := fmt.Sprintf("%d", videoID)
	outputFilePath := filepath.Join(rootPath, strVideoID, "mpeg-dash", outputFile)
	_, err = os.Stat(outputFilePath)
	assert.NoError(t, err, "Output file was not created")

	// Verify that the message was sent to the finish-conversion queue
	confirmationMsgs, err := rabbitClient.ConsumeMessages(exchangeName, finishConversionKey, finishConversionQueue)
	assert.NoError(t, err, "Failed to consume confirmation messages")

	select {
	case msg := <-confirmationMsgs:
		var confirmation converter.VideoTask
		err = json.Unmarshal(msg.Body, &confirmation)
		fmt.Println(confirmation)
		assert.NoError(t, err)
		assert.Equal(t, videoTask.VideoID, confirmation.VideoID)
		assert.Equal(t, filepath.Join(videoTask.Path, "mpeg-dash"), confirmation.Path)
		msg.Ack(false)
	case <-time.After(5 * time.Second):
		t.Fatal("Timed out waiting for the confirmation message")
	}

	// Verify that the video was marked as processed in the database
	var processed bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM processed_videos WHERE video_id = $1 AND status = 'success')", videoTask.VideoID).Scan(&processed)
	assert.NoError(t, err)
	assert.True(t, processed, "Video was not marked as processed in the database")

	// remove mpeg-dash folder
	err = os.RemoveAll(filepath.Join(rootPath, strVideoID, "mpeg-dash"))
	assert.NoError(t, err, "Failed to remove mpeg-dash folder")
}
