package service

import (
	"context"
	"gitlab.com/4uvirik/my-platform/services/order-service/internal/kafka"
	"log/slog"
	"time"

	"github.com/google/uuid"

	"gitlab.com/4uvirik/my-platform/pkg/events"
	"gitlab.com/4uvirik/my-platform/services/order-service/internal/repository/postgres"
)

type OrderService struct {
	repo     *postgres.OrderRepository
	producer *kafka.Producer
	log      *slog.Logger
}

func NewOrderService(repo *postgres.OrderRepository, producer *kafka.Producer, log *slog.Logger) *OrderService {
	return &OrderService{
		repo:     repo,
		producer: producer,
		log:      log,
	}
}

func (s *OrderService) Create(ctx context.Context, userID string, amount int64) (uuid.UUID, error) {
	id := uuid.New()

	err := s.repo.Create(ctx, postgres.Order{
		ID:     id,
		UserID: userID,
		Amount: amount,
		Status: "created",
	})
	if err != nil {
		return uuid.Nil, err
	}

	if err := s.producer.PublishOrderCreated(ctx, events.OrderCreated{
		OrderID: id.String(),
		UserID:  userID,
		Amount:  amount,
		At:      time.Now().UTC(),
	}); err != nil {
		s.log.Error(
			"failed to publish OrderCreated event",
			"order_id", id.String(),
			"user_id", userID,
			"amount", amount,
			"error", err,
		)
	}

	return id, nil
}

func (s *OrderService) Get(ctx context.Context, id uuid.UUID) (*postgres.Order, error) {
	return s.repo.Get(ctx, id)
}

func (s *OrderService) ListByUser(ctx context.Context, userID string) ([]postgres.Order, error) {
	return s.repo.ListByUser(ctx, userID)
}
