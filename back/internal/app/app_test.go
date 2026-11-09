package app

import (
	"strings"
	"testing"
	"time"

	"dd_web/back/internal/config"
)

func TestPostgresConfigRequiresConnectionFields(t *testing.T) {
	clearPostgresEnv(t)
	_, err := PostgresConfig(config.PG{})
	if err == nil || !strings.Contains(err.Error(), "host, database, user, and password are required") {
		t.Fatalf("缺少 PostgreSQL 连接参数时应拒绝启动，实际错误: %v", err)
	}
}

func TestPostgresConfigAppliesDefaults(t *testing.T) {
	clearPostgresEnv(t)
	cfg, err := PostgresConfig(config.PG{Host: "127.0.0.1", Database: "dd", User: "user", Password: "password"})
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Port != 5432 || cfg.SSLMode != "disable" || cfg.TimeZone != "UTC" || cfg.MaxOpenConns != 20 || cfg.MaxIdleConns != 0 || cfg.ConnMaxLifetime != 30*time.Minute {
		t.Fatalf("PostgreSQL 默认连接池配置错误: %+v", cfg)
	}
}

func TestPostgresConfigRejectsInvalidPoolLimits(t *testing.T) {
	clearPostgresEnv(t)
	_, err := PostgresConfig(config.PG{Host: "127.0.0.1", Database: "dd", User: "user", Password: "password", MaxOpenConns: 2, MaxIdleConns: 3})
	if err == nil {
		t.Fatal("非法 PostgreSQL 连接池配置被接受")
	}
}

func TestPostgresConfigSupportsSplitEnvironmentVariables(t *testing.T) {
	clearPostgresEnv(t)
	t.Setenv("DD_DB_HOST", "db.internal")
	t.Setenv("DD_DB_PORT", "5433")
	t.Setenv("DD_DB_NAME", "dongdong_test")
	t.Setenv("DD_DB_USER", "tester")
	t.Setenv("DD_DB_PASSWORD", "secret")
	t.Setenv("DD_DB_SSL_MODE", "require")
	t.Setenv("DD_DB_TIME_ZONE", "Asia/Shanghai")
	cfg, err := PostgresConfig(config.PG{})
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Host != "db.internal" || cfg.Port != 5433 || cfg.Database != "dongdong_test" || cfg.User != "tester" || cfg.Password != "secret" || cfg.SSLMode != "require" || cfg.TimeZone != "Asia/Shanghai" {
		t.Fatalf("PostgreSQL 环境变量覆盖失败: %+v", cfg)
	}
}

func clearPostgresEnv(t *testing.T) {
	t.Helper()
	for _, key := range []string{"DD_DB_HOST", "DD_DB_PORT", "DD_DB_NAME", "DD_DB_USER", "DD_DB_PASSWORD", "DD_DB_SSL_MODE", "DD_DB_TIME_ZONE", "DD_DB_MAX_OPEN", "DD_DB_MAX_IDLE"} {
		t.Setenv(key, "")
	}
}
