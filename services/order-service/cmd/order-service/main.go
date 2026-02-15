package main

import (
	"context"
	"gitlab.com/4uvirik/my-platform/services/order-service/internal/observability"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"

	"gitlab.com/4uvirik/my-platform/services/order-service/config"
	"gitlab.com/4uvirik/my-platform/services/order-service/internal/cache"
	grpcHandler "gitlab.com/4uvirik/my-platform/services/order-service/internal/grpc"
	"gitlab.com/4uvirik/my-platform/services/order-service/internal/idempotency"
	"gitlab.com/4uvirik/my-platform/services/order-service/internal/infrastructure/redis"
	"gitlab.com/4uvirik/my-platform/services/order-service/internal/kafka"
	"gitlab.com/4uvirik/my-platform/services/order-service/internal/ratelimit"
	"gitlab.com/4uvirik/my-platform/services/order-service/internal/repository/postgres"
	"gitlab.com/4uvirik/my-platform/services/order-service/internal/service"
	"gitlab.com/4uvirik/my-platform/services/order-service/pkg/logger"
)

func main() {
	cfg := initConfig()
	logg := initLogger(cfg)
	slog.SetDefault(logg)

	ctx := context.Background()

	redisClient, err := redis.NewClient(cfg.Redis)
	if err != nil {
		logg.Error("failed to init redis", "err", err)
		os.Exit(1)
	}

	orderCache := cache.NewOrderCache(redisClient, 30*time.Second)     // cache TTL
	idempSvc := idempotency.NewService(redisClient, 5*time.Minute)     // idempotency TTL
	rateLimitSvc := ratelimit.NewService(redisClient, 10, time.Minute) // 10 req/min

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

	svc := service.NewOrderService(repo, producer, orderCache, idempSvc, rateLimitSvc, logg)

	observability.Register()

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(observability.UnaryMetricsInterceptor()),
	)

	handler := grpcHandler.NewOrderHandler(svc)
	grpcHandler.Register(grpcServer, handler)

	lis, err := net.Listen("tcp", cfg.Server.GRPCAddr)
	if err != nil {
		log.Fatal(err)
	}

	logg.Info("order-service started", "addr", cfg.Server.GRPCAddr)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			logg.Error("grpc server stopped", "err", err)
		}
	}()

	go func() {
		http.Handle("/metrics", observability.MetricsHandler())
		if err := http.ListenAndServe(":9100", nil); err != nil {
			logg.Error("metrics server failed", "error", err)
		} // порт для Prometheus
	}()

	// graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	<-ctx.Done()
	stop()

	logg.Info("shutting down order-service")
	grpcServer.GracefulStop()

	if err := redisClient.Close(); err != nil {
		logg.Warn("failed to close redis", "err", err)
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
