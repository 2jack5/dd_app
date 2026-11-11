package biz

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"dd_web/back/internal/model"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrAdminNotFound      = errors.New("admin not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrAdminAccountExists = errors.New("admin account already exists")
	ErrAdminRoleInvalid   = errors.New("admin role is invalid")
)

// 管理员数据库层必须提供这些方法。
type AdminRepository interface {
	FindByAccount(context.Context, string) (*model.AdminUser, error)
	FindByID(context.Context, int64) (*model.AdminUser, error)
	Access(context.Context, int64) (map[string]bool, string, error)
	CreateAccount(ctx context.Context, admin *model.AdminUser, roleIDs []int64) error
}

type AdminLoginResult struct {
	Admin  *model.AdminUser
	Tokens TokenPair
}

type AdminUsecase struct {
	repo     AdminRepository
	tokens   TokenManager
	sessions SessionStore
}

func NewAdminUsecase(repo AdminRepository, tokens TokenManager, sessions SessionStore) *AdminUsecase {
	return &AdminUsecase{repo: repo, tokens: tokens, sessions: sessions}
}

func (u *AdminUsecase) FindByAccount(ctx context.Context, account string) (*model.AdminUser, error) {
	return u.repo.FindByAccount(ctx, strings.ToLower(strings.TrimSpace(account)))
}

// CreateAccount 创建启用状态的管理员账号及其角色关联。
func (u *AdminUsecase) CreateAccount(ctx context.Context, account, password, displayName string, rawRoleIDs []string) (*model.AdminUser, error) {
	account = strings.ToLower(strings.TrimSpace(account))
	displayName = strings.TrimSpace(displayName)
	// 初始密码为 8～72 字节；密码中的空格保持原样。
	if account == "" ||
		utf8.RuneCountInString(account) > 64 ||
		strings.ContainsRune(account, 0) ||
		displayName == "" ||
		utf8.RuneCountInString(displayName) > 64 ||
		strings.ContainsRune(displayName, 0) ||
		len(password) < 8 ||
		len(password) > 72 ||
		len(rawRoleIDs) == 0 {
		return nil, ErrInvalidArgument
	}
	if u.repo == nil {
		return nil, ErrAuthenticationUnavailable
	}
	// 检查管理员是否存在
	_, err := u.repo.FindByAccount(ctx, account)
	if err == nil {
		return nil, ErrAdminAccountExists
	}
	if !errors.Is(err, ErrAdminNotFound) {
		return nil, err
	}
	// 将角色 ID 转成整数，并去重。
	roleIDs := make([]int64, 0, len(rawRoleIDs))
	seen := make(map[int64]bool)
	for _, rawID := range rawRoleIDs {
		roleID, err := ParseID(rawID)
		if err != nil {
			return nil, ErrInvalidArgument
		}

		if !seen[roleID] {
			seen[roleID] = true
			roleIDs = append(roleIDs, roleID)
		}
	}
	// 密码转为 bcrypt 哈希。
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash admin password: %w", err)
	}
	admin := &model.AdminUser{
		Account:      account,
		PasswordHash: string(hash),
		DisplayName:  displayName,
		Status:       model.AdminUserStatusActive,
		AuthVersion:  1,
	}

	if err := u.repo.CreateAccount(ctx, admin, roleIDs); err != nil {
		return nil, err
	}

	return admin, nil
}

// Login 校验管理员账号、密码和状态，签发令牌并保存刷新会话。
func (u *AdminUsecase) Login(ctx context.Context, account, password string) (AdminLoginResult, error) {
	account = strings.TrimSpace(account)
	// 密码保持原样，不去除空格。
	if account == "" || password == "" || len(password) > 72 {
		return AdminLoginResult{}, ErrInvalidArgument
	}
	if u.repo == nil || u.tokens == nil || u.sessions == nil {
		return AdminLoginResult{}, ErrAuthenticationUnavailable
	}

	// 1. 查询账号。
	admin, err := u.FindByAccount(ctx, account)
	if errors.Is(err, ErrAdminNotFound) {
		return AdminLoginResult{}, ErrInvalidCredentials
	}
	if err != nil {
		return AdminLoginResult{}, err
	}

	// 2. 验证密码。
	err = bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(password))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return AdminLoginResult{}, ErrInvalidCredentials
	}
	if err != nil {
		return AdminLoginResult{}, fmt.Errorf(
			"verify admin password: %w", err,
		)
	}

	// 3. 检查账号是否允许登录。
	if admin.Status != model.AdminUserStatusActive {
		return AdminLoginResult{}, ErrAccountDisabled
	}

	// 4. 签发访问凭证和刷新凭证。
	pair, err := u.tokens.Issue(admin.ID, admin.AuthVersion)
	if err != nil {
		return AdminLoginResult{}, err
	}

	// 5. 保存刷新会话。
	session := RefreshSession{
		UserID:           admin.ID,
		AuthVersion:      admin.AuthVersion,
		ClientType:       "ADMIN",
		RefreshTokenHash: pair.RefreshTokenHash,
		ExpiresAt:        pair.RefreshExpiresAt,
	}
	if err := u.sessions.SaveRefreshSession(ctx, session); err != nil {
		return AdminLoginResult{}, err
	}

	return AdminLoginResult{
		Admin:  admin,
		Tokens: pair,
	}, nil
}

// Authenticate 供后台接口的鉴权中间件使用。
func (u *AdminUsecase) Authenticate(ctx context.Context, token, scope string) (Principal, error) {
	claims, err := u.tokens.VerifyAccess(token)
	if err != nil {
		return Principal{}, err
	}

	// 现有 TokenClaims 字段叫 UserID；
	// 在管理员专用令牌中，它承载管理员ID。
	admin, err := u.repo.FindByID(ctx, claims.UserID)
	if errors.Is(err, ErrAdminNotFound) {
		return Principal{}, ErrInvalidToken
	}
	if err != nil {
		return Principal{}, ErrAuthenticationUnavailable
	}

	if admin.Status != model.AdminUserStatusActive ||
		admin.AuthVersion != claims.AuthVersion {
		return Principal{}, ErrInvalidToken
	}

	permissions, dataScope, err := u.repo.Access(ctx, admin.ID)
	if err != nil {
		return Principal{}, ErrAuthenticationUnavailable
	}

	return Principal{
		AdminID:     admin.ID,
		Permissions: permissions,
		DataScope:   dataScope,
	}, nil
}

// 根据接口要求的身份，选择用户或管理员认证器。
type ScopedAuthenticator struct {
	User  Authenticator
	Admin Authenticator
}

func (a ScopedAuthenticator) Authenticate(ctx context.Context, token, scope string) (Principal, error) {
	auth := a.User
	if scope == "admin" {
		auth = a.Admin
	}
	if auth == nil {
		return Principal{}, ErrAuthenticationUnavailable
	}

	return auth.Authenticate(ctx, token, scope)
}
