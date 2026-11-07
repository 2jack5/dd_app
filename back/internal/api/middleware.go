package api

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	khttp "github.com/go-kratos/kratos/v3/transport/http"
)

// Principal 只能来自可信身份适配器；UserID、AdminID、权限不能从客户端Header直接构造。
// IsCoach表示陪练角色，能否接新单还须用例校验有效实名/技能/处罚；不妨碍暂停接单者查历史。
type Principal struct {
	UserID, AdminID int64
	IsCoach         bool
	Permissions     map[string]bool
	DataScope       string
}
type Authenticator interface {
	Authenticate(context.Context, string, string) (Principal, error)
}

var ErrAuthenticationUnavailable = errors.New("authentication unavailable")

type principalKey struct{}

func PrincipalFromContext(ctx context.Context) (Principal, bool) {
	p, ok := ctx.Value(principalKey{}).(Principal)
	return p, ok
}

type policy struct {
	scope, permission string
	idempotent        bool
}

// protect 是入口检查，不替代资源归属、数据范围、版本及事务幂等；未配置身份服务时拒绝执行。
func (a *API) protect(p policy, next khttp.HandlerFunc) khttp.HandlerFunc {
	return func(ctx khttp.Context) error {
		if p.scope == "channel" {
			return next(ctx)
		} // 回调处理函数固定503，真实验签接入前不返回成功。
		if p.scope != "public" {
			parts := strings.Fields(ctx.Request().Header.Get("Authorization"))
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				return reply(ctx, 401, "AUTH_REQUIRED", "请先登录", nil)
			}
			if a.auth == nil {
				return reply(ctx, 503, "AUTH_NOT_CONFIGURED", "身份验证服务尚未接入", nil)
			}
			principal, err := a.auth.Authenticate(ctx, parts[1], p.scope)
			if errors.Is(err, ErrAuthenticationUnavailable) {
				return reply(ctx, 503, "AUTH_NOT_CONFIGURED", "身份验证服务尚未接入", nil)
			}
			if err != nil {
				return reply(ctx, 401, "INVALID_TOKEN", "登录态无效或已过期", nil)
			}
			if (p.scope == "admin" && principal.AdminID <= 0) || (p.scope != "admin" && principal.UserID <= 0) {
				return reply(ctx, 403, "FORBIDDEN", "身份不匹配", nil)
			}
			if p.scope == "coach" && !principal.IsCoach {
				return reply(ctx, 403, "COACH_NOT_VERIFIED", "当前身份不具备陪练角色", nil)
			}
			if p.permission != "" && !principal.Permissions[p.permission] {
				return reply(ctx, 403, "FORBIDDEN", "没有对应操作权限", nil)
			}
			ctx.Reset(ctx.Response(), ctx.Request().WithContext(context.WithValue(ctx.Request().Context(), principalKey{}, principal)))
		}
		if p.idempotent {
			key := ctx.Request().Header.Get("Idempotency-Key")
			if strings.TrimSpace(key) == "" || len(key) > 128 {
				return reply(ctx, 400, "IDEMPOTENCY_KEY_REQUIRED", "请提供不超过128字节的Idempotency-Key", nil)
			}
		}
		return next(ctx)
	}
}
func requestFilter() khttp.FilterFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var b [16]byte
			if _, e := rand.Read(b[:]); e != nil {
				writeError(w, r, 500, "INTERNAL_ERROR", "服务暂时不可用")
				return
			}
			rid := hex.EncodeToString(b[:])
			w.Header().Set("X-Request-ID", rid)
			r = r.WithContext(context.WithValue(r.Context(), requestIDKey{}, rid))
			r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
			start := time.Now()
			defer func() {
				if recover() != nil {
					writeError(w, r, 500, "INTERNAL_ERROR", "服务暂时不可用")
				}
				slog.Info("http_request", "request_id", rid, "method", r.Method, "path", r.URL.Path, "duration", time.Since(start).String())
			}()
			next.ServeHTTP(w, r)
		})
	}
}
