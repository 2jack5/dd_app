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

type fakeTokens struct{ pair TokenPair }

type fakeSessionStore struct{ saved RefreshSession }

func (f *fakeSessionStore) SaveRefreshSession(_ context.Context, session RefreshSession) error {
	f.saved = session
	return nil
}

func (f fakeTokens) Issue(int64, int32) (TokenPair, error) { return f.pair, nil }
func (f fakeTokens) VerifyAccess(string) (TokenClaims, error) {
	return TokenClaims{UserID: 7, AuthVersion: 2}, nil
}

func TestWechatLoginCreatesSessionAndReturnsTokens(t *testing.T) {
	exchanger := &fakeWechat{identity: WechatIdentity{Provider: WechatClientMini, AppID: "app", OpenID: "openid"}}
	repo := &fakeAuthRepo{user: LoginUser{ID: 7, UserNo: "DD1", Nickname: "微信用户", Status: "ACTIVE", AuthVersion: 2}, created: true}
	sessions := &fakeSessionStore{}
	pair := TokenPair{AccessToken: "access", RefreshToken: "refresh", RefreshTokenHash: "hash", RefreshExpiresAt: time.Now().Add(time.Hour)}
	result, err := NewAuthUsecase(exchanger, repo, fakeTokens{pair}, sessions).WechatLogin(context.Background(), "wechat_mini", " code ")
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
