package infra

import (
	"fmt"
	"strings"
	"time"

	"dd_web/back/internal/config"
	"github.com/redis/go-redis/v9"
)

func NewRedis(cfg config.Redis) (*redis.Client, error) {
	cfg.Addr = strings.TrimSpace(cfg.Addr)
	if cfg.Addr == "" {
		return nil, fmt.Errorf("redis.addr is required")
	}
	if cfg.DB < 0 {
		return nil, fmt.Errorf("redis.db must be nonnegative")
	}
	if cfg.DialTimeout <= 0 {
		cfg.DialTimeout = 3 * time.Second
	}
	if cfg.ReadTimeout <= 0 {
		cfg.ReadTimeout = 2 * time.Second
	}
	if cfg.WriteTimeout <= 0 {
		cfg.WriteTimeout = 2 * time.Second
	}
	return redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}), nil
}

func CloseRedis(client *redis.Client) error {
	if client == nil {
		return nil
	}
	return client.Close()
}
