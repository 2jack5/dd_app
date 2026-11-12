package biz

import (
	"context"
	"errors"
)

var ErrAuthenticationUnavailable = errors.New("authentication unavailable")

type Principal struct {
	UserID      int64
	AdminID     int64
	SessionID   string
	IsCoach     bool
	Permissions map[string]bool
	DataScope   string
}

type Authenticator interface {
	Authenticate(context.Context, string, string) (Principal, error)
}

type principalKey struct{}

func WithPrincipal(ctx context.Context, principal Principal) context.Context {
	return context.WithValue(ctx, principalKey{}, principal)
}

func PrincipalFromContext(ctx context.Context) (Principal, bool) {
	principal, ok := ctx.Value(principalKey{}).(Principal)
	return principal, ok
}
