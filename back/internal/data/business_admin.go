package data

import (
	"context"
	"errors"

	"dd_web/back/internal/biz"
	"dd_web/back/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

// CreateAccount 在同一事务中保存管理员及其角色。
func (r *AdminRepo) CreateAccount(ctx context.Context, admin *model.AdminUser, roleIDs []int64) error {
	return r.data.db.WithContext(ctx).Transaction(
		func(tx *gorm.DB) error {
			// 查询并锁定已启用的角色。
			var roles []model.AdminRole
			err := tx.Clauses(clause.Locking{Strength: "SHARE"}).
				Where("id IN ? AND status = ?", roleIDs, model.AdminRoleStatusActive).
				Order("id ASC").Find(&roles).Error
			if err != nil {
				return err
			}
			if len(roles) != len(roleIDs) {
				return biz.ErrAdminRoleInvalid
			}
			// 保存管理员。
			if err := tx.Create(admin).Error; err != nil {
				if errors.Is(err, gorm.ErrDuplicatedKey) {
					return biz.ErrAdminAccountExists
				}
				return err
			}
			// 保存管理员的角色关联。
			links := make([]model.AdminUserRole, 0, len(roleIDs))
			for _, roleID := range roleIDs {
				links = append(links, model.AdminUserRole{
					AdminID: admin.ID,
					RoleID:  roleID,
				})
			}

			return tx.Create(&links).Error
		},
	)
}

func (r *AdminRepo) FindByAccount(ctx context.Context, account string) (*model.AdminUser, error) {
	var admin model.AdminUser
	err := r.data.db.WithContext(ctx).Where("account = ?", account).Take(&admin).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, biz.ErrAdminNotFound
	}
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func (r *AdminRepo) FindByID(ctx context.Context, id int64) (*model.AdminUser, error) {
	var admin model.AdminUser

	err := r.data.db.WithContext(ctx).
		Where("id = ?", id).
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

func (r *AdminRepo) Access(ctx context.Context, adminID int64) (map[string]bool, string, error) {
	var rows []struct {
		Code      *string
		DataScope string
	}

	err := r.data.db.WithContext(ctx).Raw(`
		SELECT p.code, r.data_scope
		FROM admin_user_roles ur
		JOIN admin_roles r ON r.id = ur.role_id
		LEFT JOIN admin_role_permissions rp ON rp.role_id = r.id
		LEFT JOIN admin_permissions p ON p.id = rp.permission_id
		WHERE ur.admin_id = ? AND r.status = ?
	`, adminID, model.AdminRoleStatusActive).Scan(&rows).Error
	if err != nil {
		return nil, "", err
	}

	permissions := make(map[string]bool)
	dataScope := "ASSIGNED"

	for _, row := range rows {
		if row.Code != nil {
			permissions[*row.Code] = true
		}
		if row.DataScope == "ALL" {
			dataScope = "ALL"
		}
	}

	return permissions, dataScope, nil
}
