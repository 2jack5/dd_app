package infra

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"dd_web/back/internal/biz"
	"dd_web/back/internal/config"
)

type WechatClient struct {
	config config.Wechat
	client *http.Client
}

func NewWechatClient(cfg config.Wechat) *WechatClient {
	if strings.TrimSpace(cfg.APIBaseURL) == "" {
		cfg.APIBaseURL = "https://api.weixin.qq.com"
	}
	return &WechatClient{config: cfg, client: &http.Client{Timeout: 5 * time.Second}}
}

type wechatCodeResponse struct {
	OpenID  string `json:"openid"`
	UnionID string `json:"unionid"`
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

func (c *WechatClient) Exchange(ctx context.Context, clientType, code string) (biz.WechatIdentity, error) {
	provider, appID, secret, path, codeKey := "", "", "", "", ""
	switch clientType {
	case biz.WechatClientMini:
		provider, appID, secret, path, codeKey = biz.WechatClientMini, c.config.MiniAppID, c.config.MiniSecret, "/sns/jscode2session", "js_code"
	case biz.WechatClientApp:
		provider, appID, secret, path, codeKey = biz.WechatClientApp, c.config.OpenAppID, c.config.OpenSecret, "/sns/oauth2/access_token", "code"
	default:
		return biz.WechatIdentity{}, biz.ErrInvalidArgument
	}
	if strings.TrimSpace(appID) == "" || strings.TrimSpace(secret) == "" {
		return biz.WechatIdentity{}, fmt.Errorf("%w: credentials for %s are not configured", biz.ErrWechatUnavailable, clientType)
	}
	endpoint, err := url.Parse(strings.TrimRight(c.config.APIBaseURL, "/") + path)
	if err != nil {
		return biz.WechatIdentity{}, fmt.Errorf("%w: invalid api endpoint", biz.ErrWechatUnavailable)
	}
	query := endpoint.Query()
	query.Set("appid", appID)
	query.Set("secret", secret)
	query.Set(codeKey, code)
	query.Set("grant_type", "authorization_code")
	endpoint.RawQuery = query.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return biz.WechatIdentity{}, fmt.Errorf("%w: %v", biz.ErrWechatUnavailable, err)
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return biz.WechatIdentity{}, fmt.Errorf("%w: %v", biz.ErrWechatUnavailable, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 4096))
		return biz.WechatIdentity{}, fmt.Errorf("%w: http status %d", biz.ErrWechatUnavailable, resp.StatusCode)
	}
	var payload wechatCodeResponse
	if err := json.NewDecoder(io.LimitReader(resp.Body, 64<<10)).Decode(&payload); err != nil {
		return biz.WechatIdentity{}, fmt.Errorf("%w: invalid response", biz.ErrWechatUnavailable)
	}
	if payload.ErrCode != 0 {
		if isInvalidWechatCode(payload.ErrCode) {
			return biz.WechatIdentity{}, fmt.Errorf("%w: wechat error %d", biz.ErrInvalidWechatCode, payload.ErrCode)
		}
		return biz.WechatIdentity{}, fmt.Errorf("%w: wechat error %d", biz.ErrWechatUnavailable, payload.ErrCode)
	}
	if payload.OpenID == "" {
		return biz.WechatIdentity{}, fmt.Errorf("%w: %v", biz.ErrWechatUnavailable, errors.New("response does not contain openid"))
	}
	return biz.WechatIdentity{Provider: provider, AppID: appID, OpenID: payload.OpenID, UnionID: payload.UnionID}, nil
}

func isInvalidWechatCode(code int) bool {
	switch code {
	case 40029, 40163, 40226:
		return true
	default:
		return false
	}
}
