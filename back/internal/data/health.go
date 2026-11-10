package data

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type DependencyHealth struct {
	data  *Data
	redis *redis.Client
}

func NewDependencyHealth(data *Data, redisClient *redis.Client) *DependencyHealth {
	return &DependencyHealth{data: data, redis: redisClient}
}

func (h *DependencyHealth) Ping(ctx context.Context) error {
	if err := h.data.Ping(ctx); err != nil {
		return fmt.Errorf("postgres readiness: %w", err)
	}
	if h.redis == nil {
		return fmt.Errorf("redis readiness: client is not configured")
	}
	if err := h.redis.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("redis readiness: %w", err)
	}
	return nil
}
