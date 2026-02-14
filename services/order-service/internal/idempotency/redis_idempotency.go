package idempotency

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Service struct {
	rdb *redis.Client
	ttl time.Duration
}

func NewService(rdb *redis.Client, ttl time.Duration) *Service {
	return &Service{
		rdb: rdb,
		ttl: ttl,
	}
}

func (s *Service) key(key string) string {
	return fmt.Sprintf("idempotency:%s", key)
}

// CheckAndSet возвращает true, если ключ уже был (дубликат)
func (s *Service) CheckAndSet(ctx context.Context, key string) (bool, error) {
	ok, err := s.rdb.SetNX(ctx, s.key(key), "1", s.ttl).Result()
	if err != nil {
		return false, err
	}

	return !ok, nil
}
