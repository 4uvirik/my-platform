package main

import (
	"context"
	"gitlab.com/4uvirik/my-platform/services/api-gateway/config"
	"gitlab.com/4uvirik/my-platform/services/api-gateway/internal/app"
	"gitlab.com/4uvirik/my-platform/services/api-gateway/pkg/logger"
	"log"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	cfg := initConfig()

	log := initLogger(cfg)

	slog.SetDefault(log)
	slog.Info("logger initialized",
		"env", cfg.App.Env,
		"port", cfg.Server.Port,
	)

	router := http.NewServeMux()
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	app := app.New(
		cfg.Server.Addr(),
		router,
		log,
	)

	log.Info("starting api-gateway")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := app.Run(ctx); err != nil {
		log.Error("application stopped wit error", "error", err)
		os.Exit(1)
	}
	log.Info("http server started", "addr", cfg.Server.Addr())

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
	log, err := logger.ConfigLogger(cfg.Logger.Level)
	if err != nil {
		slog.Error("failed to load config", "error", err)
		return nil
	}

	log.Debug("configuration reading success", slog.Any("cfg", cfg))

	return log
}
