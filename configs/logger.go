package configs

import (
	"log/slog"
	"os"
	"sync"
)

var (
	once   sync.Once
	logger *slog.Logger
)

// Init sets up global slog configuration
func LoggerInit(logFile string, level slog.Level) {
	once.Do(func() {
		// Open file for logging
		file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			panic("failed to open log file: " + err.Error())
		}

		// Configure handler options
		opts := &slog.HandlerOptions{
			Level: level, // slog.LevelDebug, slog.LevelInfo, ...
		}

		// File handler
		fileHandler := slog.NewTextHandler(file, opts)

		// Set global logger
		logger = slog.New(fileHandler)
		slog.SetDefault(logger)
	})
}

// Get returns the global logger
func Get() *slog.Logger {
	if logger == nil {
		panic("logger not initialized, call logger.Init() first")
	}
	return logger
}
