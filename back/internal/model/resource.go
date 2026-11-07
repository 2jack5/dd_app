package model

import (
	"time"
)

// MediaAsset 文件与证据资源。
type MediaAsset struct {
	Base
	OwnerUserID       *int64                `gorm:"column:owner_user_id;type:bigint" json:"owner_user_id,string"`               // →users；系统文件可空
	UploadedByAdminID *int64                `gorm:"column:uploaded_by_admin_id;type:bigint" json:"uploaded_by_admin_id,string"` // →admin_users
	StorageKey        string                `gorm:"column:storage_key;type:varchar(512);not null" json:"-"`                     // 私有对象键
	MIMEType          string                `gorm:"column:mime_type;type:varchar(128);not null" json:"mime_type"`               // 类型
	SizeBytes         int64                 `gorm:"column:size_bytes;type:bigint;not null" json:"size_bytes"`                   // 文件大小，>=0
	SHA256            string                `gorm:"column:sha256;type:varchar(64);not null" json:"sha256"`                      // 文件摘要
	Purpose           MediaAssetPurpose     `gorm:"column:purpose;type:varchar(32);not null" json:"purpose"`                    // 头像、证明、证据、账单、协议、渠道报文及游戏/运营图片；见 MediaAssetPurpose
	AccessScope       MediaAssetAccessScope `gorm:"column:access_scope;type:varchar(16);not null" json:"access_scope"`          // PRIVATE/PUBLIC
	Status            MediaAssetStatus      `gorm:"column:status;type:varchar(24);not null" json:"status"`                      // UPLOADING/PENDING/APPROVED/REJECTED/DELETED
	RetentionUntil    *time.Time            `gorm:"column:retention_until;type:timestamptz" json:"retention_until"`             // 保留到期日
}

func (MediaAsset) TableName() string { return "media_assets" }

// Notification 业务通知。
type Notification struct {
	Base
	UserID       int64               `gorm:"column:user_id;type:bigint;not null" json:"user_id,string"`             // →users
	BizType      NotificationBizType `gorm:"column:biz_type;type:varchar(32);not null" json:"biz_type"`             // APPLICATION/ORDER/REFUND/CASE/SETTLEMENT/SETTLEMENT_SUPPLEMENT
	BizID        int64               `gorm:"column:biz_id;type:bigint;not null" json:"biz_id,string"`               // 对应业务记录，逻辑关联
	TemplateCode string              `gorm:"column:template_code;type:varchar(64);not null" json:"template_code"`   // 消息模板
	Payload      JSONB               `gorm:"column:payload;type:jsonb;not null;default:'{}'::jsonb" json:"-"`       // 默认 {}；脱敏模板参数
	Channel      NotificationChannel `gorm:"column:channel;type:varchar(24);not null" json:"channel"`               // IN_APP/WECHAT_SUBSCRIBE
	Status       NotificationStatus  `gorm:"column:status;type:varchar(24);not null" json:"status"`                 // PENDING/SENT/FAILED
	ReadAt       *time.Time          `gorm:"column:read_at;type:timestamptz" json:"read_at"`                        // 站内已读时间
	SentAt       *time.Time          `gorm:"column:sent_at;type:timestamptz" json:"sent_at"`                        // 发送时间
	RetryCount   int32               `gorm:"column:retry_count;type:integer;not null;default:0" json:"retry_count"` // 默认 0
	DedupeKey    string              `gorm:"column:dedupe_key;type:varchar(160);not null" json:"dedupe_key"`        // 同一业务事件每渠道唯一
}

func (Notification) TableName() string { return "notifications" }
