package main

import (
	"fmt"
	"gitlab.com/4uvirik/my-platform/services/api-gateway/config"
	"gitlab.com/4uvirik/my-platform/services/api-gateway/pkg/logger"
	"log"
	"log/slog"
	"os"
)

func main() {
	fmt.Println("api-gateway starting...")

	cfg := initConfig()

	logger := initLogger(cfg)

	slog.SetDefault(logger)
	slog.Info("logger initialized",
		"env", cfg.App.Env,
		"port", cfg.Server.Port,
	)
}

// initConfig - загрузка и настройка конфига.
func initConfig() *config.Config {
	yamlPath := os.Getenv("YAML_PATH")

	cfg, err := config.Load(yamlPath)
	if err != nil {
		log.Fatalf("config load failed: %v", err)
	}

	if err := cfg.Validate(); err != nil {
		log.Fatalf("config validation failed: %v", err)
	}

	return cfg
}

// initLogger - загрузка и настройка логера.
func initLogger(cfg *config.Config) *slog.Logger {
	logger, err := logger.ConfigLogger(cfg.Logger.Level)
	if err != nil {
		slog.Error("failed to load config", "error", err)
		return nil
	}

	logger.Debug("configuration reading success", slog.Any("cfg", cfg))

	return logger
}
