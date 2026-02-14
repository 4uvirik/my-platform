package dedup

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

// Seen returns true if event already processed
func (s *Service) Seen(ctx context.Context, eventType, key string) (bool, error) {
	redisKey := fmt.Sprintf("dedup:%s:%s", eventType, key)

	ok, err := s.rdb.SetNX(ctx, redisKey, "1", s.ttl).Result()
	if err != nil {
		return false, err
	}

	// if SetNX returned false -> key already exists -> duplicate
	return !ok, nil
}
