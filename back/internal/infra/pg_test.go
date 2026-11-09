package infra

import (
	"net/url"
	"testing"

	"dd_web/back/internal/config"
)

func TestPostgresDSNBuildsFromSeparateFields(t *testing.T) {
	raw := postgresDSN(&config.PG{
		Host:     "2001:db8::1",
		Port:     5432,
		Database: "dongdong",
		User:     "dd_api",
		Password: "p@ss:word",
		SSLMode:  "require",
		TimeZone: "Asia/Shanghai",
	})
	parsed, err := url.Parse(raw)
	if err != nil {
		t.Fatal(err)
	}
	password, _ := parsed.User.Password()
	if parsed.Hostname() != "2001:db8::1" || parsed.Port() != "5432" || parsed.User.Username() != "dd_api" || password != "p@ss:word" || parsed.Path != "/dongdong" {
		t.Fatalf("PostgreSQL 连接参数编码错误: %s", raw)
	}
	if parsed.Query().Get("sslmode") != "require" || parsed.Query().Get("TimeZone") != "Asia/Shanghai" {
		t.Fatalf("PostgreSQL 查询参数编码错误: %s", raw)
	}
}
