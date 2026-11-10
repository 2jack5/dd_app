package infra

import (
	"testing"
	"time"

	"dd_web/back/internal/config"
)

func TestNewRedisValidatesConfiguration(t *testing.T) {
	if _, err := NewRedis(config.Redis{}); err == nil {
		t.Fatal("empty redis address was accepted")
	}
	if _, err := NewRedis(config.Redis{Addr: "127.0.0.1:6379", DB: -1}); err == nil {
		t.Fatal("negative redis database was accepted")
	}
}

func TestNewRedisAppliesTimeouts(t *testing.T) {
	client, err := NewRedis(config.Redis{Addr: "127.0.0.1:6379", DialTimeout: time.Second, ReadTimeout: 2 * time.Second, WriteTimeout: 3 * time.Second})
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()
	options := client.Options()
	if options.DialTimeout != time.Second || options.ReadTimeout != 2*time.Second || options.WriteTimeout != 3*time.Second {
		t.Fatalf("redis timeouts were not applied: %+v", options)
	}
}
