package main

import (
	"context"
	"gitlab.com/4uvirik/my-platform/services/order-service/internal/kafka"
	"log"
	"net"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"

	grpcHandler "gitlab.com/4uvirik/my-platform/services/order-service/internal/grpc"
	"gitlab.com/4uvirik/my-platform/services/order-service/internal/repository/postgres"
	"gitlab.com/4uvirik/my-platform/services/order-service/internal/service"
	orderpb "gitlab.com/4uvirik/my-platform/services/order-service/proto"
)

func main() {
	ctx := context.Background()

	dsn := os.Getenv("ORDER_SERVICE_DB_DSN")
	if dsn == "" {
		log.Fatal("ORDER_SERVICE_DB_DSN is not set")
	}

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatal(err)
	}

	producer := kafka.NewProducer(
		cfg.Kafka.Brokers,
		cfg.Kafka.TopicOrderCreated,
		logger,
	)

	repo := postgres.NewOrderRepository(pool)
	svc := service.NewOrderService(repo)
	handler := grpcHandler.NewOrderHandler(svc)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()
	orderpb.RegisterOrderServiceServer(grpcServer, handler)

	log.Println("order-service started on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
