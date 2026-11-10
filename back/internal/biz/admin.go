package biz

import (
	"context"
	"errors"
	"strings"

	"dd_web/back/internal/model"
)

var ErrAdminNotFound = errors.New("admin not found")

// 业务层需要数据库层提供这些方法。
// 接口只约定方法，不写 SQL。
type AdminRepository interface {
	FindByAccount(ctx context.Context, account string) (*model.AdminUser, error)
}

type AdminUsecase struct {
	repo AdminRepository
}

func NewAdminUsecase(repo AdminRepository) *AdminUsecase {
	return &AdminUsecase{repo: repo}
}

func (u *AdminUsecase) FindByAccount(ctx context.Context, account string) (*model.AdminUser, error) {
	// 管理员模型约定：账号统一转小写。
	account = strings.ToLower(strings.TrimSpace(account))

	return u.repo.FindByAccount(ctx, account)
}
