package main

import (
	"context"
	"log"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"

	"gitlab.com/4uvirik/my-platform/services/order-service/config"
	grpcHandler "gitlab.com/4uvirik/my-platform/services/order-service/internal/grpc"
	"gitlab.com/4uvirik/my-platform/services/order-service/internal/kafka"
	"gitlab.com/4uvirik/my-platform/services/order-service/internal/repository/postgres"
	"gitlab.com/4uvirik/my-platform/services/order-service/internal/service"
	"gitlab.com/4uvirik/my-platform/services/order-service/pkg/logger"
)

func main() {
	cfg := initConfig()
	logg := initLogger(cfg)
	slog.SetDefault(logg)

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, cfg.DB.DSN)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	repo := postgres.NewOrderRepository(pool)

	producer, err := kafka.NewProducer(cfg.Kafka.Brokers, cfg.Kafka.TopicOrderCreated, logg)
	if err != nil {
		log.Fatal(err)
	}
	defer producer.Close()

	svc := service.NewOrderService(repo, producer, logg)

	grpcServer := grpc.NewServer()
	handler := grpcHandler.NewOrderHandler(svc)
	grpcHandler.Register(grpcServer, handler)

	lis, err := net.Listen("tcp", cfg.Server.GRPCAddr)
	if err != nil {
		log.Fatal(err)
	}

	logg.Info("order-service started", "addr", cfg.Server.GRPCAddr)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}

	// graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	<-ctx.Done()
	stop()

	logg.Info("shutting down order-service")
	grpcServer.GracefulStop()
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
