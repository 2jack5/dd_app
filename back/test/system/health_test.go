package system_test

import (
	"dd_web/back/internal/ecode"
	"dd_web/back/test/testkit"
	"net/http"
	"testing"
)

func TestHealthEndpoint(t *testing.T) {
	response := testkit.Request(t, http.MethodGet, "/healthz", nil)
	testkit.Require(t, response, http.StatusOK, ecode.Success)

	var data struct {
		Status string `json:"status"`
	}
	testkit.DecodeData(t, response, &data)
	if data.Status != "alive" {
		t.Fatalf("健康状态错误: %q", data.Status)
	}
}

func TestReadinessEndpoint(t *testing.T) {
	response := testkit.Request(t, http.MethodGet, "/readyz", nil)
	switch response.StatusCode {
	case http.StatusOK:
		testkit.Require(t, response, http.StatusOK, ecode.Success)
	case http.StatusInternalServerError:
		testkit.Require(t, response, http.StatusInternalServerError, ecode.InternalError)
	default:
		t.Fatalf("就绪接口返回了未定义状态: status=%d body=%s", response.StatusCode, response.Body)
	}
}
