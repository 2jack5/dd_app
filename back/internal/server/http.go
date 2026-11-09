package server

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	commonv1 "dd_web/back/api/common/v1"
	"dd_web/back/internal/biz"
	"dd_web/back/internal/config"
	"dd_web/back/internal/infra"
	kerrors "github.com/go-kratos/kratos/v3/errors"
	"github.com/go-kratos/kratos/v3/middleware"
	"github.com/go-kratos/kratos/v3/middleware/recovery"
	"github.com/go-kratos/kratos/v3/transport"
	khttp "github.com/go-kratos/kratos/v3/transport/http"
)

type requestIDKey struct{}

type routePolicy struct {
	Scope      string
	Permission string
	Idempotent bool
}

type envelope struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	Data      any    `json:"data"`
	RequestID string `json:"request_id"`
}

func newHTTPServer(httpCfg config.HTTP, defaultAddress, addressEnv, timeoutEnv string, auth biz.Authenticator, policies map[string]routePolicy) (*khttp.Server, error) {
	address, timeout, err := httpConfig(httpCfg, defaultAddress, addressEnv, timeoutEnv)
	if err != nil {
		return nil, err
	}
	return khttp.NewServer(
		khttp.Address(address), khttp.Timeout(timeout),
		khttp.Filter(requestFilter()), khttp.Middleware(recovery.Recovery(), accessMiddleware(auth, policies)),
		khttp.ResponseEncoder(encodeResponse), khttp.ErrorEncoder(encodeError),
		khttp.NotFoundHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			writeEnvelope(w, r, 404, "NOT_FOUND", "接口不存在", nil)
		})),
		khttp.MethodNotAllowedHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			writeEnvelope(w, r, 405, "METHOD_NOT_ALLOWED", "请求方法不支持", nil)
		})),
	), nil
}

func httpConfig(cfg config.HTTP, defaultAddress, addressEnv, timeoutEnv string) (string, time.Duration, error) {
	address, timeout := defaultAddress, 5*time.Second
	if value := strings.TrimSpace(cfg.Port); value != "" {
		if _, _, err := net.SplitHostPort(value); err == nil {
			address = value
		} else {
			port, err := strconv.Atoi(value)
			if err != nil || port < 1 || port > 65535 {
				return "", 0, fmt.Errorf("http.port must be a port number or host:port")
			}
			address = net.JoinHostPort("127.0.0.1", value)
		}
	}
	if cfg.ShutdownTimeout < 0 {
		return "", 0, fmt.Errorf("http.shutdown_timeout must be nonnegative")
	}
	if value, ok := os.LookupEnv(addressEnv); ok {
		address = value
	}
	if value := os.Getenv(timeoutEnv); value != "" {
		parsed, err := time.ParseDuration(value)
		if err != nil || parsed <= 0 {
			return "", 0, fmt.Errorf("%s must be a positive duration", timeoutEnv)
		}
		timeout = parsed
	}
	return address, timeout, nil
}

func encodeResponse(w http.ResponseWriter, r *http.Request, value any) error {
	reply, ok := value.(*commonv1.ApiResponse)
	if !ok {
		return fmt.Errorf("unsupported response type %T", value)
	}
	var data any
	if reply.Data != nil {
		data = reply.Data.AsInterface()
	}
	writeEnvelope(w, r, http.StatusOK, reply.Code, reply.Message, data)
	return nil
}

func encodeError(w http.ResponseWriter, r *http.Request, err error) {
	e := kerrors.FromError(err)
	var data any
	if operation := e.Metadata["operation"]; operation != "" {
		data = map[string]string{"operation": operation}
	}
	writeEnvelope(w, r, int(e.Code), e.Reason, e.Message, data)
}

func writeEnvelope(w http.ResponseWriter, r *http.Request, status int, code, message string, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(envelope{code, message, data, requestID(r)})
}

func requestID(r *http.Request) string {
	value, _ := r.Context().Value(requestIDKey{}).(string)
	return value
}

func requestFilter() khttp.FilterFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var bytes [16]byte
			if _, err := rand.Read(bytes[:]); err != nil {
				writeEnvelope(w, r, 500, "INTERNAL_ERROR", "服务暂时不可用", nil)
				return
			}
			id := hex.EncodeToString(bytes[:])
			w.Header().Set("X-Request-ID", id)
			r = r.WithContext(context.WithValue(r.Context(), requestIDKey{}, id))
			r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
			started := time.Now()
			defer func() {
				infra.L().Info("http_request", "request_id", id, "method", r.Method, "path", r.URL.Path, "duration", time.Since(started).String())
			}()
			next.ServeHTTP(w, r)
		})
	}
}

func accessMiddleware(auth biz.Authenticator, policies map[string]routePolicy) middleware.Middleware {
	return func(next middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (any, error) {
			tr, ok := transport.FromServerContext(ctx)
			if !ok {
				return next(ctx, req)
			}
			policy, ok := policies[tr.Operation()]
			if !ok {
				return nil, kerrors.InternalServer("INTERNAL_ERROR", "接口策略未配置")
			}
			if policy.Scope != "public" && policy.Scope != "channel" {
				parts := strings.Fields(tr.RequestHeader().Get("Authorization"))
				if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
					return nil, kerrors.Unauthorized("AUTH_REQUIRED", "请先登录")
				}
				if auth == nil {
					return nil, kerrors.ServiceUnavailable("AUTH_NOT_CONFIGURED", "身份验证服务尚未接入")
				}
				principal, err := auth.Authenticate(ctx, parts[1], policy.Scope)
				if errors.Is(err, biz.ErrAuthenticationUnavailable) {
					return nil, kerrors.ServiceUnavailable("AUTH_NOT_CONFIGURED", "身份验证服务尚未接入")
				}
				if err != nil {
					return nil, kerrors.Unauthorized("INVALID_TOKEN", "登录态无效或已过期")
				}
				if (policy.Scope == "admin" && principal.AdminID <= 0) || (policy.Scope != "admin" && principal.UserID <= 0) {
					return nil, kerrors.Forbidden("FORBIDDEN", "身份不匹配")
				}
				if policy.Scope == "coach" && !principal.IsCoach {
					return nil, kerrors.Forbidden("COACH_NOT_VERIFIED", "当前身份不具备陪练角色")
				}
				if policy.Permission != "" && !principal.Permissions[policy.Permission] {
					return nil, kerrors.Forbidden("FORBIDDEN", "没有对应操作权限")
				}
				ctx = biz.WithPrincipal(ctx, principal)
			}
			if policy.Idempotent {
				key := strings.TrimSpace(tr.RequestHeader().Get("Idempotency-Key"))
				if key == "" || len(key) > 128 {
					return nil, kerrors.BadRequest("IDEMPOTENCY_KEY_REQUIRED", "请提供不超过128字节的Idempotency-Key")
				}
			}
			return next(ctx, req)
		}
	}
}
