package logging

import (
	"log/slog"
	"os"
	"strings"
	"todolist/internal/configs"
)

func GetLogger(loggerConfig configs.LoggerConfig) (*slog.Logger, error) {
	opts := &slog.HandlerOptions{
		Level: getLoggerLevel(loggerConfig),
	}
	jsonHandler := slog.NewJSONHandler(os.Stdout, opts)
	handler := &ContextJsonHandler{Handler: jsonHandler}
	logger := slog.New(handler)
	return logger, nil
}

func getLoggerLevel(loggerConfig configs.LoggerConfig) slog.Level {
	var logLevel slog.Level
	switch strings.ToLower(loggerConfig.Level) {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}
	return logLevel
}
