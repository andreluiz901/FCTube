package log

import (
	"bytes"
	"io"
	"os"
	"testing"
)

// TestLoggerOutput checks if the logger properly outputs to os.Stdout.
func TestLoggerOutput(t *testing.T) {
	// Create a pipe to capture the output
	r, w, _ := os.Pipe()
	stdout := os.Stdout
	os.Stdout = w

	// Create a logger and log a message
	log := NewLogger(false)
	log.Error("Testing logger output")

	// Close the writer and restore os.Stdout
	w.Close()
	os.Stdout = stdout

	// Read the captured output
	var buf bytes.Buffer
	io.Copy(&buf, r)

	// Verify that the log message was written to os.Stdout
	if !containsLogLevel(buf.String(), "ERROR") {
		t.Errorf("Expected log level ERROR, but got: %s", buf.String())
	}
}

// containsLogLevel is a helper function that checks if the log output contains the expected log level.
func containsLogLevel(output, level string) bool {
	return bytes.Contains([]byte(output), []byte(level))
}
