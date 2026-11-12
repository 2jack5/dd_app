package infra

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"dd_web/back/internal/biz"
	"dd_web/back/internal/config"
)

type IdentityClient struct {
	provider  string
	verifyURL string
	appID     string
	apiKey    string
	http      *http.Client
}

func NewIdentityClient(cfg config.Identity) *IdentityClient {
	timeout := cfg.Timeout
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	return &IdentityClient{
		provider: strings.TrimSpace(strings.ToUpper(cfg.Provider)), verifyURL: strings.TrimSpace(cfg.VerifyURL),
		appID: strings.TrimSpace(cfg.AppID), apiKey: strings.TrimSpace(cfg.APIKey), http: &http.Client{Timeout: timeout},
	}
}

type identityVerifyRequest struct {
	RealName       string `json:"real_name"`
	IdentityNumber string `json:"identity_number"`
	ConsentVersion string `json:"consent_version"`
}

type identityVerifyResponse struct {
	RequestID string     `json:"request_id"`
	Matched   bool       `json:"matched"`
	ResultRef string     `json:"result_ref"`
	ExpiresAt *time.Time `json:"expires_at"`
}

func (c *IdentityClient) Verify(ctx context.Context, input biz.IdentityVerifyInput) (biz.IdentityProviderResult, error) {
	if c == nil || c.http == nil || c.provider == "" || c.verifyURL == "" || c.apiKey == "" {
		return biz.IdentityProviderResult{}, biz.ErrIdentityUnavailable
	}
	if strings.TrimSpace(strings.ToUpper(input.Provider)) != c.provider {
		return biz.IdentityProviderResult{}, biz.ErrInvalidArgument
	}
	body, err := json.Marshal(identityVerifyRequest{RealName: input.RealName, IdentityNumber: input.IdentityNumber, ConsentVersion: input.ConsentVersion})
	if err != nil {
		return biz.IdentityProviderResult{}, fmt.Errorf("encode identity request: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.verifyURL, bytes.NewReader(body))
	if err != nil {
		return biz.IdentityProviderResult{}, fmt.Errorf("create identity request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	if c.appID != "" {
		req.Header.Set("X-App-ID", c.appID)
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return biz.IdentityProviderResult{}, fmt.Errorf("%w: provider request failed", biz.ErrIdentityUnavailable)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return biz.IdentityProviderResult{}, fmt.Errorf("%w: provider returned status %d", biz.ErrIdentityUnavailable, resp.StatusCode)
	}
	var result identityVerifyResponse
	decoder := json.NewDecoder(io.LimitReader(resp.Body, 1<<20))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&result); err != nil {
		return biz.IdentityProviderResult{}, fmt.Errorf("%w: invalid provider response", biz.ErrIdentityUnavailable)
	}
	if strings.TrimSpace(result.RequestID) == "" {
		return biz.IdentityProviderResult{}, fmt.Errorf("%w: provider response missing request_id", biz.ErrIdentityUnavailable)
	}
	return biz.IdentityProviderResult{RequestID: result.RequestID, Matched: result.Matched, ResultRef: result.ResultRef, ExpiresAt: result.ExpiresAt}, nil
}
