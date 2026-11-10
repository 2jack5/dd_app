package infra

import (
	"errors"
	"testing"
	"time"

	"dd_web/back/internal/biz"
	"dd_web/back/internal/config"
)

func TestTokenManagerIssuesAndVerifiesAccessToken(t *testing.T) {
	manager := NewTokenManager(config.Auth{Issuer: "test", TokenSecret: "01234567890123456789012345678901", AccessTokenTTL: time.Hour, RefreshTokenTTL: 24 * time.Hour})
	pair, err := manager.Issue(42, 3)
	if err != nil {
		t.Fatal(err)
	}
	claims, err := manager.VerifyAccess(pair.AccessToken)
	if err != nil || claims.UserID != 42 || claims.AuthVersion != 3 || pair.RefreshTokenHash == "" {
		t.Fatalf("claims=%+v pair=%+v err=%v", claims, pair, err)
	}
	if _, err := manager.VerifyAccess(pair.RefreshToken); !errors.Is(err, biz.ErrInvalidToken) {
		t.Fatalf("refresh token accepted as access token: %v", err)
	}
}

func TestTokenManagerRequiresStrongSecret(t *testing.T) {
	manager := NewTokenManager(config.Auth{TokenSecret: "short"})
	if _, err := manager.Issue(1, 1); !errors.Is(err, biz.ErrAuthenticationUnavailable) {
		t.Fatalf("short secret accepted: %v", err)
	}
}
