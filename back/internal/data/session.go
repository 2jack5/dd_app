package data

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"dd_web/back/internal/biz"
	"github.com/redis/go-redis/v9"
)

const (
	refreshSessionPrefix = "auth:refresh:"
	userSessionsPrefix   = "auth:user_sessions:"
)

type RedisSessionStore struct {
	client *redis.Client
	now    func() time.Time
}

func NewRedisSessionStore(client *redis.Client) *RedisSessionStore {
	return &RedisSessionStore{client: client, now: time.Now}
}

type refreshSessionValue struct {
	UserID      string `json:"user_id"`
	AuthVersion int32  `json:"auth_version"`
	ClientType  string `json:"client_type"`
	CreatedAt   int64  `json:"created_at"`
}

func (s *RedisSessionStore) SaveRefreshSession(ctx context.Context, session biz.RefreshSession) error {
	if s == nil || s.client == nil {
		return fmt.Errorf("redis session store is not configured")
	}
	ttl := session.ExpiresAt.Sub(s.now())
	if ttl <= 0 || session.UserID <= 0 || session.RefreshTokenHash == "" {
		return fmt.Errorf("invalid refresh session")
	}
	value, err := json.Marshal(refreshSessionValue{
		UserID: strconv.FormatInt(session.UserID, 10), AuthVersion: session.AuthVersion,
		ClientType: session.ClientType, CreatedAt: s.now().UTC().Unix(),
	})
	if err != nil {
		return fmt.Errorf("marshal refresh session: %w", err)
	}
	refreshKey := refreshSessionPrefix + session.RefreshTokenHash
	userKey := userSessionsPrefix + strconv.FormatInt(session.UserID, 10)
	_, err = s.client.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.Set(ctx, refreshKey, value, ttl)
		pipe.SAdd(ctx, userKey, session.RefreshTokenHash)
		pipe.Expire(ctx, userKey, ttl)
		return nil
	})
	if err != nil {
		return fmt.Errorf("save refresh session to redis: %w", err)
	}
	return nil
}
