package config

import "testing"

func TestBlankConfigurationTemplateLoads(t *testing.T) {
	for _, key := range []string{"NX_ENV", "NX_HTTP__PORT", "NX_LOG__LEVEL", "NX_LOG__FORMAT", "NX_PG__DSN", "NX_PG__MAX_OPEN_CONNS", "NX_PG__MAX_IDLE_CONNS", "NX_PG__CONN_MAX_LIFETIME", "NX_REDIS__ADDR", "NX_REDIS__PASSWORD", "NX_REDIS__DB"} {
		t.Setenv(key, "")
	}
	cfg, err := Load("../../config")
	if err != nil {
		t.Fatal(err)
	}
	if cfg.PG.DSN != "" || cfg.Log.Level != "" || cfg.Log.Format != "" {
		t.Fatal("template unexpectedly contains configured values")
	}
}

func TestBlankTemplateSupportsEnvironmentOverrides(t *testing.T) {
	t.Setenv("NX_ENV", "dev")
	t.Setenv("NX_LOG__LEVEL", "warn")
	t.Setenv("NX_LOG__FORMAT", "console")
	cfg, err := Load("../../config")
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Log.Level != "warn" || cfg.Log.Format != "console" {
		t.Fatal("log configuration overrides were ignored")
	}
}
