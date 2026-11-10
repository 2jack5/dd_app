package access_test

import (
	"dd_web/back/internal/ecode"
	"dd_web/back/test/testkit"
	"net/http"
	"testing"
)

func TestUserEndpointRequiresAuthentication(t *testing.T) {
	response := testkit.Request(t, http.MethodGet, "/api/v1/users/me", nil)
	testkit.Require(t, response, http.StatusUnauthorized, ecode.AuthRequired)
}

func TestCoachEndpointRequiresAuthentication(t *testing.T) {
	response := testkit.Request(t, http.MethodGet, "/api/v1/coach/me", nil)
	testkit.Require(t, response, http.StatusUnauthorized, ecode.AuthRequired)
}

func TestAdminEndpointRequiresAuthentication(t *testing.T) {
	response := testkit.Request(t, http.MethodGet, "/api/v1/admin/me", nil)
	testkit.Require(t, response, http.StatusUnauthorized, ecode.AuthRequired)
}
