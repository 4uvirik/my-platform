package main

import (
	"context"
	"gitlab.com/4uvirik/my-platform/services/notification-service/internal/dedup"
	"gitlab.com/4uvirik/my-platform/services/notification-service/internal/infrastructure/redis"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	redisClient, err := redis.NewClient(redis.Config{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	if err != nil {
		logg.Error("failed to init redis", "err", err)
		os.Exit(1)
	}

	dedupSvc := dedup.NewService(redisClient, 24*time.Hour)

	consumer, err := kafka.NewConsumer(
		cfg.Kafka.Brokers,
		cfg.Kafka.TopicOrderCreated,
		cfg.Kafka.GroupID,
		dedupSvc,
		logg,
	)
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

	if err := redisClient.Close(); err != nil {
		logg.Warn("failed to close redis", "err", err)
	}

	if err := consumer.Close(); err != nil {
		logg.Warn("failed to close kafka consumer", "err", err)
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
