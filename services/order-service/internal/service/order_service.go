package service

import (
	"context"

	"github.com/google/uuid"

	"gitlab.com/4uvirik/my-platform/services/order-service/internal/repository/postgres"
)

type OrderService struct {
	repo *postgres.OrderRepository
}

func NewOrderService(repo *postgres.OrderRepository) *OrderService {
	return &OrderService{repo: repo}
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

	return id, nil
}

func (s *OrderService) Get(ctx context.Context, id uuid.UUID) (*postgres.Order, error) {
	return s.repo.Get(ctx, id)
}

func (s *OrderService) ListByUser(ctx context.Context, userID string) ([]postgres.Order, error) {
	return s.repo.ListByUser(ctx, userID)
}
