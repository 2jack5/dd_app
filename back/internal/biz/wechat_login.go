package biz

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	ErrWechatUnavailable = errors.New("wechat unavailable")
	ErrInvalidWechatCode = errors.New("invalid wechat code")
	ErrAccountDisabled   = errors.New("account disabled")
	ErrInvalidToken      = errors.New("invalid token")
)

const (
	WechatClientMini = "WECHAT_MINI"
	WechatClientApp  = "WECHAT_APP"
)

type WechatIdentity struct {
	Provider string
	AppID    string
	OpenID   string
	UnionID  string
}

type LoginUser struct {
	ID          int64
	UserNo      string
	Nickname    string
	Status      string
	AuthVersion int32
}

type TokenPair struct {
	AccessToken      string
	RefreshToken     string
	AccessExpiresAt  time.Time
	RefreshExpiresAt time.Time
	RefreshTokenHash string
}

type TokenClaims struct {
	UserID      int64
	AuthVersion int32
}

type WechatCodeExchanger interface {
	Exchange(context.Context, string, string) (WechatIdentity, error)
}

type AuthRepository interface {
	FindOrCreateWechatUser(context.Context, WechatIdentity) (LoginUser, bool, error)
	UserForAuthentication(context.Context, int64) (LoginUser, error)
}

type RefreshSession struct {
	UserID           int64
	AuthVersion      int32
	ClientType       string
	RefreshTokenHash string
	ExpiresAt        time.Time
}

type SessionStore interface {
	SaveRefreshSession(context.Context, RefreshSession) error
}

type TokenManager interface {
	Issue(int64, int32) (TokenPair, error)
	VerifyAccess(string) (TokenClaims, error)
}

type WechatLoginResult struct {
	User      LoginUser
	Tokens    TokenPair
	IsNewUser bool
}

type AuthUsecase struct {
	exchanger WechatCodeExchanger
	repo      AuthRepository
	tokens    TokenManager
	sessions  SessionStore
}

func NewAuthUsecase(exchanger WechatCodeExchanger, repo AuthRepository, tokens TokenManager, sessions SessionStore) *AuthUsecase {
	return &AuthUsecase{exchanger: exchanger, repo: repo, tokens: tokens, sessions: sessions}
}

func (u *AuthUsecase) WechatLogin(ctx context.Context, clientType, code string) (WechatLoginResult, error) {
	clientType, code = strings.TrimSpace(strings.ToUpper(clientType)), strings.TrimSpace(code)
	if clientType != WechatClientMini && clientType != WechatClientApp || code == "" || len(code) > 256 {
		return WechatLoginResult{}, ErrInvalidArgument
	}
	identity, err := u.exchanger.Exchange(ctx, clientType, code)
	if err != nil {
		return WechatLoginResult{}, err
	}
	if identity.AppID == "" || identity.OpenID == "" {
		return WechatLoginResult{}, fmt.Errorf("%w: empty identity", ErrWechatUnavailable)
	}
	user, created, err := u.repo.FindOrCreateWechatUser(ctx, identity)
	if err != nil {
		return WechatLoginResult{}, err
	}
	if user.Status != "ACTIVE" {
		return WechatLoginResult{}, ErrAccountDisabled
	}
	pair, err := u.tokens.Issue(user.ID, user.AuthVersion)
	if err != nil {
		return WechatLoginResult{}, err
	}
	if u.sessions == nil {
		return WechatLoginResult{}, fmt.Errorf("save refresh session: %w", ErrAuthenticationUnavailable)
	}
	session := RefreshSession{UserID: user.ID, AuthVersion: user.AuthVersion, ClientType: clientType, RefreshTokenHash: pair.RefreshTokenHash, ExpiresAt: pair.RefreshExpiresAt}
	if err := u.sessions.SaveRefreshSession(ctx, session); err != nil {
		return WechatLoginResult{}, err
	}
	return WechatLoginResult{User: user, Tokens: pair, IsNewUser: created}, nil
}

func (u *AuthUsecase) Authenticate(ctx context.Context, bearerToken, _ string) (Principal, error) {
	claims, err := u.tokens.VerifyAccess(strings.TrimSpace(bearerToken))
	if err != nil {
		return Principal{}, err
	}
	user, err := u.repo.UserForAuthentication(ctx, claims.UserID)
	if err != nil || user.Status != "ACTIVE" || user.AuthVersion != claims.AuthVersion {
		return Principal{}, ErrAuthenticationUnavailable
	}
	return Principal{UserID: user.ID}, nil
}
