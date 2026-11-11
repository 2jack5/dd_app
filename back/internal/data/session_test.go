package data

import (
	"context"
	"net"
	"testing"
	"time"

	"dd_web/back/internal/biz"
	"github.com/redis/go-redis/v9"
)

type captureRedisHook struct{ commands []redis.Cmder }

func (h *captureRedisHook) DialHook(next redis.DialHook) redis.DialHook          { return next }
func (h *captureRedisHook) ProcessHook(next redis.ProcessHook) redis.ProcessHook { return next }
func (h *captureRedisHook) ProcessPipelineHook(redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(_ context.Context, commands []redis.Cmder) error {
		h.commands = append(h.commands, commands...)
		return nil
	}
}

func TestRedisSessionStoreWritesTokenAndUserIndex(t *testing.T) {
	client := redis.NewClient(&redis.Options{Addr: "unused:6379", Dialer: func(context.Context, string, string) (net.Conn, error) {
		t.Fatal("capture hook should prevent network access")
		return nil, nil
	}})
	defer client.Close()
	hook := &captureRedisHook{}
	client.AddHook(hook)
	store := NewRedisSessionStore(client)
	store.now = func() time.Time { return time.Unix(1_700_000_000, 0) }
	err := store.SaveRefreshSession(context.Background(), biz.RefreshSession{
		UserID: 7, AuthVersion: 2, ClientType: biz.WechatClientMini,
		RefreshTokenHash: "token-hash", ExpiresAt: store.now().Add(time.Hour),
	})
	if err != nil {
		t.Fatal(err)
	}
	names := map[string]bool{}
	for _, command := range hook.commands {
		names[command.Name()] = true
	}
	for _, name := range []string{"set", "sadd", "expire"} {
		if !names[name] {
			t.Fatalf("redis transaction missing %s: %#v", name, names)
		}
	}
}

func TestRedisSessionStoreRejectsUnavailableClient(t *testing.T) {
	client := redis.NewClient(&redis.Options{Addr: "unused:6379"})
	_ = client.Close()
	store := NewRedisSessionStore(client)
	err := store.SaveRefreshSession(context.Background(), biz.RefreshSession{UserID: 1, RefreshTokenHash: "hash", ExpiresAt: time.Now().Add(time.Hour)})
	if err == nil {
		t.Fatal("closed redis client was accepted")
	}
}

func TestRedisSessionNamespacesSeparateAdminAndUser(t *testing.T) {
	client := redis.NewClient(&redis.Options{Addr: "unused:6379", Dialer: func(context.Context, string, string) (net.Conn, error) {
		t.Fatal("capture hook should prevent network access")
		return nil, nil
	}})
	defer client.Close()
	hook := &captureRedisHook{}
	client.AddHook(hook)
	session := biz.RefreshSession{UserID: 7, AuthVersion: 1, ClientType: "ADMIN", RefreshTokenHash: "same-hash", ExpiresAt: time.Now().Add(time.Hour)}
	for _, store := range []*RedisSessionStore{NewRedisSessionStore(client), NewRedisSessionStore(client, "admin")} {
		if err := store.SaveRefreshSession(context.Background(), session); err != nil {
			t.Fatal(err)
		}
	}
	keys := map[string]bool{}
	for _, command := range hook.commands {
		if command.Name() == "set" || command.Name() == "sadd" {
			key, ok := command.Args()[1].(string)
			if !ok {
				t.Fatalf("unexpected key: %#v", command.Args())
			}
			keys[key] = true
		}
	}
	for _, key := range []string{"auth:refresh:same-hash", "auth:admin_refresh:same-hash", "auth:user_sessions:7", "auth:admin_sessions:7"} {
		if !keys[key] {
			t.Errorf("missing isolated session key %s", key)
		}
	}
}
