package converter_test

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"path/filepath"

	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	postgresUser     = "testuser"
	postgresPassword = "testpassword"
	postgresDB       = "testdb"
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

// setupPostgresContainer creates and sets up a PostgreSQL container
func setupPostgresContainer(ctx context.Context) (testcontainers.Container, *sql.DB, error) {
	req := testcontainers.ContainerRequest{
		Image:        "postgres:13",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     postgresUser,
			"POSTGRES_PASSWORD": postgresPassword,
			"POSTGRES_DB":       postgresDB,
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}

	postgresContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, nil, err
	}

	host, err := postgresContainer.Host(ctx)
	if err != nil {
		return nil, nil, err
	}

	port, err := postgresContainer.MappedPort(ctx, "5432")
	if err != nil {
		return nil, nil, err
	}

	// Create DB connection
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", postgresUser, postgresPassword, host, port.Port(), postgresDB)
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, nil, err
	}

	// Wait for the DB to be ready
	err = db.Ping()
	if err != nil {
		return nil, nil, err
	}

	// Execute the SQL file
	sqlFilePath := filepath.Join("../../", "db.sql")
	sqlContent, err := ioutil.ReadFile(sqlFilePath)
	if err != nil {
		return nil, nil, err
	}

	_, err = db.Exec(string(sqlContent))
	if err != nil {
		return nil, nil, err
	}

	return postgresContainer, db, nil
}
