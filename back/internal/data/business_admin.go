package data

import (
	"context"
	"errors"

	"dd_web/back/internal/biz"
	"dd_web/back/internal/model"

	"gorm.io/gorm"
)

// 管理员数据库访问对象。
// data 中保存了程序启动时建立的数据库连接。
type AdminRepo struct {
	data *Data
}

func NewAdminRepo(data *Data) *AdminRepo {
	return &AdminRepo{data: data}
}

// 检查 AdminRepo 是否实现业务层要求的接口。
var _ biz.AdminRepository = (*AdminRepo)(nil)

func (r *AdminRepo) FindByAccount(ctx context.Context, account string) (*model.AdminUser, error) {
	var admin model.AdminUser

	err := r.data.db.WithContext(ctx).
		Where("account = ?", account).
		Take(&admin).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, biz.ErrAdminNotFound
	}
	if err != nil {
		return nil, err
	}

	return &admin, nil
}
