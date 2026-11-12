package server

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"strings"
	"testing"

	"dd_web/back/internal/biz"
	"dd_web/back/internal/config"
	"dd_web/back/internal/ecode"
	"dd_web/back/internal/model"
	"dd_web/back/internal/service"
)

type routeContract struct {
	Method, Path, Operation, Auth, Permission string
	Idempotent, Implemented                   bool
	Body                                      any
}
type fakeRepo struct{}

func (fakeRepo) Ping(context.Context) error                              { return nil }
func (fakeRepo) Games(context.Context, int64, int) ([]model.Game, error) { return nil, nil }
func (fakeRepo) Regions(context.Context, int64, int64, int) ([]model.GameRegion, error) {
	return nil, nil
}
func (fakeRepo) Ranks(context.Context, int64, int64, int) ([]model.GameRank, error) { return nil, nil }
func (fakeRepo) Categories(context.Context, int64, int64, int) ([]model.ServiceCategory, error) {
	return nil, nil
}

type allowAuth struct{ permissions map[string]bool }

func (a allowAuth) Authenticate(context.Context, string, string) (biz.Principal, error) {
	return biz.Principal{UserID: 1, AdminID: 1, IsCoach: true, Permissions: a.permissions}, nil
}

func testServers(t *testing.T, auth biz.Authenticator) (http.Handler, http.Handler) {
	t.Helper()
	t.Setenv("DD_BUSINESS_HTTP_ADDR", "127.0.0.1:0")
	t.Setenv("DD_ORDER_HTTP_ADDR", "127.0.0.1:0")
	repo := fakeRepo{}
	system := biz.NewSystemUsecase(repo)
	business, err := NewBusinessHTTPServer(&config.Config{}, service.NewBusinessService(system, biz.NewCatalogUsecase(repo), nil, nil), auth)
	if err != nil {
		t.Fatal(err)
	}
	order, err := NewOrderHTTPServer(&config.Config{}, service.NewOrderService(system), auth)
	if err != nil {
		t.Fatal(err)
	}
	return business, order
}
func readContract(t *testing.T) []routeContract {
	t.Helper()
	raw, err := os.ReadFile("../../api/v1/routes.json")
	if err != nil {
		t.Fatal(err)
	}
	var routes []routeContract
	if err := json.Unmarshal(raw, &routes); err != nil {
		t.Fatal(err)
	}
	return routes
}
func invoke(h http.Handler, method, path, token, key string, hasBody bool) *httptest.ResponseRecorder {
	body := bytes.NewReader(nil)
	if hasBody {
		body = bytes.NewReader([]byte(`{}`))
	}
	req := httptest.NewRequest(method, path, body)
	if hasBody {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	if key != "" {
		req.Header.Set("Idempotency-Key", key)
	}
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	return w
}

var orderOperations = map[string]bool{
	"coach.orders": true, "coach.orderDetail": true, "coach.orderEvents": true, "coach.acceptOrder": true, "coach.rejectOrder": true, "coach.startOrder": true, "coach.completeOrder": true, "coach.earningsSummary": true, "coach.earningsItems": true, "coach.settlements": true, "coach.settlementDetail": true, "coach.supplements": true, "coach.settlementAccount": true, "coach.onboarding": true,
	"admin.orders": true, "admin.orderDetail": true, "admin.orderEvents": true, "admin.refunds": true, "admin.refundDetail": true, "admin.decideRefund": true, "admin.cases": true, "admin.caseDetail": true, "admin.assignCase": true, "admin.caseMessage": true, "admin.resolveCase": true, "admin.reopenCase": true, "admin.dashboard": true, "admin.payments": true, "admin.settlementItems": true, "admin.batches": true, "admin.createBatch": true, "admin.adjustments": true, "admin.createAdjustment": true, "admin.supplements": true, "admin.createSupplement": true, "admin.retrySupplement": true, "admin.createReconciliation": true, "admin.reconciliationRuns": true, "admin.reconciliationDetail": true, "admin.reconciliationDiffs": true, "admin.resolveDiff": true, "admin.commissionRules": true, "admin.createCommission": true, "admin.publishCommission": true, "admin.businessRules": true, "admin.createBusinessRule": true, "admin.publishBusinessRule": true,
	"callback.payment": true, "callback.refund": true, "callback.settlement": true,
	"callback.walletRecharge":      true,
	"admin.walletRechargePackages": true, "admin.createRechargePackage": true, "admin.updateRechargePackage": true,
}

func isOrder(op string) bool {
	prefix := strings.Split(op, ".")[0]
	return orderOperations[op] || prefix == "trade" || prefix == "payment" || prefix == "refund" || prefix == "review" || prefix == "case" || prefix == "wallet"
}

func TestGeneratedRoutesAreSplitAcrossServices(t *testing.T) {
	routes := readContract(t)
	if len(routes) != 190 {
		t.Fatalf("路由契约数量错误: %d", len(routes))
	}
	permissions := map[string]bool{}
	for _, r := range routes {
		if r.Permission != "" {
			permissions[r.Permission] = true
		}
	}
	business, order := testServers(t, allowAuth{permissions})
	replace := regexp.MustCompile(`\{[^}]+\}`)
	for _, r := range routes {
		t.Run(r.Operation, func(t *testing.T) {
			h := business
			if isOrder(r.Operation) {
				h = order
			}
			path := replace.ReplaceAllString(r.Path, "1")
			if r.Operation == "catalog.categories" {
				path += "?game_id=1"
			}
			w := invoke(h, r.Method, path, "token", "test-key", r.Body != nil)
			if w.Code == 404 || w.Code == 405 {
				t.Fatalf("目标微服务未注册路由: %d %s", w.Code, w.Body.String())
			}
			// 每个微服务都必须有自己的存活和就绪探针。
			if r.Operation == "system.health" || r.Operation == "system.ready" {
				return
			}
			other := order
			if isOrder(r.Operation) {
				other = business
			}
			if got := invoke(other, r.Method, path, "token", "test-key", r.Body != nil); got.Code != 404 {
				t.Fatalf("路由同时暴露在错误微服务: %d %s", got.Code, got.Body.String())
			}
		})
	}
}
func TestBothServicesExposeHealth(t *testing.T) {
	business, order := testServers(t, nil)
	for name, h := range map[string]http.Handler{"business": business, "order": order} {
		if got := invoke(h, "GET", "/healthz", "", "", false); got.Code != 200 {
			t.Errorf("%s health: %d", name, got.Code)
		}
	}
}
func TestGeneratedMiddlewareProtectsRoutes(t *testing.T) {
	business, order := testServers(t, nil)
	for _, tc := range []struct {
		h    http.Handler
		path string
	}{{business, "/api/v1/users/me"}, {order, "/api/v1/orders"}} {
		if got := invoke(tc.h, "GET", tc.path, "", "", false); got.Code != 401 {
			t.Fatalf("匿名访问状态: %d", got.Code)
		}
	}
}
func TestGeneratedCatalogValidation(t *testing.T) {
	business, _ := testServers(t, allowAuth{map[string]bool{}})
	for _, path := range []string{"/api/v1/games?page_size=101", "/api/v1/games?cursor=-1", "/api/v1/games/not-id/regions", "/api/v1/service-categories"} {
		if got := invoke(business, "GET", path, "", "", false); got.Code != 400 {
			t.Errorf("%s: %d %s", path, got.Code, got.Body.String())
		}
	}
}

func TestUnifiedBusinessEnvelope(t *testing.T) {
	business, _ := testServers(t, nil)
	for _, tc := range []struct {
		name, path string
		status     int
		code       ecode.Code
		errorMsg   string
	}{
		{"success", "/healthz", http.StatusOK, ecode.Success, ""},
		{"not-found", "/__missing__", http.StatusNotFound, ecode.NotFound, ecode.Message(ecode.NotFound)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			response := invoke(business, http.MethodGet, tc.path, "", "", false)
			var body struct {
				Code      int32  `json:"code"`
				ErrorMsg  string `json:"error_msg"`
				RequestID string `json:"request_id"`
			}
			if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
				t.Fatal(err)
			}
			if response.Code != tc.status || body.Code != int32(tc.code) || body.ErrorMsg != tc.errorMsg || body.RequestID == "" {
				t.Fatalf("unexpected envelope: status=%d body=%s", response.Code, response.Body.String())
			}
		})
	}
}

func TestUnknownErrorUsesInternalCode5000(t *testing.T) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	encodeError(recorder, request, errors.New("database details must not leak"))
	var body map[string]any
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if recorder.Code != http.StatusInternalServerError || body["code"] != float64(ecode.InternalError) || body["error_msg"] != ecode.Message(ecode.InternalError) {
		t.Fatalf("unexpected internal envelope: status=%d body=%s", recorder.Code, recorder.Body.String())
	}
	if strings.Contains(recorder.Body.String(), "database details") {
		t.Fatal("internal error leaked into response")
	}
}
