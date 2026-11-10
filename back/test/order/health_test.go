package order_test

import (
	"dd_web/back/internal/ecode"
	"dd_web/back/test/testkit"
	"net/http"
	"testing"
)

func TestOrderServiceHealth(t *testing.T) {
	response := testkit.RequestOrder(t, http.MethodGet, "/healthz", nil)
	testkit.Require(t, response, http.StatusOK, ecode.Success)
}
