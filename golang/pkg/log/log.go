package log

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

func NewLogger(isDebug bool) *slog.Logger {
	if isDebug {
		// Logger in debug mode with colorized text output
		handler := tint.NewHandler(os.Stdout, &tint.Options{
			Level: slog.LevelDebug, // Set the log level to DEBUG
		})
		return slog.New(handler)
	}

	// Logger in non-debug mode showing only errors
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelError, // Set the log level to ERROR
	})
	return slog.New(handler)
}
