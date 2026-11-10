package config

import (
	"os"
	"path/filepath"
	"testing"
)

func blankConfig(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	content := []byte("env:\nhttp:\n  port:\nlog:\n  level:\n  format:\npg:\n  host:\n  port:\n  database:\n  user:\n  password:\n  ssl_mode:\n  time_zone:\n  max_open_conns:\n  max_idle_conns:\n  conn_max_lifetime:\nredis:\n  addr:\n  password:\n  db:\n")
	if err := os.WriteFile(filepath.Join(dir, "config.yaml"), content, 0o600); err != nil {
		t.Fatal(err)
	}
	return dir
}

func clearNXEnv(t *testing.T) {
	t.Helper()
	for _, key := range []string{"NX_ENV", "NX_HTTP__PORT", "NX_LOG__LEVEL", "NX_LOG__FORMAT", "NX_PG__HOST", "NX_PG__PORT", "NX_PG__DATABASE", "NX_PG__USER", "NX_PG__PASSWORD", "NX_PG__SSL_MODE", "NX_PG__TIME_ZONE", "NX_PG__MAX_OPEN_CONNS", "NX_PG__MAX_IDLE_CONNS", "NX_PG__CONN_MAX_LIFETIME", "NX_REDIS__ADDR", "NX_REDIS__PASSWORD", "NX_REDIS__DB", "NX_REDIS__DIAL_TIMEOUT", "NX_REDIS__READ_TIMEOUT", "NX_REDIS__WRITE_TIMEOUT", "NX_AUTH__TOKEN_SECRET", "NX_WECHAT__MINI_APP_ID", "NX_WECHAT__MINI_SECRET", "NX_WECHAT__OPEN_APP_ID", "NX_WECHAT__OPEN_SECRET"} {
		t.Setenv(key, "")
	}
}

func TestBlankConfigurationTemplateLoads(t *testing.T) {
	clearNXEnv(t)
	cfg, err := Load(blankConfig(t))
	if err != nil {
		t.Fatal(err)
	}
	if cfg.PG.Host != "" || cfg.PG.Database != "" || cfg.PG.User != "" || cfg.PG.Password != "" || cfg.Log.Level != "" || cfg.Log.Format != "" {
		t.Fatal("blank fixture unexpectedly contains configured values")
	}
}

func TestBlankTemplateSupportsEnvironmentOverrides(t *testing.T) {
	clearNXEnv(t)
	t.Setenv("NX_ENV", "dev")
	t.Setenv("NX_LOG__LEVEL", "warn")
	t.Setenv("NX_LOG__FORMAT", "console")
	cfg, err := Load(blankConfig(t))
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Log.Level != "warn" || cfg.Log.Format != "console" {
		t.Fatal("log configuration overrides were ignored")
	}
}

func TestPostgresFieldsSupportEnvironmentOverrides(t *testing.T) {
	clearNXEnv(t)
	t.Setenv("NX_PG__HOST", "db.internal")
	t.Setenv("NX_PG__PORT", "5433")
	t.Setenv("NX_PG__DATABASE", "dongdong_test")
	t.Setenv("NX_PG__USER", "tester")
	t.Setenv("NX_PG__PASSWORD", "secret")
	t.Setenv("NX_PG__SSL_MODE", "require")
	t.Setenv("NX_PG__TIME_ZONE", "Asia/Shanghai")
	cfg, err := Load(blankConfig(t))
	if err != nil {
		t.Fatal(err)
	}
	if cfg.PG.Host != "db.internal" || cfg.PG.Port != 5433 || cfg.PG.Database != "dongdong_test" || cfg.PG.User != "tester" || cfg.PG.Password != "secret" || cfg.PG.SSLMode != "require" || cfg.PG.TimeZone != "Asia/Shanghai" {
		t.Fatalf("PostgreSQL 环境变量覆盖失败: %+v", cfg.PG)
	}
}
