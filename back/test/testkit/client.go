package testkit

import (
	"bytes"
	"dd_web/back/internal/ecode"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

const (
	DefaultBusinessBaseURL = "http://127.0.0.1:8000"
	DefaultOrderBaseURL    = "http://127.0.0.1:8001"
)

var httpClient = &http.Client{Timeout: 5 * time.Second}

type Envelope struct {
	Code      int32           `json:"code"`
	ErrorMsg  string          `json:"error_msg"`
	Data      json.RawMessage `json:"data"`
	RequestID string          `json:"request_id"`
}

type Response struct {
	StatusCode int
	Header     http.Header
	Body       []byte
	Envelope   Envelope
}

func BaseURL() string {
	return BusinessBaseURL()
}

func BusinessBaseURL() string {
	if value := strings.TrimSpace(os.Getenv("DD_TEST_BUSINESS_BASE_URL")); value != "" {
		return strings.TrimRight(value, "/")
	}
	if value := strings.TrimSpace(os.Getenv("DD_TEST_BASE_URL")); value != "" {
		return strings.TrimRight(value, "/")
	}
	return DefaultBusinessBaseURL
}

func OrderBaseURL() string {
	if value := strings.TrimSpace(os.Getenv("DD_TEST_ORDER_BASE_URL")); value != "" {
		return strings.TrimRight(value, "/")
	}
	return DefaultOrderBaseURL
}

func Request(t *testing.T, method, path string, body io.Reader) Response {
	return RequestWithHeaders(t, method, path, body, nil)
}

// RequestWithHeaders 请求已经启动的真实服务。每次调用都完整校验统一响应格式，
// 因此任意顶层 TestXxx 都可以从 IDE 的运行按钮单独执行。
func RequestWithHeaders(t *testing.T, method, path string, body io.Reader, headers http.Header) Response {
	return requestWithHeaders(t, BusinessBaseURL(), "business", method, path, body, headers)
}

func RequestOrder(t *testing.T, method, path string, body io.Reader) Response {
	return requestWithHeaders(t, OrderBaseURL(), "order", method, path, body, nil)
}

func RequestOrderWithHeaders(t *testing.T, method, path string, body io.Reader, headers http.Header) Response {
	return requestWithHeaders(t, OrderBaseURL(), "order", method, path, body, headers)
}

func requestWithHeaders(t *testing.T, baseURL, service, method, path string, body io.Reader, headers http.Header) Response {
	t.Helper()
	target := baseURL + "/" + strings.TrimLeft(path, "/")
	req, err := http.NewRequest(method, target, body)
	if err != nil {
		t.Fatalf("创建请求失败: %v", err)
	}
	req.Header.Set("Accept", "application/json")
	for key, values := range headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	t.Logf("请求真实服务: %s %s", method, target)
	resp, err := httpClient.Do(req)
	if err != nil {
		t.Fatalf("请求 %s 失败: %v；请先启动 %s 服务；非默认地址请设置对应的 DD_TEST_*_BASE_URL", target, err, service)
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("读取响应失败: %v", err)
	}

	result := Response{StatusCode: resp.StatusCode, Header: resp.Header.Clone(), Body: raw}
	if err := json.Unmarshal(raw, &result.Envelope); err != nil {
		t.Fatalf("接口未返回统一 JSON 响应: status=%d body=%q error=%v", resp.StatusCode, raw, err)
	}
	if result.Envelope.RequestID == "" {
		t.Fatalf("响应缺少 request_id: %s", raw)
	}
	if headerID := resp.Header.Get("X-Request-ID"); headerID != result.Envelope.RequestID {
		t.Fatalf("请求 ID 不一致: header=%q body=%q", headerID, result.Envelope.RequestID)
	}
	if result.Envelope.Code == int32(ecode.Success) && result.Envelope.ErrorMsg != "" {
		t.Fatalf("成功响应的error_msg必须为空: %s", raw)
	}
	if result.Envelope.Code != int32(ecode.Success) && result.Envelope.ErrorMsg == "" {
		t.Fatalf("失败响应必须包含error_msg: %s", raw)
	}
	t.Logf("真实服务响应: status=%d code=%d request_id=%s", result.StatusCode, result.Envelope.Code, result.Envelope.RequestID)
	return result
}

// RequestJSON 编码 JSON 请求体并设置 Content-Type，供需要请求体的 case 直接使用。
func RequestJSON(t *testing.T, method, path string, payload any) Response {
	return requestJSON(t, false, method, path, payload)
}

func RequestOrderJSON(t *testing.T, method, path string, payload any) Response {
	return requestJSON(t, true, method, path, payload)
}

func requestJSON(t *testing.T, order bool, method, path string, payload any) Response {
	t.Helper()
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("编码 JSON 请求失败: %v", err)
	}
	header := make(http.Header)
	header.Set("Content-Type", "application/json")
	if order {
		return RequestOrderWithHeaders(t, method, path, bytes.NewReader(body), header)
	}
	return RequestWithHeaders(t, method, path, bytes.NewReader(body), header)
}

func Require(t *testing.T, response Response, status int, code ecode.Code) {
	t.Helper()
	if response.StatusCode != status || response.Envelope.Code != int32(code) {
		t.Fatalf("接口结果不符合预期: want_status=%d want_code=%d got_status=%d got_code=%d body=%s",
			status, code, response.StatusCode, response.Envelope.Code, response.Body)
	}
}

func DecodeData(t *testing.T, response Response, target any) {
	t.Helper()
	if err := json.Unmarshal(response.Envelope.Data, target); err != nil {
		t.Fatalf("解析 data 失败: %v; data=%s", err, response.Envelope.Data)
	}
}

func URL(path string) string {
	return fmt.Sprintf("%s/%s", BusinessBaseURL(), strings.TrimLeft(path, "/"))
}
