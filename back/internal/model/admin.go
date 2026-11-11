package model

import (
	"database/sql"
	"time"
)

// AdminUser 后台管理员。
type AdminUser struct {
	Base
	Account         string          `gorm:"column:account;type:varchar(64);not null;uniqueIndex:uk_admin_users_account;comment:统一转小写的登录名" json:"account"` // 统一转小写的登录名
	PasswordHash    string          `gorm:"column:password_hash;type:varchar(256);not null;comment:强密码哈希；不存明文" json:"-"`                                  // 强密码哈希；不存明文
	DisplayName     string          `gorm:"column:display_name;type:varchar(64);not null;comment:名称" json:"display_name"`                                 // 名称
	Status          AdminUserStatus `gorm:"column:status;type:varchar(16);not null;comment:ACTIVE/DISABLED" json:"status"`                                // ACTIVE/DISABLED
	MFASecretCipher []byte          `gorm:"column:mfa_secret_cipher;type:bytea;comment:二次认证密钥密文" json:"-"`                                                // 二次认证密钥密文
	AuthVersion     int32           `gorm:"column:auth_version;type:integer;not null;default:1;comment:默认 1" json:"auth_version"`                         // 默认 1
	LastLoginAt     *time.Time      `gorm:"column:last_login_at;type:timestamptz;comment:最近登录" json:"last_login_at"`                                      // 最近登录
}

func (AdminUser) TableName() string { return "admin_users" }

// AdminRole 后台角色。
type AdminRole struct {
	Base
	Code      string             `gorm:"column:code;type:varchar(48);not null;comment:角色编码" json:"code"`                              // 角色编码
	Name      string             `gorm:"column:name;type:varchar(64);not null;comment:名称" json:"name"`                                // 名称
	DataScope AdminRoleDataScope `gorm:"column:data_scope;type:varchar(24);not null;comment:ALL/ASSIGNED；首期只支持这两种" json:"data_scope"` // ALL/ASSIGNED；首期只支持这两种
	Status    AdminRoleStatus    `gorm:"column:status;type:varchar(16);not null;comment:ACTIVE/DISABLED" json:"status"`               // ACTIVE/DISABLED
}

func (AdminRole) TableName() string { return "admin_roles" }

// AdminPermission 后台权限字典。
type AdminPermission struct {
	Base
	Code      string `gorm:"column:code;type:varchar(96);not null;comment:如 refund.review/settlement.submit" json:"code"` // 如 refund.review/settlement.submit
	Name      string `gorm:"column:name;type:varchar(64);not null;comment:权限名称" json:"name"`                              // 权限名称
	Resource  string `gorm:"column:resource;type:varchar(48);not null;comment:资源类别" json:"resource"`                      // 资源类别
	Action    string `gorm:"column:action;type:varchar(48);not null;comment:动作" json:"action"`                            // 动作
	Sensitive bool   `gorm:"column:sensitive;type:boolean;not null;default:false;comment:默认 false" json:"sensitive"`      // 默认 false
}

func (AdminPermission) TableName() string { return "admin_permissions" }

// AdminUserRole 管理员角色关联。
type AdminUserRole struct {
	Base
	AdminID int64 `gorm:"column:admin_id;type:bigint;not null;comment:→admin_users" json:"admin_id,string"` // →admin_users
	RoleID  int64 `gorm:"column:role_id;type:bigint;not null;comment:→admin_roles" json:"role_id,string"`   // →admin_roles
}

func (AdminUserRole) TableName() string { return "admin_user_roles" }

// AdminRolePermission 角色权限关联。
type AdminRolePermission struct {
	Base
	RoleID       int64 `gorm:"column:role_id;type:bigint;not null;comment:→admin_roles" json:"role_id,string"`                   // →admin_roles
	PermissionID int64 `gorm:"column:permission_id;type:bigint;not null;comment:→admin_permissions" json:"permission_id,string"` // →admin_permissions
}

func (AdminRolePermission) TableName() string { return "admin_role_permissions" }

// AdminApproval 敏感操作双人复核。
type AdminApproval struct {
	Base
	Operation       string              `gorm:"column:operation;type:varchar(48);not null;comment:操作代码" json:"operation"`                                 // 操作代码
	ResourceType    string              `gorm:"column:resource_type;type:varchar(32);not null;comment:资源类别" json:"resource_type"`                         // 资源类别
	ResourceID      int64               `gorm:"column:resource_id;type:bigint;not null;comment:业务资源 ID，逻辑关联" json:"resource_id,string"`                   // 业务资源 ID，逻辑关联
	RequesterID     int64               `gorm:"column:requester_id;type:bigint;not null;comment:→admin_users" json:"requester_id,string"`                 // →admin_users
	ReviewerID      *int64              `gorm:"column:reviewer_id;type:bigint;comment:→admin_users" json:"reviewer_id,string"`                            // →admin_users
	PayloadSnapshot JSONB               `gorm:"column:payload_snapshot;type:jsonb;not null;comment:待批准参数快照" json:"-"`                                     // 待批准参数快照
	PayloadHash     string              `gorm:"column:payload_hash;type:varchar(64);not null;comment:快照哈希" json:"-"`                                      // 快照哈希
	ResourceVersion int32               `gorm:"column:resource_version;type:integer;not null;comment:提交时业务版本" json:"resource_version"`                    // 提交时业务版本
	Status          AdminApprovalStatus `gorm:"column:status;type:varchar(16);not null;comment:PENDING/APPROVED/REJECTED/EXPIRED/EXECUTED" json:"status"` // PENDING/APPROVED/REJECTED/EXPIRED/EXECUTED
	ExpiresAt       time.Time           `gorm:"column:expires_at;type:timestamptz;not null;comment:批准有效期" json:"expires_at"`                              // 批准有效期
	ReviewedAt      *time.Time          `gorm:"column:reviewed_at;type:timestamptz;comment:复核时间" json:"reviewed_at"`                                      // 复核时间
	ExecutedAt      *time.Time          `gorm:"column:executed_at;type:timestamptz;comment:实际执行时间" json:"executed_at"`                                    // 实际执行时间
}

func (AdminApproval) TableName() string { return "admin_approvals" }

// AdminAuditLog 后台与系统敏感操作审计。
type AdminAuditLog struct {
	EventBase
	ActorType      AdminAuditLogActorType `gorm:"column:actor_type;type:varchar(16);not null;comment:ADMIN/SYSTEM" json:"actor_type"`         // ADMIN/SYSTEM
	AdminID        *int64                 `gorm:"column:admin_id;type:bigint;comment:→admin_users；系统任务为空" json:"admin_id,string"`             // →admin_users；系统任务为空
	Operation      string                 `gorm:"column:operation;type:varchar(64);not null;comment:动作" json:"operation"`                     // 动作
	ResourceType   string                 `gorm:"column:resource_type;type:varchar(32);not null;comment:对象类型" json:"resource_type"`           // 对象类型
	ResourceID     int64                  `gorm:"column:resource_id;type:bigint;not null;comment:对象 ID，逻辑关联" json:"resource_id,string"`       // 对象 ID，逻辑关联
	BeforeRedacted JSONB                  `gorm:"column:before_redacted;type:jsonb;not null;default:'{}'::jsonb;comment:默认 {}；脱敏前值" json:"-"` // 默认 {}；脱敏前值
	AfterRedacted  JSONB                  `gorm:"column:after_redacted;type:jsonb;not null;default:'{}'::jsonb;comment:默认 {}；脱敏后值" json:"-"`  // 默认 {}；脱敏后值
	PayloadAssetID *int64                 `gorm:"column:payload_asset_id;type:bigint;comment:→media_assets；加密渠道报文引用" json:"-"`                // →media_assets；加密渠道报文引用
	RequestID      string                 `gorm:"column:request_id;type:varchar(64);not null;comment:链路追踪" json:"request_id"`                 // 链路追踪
	IP             sql.NullString         `gorm:"column:ip;type:inet;comment:来源 IP" json:"-"`                                                 // 来源 IP
	EventKey       string                 `gorm:"column:event_key;type:varchar(160);not null;comment:唯一审计事件键" json:"event_key"`               // 唯一审计事件键
	Outcome        AdminAuditLogOutcome   `gorm:"column:outcome;type:varchar(16);not null;comment:SUCCESS/FAILED" json:"outcome"`             // SUCCESS/FAILED
}

func (AdminAuditLog) TableName() string { return "admin_audit_logs" }
