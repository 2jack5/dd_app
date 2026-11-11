package data

import (
	"context"
	"encoding/json"
	"errors"
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

	refreshPrefix string
	indexPrefix   string
}

func NewRedisSessionStore(client *redis.Client, namespace ...string) *RedisSessionStore {
	refreshPrefix := refreshSessionPrefix
	indexPrefix := userSessionsPrefix

	if len(namespace) > 0 && namespace[0] != "" {
		refreshPrefix = "auth:" + namespace[0] + "_refresh:"
		indexPrefix = "auth:" + namespace[0] + "_sessions:"
	}
	return &RedisSessionStore{
		client:        client,
		now:           time.Now,
		refreshPrefix: refreshPrefix,
		indexPrefix:   indexPrefix,
	}
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
	refreshKey := s.refreshPrefix + session.RefreshTokenHash
	userKey := s.indexPrefix + strconv.FormatInt(session.UserID, 10)
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

func (s *RedisSessionStore) RotateRefreshSession(ctx context.Context, currentHash string, next biz.RefreshSession) error {
	if s == nil || s.client == nil {
		return fmt.Errorf("redis session store is not configured")
	}
	ttl := next.ExpiresAt.Sub(s.now())
	if currentHash == "" || ttl <= 0 || next.UserID <= 0 || next.RefreshTokenHash == "" {
		return biz.ErrInvalidToken
	}
	currentKey := s.refreshPrefix + currentHash
	nextKey := s.refreshPrefix + next.RefreshTokenHash
	userKey := s.indexPrefix + strconv.FormatInt(next.UserID, 10)
	err := s.client.Watch(ctx, func(tx *redis.Tx) error {
		raw, err := tx.Get(ctx, currentKey).Bytes()
		if err == redis.Nil {
			return biz.ErrInvalidToken
		}
		if err != nil {
			return fmt.Errorf("load refresh session from redis: %w", err)
		}
		var current refreshSessionValue
		if json.Unmarshal(raw, &current) != nil || current.UserID != strconv.FormatInt(next.UserID, 10) || current.AuthVersion != next.AuthVersion {
			return biz.ErrInvalidToken
		}
		value, err := json.Marshal(refreshSessionValue{
			UserID: current.UserID, AuthVersion: current.AuthVersion,
			ClientType: current.ClientType, CreatedAt: s.now().UTC().Unix(),
		})
		if err != nil {
			return fmt.Errorf("marshal refresh session: %w", err)
		}
		_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			pipe.Del(ctx, currentKey)
			pipe.Set(ctx, nextKey, value, ttl)
			pipe.SRem(ctx, userKey, currentHash)
			pipe.SAdd(ctx, userKey, next.RefreshTokenHash)
			pipe.Expire(ctx, userKey, ttl)
			return nil
		})
		return err
	}, currentKey)
	if errors.Is(err, redis.TxFailedErr) {
		return biz.ErrInvalidToken
	}
	return err
}
