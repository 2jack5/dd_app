package infra

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"dd_web/back/internal/biz"
	"dd_web/back/internal/config"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func wechatHTTPClient(t *testing.T, body string, assert func(*http.Request)) *http.Client {
	t.Helper()
	return &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		assert(r)
		return &http.Response{StatusCode: http.StatusOK, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(body)), Request: r}, nil
	})}
}

func TestWechatClientExchangesMiniProgramCode(t *testing.T) {
	client := NewWechatClient(config.Wechat{APIBaseURL: "https://wechat.test", MiniAppID: "mini-id", MiniSecret: "mini-secret"})
	client.client = wechatHTTPClient(t, `{"openid":"open-id","unionid":"union-id"}`, func(r *http.Request) {
		if r.URL.Path != "/sns/jscode2session" || r.URL.Query().Get("appid") != "mini-id" || r.URL.Query().Get("secret") != "mini-secret" || r.URL.Query().Get("js_code") != "login-code" {
			t.Fatalf("unexpected request: %s", r.URL.String())
		}
	})
	identity, err := client.Exchange(context.Background(), biz.WechatClientMini, "login-code")
	if err != nil || identity.OpenID != "open-id" || identity.UnionID != "union-id" || identity.Provider != biz.WechatClientMini {
		t.Fatalf("identity=%+v err=%v", identity, err)
	}
}

func TestWechatClientMapsConsumedCodeToUnauthorized(t *testing.T) {
	client := NewWechatClient(config.Wechat{APIBaseURL: "https://wechat.test", OpenAppID: "app-id", OpenSecret: "app-secret"})
	client.client = wechatHTTPClient(t, `{"errcode":40163,"errmsg":"code been used"}`, func(*http.Request) {})
	_, err := client.Exchange(context.Background(), biz.WechatClientApp, "used-code")
	if !errors.Is(err, biz.ErrInvalidWechatCode) {
		t.Fatalf("unexpected error: %v", err)
	}
}
