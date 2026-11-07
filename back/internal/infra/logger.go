package infra

import (
	"io"
	"log/slog"
	"os"
	"strings"
)

// Init 初始化项目统一日志，供main.go在加载配置后调用。
// level: debug/info/warn/error；format: json/console/text；留空时使用info/json。
func Init(level, format string) { slog.SetDefault(newLogger(os.Stdout, level, format)) }

func newLogger(w io.Writer, level, format string) *slog.Logger {
	var l slog.Level
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "debug":
		l = slog.LevelDebug
	case "warn", "warning":
		l = slog.LevelWarn
	case "error":
		l = slog.LevelError
	default:
		l = slog.LevelInfo
	}
	opts := &slog.HandlerOptions{Level: l, ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey {
			a.Key = "ts"
		}
		return a
	}}
	switch strings.ToLower(strings.TrimSpace(format)) {
	case "console", "text":
		return slog.New(slog.NewTextHandler(w, opts))
	default:
		return slog.New(slog.NewJSONHandler(w, opts))
	}
}

// L 与slog.Default返回同一个实例，Kratos、HTTP与GORM共用日志级别和格式。
func L() *slog.Logger               { return slog.Default() }
func Debug(msg string, args ...any) { L().Debug(msg, args...) }
func Info(msg string, args ...any)  { L().Info(msg, args...) }
func Warn(msg string, args ...any)  { L().Warn(msg, args...) }
func Error(msg string, args ...any) { L().Error(msg, args...) }
func TraceID(id string) slog.Attr   { return slog.String("trace_id", id) }
