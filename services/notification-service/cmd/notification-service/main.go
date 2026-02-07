package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"gitlab.com/4uvirik/my-platform/services/notification-service/internal/kafka"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	consumer := kafka.NewConsumer(
		[]string{"localhost:9092"},
		"order.created",
		"notification-service",
		logger,
	)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	logger.Info("notification-service started")

	if err := consumer.Run(ctx); err != nil {
		logger.Error("consumer stopped", "error", err)
	}
}
