package model

import (
	"time"
)

// User 用户。
type User struct {
	Base
	UserNo         string     `gorm:"column:user_no;type:varchar(32);not null" json:"user_no"`                       // 唯一公开编号
	Nickname       string     `gorm:"column:nickname;type:varchar(64);not null" json:"nickname"`                     // 已审核昵称
	AvatarAssetID  *int64     `gorm:"column:avatar_asset_id;type:bigint" json:"avatar_asset_id,string"`              // →media_assets
	PhoneCipher    []byte     `gorm:"column:phone_cipher;type:bytea" json:"-"`                                       // 手机号密文
	PhoneDigest    *string    `gorm:"column:phone_digest;type:varchar(64)" json:"-"`                                 // 规范化手机号的密钥摘要，供精确查找
	Status         UserStatus `gorm:"column:status;type:varchar(24);not null;default:'ACTIVE'" json:"status"`        // 默认 ACTIVE；ACTIVE/BANNED/CANCEL_PENDING/CANCELLED
	AuthVersion    int32      `gorm:"column:auth_version;type:integer;not null;default:1" json:"auth_version"`       // 默认 1；注销、封禁时递增使旧 token 失效
	LastLoginAt    *time.Time `gorm:"column:last_login_at;type:timestamptz" json:"last_login_at"`                    // 最近登录
	CancelledAt    *time.Time `gorm:"column:cancelled_at;type:timestamptz" json:"cancelled_at"`                      // 注销完成时间
	ContentVersion int32      `gorm:"column:content_version;type:integer;not null;default:1" json:"content_version"` // 默认 1；内容审核版本，提交新版本时递增；独立于登录或价格版本
}

func (User) TableName() string { return "users" }

// AuthIdentity 第三方登录身份。
type AuthIdentity struct {
	Base
	UserID   int64                `gorm:"column:user_id;type:bigint;not null" json:"user_id,string"` // →users
	Provider AuthIdentityProvider `gorm:"column:provider;type:varchar(24);not null" json:"provider"` // WECHAT_MINI/WECHAT_APP
	AppID    string               `gorm:"column:app_id;type:varchar(64);not null" json:"app_id"`     // 微信应用标识
	Subject  string               `gorm:"column:subject;type:varchar(128);not null" json:"-"`        // 该应用下 OpenID
	UnionID  *string              `gorm:"column:union_id;type:varchar(128)" json:"-"`                // 开放平台关联标识；不作为实名认证凭据
	BoundAt  time.Time            `gorm:"column:bound_at;type:timestamptz;not null" json:"bound_at"` // 绑定时间
}

func (AuthIdentity) TableName() string { return "auth_identities" }

// AgreementAcceptance 协议同意记录。
type AgreementAcceptance struct {
	Base
	UserID           int64                            `gorm:"column:user_id;type:bigint;not null" json:"user_id,string"`                   // →users
	AgreementCode    AgreementAcceptanceAgreementCode `gorm:"column:agreement_code;type:varchar(48);not null" json:"agreement_code"`       // USER/PRIVACY/COACH/ORDER
	AgreementVersion string                           `gorm:"column:agreement_version;type:varchar(32);not null" json:"agreement_version"` // 版本
	ContentHash      string                           `gorm:"column:content_hash;type:varchar(64);not null" json:"content_hash"`           // 协议正文摘要
	ContentAssetID   int64                            `gorm:"column:content_asset_id;type:bigint;not null" json:"content_asset_id,string"` // →media_assets；不可变协议文本
	OrderID          *int64                           `gorm:"column:order_id;type:bigint" json:"order_id,string"`                          // →orders；下单协议可关联订单
	AcceptedAt       time.Time                        `gorm:"column:accepted_at;type:timestamptz;not null" json:"accepted_at"`             // 同意时间
	RequestID        string                           `gorm:"column:request_id;type:varchar(64);not null" json:"request_id"`               // 请求追踪号
}

func (AgreementAcceptance) TableName() string { return "agreement_acceptances" }
