//go:build testcontainers

package converter_test

import (
	"context"
	"encoding/json"
	"fmt"
	"imersaofc/internal/converter"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestIsProcessed(t *testing.T) {
	ctx := context.Background()

	// Setup do container PostgreSQL
	postgresContainer, db, err := setupPostgresContainer(ctx)
	if err != nil {
		t.Fatalf("Failed to setup PostgreSQL container: %v", err)
	}
	defer postgresContainer.Terminate(ctx)
	defer db.Close()

	// Simular a inserção de um vídeo processado
	_, err = db.Exec("INSERT INTO processed_videos (video_id, status, processed_at) VALUES ($1, $2, $3)", 1, "success", time.Now())
	assert.NoError(t, err)

	// Verificar se o vídeo foi marcado como processado
	isProcessed := converter.IsProcessed(db, 1)
	assert.True(t, isProcessed)

	// Verificar se um vídeo não processado retorna `false`
	isProcessed = converter.IsProcessed(db, 999)
	assert.False(t, isProcessed)
}

func TestMarkProcessed(t *testing.T) {
	ctx := context.Background()

	// Setup do container PostgreSQL
	postgresContainer, db, err := setupPostgresContainer(ctx)
	if err != nil {
		t.Fatalf("Failed to setup PostgreSQL container: %v", err)
	}
	defer postgresContainer.Terminate(ctx)
	defer db.Close()

	// Marcar um vídeo como processado
	err = converter.MarkProcessed(db, 2)
	assert.NoError(t, err)

	// Verificar se o vídeo foi processado corretamente
	var status string
	err = db.QueryRow("SELECT status FROM processed_videos WHERE video_id = $1", 2).Scan(&status)
	assert.NoError(t, err)
	assert.Equal(t, "success", status)
}

func TestRegisterError(t *testing.T) {
	ctx := context.Background()

	// Setup do container PostgreSQL
	postgresContainer, db, err := setupPostgresContainer(ctx)
	if err != nil {
		t.Fatalf("Failed to setup PostgreSQL container: %v", err)
	}
	defer postgresContainer.Terminate(ctx)
	defer db.Close()

	// Dados de erro simulados
	errorData := map[string]interface{}{
		"video_id":  1,
		"phases":    []string{"Phase 1", "Phase 2"},
		"error_msg": "Test error", // Certifique-se de incluir "error_msg"
	}

	// Registrar o erro
	converter.RegisterError(db, errorData, fmt.Errorf("Test error"))

	// Verificar se o erro foi registrado corretamente no banco de dados
	var errorDetails []byte
	err = db.QueryRow("SELECT error_details FROM process_errors_log WHERE id = 1").Scan(&errorDetails)
	assert.NoError(t, err)

	// Verificar se os detalhes do erro estão corretos
	var loggedError map[string]interface{}
	err = json.Unmarshal(errorDetails, &loggedError)
	assert.NoError(t, err)

	assert.Equal(t, float64(1), loggedError["video_id"].(float64)) // A conversão de int para float64 é normal no JSON
	assert.Equal(t, "Test error", loggedError["error_msg"])
}
