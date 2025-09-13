package configs

import (
	"log"
	"log/slog"
	"os"
)

// Config func to get env value
func Config(key string) string {
	return os.Getenv(key)
}

func InitializeConfig() {
	// Open or create a log file
	logFile, err := os.OpenFile("auth_server.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()

	// Create a slog handler that writes to the log file
	fileHandler := slog.NewJSONHandler(logFile, nil)

	// Set the global logger to use the file handler
	slog.SetDefault(slog.New(fileHandler))

	// Example log message
	slog.Info("Starting the application")
}
