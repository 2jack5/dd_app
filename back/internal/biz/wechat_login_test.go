package biz

import (
	"context"
	"errors"
	"testing"
	"time"
)

type fakeWechat struct {
	identity WechatIdentity
	err      error
	calls    int
}

func (f *fakeWechat) Exchange(context.Context, string, string) (WechatIdentity, error) {
	f.calls++
	return f.identity, f.err
}

type fakeAuthRepo struct {
	user    LoginUser
	created bool
}

func (f *fakeAuthRepo) FindOrCreateWechatUser(context.Context, WechatIdentity) (LoginUser, bool, error) {
	return f.user, f.created, nil
}
func (f *fakeAuthRepo) UserForAuthentication(context.Context, int64) (LoginUser, error) {
	return f.user, nil
}

type fakeTokens struct {
	pair   TokenPair
	claims TokenClaims
	err    error
}

type fakeSessionStore struct {
	saved       RefreshSession
	rotatedFrom string
	rotateErr   error
	revokedUser int64
	revokedSID  string
	inactive    bool
}

func (f *fakeSessionStore) SaveRefreshSession(_ context.Context, session RefreshSession) error {
	f.saved = session
	return nil
}
func (f *fakeSessionStore) RotateRefreshSession(_ context.Context, currentHash string, session RefreshSession) error {
	f.rotatedFrom = currentHash
	f.saved = session
	return f.rotateErr
}
func (f *fakeSessionStore) RevokeSession(_ context.Context, userID int64, sessionID string) error {
	f.revokedUser, f.revokedSID = userID, sessionID
	return nil
}
func (f *fakeSessionStore) SessionActive(context.Context, int64, string) (bool, error) {
	return !f.inactive, nil
}

func (f fakeTokens) Issue(int64, int32) (TokenPair, error) { return f.pair, nil }
func (f fakeTokens) IssueForSession(_ int64, _ int32, sessionID string) (TokenPair, error) {
	pair := f.pair
	pair.SessionID = sessionID
	return pair, nil
}
func (f fakeTokens) VerifyAccess(string) (TokenClaims, error) {
	return TokenClaims{UserID: 7, AuthVersion: 2, SessionID: "session-1"}, nil
}
func (f fakeTokens) VerifyRefresh(string) (TokenClaims, error) {
	if f.err != nil {
		return TokenClaims{}, f.err
	}
	return f.claims, nil
}

func TestWechatLoginCreatesSessionAndReturnsTokens(t *testing.T) {
	exchanger := &fakeWechat{identity: WechatIdentity{Provider: WechatClientMini, AppID: "app", OpenID: "openid"}}
	repo := &fakeAuthRepo{user: LoginUser{ID: 7, UserNo: "DD1", Nickname: "微信用户", Status: "ACTIVE", AuthVersion: 2}, created: true}
	sessions := &fakeSessionStore{}
	pair := TokenPair{AccessToken: "access", RefreshToken: "refresh", SessionID: "session-1", RefreshTokenHash: "hash", RefreshExpiresAt: time.Now().Add(time.Hour)}
	result, err := NewAuthUsecase(exchanger, repo, fakeTokens{pair: pair}, sessions).WechatLogin(context.Background(), "wechat_mini", " code ")
	if err != nil {
		t.Fatal(err)
	}
	if !result.IsNewUser || result.Tokens.AccessToken != "access" || sessions.saved.UserID != 7 || sessions.saved.RefreshTokenHash != "hash" {
		t.Fatalf("unexpected login result: %+v session=%+v", result, sessions.saved)
	}
}

func TestWechatLoginRejectsInvalidInputBeforeCallingWechat(t *testing.T) {
	exchanger := &fakeWechat{}
	_, err := NewAuthUsecase(exchanger, &fakeAuthRepo{}, fakeTokens{}, &fakeSessionStore{}).WechatLogin(context.Background(), "browser", "")
	if !errors.Is(err, ErrInvalidArgument) || exchanger.calls != 0 {
		t.Fatalf("invalid input result: err=%v calls=%d", err, exchanger.calls)
	}
}

func TestAuthUsecaseAuthenticatesActiveMatchingVersion(t *testing.T) {
	repo := &fakeAuthRepo{user: LoginUser{ID: 7, Status: "ACTIVE", AuthVersion: 2}}
	principal, err := NewAuthUsecase(&fakeWechat{}, repo, fakeTokens{}, &fakeSessionStore{}).Authenticate(context.Background(), "token", "user")
	if err != nil || principal.UserID != 7 {
		t.Fatalf("authenticate: principal=%+v err=%v", principal, err)
	}
}

func TestAuthUsecaseRejectsRevokedSession(t *testing.T) {
	repo := &fakeAuthRepo{user: LoginUser{ID: 7, Status: "ACTIVE", AuthVersion: 2}}
	_, err := NewAuthUsecase(&fakeWechat{}, repo, fakeTokens{}, &fakeSessionStore{inactive: true}).Authenticate(context.Background(), "token", "user")
	if !errors.Is(err, ErrInvalidToken) {
		t.Fatalf("revoked session accepted: %v", err)
	}
}

func TestAuthUsecaseRefreshTokenRotatesSession(t *testing.T) {
	repo := &fakeAuthRepo{user: LoginUser{ID: 7, Status: "ACTIVE", AuthVersion: 2}}
	sessions := &fakeSessionStore{}
	pair := TokenPair{AccessToken: "new-access", RefreshToken: "new-refresh", RefreshTokenHash: "new-hash", RefreshExpiresAt: time.Now().Add(time.Hour)}
	result, err := NewAuthUsecase(&fakeWechat{}, repo, fakeTokens{pair: pair, claims: TokenClaims{UserID: 7, AuthVersion: 2, SessionID: "session-1"}}, sessions).RefreshToken(context.Background(), " old-refresh ")
	if err != nil {
		t.Fatal(err)
	}
	if result.Tokens.AccessToken != "new-access" || sessions.saved.RefreshTokenHash != "new-hash" || sessions.rotatedFrom == "" {
		t.Fatalf("unexpected refresh result: %+v session=%+v", result, sessions)
	}
}

func TestAuthUsecaseRefreshTokenRejectsAuthVersionMismatch(t *testing.T) {
	repo := &fakeAuthRepo{user: LoginUser{ID: 7, Status: "ACTIVE", AuthVersion: 3}}
	_, err := NewAuthUsecase(&fakeWechat{}, repo, fakeTokens{claims: TokenClaims{UserID: 7, AuthVersion: 2, SessionID: "session-1"}}, &fakeSessionStore{}).RefreshToken(context.Background(), "refresh")
	if !errors.Is(err, ErrInvalidToken) {
		t.Fatalf("auth version mismatch: %v", err)
	}
}

func TestAuthUsecaseLogoutRevokesCurrentSession(t *testing.T) {
	sessions := &fakeSessionStore{}
	usecase := NewAuthUsecase(&fakeWechat{}, &fakeAuthRepo{}, fakeTokens{}, sessions)
	if err := usecase.Logout(context.Background(), 7, "session-1"); err != nil {
		t.Fatal(err)
	}
	if sessions.revokedUser != 7 || sessions.revokedSID != "session-1" {
		t.Fatalf("unexpected revoked session: user=%d sid=%q", sessions.revokedUser, sessions.revokedSID)
	}
}
