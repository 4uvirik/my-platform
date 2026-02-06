package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Order struct {
	ID     uuid.UUID
	UserID string
	Amount int64
	Status string
}

type OrderRepository struct {
	pool *pgxpool.Pool
}

func NewOrderRepository(pool *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{pool: pool}
}

func (r *OrderRepository) Create(ctx context.Context, order Order) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO orders (id, user_id, amount, status) VALUES ($1, $2, $3, $4)`,
		order.ID, order.UserID, order.Amount, order.Status,
	)
	return err
}

func (r *OrderRepository) Get(ctx context.Context, id uuid.UUID) (*Order, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id, user_id, amount, status FROM orders WHERE id = $1`, id)

	var o Order
	if err := row.Scan(&o.ID, &o.UserID, &o.Amount, &o.Status); err != nil {
		return nil, err
	}
	return &o, nil
}

func (r *OrderRepository) ListByUser(ctx context.Context, userID string) ([]Order, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, amount, status FROM orders WHERE user_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var o Order
		if err := rows.Scan(&o.ID, &o.UserID, &o.Amount, &o.Status); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}
