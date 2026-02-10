package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"gitlab.com/4uvirik/my-platform/pkg/logger"
	"gitlab.com/4uvirik/my-platform/services/order-service/config"
	grpcHandler "gitlab.com/4uvirik/my-platform/services/order-service/internal/grpc"
	"gitlab.com/4uvirik/my-platform/services/order-service/internal/kafka"
	"gitlab.com/4uvirik/my-platform/services/order-service/internal/repository/postgres"
	"gitlab.com/4uvirik/my-platform/services/order-service/internal/service"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"net"
)

func main() {
	cfgPath := os.Getenv("CONFIG_PATH")
	cfg, err := config.Load(cfgPath)
	if err != nil {
		panic(err)
	}

	log := logger.New(logger.Config{
		Level: cfg.Logger.Level,
		JSON:  cfg.Logger.JSON,
	})
	slog.SetDefault(log)

	pool, err := pgxpool.New(context.Background(), cfg.DB.DSN)
	if err != nil {
		log.Error("db connect failed", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	repo := postgres.NewOrderRepository(pool)

	producer, err := kafka.NewProducer(cfg.Kafka.Brokers, cfg.Kafka.TopicOrderCreated, log)
	if err != nil {
		log.Error("kafka producer init failed", "error", err)
		os.Exit(1)
	}
	defer producer.Close()

	svc := service.NewOrderService(repo, producer, log)

	lis, err := net.Listen("tcp", cfg.Server.GRPCAddr)
	if err != nil {
		log.Error("listen failed", "error", err)
		os.Exit(1)
	}

	grpcServer := grpc.NewServer()
	grpcHandler.Register(grpcServer, svc)

	go func() {
		log.Info("grpc server started", "addr", cfg.Server.GRPCAddr)
		if err := grpcServer.Serve(lis); err != nil {
			log.Error("grpc serve failed", "error", err)
		}
	}()

	// graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	<-ctx.Done()
	stop()

	log.Info("shutting down order-service")
	grpcServer.GracefulStop()
}
