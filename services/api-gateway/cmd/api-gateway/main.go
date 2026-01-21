package main

import (
	"fmt"
	"gitlab.com/4uvirik/my-platform/services/api-gateway/pkg/logger"
	"log/slog"
	"os"
)

func main() {
	fmt.Println("api-gateway starting...")

	cfg := initConfig()

	logger := initLogger(cfg)
	slog.SetDefault(logger)

	slog.Info("logger initialized")
}

// initLogger - загрузка и настройка логера.
func initLogger(cfg *config.Config) *slog.Logger {
	logger, err := logger.ConfigLogger(cfg.Logger.Level)
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	logger.Debug("configuration reading success", slog.Any("cfg", cfg))

	return logger
}
