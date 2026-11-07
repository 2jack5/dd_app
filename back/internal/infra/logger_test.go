package infra

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"strings"
	"testing"
)

func TestStructuredLoggerLevelAndFields(t *testing.T) {
	var buf bytes.Buffer
	l := newLogger(&buf, "warn", "json")
	l.Info("hidden")
	l.Warn("visible", TraceID("request-1"), "order_id", "42")
	var entry map[string]any
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatal(err)
	}
	if entry["level"] != "WARN" || entry["msg"] != "visible" || entry["trace_id"] != "request-1" || entry["order_id"] != "42" || entry["ts"] == nil {
		t.Fatalf("invalid structured log: %s", buf.String())
	}
}
func TestInitSharesStandardDefaultLogger(t *testing.T) {
	original := slog.Default()
	t.Cleanup(func() { slog.SetDefault(original) })
	Init("error", "console")
	if L() != slog.Default() {
		t.Fatal("infra and standard logger diverged")
	}
	var buf bytes.Buffer
	slog.SetDefault(newLogger(&buf, "warn", "console"))
	gormWriter{}.Printf("slow query %s", "sample")
	if !strings.Contains(buf.String(), "level=WARN") || !strings.Contains(buf.String(), "slow query sample") {
		t.Fatalf("GORM did not use configured slog handler: %s", buf.String())
	}
}
