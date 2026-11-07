package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	khttp "github.com/go-kratos/kratos/v3/transport/http"
)

// HTTP错误分类保留在api层，数据库初始化不再放在api包。
var errDatabaseUnavailable = errors.New("database unavailable")
var errNotFound = errors.New("resource not found")
var errInvalidArgument = errors.New("invalid argument")

type requestIDKey struct{}
type response struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	Data      any    `json:"data"`
	RequestID string `json:"request_id"`
}

func requestID(ctx context.Context) string { v, _ := ctx.Value(requestIDKey{}).(string); return v }
func reply(ctx khttp.Context, status int, code, message string, data any) error {
	return ctx.JSON(status, response{code, message, data, requestID(ctx)})
}
func notImplemented(ctx khttp.Context, operation string) error {
	return reply(ctx, 501, "NOT_IMPLEMENTED", "接口已设计，业务逻辑尚未实现", map[string]string{"operation": operation})
}
func dataError(ctx khttp.Context, err error) error {
	switch {
	case errors.Is(err, errDatabaseUnavailable):
		return reply(ctx, 503, "DATABASE_UNAVAILABLE", "数据库尚未配置或不可用", nil)
	case errors.Is(err, errNotFound):
		return reply(ctx, 404, "NOT_FOUND", "资源不存在或未启用", nil)
	case errors.Is(err, errInvalidArgument):
		return reply(ctx, 400, "INVALID_ARGUMENT", "ID、日期或分页参数无效", nil)
	default:
		return reply(ctx, 500, "INTERNAL_ERROR", "服务暂时不可用", nil)
	}
}
func writeError(w http.ResponseWriter, r *http.Request, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(response{Code: code, Message: message, Data: nil, RequestID: requestID(r.Context())})
}
