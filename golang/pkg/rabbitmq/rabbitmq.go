package rabbitmq

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/streadway/amqp"
)

// RabbitClientInterface defines the interface for RabbitMQ operations, making it easier to mock in tests
type RabbitClientInterface interface {
	ConsumeMessages(exchange, routingKey, queueName string) (<-chan amqp.Delivery, error)
	PublishMessage(exchange, routingKey, queueName string, message []byte) error // Alterado para aceitar o nome da fila
	Close() error
	IsClosed() bool
}

// RabbitClient manages RabbitMQ connections
type RabbitClient struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	url     string
}

// newConnection establishes a new connection and channel with RabbitMQ
func newConnection(url string) (*amqp.Connection, *amqp.Channel, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close() // Ensure the connection is closed if channel creation fails
		return nil, nil, fmt.Errorf("failed to open channel: %v", err)
	}

	return conn, channel, nil
}

// NewRabbitClient creates a new RabbitMQ client with the given connection URL
func NewRabbitClient(ctx context.Context, connectionURL string) (*RabbitClient, error) {
	conn, channel, err := newConnection(connectionURL)
	if err != nil {
		return nil, err
	}

	return &RabbitClient{
		conn:    conn,
		channel: channel,
		url:     connectionURL,
	}, nil
}

// ConsumeMessages consumes messages from a specified exchange using a custom queue name and routing key
func (client *RabbitClient) ConsumeMessages(exchange, routingKey, queueName string) (<-chan amqp.Delivery, error) {
	err := client.channel.ExchangeDeclare(
		exchange, "direct", true, true, false, false, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to declare exchange: %v", err)
	}

	queue, err := client.channel.QueueDeclare(
		queueName, true, true, false, false, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to declare queue: %v", err)
	}

	err = client.channel.QueueBind(queue.Name, routingKey, exchange, false, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to bind queue: %v", err)
	}

	msgs, err := client.channel.Consume(queue.Name, "", false, false, false, false, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to consume messages: %v", err)
	}

	return msgs, nil
}

// PublishMessage publishes a message to a specified exchange and binds it to a queue
func (client *RabbitClient) PublishMessage(exchange, routingKey, queueName string, message []byte) error {
	// Ensure the exchange exists before publishing
	err := client.channel.ExchangeDeclare(
		exchange, "direct", true, true, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to declare exchange: %v", err)
	}

	// Ensure the queue exists
	_, err = client.channel.QueueDeclare(
		queueName, true, true, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %v", err)
	}

	// Bind the queue to the exchange with the routing key
	err = client.channel.QueueBind(queueName, routingKey, exchange, false, nil)
	if err != nil {
		return fmt.Errorf("failed to bind queue to exchange: %v", err)
	}

	// Publish the message to the exchange with the routing key
	err = client.channel.Publish(
		exchange, routingKey, false, false, amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
		})
	if err != nil {
		return fmt.Errorf("failed to publish message: %v", err)
	}
	return nil
}

// IsClosed checks if the RabbitMQ connection is closed
func (client *RabbitClient) IsClosed() bool {
	return client.conn.IsClosed()
}

// Close terminates the RabbitMQ connection and channel
func (client *RabbitClient) Close() error {
	err := client.channel.Close()
	if err != nil {
		return fmt.Errorf("failed to close channel: %v", err)
	}
	err = client.conn.Close()
	if err != nil {
		return fmt.Errorf("failed to close connection: %v", err)
	}
	return nil
}

// Reconnect attempts to reconnect to RabbitMQ in case of a lost connection
func (client *RabbitClient) Reconnect(ctx context.Context) error {
	var err error
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context canceled while trying to reconnect")
		default:
			slog.Info("Attempting to reconnect to RabbitMQ...")
			client.conn, client.channel, err = newConnection(client.url)
			if err == nil {
				slog.Info("Reconnected to RabbitMQ successfully")
				return nil
			}
			slog.Error("Failed to reconnect to RabbitMQ", slog.String("error", err.Error()))
			time.Sleep(5 * time.Second) // Retry every 5 seconds
		}
	}
}
