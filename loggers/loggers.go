package loggers

import (
	"log/slog"
	"os"
)

var (
	API_LOGGER  *slog.Logger
	CRON_LOGGER *slog.Logger
)

// InitSlog initializes a base structured logger with a JSON handler
// It uses the LOG_LEVEL environment variable (default DEBUG) and then
// creates sub-loggers for each component
func InitSlog() {
	// Determine the log level
	logLevel := slog.LevelDebug
	switch os.Getenv("LOG_LEVEL") {
	case "INFO":
		logLevel = slog.LevelInfo
	case "ERROR":
		logLevel = slog.LevelError
	case "DEBUG":
		logLevel = slog.LevelDebug
	}

	// Create a JSON handler that writes to stdout
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     logLevel,
		AddSource: false,
	})

	// Create the base logger and set it as the default
	baseLogger := slog.New(handler)
	slog.SetDefault(baseLogger)

	// Create sub-loggers with additional context for each component
	API_LOGGER = baseLogger.With("component", "API")
	CRON_LOGGER = baseLogger.With("component", "CRON")
}
