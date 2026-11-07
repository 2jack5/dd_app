package model

import (
	"time"
)

// CoachService 陪练服务商品。
type CoachService struct {
	Base
	CoachID               int64                `gorm:"column:coach_id;type:bigint;not null" json:"coach_id,string"`                           // →coaches
	GameID                int64                `gorm:"column:game_id;type:bigint;not null" json:"game_id,string"`                             // →games
	CategoryID            int64                `gorm:"column:category_id;type:bigint;not null" json:"category_id,string"`                     // →service_categories；复合外键保证游戏一致
	Title                 string               `gorm:"column:title;type:varchar(128);not null" json:"title"`                                  // 已审核标题
	Description           string               `gorm:"column:description;type:text;not null" json:"description"`                              // 已审核服务介绍
	UnitType              CoachServiceUnitType `gorm:"column:unit_type;type:varchar(16);not null" json:"unit_type"`                           // HOUR/ROUND；计价单位，从分类复制，发布后不可直接改变含义
	BillingMinutes        *int32               `gorm:"column:billing_minutes;type:integer" json:"billing_minutes"`                            // 从分类复制，订单再快照
	BookingMinutesPerUnit int32                `gorm:"column:booking_minutes_per_unit;type:integer;not null" json:"booking_minutes_per_unit"` // 每单位预计占用分钟，>0；按局也必须配置
	PriceCent             int64                `gorm:"column:price_cent;type:bigint;not null" json:"price_cent"`                              // 单位价格，>0
	MinQty                int32                `gorm:"column:min_qty;type:integer;not null;default:1" json:"min_qty"`                         // 默认 1，>0
	MaxQty                int32                `gorm:"column:max_qty;type:integer;not null" json:"max_qty"`                                   // >=min_qty
	Status                CoachServiceStatus   `gorm:"column:status;type:varchar(24);not null" json:"status"`                                 // DRAFT/ON_SHELF/OFF_SHELF/BLOCKED
	Version               int32                `gorm:"column:version;type:integer;not null;default:1" json:"version"`                         // 默认 1；改价、可售范围修改递增
	ContentVersion        int32                `gorm:"column:content_version;type:integer;not null;default:1" json:"content_version"`         // 默认 1；内容审核版本，提交新版本时递增；独立于登录或价格版本
}

func (CoachService) TableName() string { return "coach_services" }

// CoachServiceRegion 服务可售区服。
type CoachServiceRegion struct {
	Base
	ServiceID int64 `gorm:"column:service_id;type:bigint;not null" json:"service_id,string"` // →coach_services
	GameID    int64 `gorm:"column:game_id;type:bigint;not null" json:"game_id,string"`       // →games；用于同游戏约束
	RegionID  int64 `gorm:"column:region_id;type:bigint;not null" json:"region_id,string"`   // →game_regions
}

func (CoachServiceRegion) TableName() string { return "coach_service_regions" }

// CoachAvailability 可预约排期。
type CoachAvailability struct {
	Base
	CoachID int64                   `gorm:"column:coach_id;type:bigint;not null" json:"coach_id,string"`   // →coaches
	StartAt time.Time               `gorm:"column:start_at;type:timestamptz;not null" json:"start_at"`     // 可用开始
	EndAt   time.Time               `gorm:"column:end_at;type:timestamptz;not null" json:"end_at"`         // 可用结束
	Status  CoachAvailabilityStatus `gorm:"column:status;type:varchar(16);not null" json:"status"`         // OPEN/CLOSED
	Version int32                   `gorm:"column:version;type:integer;not null;default:1" json:"version"` // 默认 1
}

func (CoachAvailability) TableName() string { return "coach_availability" }

// BookingHold 预约预占与已占库存。
type BookingHold struct {
	Base
	CoachID    int64             `gorm:"column:coach_id;type:bigint;not null" json:"coach_id,string"` // →coaches
	OrderID    int64             `gorm:"column:order_id;type:bigint;not null" json:"order_id,string"` // →orders
	StartAt    time.Time         `gorm:"column:start_at;type:timestamptz;not null" json:"start_at"`   // 预约开始
	EndAt      time.Time         `gorm:"column:end_at;type:timestamptz;not null" json:"end_at"`       // 预约结束；半开区间
	Status     BookingHoldStatus `gorm:"column:status;type:varchar(16);not null" json:"status"`       // HELD/CONFIRMED/RELEASED/EXPIRED
	ExpiresAt  *time.Time        `gorm:"column:expires_at;type:timestamptz" json:"expires_at"`        // 待支付占用的失效检查时间
	ReleasedAt *time.Time        `gorm:"column:released_at;type:timestamptz" json:"released_at"`      // 释放时间
}

func (BookingHold) TableName() string { return "booking_holds" }
