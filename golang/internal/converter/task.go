package converter

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"time"
	"imersaofc/pkg/rabbitmq"
	"github.com/streadway/amqp"
)

// VideoConverter handles video conversion tasks
type VideoConverter struct {
	rabbitClient *rabbitmq.RabbitClient
	db           *sql.DB
	rootPath     string
}

// VideoTask represents a video conversion task
type VideoTask struct {
	VideoID int    `json:"video_id"`
	Path    string `json:"path"`
}

// NewVideoConverter creates a new instance of VideoConverter
func NewVideoConverter(rabbitClient *rabbitmq.RabbitClient, db *sql.DB, rootPath string) *VideoConverter {
	return &VideoConverter{
		rabbitClient: rabbitClient,
		db:           db,
		rootPath:     rootPath,
	}
}

// HandleMessage processes a video conversion message
func (vc *VideoConverter) HandleMessage(ctx context.Context, d amqp.Delivery, conversionExch, confirmationKey, confirmationQueue string) {
	var task VideoTask

	if err := json.Unmarshal(d.Body, &task); err != nil {
		vc.logError(task, "Failed to deserialize message", err)
		d.Ack(false)
		return
	}

	// Check if the video has already been processed
	if IsProcessed(vc.db, task.VideoID) {
		slog.Warn("Video already processed", slog.Int("video_id", task.VideoID))
		d.Ack(false)
		return
	}

	// Process the video
	err := vc.processVideo(&task)
	if err != nil {
		vc.logError(task, "Error during video conversion", err)
		d.Ack(false)
		return
	}
	slog.Info("Video conversion processed", slog.Int("video_id", task.VideoID))

	// Mark as processed
	err = MarkProcessed(vc.db, task.VideoID)
	if err != nil {
		vc.logError(task, "Failed to mark video as processed", err)
	}
	d.Ack(false)
	slog.Info("Video marked as processed", slog.Int("video_id", task.VideoID))

	// Publicar a mensagem de confirmação
	confirmationMessage := []byte(fmt.Sprintf(`{"video_id": %d, "path":"%s"}`, task.VideoID, task.Path))
	err = vc.rabbitClient.PublishMessage(conversionExch, confirmationKey, confirmationQueue, confirmationMessage)
	if err != nil {
		slog.Error("Failed to publish confirmation message", slog.String("error", err.Error()))
	}
	slog.Info("Published confirmation message", slog.Int("video_id", task.VideoID))
}

// processVideo handles video processing (merging chunks and converting)
func (vc *VideoConverter) processVideo(task *VideoTask) error {
	chunkPath := filepath.Join(vc.rootPath, fmt.Sprintf("%d", task.VideoID))
	mergedFile := filepath.Join(chunkPath, "merged.mp4")
	mpegDashPath := filepath.Join(chunkPath, "mpeg-dash")

	// Merge chunks
	slog.Info("Merging chunks", slog.String("path", chunkPath))
	if err := vc.mergeChunks(chunkPath, mergedFile); err != nil {
		return fmt.Errorf("failed to merge chunks: %v", err)
	}

	// Create directory for MPEG-DASH output
	if err := os.MkdirAll(mpegDashPath, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	// Convert to MPEG-DASH
	ffmpegCmd := exec.Command(
		"ffmpeg", "-i", mergedFile, // Arquivo de entrada
		"-f", "dash", // Formato de saída
		filepath.Join(mpegDashPath, "output.mpd"), // Caminho para salvar o arquivo .mpd
	)

	output, err := ffmpegCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to convert to MPEG-DASH: %v, output: %s", err, string(output))
	}
	slog.Info("Converted to MPEG-DASH", slog.String("path", mpegDashPath))

	// Remove merged file after processing
	if err := os.Remove(mergedFile); err != nil {
		slog.Warn("Failed to remove merged file", slog.String("file", mergedFile), slog.String("error", err.Error()))
	}
	slog.Info("Removed merged file", slog.String("file", mergedFile))

	return nil
}

// Método para extrair o número do nome do arquivo
func (vc *VideoConverter) extractNumber(fileName string) int {
	re := regexp.MustCompile(`\d+`)
	numStr := re.FindString(filepath.Base(fileName)) // Pega o nome do arquivo, sem o caminho
	num, _ := strconv.Atoi(numStr)
	return num
}

// Método para mesclar os chunks
func (vc *VideoConverter) mergeChunks(inputDir, outputFile string) error {
	// Buscar todos os arquivos .chunk no diretório
	chunks, err := filepath.Glob(filepath.Join(inputDir, "*.chunk"))
	if err != nil {
		return fmt.Errorf("failed to find chunks: %v", err)
	}

	// Ordenar os chunks numericamente
	sort.Slice(chunks, func(i, j int) bool {
		return vc.extractNumber(chunks[i]) < vc.extractNumber(chunks[j])
	})

	// Criar arquivo de saída
	output, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create merged file: %v", err)
	}
	defer output.Close()

	// Ler cada chunk e escrever no arquivo final
	for _, chunk := range chunks {
		input, err := os.Open(chunk)
		if err != nil {
			return fmt.Errorf("failed to open chunk %s: %v", chunk, err)
		}

		// Copiar dados do chunk para o arquivo de saída
		_, err = output.ReadFrom(input)
		if err != nil {
			return fmt.Errorf("failed to write chunk %s to merged file: %v", chunk, err)
		}
		input.Close()
	}
	return nil
}

// logError handles logging the error in JSON format
func (vc *VideoConverter) logError(task VideoTask, message string, err error) {
	errorData := map[string]interface{}{
		"video_id": task.VideoID,
		"error":    message,
		"details":  err.Error(),
		"time":     time.Now(),
	}

	serializedError, _ := json.Marshal(errorData)
	slog.Error("Processing error", slog.String("error_details", string(serializedError)))

	RegisterError(vc.db, errorData, err)
}
