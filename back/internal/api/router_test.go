package api

import (
	"context"
	"dd_web/back/internal/config"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"testing"
	"time"

	khttp "github.com/go-kratos/kratos/v3/transport/http"
)

type contractRoute struct {
	Method, Path, Operation, Auth, Permission string
	Idempotent, Implemented                   bool
}

func readContract(t *testing.T) []contractRoute {
	t.Helper()
	raw, e := os.ReadFile("../../api/v1/routes.json")
	if e != nil {
		t.Fatal(e)
	}
	var r []contractRoute
	if e = json.Unmarshal(raw, &r); e != nil {
		t.Fatal(e)
	}
	return r
}
func testRouter(t *testing.T, auth Authenticator) *khttp.Server {
	t.Helper()
	s, e := newRouter(routerConfig{address: "127.0.0.1:0", timeout: time.Second, auth: auth})
	if e != nil {
		t.Fatal(e)
	}
	return s
}
func invoke(s http.Handler, method, path, token, key string) *httptest.ResponseRecorder {
	r := httptest.NewRequest(method, path, nil)
	if token != "" {
		r.Header.Set("Authorization", "Bearer "+token)
	}
	if key != "" {
		r.Header.Set("Idempotency-Key", key)
	}
	w := httptest.NewRecorder()
	s.ServeHTTP(w, r)
	return w
}
func TestExplicitRoutesMatchOpenAPI(t *testing.T) {
	s := testRouter(t, nil)
	actual := map[string]bool{}
	if e := s.WalkRoute(func(info khttp.RouteInfo) error {
		key := info.Method + " " + info.Path
		if actual[key] {
			t.Errorf("duplicate registered route: %s", key)
		}
		actual[key] = true
		return nil
	}); e != nil {
		t.Fatal(e)
	}
	raw, e := os.ReadFile("../../api/v1/openapi.json")
	if e != nil {
		t.Fatal(e)
	}
	var doc struct {
		Paths map[string]map[string]json.RawMessage `json:"paths"`
	}
	if e = json.Unmarshal(raw, &doc); e != nil {
		t.Fatal(e)
	}
	for path, methods := range doc.Paths {
		for method := range methods {
			key := stringUpper(method) + " " + path
			if !actual[key] {
				t.Errorf("documented route not registered: %s", key)
			}
			delete(actual, key)
		}
	}
	if len(actual) > 0 {
		t.Fatalf("undocumented routes: %v", actual)
	}
}
func stringUpper(s string) string {
	b := []byte(s)
	for i, c := range b {
		if c >= 'a' && c <= 'z' {
			b[i] = c - 32
		}
	}
	return string(b)
}

type allowAuth struct{ permissions map[string]bool }

func (a allowAuth) Authenticate(_ context.Context, _ string, scope string) (Principal, error) {
	if scope == "admin" {
		return Principal{AdminID: 1, Permissions: a.permissions}, nil
	}
	return Principal{UserID: 1, IsCoach: true}, nil
}
func TestEveryRouteReachesItsOwnHandler(t *testing.T) {
	contract := readContract(t)
	permissions := map[string]bool{}
	for _, r := range contract {
		if r.Permission != "" {
			permissions[r.Permission] = true
		}
	}
	s := testRouter(t, allowAuth{permissions})
	replace := regexp.MustCompile(`\{[^}]+\}`)
	for _, r := range contract {
		t.Run(r.Operation, func(t *testing.T) {
			path := replace.ReplaceAllString(r.Path, "1")
			if r.Operation == "catalog.categories" {
				path += "?game_id=1"
			}
			w := invoke(s, r.Method, path, "test-token", "test-key")
			var got response
			if e := json.Unmarshal(w.Body.Bytes(), &got); e != nil {
				t.Fatal(e)
			}
			if got.RequestID == "" || got.RequestID != w.Header().Get("X-Request-ID") {
				t.Fatal("inconsistent request ID")
			}
			switch {
			case r.Operation == "system.health":
				if w.Code != 200 {
					t.Fatalf("health: %d", w.Code)
				}
			case r.Operation == "system.ready":
				if w.Code != 503 {
					t.Fatalf("readiness without DB: %d", w.Code)
				}
			case r.Implemented:
				if w.Code != 503 || got.Code != "DATABASE_UNAVAILABLE" {
					t.Fatalf("catalog without DB: %d %s", w.Code, w.Body.String())
				}
			case r.Auth == "channel":
				if w.Code != 503 || got.Code != "CHANNEL_NOT_CONFIGURED" {
					t.Fatalf("callback returned fake success: %d", w.Code)
				}
			default:
				if w.Code != 501 || got.Code != "NOT_IMPLEMENTED" {
					t.Fatalf("unimplemented handler: %d %s", w.Code, w.Body.String())
				}
				data, ok := got.Data.(map[string]any)
				if !ok || data["operation"] != r.Operation {
					t.Fatalf("route reached wrong handler: %s", w.Body.String())
				}
			}
		})
	}
}
func TestProtectedRoutesRemainClosed(t *testing.T) {
	s := testRouter(t, nil)
	for _, r := range readContract(t) {
		if r.Auth == "public" || r.Auth == "channel" {
			continue
		}
		path := regexp.MustCompile(`\{[^}]+\}`).ReplaceAllString(r.Path, "1")
		if w := invoke(s, r.Method, path, "", "key"); w.Code != 401 {
			t.Errorf("anonymous protected route %s: %d", r.Operation, w.Code)
		}
		if w := invoke(s, r.Method, path, "unverified", "key"); w.Code != 503 {
			t.Errorf("unconfigured auth %s: %d", r.Operation, w.Code)
		}
	}
	s = testRouter(t, allowAuth{map[string]bool{}})
	if w := invoke(s, "GET", "/api/v1/admin/users", "token", ""); w.Code != 403 {
		t.Fatalf("missing admin permission: %d", w.Code)
	}
	if w := invoke(s, "POST", "/api/v1/orders", "token", ""); w.Code != 400 {
		t.Fatalf("missing idempotency key: %d", w.Code)
	}
}
func TestCatalogAndHTTPValidation(t *testing.T) {
	s := testRouter(t, nil)
	for _, tc := range []struct {
		method, path string
		status       int
	}{
		{"GET", "/api/v1/games?page_size=101", 400}, {"GET", "/api/v1/games?cursor=-1", 400},
		{"GET", "/api/v1/games/not-id/regions", 400}, {"GET", "/api/v1/service-categories", 400},
		{"GET", "/missing", 404}, {"DELETE", "/healthz", 405},
	} {
		w := invoke(s, tc.method, tc.path, "", "")
		if w.Code != tc.status {
			t.Errorf("%s got %d", tc.path, w.Code)
		}
	}
	// Cursor truncation must use the last emitted row, not the extra lookahead row.
	rows, next := trimCatalog([]int64{3, 4}, catalogPage{limit: 1}, func(n int64) int64 { return n })
	if len(rows) != 1 || rows[0] != 3 || next != "3" {
		t.Fatal("cursor skipped a row")
	}
}

func TestRouterAcceptsLoadedConfiguration(t *testing.T) {
	t.Setenv("DD_DATABASE_URL", "")
	t.Setenv("DD_HTTP_ADDR", "")
	t.Setenv("DD_HTTP_TIMEOUT", "")
	t.Setenv("DD_DB_MAX_OPEN", "")
	t.Setenv("DD_DB_MAX_IDLE", "")
	srv, err := NewRouter(&config.Config{HTTP: config.HTTP{Port: "9123"}}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if w := invoke(srv, "GET", "/healthz", "", ""); w.Code != 200 {
		t.Fatalf("configured router did not respond: %d", w.Code)
	}
}
