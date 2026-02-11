package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"gitlab.com/4uvirik/my-platform/services/notification-service/config"
	"gitlab.com/4uvirik/my-platform/services/notification-service/internal/kafka"
	"gitlab.com/4uvirik/my-platform/services/notification-service/pkg/logger"
)

func main() {
	cfg := initConfig()
	logg := initLogger(cfg)
	slog.SetDefault(logg)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	consumer, err := kafka.NewConsumer(cfg.Kafka.Brokers, cfg.Kafka.TopicOrderCreated, cfg.Kafka.GroupID, logg)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		if err := consumer.Run(ctx); err != nil {
			logg.Error("kafka consumer stopped with error", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logg.Info("notification-service shutting down")
	cancel()
	consumer.Close()
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
