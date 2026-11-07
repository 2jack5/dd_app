package model

import (
	"time"
)

// UserActivityEvent 用户活跃事件。
type UserActivityEvent struct {
	EventBase
	UserID     int64                      `gorm:"column:user_id;type:bigint;not null" json:"user_id,string"`             // →users
	EventKey   string                     `gorm:"column:event_key;type:varchar(128);not null" json:"event_key"`          // 客户端事件唯一号
	EventType  UserActivityEventEventType `gorm:"column:event_type;type:varchar(32);not null" json:"event_type"`         // APP_ACTIVE/VIEW_SERVICE/CREATE_ORDER
	OccurredAt time.Time                  `gorm:"column:occurred_at;type:timestamptz;not null" json:"occurred_at"`       // 事件发生时间
	ReceivedAt time.Time                  `gorm:"column:received_at;type:timestamptz;not null" json:"received_at"`       // 服务端接收时间
	Context    JSONB                      `gorm:"column:context;type:jsonb;not null;default:'{}'::jsonb" json:"context"` // 默认 {}；仅白名单维度
}

func (UserActivityEvent) TableName() string { return "user_activity_events" }

// DailyMetric 日统计投影。
type DailyMetric struct {
	Base
	StatDate        Date                  `gorm:"column:stat_date;type:date;not null" json:"stat_date"`                      // 业务日期，Asia/Shanghai
	MetricCode      DailyMetricMetricCode `gorm:"column:metric_code;type:varchar(48);not null" json:"metric_code"`           // DAU/PAID_GMV/REFUND/COMMISSION/COACH_ACTIVE 等
	DimensionKey    string                `gorm:"column:dimension_key;type:varchar(128);not null" json:"dimension_key"`      // ALL 或规范化游戏维度
	ValueBigint     int64                 `gorm:"column:value_bigint;type:bigint;not null" json:"value_bigint"`              // 计数或分金额
	SourceWatermark time.Time             `gorm:"column:source_watermark;type:timestamptz;not null" json:"source_watermark"` // 源数据水位
	CalculatedAt    time.Time             `gorm:"column:calculated_at;type:timestamptz;not null" json:"calculated_at"`       // 计算时间
}

func (DailyMetric) TableName() string { return "daily_metrics" }
