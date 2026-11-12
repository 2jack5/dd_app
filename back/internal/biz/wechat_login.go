package biz

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
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
	SessionID        string
	AccessExpiresAt  time.Time
	RefreshExpiresAt time.Time
	RefreshTokenHash string
}

type TokenClaims struct {
	UserID      int64
	AuthVersion int32
	SessionID   string
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
	SessionID        string
	ClientType       string
	RefreshTokenHash string
	ExpiresAt        time.Time
}

type SessionStore interface {
	SaveRefreshSession(context.Context, RefreshSession) error
	RotateRefreshSession(context.Context, string, RefreshSession) error
	RevokeSession(context.Context, int64, string) error
	SessionActive(context.Context, int64, string) (bool, error)
}

type TokenManager interface {
	Issue(int64, int32) (TokenPair, error)
	IssueForSession(int64, int32, string) (TokenPair, error)
	VerifyAccess(string) (TokenClaims, error)
	VerifyRefresh(string) (TokenClaims, error)
}

type WechatLoginResult struct {
	User      LoginUser
	Tokens    TokenPair
	IsNewUser bool
}

type RefreshTokenResult struct {
	Tokens TokenPair
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
	session := RefreshSession{UserID: user.ID, AuthVersion: user.AuthVersion, SessionID: pair.SessionID, ClientType: clientType, RefreshTokenHash: pair.RefreshTokenHash, ExpiresAt: pair.RefreshExpiresAt}
	if err := u.sessions.SaveRefreshSession(ctx, session); err != nil {
		return WechatLoginResult{}, err
	}
	return WechatLoginResult{User: user, Tokens: pair, IsNewUser: created}, nil
}

func (u *AuthUsecase) RefreshToken(ctx context.Context, refreshToken string) (RefreshTokenResult, error) {
	refreshToken = strings.TrimSpace(refreshToken)
	if refreshToken == "" || len(refreshToken) > 4096 {
		return RefreshTokenResult{}, ErrInvalidArgument
	}
	claims, err := u.tokens.VerifyRefresh(refreshToken)
	if err != nil {
		return RefreshTokenResult{}, err
	}
	user, err := u.repo.UserForAuthentication(ctx, claims.UserID)
	if err != nil {
		return RefreshTokenResult{}, err
	}
	if user.Status != "ACTIVE" {
		return RefreshTokenResult{}, ErrAccountDisabled
	}
	if user.AuthVersion != claims.AuthVersion {
		return RefreshTokenResult{}, ErrInvalidToken
	}
	pair, err := u.tokens.IssueForSession(user.ID, user.AuthVersion, claims.SessionID)
	if err != nil {
		return RefreshTokenResult{}, err
	}
	if u.sessions == nil {
		return RefreshTokenResult{}, fmt.Errorf("rotate refresh session: %w", ErrAuthenticationUnavailable)
	}
	currentHash := sha256.Sum256([]byte(refreshToken))
	next := RefreshSession{
		UserID: user.ID, AuthVersion: user.AuthVersion, SessionID: pair.SessionID,
		RefreshTokenHash: pair.RefreshTokenHash, ExpiresAt: pair.RefreshExpiresAt,
	}
	if err := u.sessions.RotateRefreshSession(ctx, hex.EncodeToString(currentHash[:]), next); err != nil {
		return RefreshTokenResult{}, err
	}
	return RefreshTokenResult{Tokens: pair}, nil
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
	if u.sessions == nil {
		return Principal{}, ErrAuthenticationUnavailable
	}
	active, err := u.sessions.SessionActive(ctx, user.ID, claims.SessionID)
	if err != nil {
		return Principal{}, ErrAuthenticationUnavailable
	}
	if !active {
		return Principal{}, ErrInvalidToken
	}
	return Principal{UserID: user.ID, SessionID: claims.SessionID}, nil
}

func (u *AuthUsecase) Logout(ctx context.Context, userID int64, sessionID string) error {
	if userID <= 0 || strings.TrimSpace(sessionID) == "" {
		return ErrInvalidToken
	}
	if u.sessions == nil {
		return ErrAuthenticationUnavailable
	}
	return u.sessions.RevokeSession(ctx, userID, sessionID)
}
