package infra

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"dd_web/back/internal/biz"
	"dd_web/back/internal/config"
)

func TestIdentityClientCallsConfiguredProvider(t *testing.T) {
	client := &IdentityClient{provider: "GENERIC_HTTP", verifyURL: "https://identity.example/verify", appID: "app", apiKey: "secret"}
	client.http = &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.Header.Get("Authorization") != "Bearer secret" || r.Header.Get("X-App-ID") != "app" {
			t.Fatalf("missing provider credentials")
		}
		var body identityVerifyRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body.RealName != "张三" || body.IdentityNumber != "11010519491231002X" {
			t.Fatalf("unexpected provider request: %+v", body)
		}
		return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(`{"request_id":"request-1","matched":true}`)), Header: make(http.Header)}, nil
	})}
	result, err := client.Verify(context.Background(), biz.IdentityVerifyInput{Provider: "GENERIC_HTTP", RealName: "张三", IdentityNumber: "11010519491231002X"})
	if err != nil || !result.Matched || result.RequestID != "request-1" {
		t.Fatalf("result=%+v err=%v", result, err)
	}
}

func TestIdentityClientRequiresConfiguration(t *testing.T) {
	_, err := NewIdentityClient(config.Identity{}).Verify(context.Background(), biz.IdentityVerifyInput{})
	if err == nil {
		t.Fatal("unconfigured identity client was accepted")
	}
}
