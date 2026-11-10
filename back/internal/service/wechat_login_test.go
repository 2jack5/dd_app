package service

import (
	"context"
	"errors"
	"testing"
	"time"

	businessv1 "dd_web/back/api/business/v1"
	commonv1 "dd_web/back/api/common/v1"
	"dd_web/back/internal/biz"
	"dd_web/back/internal/ecode"
)

type loginWechatStub struct{}

func (loginWechatStub) Exchange(context.Context, string, string) (biz.WechatIdentity, error) {
	return biz.WechatIdentity{Provider: biz.WechatClientMini, AppID: "app", OpenID: "openid"}, nil
}

type loginRepoStub struct{}

func (loginRepoStub) FindOrCreateWechatUser(context.Context, biz.WechatIdentity) (biz.LoginUser, bool, error) {
	return biz.LoginUser{ID: 9, UserNo: "DD9", Nickname: "微信用户", Status: "ACTIVE", AuthVersion: 1}, true, nil
}
func (loginRepoStub) UserForAuthentication(context.Context, int64) (biz.LoginUser, error) {
	return biz.LoginUser{}, nil
}

type loginTokenStub struct{}
type loginSessionStub struct{ err error }

func (s loginSessionStub) SaveRefreshSession(context.Context, biz.RefreshSession) error { return s.err }

func (loginTokenStub) Issue(int64, int32) (biz.TokenPair, error) {
	return biz.TokenPair{AccessToken: "access", RefreshToken: "refresh", RefreshTokenHash: "hash", AccessExpiresAt: time.Now().Add(time.Hour), RefreshExpiresAt: time.Now().Add(24 * time.Hour)}, nil
}
func (loginTokenStub) VerifyAccess(string) (biz.TokenClaims, error) { return biz.TokenClaims{}, nil }

func TestAuthWechatLoginReturnsTokenAndUser(t *testing.T) {
	auth := biz.NewAuthUsecase(loginWechatStub{}, loginRepoStub{}, loginTokenStub{}, loginSessionStub{})
	svc := NewBusinessService(nil, nil, auth)
	response, err := svc.AuthWechatLogin(context.Background(), &businessv1.AuthWechatLoginRequest{Body: &commonv1.WechatLogin{ClientType: biz.WechatClientMini, Code: "code"}})
	if err != nil {
		t.Fatal(err)
	}
	data, ok := response.Data.AsInterface().(map[string]any)
	if !ok || data["access_token"] != "access" || data["refresh_token"] != "refresh" || data["is_new_user"] != true {
		t.Fatalf("unexpected response: %#v", response.Data.AsInterface())
	}
}

func TestAuthWechatLoginValidatesBody(t *testing.T) {
	svc := NewBusinessService(nil, nil, biz.NewAuthUsecase(loginWechatStub{}, loginRepoStub{}, loginTokenStub{}, loginSessionStub{}))
	if _, err := svc.AuthWechatLogin(context.Background(), &businessv1.AuthWechatLoginRequest{}); err == nil {
		t.Fatal("missing body was accepted")
	}
}

func TestAuthWechatLoginMapsRedisFailureToInternalCode5000(t *testing.T) {
	auth := biz.NewAuthUsecase(loginWechatStub{}, loginRepoStub{}, loginTokenStub{}, loginSessionStub{err: errors.New("redis unavailable")})
	svc := NewBusinessService(nil, nil, auth)
	_, err := svc.AuthWechatLogin(context.Background(), &businessv1.AuthWechatLoginRequest{Body: &commonv1.WechatLogin{ClientType: biz.WechatClientMini, Code: "code"}})
	businessError, ok := ecode.FromError(err)
	if !ok || businessError.Code() != ecode.InternalError || businessError.Code() != 5000 {
		t.Fatalf("redis failure must map to code 5000: %v", err)
	}
}
