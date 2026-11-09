package testkit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

const DefaultBaseURL = "http://127.0.0.1:8000"

var httpClient = &http.Client{Timeout: 5 * time.Second}

type Envelope struct {
	Code      string          `json:"code"`
	Message   string          `json:"message"`
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
	if value := strings.TrimSpace(os.Getenv("DD_TEST_BASE_URL")); value != "" {
		return strings.TrimRight(value, "/")
	}
	return DefaultBaseURL
}

func Request(t *testing.T, method, path string, body io.Reader) Response {
	return RequestWithHeaders(t, method, path, body, nil)
}

// RequestWithHeaders 请求已经启动的真实服务。每次调用都完整校验统一响应格式，
// 因此任意顶层 TestXxx 都可以从 IDE 的运行按钮单独执行。
func RequestWithHeaders(t *testing.T, method, path string, body io.Reader, headers http.Header) Response {
	t.Helper()
	target := BaseURL() + "/" + strings.TrimLeft(path, "/")
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
		t.Fatalf("请求 %s 失败: %v；请先在 back 目录运行 go run ./cmd/server；如果服务不是默认地址，请设置 DD_TEST_BASE_URL", target, err)
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
	t.Logf("真实服务响应: status=%d code=%s request_id=%s", result.StatusCode, result.Envelope.Code, result.Envelope.RequestID)
	return result
}

// RequestJSON 编码 JSON 请求体并设置 Content-Type，供需要请求体的 case 直接使用。
func RequestJSON(t *testing.T, method, path string, payload any) Response {
	t.Helper()
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("编码 JSON 请求失败: %v", err)
	}
	header := make(http.Header)
	header.Set("Content-Type", "application/json")
	return RequestWithHeaders(t, method, path, bytes.NewReader(body), header)
}

func Require(t *testing.T, response Response, status int, code string) {
	t.Helper()
	if response.StatusCode != status || response.Envelope.Code != code {
		t.Fatalf("接口结果不符合预期: want_status=%d want_code=%q got_status=%d got_code=%q body=%s",
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
	return fmt.Sprintf("%s/%s", BaseURL(), strings.TrimLeft(path, "/"))
}
