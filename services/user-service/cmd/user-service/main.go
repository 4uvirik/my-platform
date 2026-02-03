package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	"gitlab.com/4uvirik/my-platform/services/user-service/config"
	"gitlab.com/4uvirik/my-platform/services/user-service/internal/app"
	http2 "gitlab.com/4uvirik/my-platform/services/user-service/internal/http"
	"gitlab.com/4uvirik/my-platform/services/user-service/internal/repository/memory"
	"gitlab.com/4uvirik/my-platform/services/user-service/pkg/logger"
)

func main() {
	cfg := initConfig()
	logg := initLogger(cfg)

	slog.SetDefault(logg)
	logg.Info("user-service starting", "env", cfg.App.Env, "addr", cfg.Addr())

	repo := memory.New()
	router := http2.NewRouter(repo)

	application := app.New(cfg.Addr(), router, logg)

	if err := application.Run(context.Background()); err != nil {
		logg.Error("service stopped with error", "error", err)
		os.Exit(1)
	}
}

func initConfig() *config.Config {
	yamlPath := os.Getenv("YAML_PATH")

	cfg, err := config.Load(yamlPath)
	if err != nil {
		log.Fatalf("config load failed: %v", err)
	}

	return cfg
}

func initLogger(cfg *config.Config) *slog.Logger {
	logg, err := logger.ConfigLogger(cfg.Logger.Level)
	if err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}

	return logg
}
