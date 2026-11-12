package server

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"dd_web/back/internal/biz"
	"dd_web/back/internal/config"
	"dd_web/back/internal/ecode"
	"dd_web/back/internal/infra"
	"dd_web/back/internal/model"
	"dd_web/back/internal/service"

	"golang.org/x/crypto/bcrypt"
)

// 使用真实 HTTP 路由、业务用例、bcrypt 和令牌签发，替换外部存储。
type adminPasswordRepo struct {
	users  map[int64]*model.AdminUser
	roles  map[int64][]int64
	nextID int64
}

func (r *adminPasswordRepo) FindByAccount(_ context.Context, account string) (*model.AdminUser, error) {
	for _, user := range r.users {
		if user.Account == account {
			return user, nil
		}
	}
	return nil, biz.ErrAdminNotFound
}
func (r *adminPasswordRepo) FindByID(_ context.Context, id int64) (*model.AdminUser, error) {
	if user, ok := r.users[id]; ok {
		return user, nil
	}
	return nil, biz.ErrAdminNotFound
}
func (r *adminPasswordRepo) Access(_ context.Context, id int64) (map[string]bool, string, error) {
	permissions := map[string]bool{}
	for _, role := range r.roles[id] {
		if role == 1 {
			permissions["admin.account.write"] = true
		}
	}
	return permissions, "ALL", nil
}
func (r *adminPasswordRepo) CreateAccount(ctx context.Context, admin *model.AdminUser, roleIDs []int64) error {
	for _, roleID := range roleIDs {
		if roleID != 1 {
			return biz.ErrAdminRoleInvalid
		}
	}
	if _, err := r.FindByAccount(ctx, admin.Account); err == nil {
		return biz.ErrAdminAccountExists
	}
	r.nextID++
	admin.ID = r.nextID
	r.users[admin.ID] = admin
	r.roles[admin.ID] = append([]int64(nil), roleIDs...)
	return nil
}

type adminPasswordSessions struct {
	saved biz.RefreshSession
	err   error
}

func (s *adminPasswordSessions) SaveRefreshSession(_ context.Context, value biz.RefreshSession) error {
	if s.err != nil {
		return s.err
	}
	s.saved = value
	return nil
}
func (s *adminPasswordSessions) RotateRefreshSession(context.Context, string, biz.RefreshSession) error {
	return s.err
}

func adminPasswordRequest(t *testing.T, h http.Handler, path string, body any, token, key string) (*httptest.ResponseRecorder, map[string]any) {
	t.Helper()
	raw, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(string(raw)))
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	if key != "" {
		req.Header.Set("Idempotency-Key", key)
	}
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	var reply map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &reply); err != nil {
		t.Fatal(err)
	}
	return w, reply
}

func TestAdminCreateThenPasswordLoginWithoutMFA(t *testing.T) {
	t.Setenv("DD_BUSINESS_HTTP_ADDR", "127.0.0.1:0")
	repo := &adminPasswordRepo{
		users: map[int64]*model.AdminUser{
			1: {Base: model.Base{ID: 1}, Account: "root", Status: model.AdminUserStatusActive, AuthVersion: 1},
			2: {Base: model.Base{ID: 2}, Account: "viewer", Status: model.AdminUserStatusActive, AuthVersion: 1},
		},
		roles:  map[int64][]int64{1: {1}},
		nextID: 2,
	}
	sessions := &adminPasswordSessions{}
	tokenCfg := config.Auth{Issuer: "admin-test/admin", TokenSecret: strings.Repeat("test-key", 8), AccessTokenTTL: time.Hour, RefreshTokenTTL: 24 * time.Hour}
	tokens := infra.NewTokenManager(tokenCfg)
	admin := biz.NewAdminUsecase(repo, tokens, sessions)
	svc := service.NewBusinessService(nil, nil, nil, admin)
	h, err := NewBusinessHTTPServer(&config.Config{}, svc, biz.ScopedAuthenticator{Admin: admin})
	if err != nil {
		t.Fatal(err)
	}
	rootTokens, err := tokens.Issue(1, 1)
	if err != nil {
		t.Fatal(err)
	}
	viewerTokens, err := tokens.Issue(2, 1)
	if err != nil {
		t.Fatal(err)
	}
	password := " Password123! "
	create := map[string]any{"account": " Operator01 ", "initial_password": password, "display_name": " 运营人员 ", "role_ids": []string{"1", "1"}}
	createPath := "/api/v1/admin/accounts"
	loginPath := "/api/v1/admin/auth/login"

	for _, tc := range []struct {
		name, token, key string
		status           int
	}{
		{"anonymous", "", "key", http.StatusUnauthorized},
		{"no permission", viewerTokens.AccessToken, "key", http.StatusForbidden},
		{"missing idempotency key", rootTokens.AccessToken, "", http.StatusBadRequest},
	} {
		t.Run(tc.name, func(t *testing.T) {
			w, _ := adminPasswordRequest(t, h, createPath, create, tc.token, tc.key)
			if w.Code != tc.status {
				t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
			}
			if len(repo.users) != 2 {
				t.Fatal("unauthorized request created an account")
			}
		})
	}

	w, reply := adminPasswordRequest(t, h, createPath, create, rootTokens.AccessToken, "create-1")
	if w.Code != http.StatusOK || reply["code"] != float64(ecode.Success) {
		t.Fatalf("create: %d %s", w.Code, w.Body.String())
	}
	data := reply["data"].(map[string]any)
	id, err := strconv.ParseInt(data["id"].(string), 10, 64)
	if err != nil {
		t.Fatal(err)
	}
	created := repo.users[id]
	if data["status"] != "ACTIVE" || data["account"] != "operator01" || created.DisplayName != "运营人员" || len(created.MFASecretCipher) != 0 {
		t.Fatalf("unexpected created account: %#v", data)
	}
	if len(repo.roles[id]) != 1 {
		t.Fatalf("role IDs were not deduplicated: %v", repo.roles[id])
	}
	if created.PasswordHash == password || bcrypt.CompareHashAndPassword([]byte(created.PasswordHash), []byte(password)) != nil {
		t.Fatal("password was not stored as a matching bcrypt hash")
	}
	for _, secret := range []string{password, created.PasswordHash, "password_hash", "mfa_secret"} {
		if strings.Contains(w.Body.String(), secret) {
			t.Fatal("creation response exposed credentials")
		}
	}

	w, reply = adminPasswordRequest(t, h, createPath, create, rootTokens.AccessToken, "create-duplicate")
	if w.Code != http.StatusConflict || reply["code"] != float64(ecode.AdminAccountExists) {
		t.Fatalf("duplicate: %d %s", w.Code, w.Body.String())
	}

	for _, tc := range []struct{ name, account, password string }{
		{"wrong password", "operator01", "wrong-password"},
		{"unknown account", "unknown", password},
		{"password whitespace matters", "operator01", strings.TrimSpace(password)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			w, reply := adminPasswordRequest(t, h, loginPath, map[string]any{"account": tc.account, "password": tc.password}, "", "")
			if w.Code != http.StatusUnauthorized || reply["code"] != float64(ecode.InvalidCredentials) {
				t.Fatalf("login: %d %s", w.Code, w.Body.String())
			}
		})
	}

	// 请求没有 mfa_code，账号也没有 MFASecretCipher。
	w, reply = adminPasswordRequest(t, h, loginPath, map[string]any{"account": " OPERATOR01 ", "password": password}, "", "")
	if w.Code != http.StatusOK {
		t.Fatalf("password login: %d %s", w.Code, w.Body.String())
	}
	data = reply["data"].(map[string]any)
	access := data["access_token"].(string)
	if data["refresh_token"] == "" || data["token_type"] != "Bearer" {
		t.Fatalf("missing tokens: %#v", data)
	}
	claims, err := tokens.VerifyAccess(access)
	if err != nil || claims.UserID != id || claims.AuthVersion != 1 {
		t.Fatalf("invalid admin token: %+v %v", claims, err)
	}
	if sessions.saved.UserID != id || sessions.saved.ClientType != "ADMIN" || sessions.saved.RefreshTokenHash == "" {
		t.Fatalf("admin refresh session missing: %+v", sessions.saved)
	}
	userCfg := tokenCfg
	userCfg.Issuer = "admin-test"
	if _, err := infra.NewTokenManager(userCfg).VerifyAccess(access); err == nil {
		t.Fatal("admin token accepted as a user token")
	}

	sessions.err = errors.New("private redis failure")
	w, reply = adminPasswordRequest(t, h, loginPath, map[string]any{"account": "operator01", "password": password}, "", "")
	if w.Code != http.StatusInternalServerError || reply["code"] != float64(ecode.InternalError) || reply["data"] != nil || strings.Contains(w.Body.String(), "private redis failure") {
		t.Fatalf("session failure leaked success or internals: %d %s", w.Code, w.Body.String())
	}
	sessions.err = nil

	created.Status = model.AdminUserStatusDisabled
	w, reply = adminPasswordRequest(t, h, loginPath, map[string]any{"account": "operator01", "password": password}, "", "")
	if w.Code != http.StatusForbidden || reply["code"] != float64(ecode.AccountDisabled) {
		t.Fatalf("disabled login: %d %s", w.Code, w.Body.String())
	}
	w, _ = adminPasswordRequest(t, h, createPath, create, access, "disabled-token")
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("disabled account token accepted: %d", w.Code)
	}
	created.Status = model.AdminUserStatusActive
	created.AuthVersion++
	w, _ = adminPasswordRequest(t, h, createPath, create, access, "old-token")
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("old auth version token accepted: %d", w.Code)
	}

	invalidRole := map[string]any{"account": "other", "initial_password": password, "display_name": "Other", "role_ids": []string{"99"}}
	w, reply = adminPasswordRequest(t, h, createPath, invalidRole, rootTokens.AccessToken, "bad-role")
	if w.Code != http.StatusBadRequest || reply["code"] != float64(ecode.AdminRoleInvalid) || len(repo.users) != 3 {
		t.Fatalf("invalid role: %d %s", w.Code, w.Body.String())
	}
	shortPassword := map[string]any{"account": "other", "initial_password": "short", "display_name": "Other", "role_ids": []string{"1"}}
	w, reply = adminPasswordRequest(t, h, createPath, shortPassword, rootTokens.AccessToken, "bad-password")
	if w.Code != http.StatusBadRequest || reply["code"] != float64(ecode.InvalidArgument) || len(repo.users) != 3 {
		t.Fatalf("invalid password: %d %s", w.Code, w.Body.String())
	}
}
