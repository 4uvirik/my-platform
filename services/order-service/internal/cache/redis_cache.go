package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"gitlab.com/4uvirik/my-platform/services/order-service/internal/repository/postgres"
)

type OrderCache struct {
	rdb *redis.Client
	ttl time.Duration
}

func NewOrderCache(rdb *redis.Client, ttl time.Duration) *OrderCache {
	return &OrderCache{
		rdb: rdb,
		ttl: ttl,
	}
}

func (c *OrderCache) key(id uuid.UUID) string {
	return fmt.Sprintf("order:%s", id.String())
}

func (c *OrderCache) Get(ctx context.Context, id uuid.UUID) (*postgres.Order, bool, error) {
	val, err := c.rdb.Get(ctx, c.key(id)).Result()
	if err == redis.Nil {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}

	var order postgres.Order
	if err := json.Unmarshal([]byte(val), &order); err != nil {
		return nil, false, err
	}

	return &order, true, nil
}

func (c *OrderCache) Set(ctx context.Context, order postgres.Order) error {
	b, err := json.Marshal(order)
	if err != nil {
		return err
	}

	return c.rdb.Set(ctx, c.key(order.ID), b, c.ttl).Err()
}
