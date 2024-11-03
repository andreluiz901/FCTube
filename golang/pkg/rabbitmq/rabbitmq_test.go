//go:build testcontainers

package rabbitmq_test

import (
	"context"
	"fmt"
	"imersaofc/pkg/rabbitmq"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// Helper function to start a RabbitMQ container for testing
func startRabbitMQContainer(ctx context.Context) (testcontainers.Container, string, error) {
	req := testcontainers.ContainerRequest{
		Image:        "rabbitmq:3-management",
		ExposedPorts: []string{"5672/tcp", "15672/tcp"},
		WaitingFor:   wait.ForLog("Server startup complete"),
	}

	rabbitmqC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, "", err
	}

	// Wait for RabbitMQ to be ready
	host, err := rabbitmqC.Host(ctx)
	if err != nil {
		return nil, "", err
	}

	port, err := rabbitmqC.MappedPort(ctx, "5672")
	if err != nil {
		return nil, "", err
	}

	rabbitmqURL := fmt.Sprintf("amqp://guest:guest@%s:%s/", host, port.Port())
	fmt.Println("RabbitMQ URL:", rabbitmqURL)
	return rabbitmqC, rabbitmqURL, nil
}

// TestRabbitMQPublishAndConsume tests publishing and consuming messages using RabbitMQ
func TestRabbitMQPublishAndConsume(t *testing.T) {
	// Set up RabbitMQ container
	ctx := context.Background()
	rabbitmqC, rabbitMQURL, err := startRabbitMQContainer(ctx)
	assert.NoError(t, err, "Failed to start RabbitMQ container")
	defer rabbitmqC.Terminate(ctx)

	// Set up RabbitMQ client
	client, err := rabbitmq.NewRabbitClient(ctx, rabbitMQURL)
	assert.NoError(t, err, "Failed to connect to RabbitMQ")
	defer client.Close()

	// Declare exchange, routing key, and queue
	exchange := "test_exchange"
	routingKey := "test_key"
	queueName := "test_queue"

	// Start consuming messages in the background
	msgs, err := client.ConsumeMessages(exchange, routingKey, queueName)
	assert.NoError(t, err, "Failed to consume messages")

	// Test case 1: Successfully publish a message
	t.Run("Publish and consume a valid message", func(t *testing.T) {
		testMessage := []byte("Hello, RabbitMQ!")
		err = client.PublishMessage(exchange, routingKey, queueName, testMessage) // Adjusted to include queueName
		assert.NoError(t, err, "Failed to publish message")

		// Wait for the message or timeout
		select {
		case msg := <-msgs:
			assert.Equal(t, string(testMessage), string(msg.Body), "Message content mismatch")
			fmt.Println("Received message:", string(msg.Body))
			msg.Ack(false)
		case <-time.After(5 * time.Second):
			t.Fatal("Timed out waiting for message")
		}
	})

	// Test case 2: Attempt to consume from a non-existent queue (should not error but no messages expected)
	t.Run("Consume from non-existent queue", func(t *testing.T) {
		invalidQueue := "invalid_queue"
		msgs, err := client.ConsumeMessages(exchange, routingKey, invalidQueue)
		assert.NoError(t, err, "Consuming from a non-existent queue should not return an error immediately")

		select {
		case <-msgs:
			t.Fatal("Did not expect to receive any messages from non-existent queue")
		case <-time.After(2 * time.Second):
			// Pass the test, since we do not expect any messages
		}
	})

	// Test case 3: Check reconnection logic
	t.Run("Reconnect on connection failure", func(t *testing.T) {
		client.Close()

		err := client.Reconnect(ctx)
		assert.NoError(t, err, "Failed to reconnect to RabbitMQ")

		// Re-declare exchange and queue after reconnect
		msgs, err = client.ConsumeMessages(exchange, routingKey, queueName)
		assert.NoError(t, err, "Failed to consume messages after reconnect")

		err = client.PublishMessage(exchange, routingKey, queueName, []byte("Reconnected Message")) // Adjusted to include queueName
		assert.NoError(t, err, "Failed to publish message after reconnect")

		// Wait for the message or timeout
		select {
		case msg := <-msgs:
			assert.Equal(t, "Reconnected Message", string(msg.Body), "Message mismatch after reconnect")
			msg.Ack(false)
		case <-time.After(5 * time.Second):
			t.Fatal("Timed out waiting for message after reconnect")
		}
	})
}
