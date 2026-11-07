package model

import (
	"time"
)

// Review 订单评价主记录。
type Review struct {
	Base
	OrderID     int64            `gorm:"column:order_id;type:bigint;not null" json:"order_id,string"`   // →orders
	BuyerID     int64            `gorm:"column:buyer_id;type:bigint;not null" json:"buyer_id,string"`   // →users
	CoachID     int64            `gorm:"column:coach_id;type:bigint;not null" json:"coach_id,string"`   // →coaches
	Rating      int16            `gorm:"column:rating;type:smallint;not null" json:"rating"`            // 1..5
	Visibility  ReviewVisibility `gorm:"column:visibility;type:varchar(16);not null" json:"visibility"` // PENDING/PUBLISHED/HIDDEN
	PublishedAt *time.Time       `gorm:"column:published_at;type:timestamptz" json:"published_at"`      // 首评公开时间
}

func (Review) TableName() string { return "reviews" }

// ReviewContent 评价内容版本。
type ReviewContent struct {
	Base
	ReviewID      int64                    `gorm:"column:review_id;type:bigint;not null" json:"review_id,string"`           // →reviews
	Kind          ReviewContentKind        `gorm:"column:kind;type:varchar(16);not null" json:"kind"`                       // PRIMARY/FOLLOW_UP
	Revision      int32                    `gorm:"column:revision;type:integer;not null" json:"revision"`                   // 版本
	Content       string                   `gorm:"column:content;type:text;not null" json:"content"`                        // 评价正文
	ImageAssetIds JSONB                    `gorm:"column:image_asset_ids;type:jsonb;not null;default:'[]'::jsonb" json:"-"` // 默认 []
	AuditStatus   ReviewContentAuditStatus `gorm:"column:audit_status;type:varchar(16);not null" json:"audit_status"`       // PENDING/APPROVED/REJECTED
	IsCurrent     bool                     `gorm:"column:is_current;type:boolean;not null;default:false" json:"is_current"` // 默认 false；该类型当前公开版本
	PublishedAt   *time.Time               `gorm:"column:published_at;type:timestamptz" json:"published_at"`                // 公开时间
}

func (ReviewContent) TableName() string { return "review_contents" }
