package service

import (
	"context"
	"errors"
	"gitlab.com/4uvirik/my-platform/services/order-service/internal/kafka"
	"log/slog"
	"time"

	"github.com/google/uuid"

	"gitlab.com/4uvirik/my-platform/pkg/events"
	"gitlab.com/4uvirik/my-platform/services/order-service/internal/cache"
	"gitlab.com/4uvirik/my-platform/services/order-service/internal/idempotency"
	"gitlab.com/4uvirik/my-platform/services/order-service/internal/ratelimit"
	"gitlab.com/4uvirik/my-platform/services/order-service/internal/repository/postgres"
)

type OrderService struct {
	repo       *postgres.OrderRepository
	producer   *kafka.Producer
	cache      *cache.OrderCache
	idempotent *idempotency.Service
	rateLimit  *ratelimit.Service
	logger     *slog.Logger
}

func NewOrderService(
	repo *postgres.OrderRepository,
	producer *kafka.Producer,
	cache *cache.OrderCache,
	idempotent *idempotency.Service,
	rateLimit *ratelimit.Service,
	logger *slog.Logger,
) *OrderService {
	return &OrderService{
		repo:       repo,
		producer:   producer,
		cache:      cache,
		idempotent: idempotent,
		rateLimit:  rateLimit,
		logger:     logger,
	}
}

func (s *OrderService) Create(ctx context.Context, userID string, amount int64, idempotencyKey string) (uuid.UUID, error) {
	allowed, err := s.rateLimit.Allow(ctx, userID)
	if err != nil {
		return uuid.Nil, err
	}
	if !allowed {
		return uuid.Nil, errors.New("rate limit exceeded")
	}

	duplicated, err := s.idempotent.CheckAndSet(ctx, idempotencyKey)
	if err != nil {
		return uuid.Nil, err
	}
	if duplicated {
		return uuid.Nil, errors.New("duplicate request")
	}

	id := uuid.New()

	err = s.repo.Create(ctx, postgres.Order{
		ID:     id,
		UserID: userID,
		Amount: amount,
		Status: "created",
	})
	if err != nil {
		return uuid.Nil, err
	}

	if err := s.cache.Set(ctx, postgres.Order{
		ID:     id,
		UserID: userID,
		Amount: amount,
		Status: "created",
	}); err != nil {
		s.logger.Error(
			"failed to Set cache event",
			"id", id,
			"user_id", userID,
			"error", err,
		)
	}

	if err := s.producer.PublishOrderCreated(ctx, events.OrderCreated{
		OrderID: id.String(),
		UserID:  userID,
		Amount:  amount,
		At:      time.Now().UTC(),
	}); err != nil {
		s.logger.Error(
			"failed to publish OrderCreated event",
			"order_id", id.String(),
			"user_id", userID,
			"error", err,
		)
	}

	return id, nil
}

func (s *OrderService) Get(ctx context.Context, id uuid.UUID) (*postgres.Order, error) {
	if order, ok, err := s.cache.Get(ctx, id); err == nil && ok {
		return order, nil
	}

	order, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := s.cache.Set(ctx, *order); err != nil {
		s.logger.Error(
			"failed to set order cache",
			"order_id", id.String(),
			"error", err,
		)
	}

	return order, nil
}

func (s *OrderService) ListByUser(ctx context.Context, userID string) ([]postgres.Order, error) {
	return s.repo.ListByUser(ctx, userID)
}
