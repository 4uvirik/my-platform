package ratelimit

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Service struct {
	rdb        *redis.Client
	limit      int64
	windowTime time.Duration
}

func NewService(rdb *redis.Client, limit int64, window time.Duration) *Service {
	return &Service{
		rdb:        rdb,
		limit:      limit,
		windowTime: window,
	}
}

func (s *Service) key(key string) string {
	return fmt.Sprintf("ratelimit:%s", key)
}

func (s *Service) Allow(ctx context.Context, key string) (bool, error) {
	redisKey := s.key(key)

	cnt, err := s.rdb.Incr(ctx, redisKey).Result()
	if err != nil {
		return false, err
	}

	if cnt == 1 {
		if err := s.rdb.Expire(ctx, redisKey, s.windowTime).Err(); err != nil {
			return false, err
		}
	}

	return cnt <= s.limit, nil
}
