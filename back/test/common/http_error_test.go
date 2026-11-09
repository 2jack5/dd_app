package common_test

import (
	"dd_web/back/test/testkit"
	"net/http"
	"testing"
)

func TestMissingRouteReturnsNotFound(t *testing.T) {
	response := testkit.Request(t, http.MethodGet, "/__test_missing_route__", nil)
	testkit.Require(t, response, http.StatusNotFound, "NOT_FOUND")
}

func TestUnsupportedMethodReturnsMethodNotAllowed(t *testing.T) {
	response := testkit.Request(t, http.MethodDelete, "/healthz", nil)
	testkit.Require(t, response, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED")
}
