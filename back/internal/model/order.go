package model

import (
	"time"
)

// OrderQuote 订单报价。
type OrderQuote struct {
	Base
	QuoteNo        string           `gorm:"column:quote_no;type:varchar(64);not null" json:"quote_no"`           // 不可猜报价号
	BuyerID        int64            `gorm:"column:buyer_id;type:bigint;not null" json:"buyer_id,string"`         // →users
	ServiceID      int64            `gorm:"column:service_id;type:bigint;not null" json:"service_id,string"`     // →coach_services
	ServiceVersion int32            `gorm:"column:service_version;type:integer;not null" json:"service_version"` // 报价时服务版本
	RegionID       int64            `gorm:"column:region_id;type:bigint;not null" json:"region_id,string"`       // →game_regions
	Quantity       int32            `gorm:"column:quantity;type:integer;not null" json:"quantity"`               // 计价单位数量，>0
	StartAt        time.Time        `gorm:"column:start_at;type:timestamptz;not null" json:"start_at"`           // 预约开始
	EndAt          time.Time        `gorm:"column:end_at;type:timestamptz;not null" json:"end_at"`               // 服务端算出的预计结束
	AmountCent     int64            `gorm:"column:amount_cent;type:bigint;not null" json:"amount_cent"`          // 报价总价，>0
	Snapshot       JSONB            `gorm:"column:snapshot;type:jsonb;not null" json:"-"`                        // 服务、单价、计价规则、佣金及业务规则版本
	Status         OrderQuoteStatus `gorm:"column:status;type:varchar(16);not null" json:"status"`               // VALID/CONSUMED/EXPIRED
	ExpiresAt      time.Time        `gorm:"column:expires_at;type:timestamptz;not null" json:"expires_at"`       // 报价过期时间
}

func (OrderQuote) TableName() string { return "order_quotes" }

// Order 业务订单。
type Order struct {
	Base
	OrderNo             string             `gorm:"column:order_no;type:varchar(40);not null" json:"order_no"`                            // 业务订单号
	QuoteID             int64              `gorm:"column:quote_id;type:bigint;not null" json:"quote_id,string"`                          // →order_quotes
	BuyerID             int64              `gorm:"column:buyer_id;type:bigint;not null" json:"buyer_id,string"`                          // →users
	CoachID             int64              `gorm:"column:coach_id;type:bigint;not null" json:"coach_id,string"`                          // →coaches
	ServiceID           int64              `gorm:"column:service_id;type:bigint;not null" json:"service_id,string"`                      // →coach_services
	GameID              int64              `gorm:"column:game_id;type:bigint;not null" json:"game_id,string"`                            // →games
	RegionID            int64              `gorm:"column:region_id;type:bigint;not null" json:"region_id,string"`                        // →game_regions
	Status              OrderStatus        `gorm:"column:status;type:varchar(24);not null;default:'PENDING_PAYMENT'" json:"status"`      // 默认 PENDING_PAYMENT；PENDING_PAYMENT/WAIT_ACCEPT/WAIT_START/IN_SERVICE/WAIT_CONFIRM/COMPLETED/CLOSED
	Quantity            int32              `gorm:"column:quantity;type:integer;not null" json:"quantity"`                                // >0
	UnitPriceCent       int64              `gorm:"column:unit_price_cent;type:bigint;not null" json:"unit_price_cent"`                   // 成交单价，>0
	TotalCent           int64              `gorm:"column:total_cent;type:bigint;not null" json:"total_cent"`                             // 总价；首期无优惠=单价×数量
	Currency            string             `gorm:"column:currency;type:char(3);not null;default:'CNY'" json:"currency"`                  // 默认 CNY
	CommissionBPS       int32              `gorm:"column:commission_bps;type:integer;not null" json:"commission_bps"`                    // 0..10000；万分比快照
	CommissionCent      int64              `gorm:"column:commission_cent;type:bigint;not null" json:"commission_cent"`                   // 订单初始佣金
	CoachIncomeCent     int64              `gorm:"column:coach_income_cent;type:bigint;not null" json:"coach_income_cent"`               // 初始理论应得=总额-佣金
	CommissionRuleID    int64              `gorm:"column:commission_rule_id;type:bigint;not null" json:"commission_rule_id,string"`      // →commission_rules
	BusinessRuleID      int64              `gorm:"column:business_rule_id;type:bigint;not null" json:"business_rule_id,string"`          // →business_rule_versions
	ServiceSnapshot     JSONB              `gorm:"column:service_snapshot;type:jsonb;not null" json:"-"`                                 // 不可变展示和计价快照
	AppointmentStartAt  time.Time          `gorm:"column:appointment_start_at;type:timestamptz;not null" json:"appointment_start_at"`    // 预约开始
	AppointmentEndAt    time.Time          `gorm:"column:appointment_end_at;type:timestamptz;not null" json:"appointment_end_at"`        // 预约结束
	GamePlayerIDCipher  []byte             `gorm:"column:game_player_id_cipher;type:bytea;not null" json:"-"`                            // 游戏 ID/角色名密文
	Remark              string             `gorm:"column:remark;type:varchar(512);not null;default:''" json:"remark"`                    // 默认空；限制隐私内容
	PayDeadlineAt       time.Time          `gorm:"column:pay_deadline_at;type:timestamptz;not null" json:"pay_deadline_at"`              // 支付截止
	AcceptDeadlineAt    *time.Time         `gorm:"column:accept_deadline_at;type:timestamptz" json:"accept_deadline_at"`                 // 支付后生成接单截止
	AcceptedAt          *time.Time         `gorm:"column:accepted_at;type:timestamptz" json:"accepted_at"`                               // 接单时间
	StartedAt           *time.Time         `gorm:"column:started_at;type:timestamptz" json:"started_at"`                                 // 真实开始时间
	SubmittedCompleteAt *time.Time         `gorm:"column:submitted_complete_at;type:timestamptz" json:"submitted_complete_at"`           // 提交完成时间
	AutoConfirmAt       *time.Time         `gorm:"column:auto_confirm_at;type:timestamptz" json:"auto_confirm_at"`                       // 允许自动确认时间
	CompletedAt         *time.Time         `gorm:"column:completed_at;type:timestamptz" json:"completed_at"`                             // 确认完成时间
	ClosedAt            *time.Time         `gorm:"column:closed_at;type:timestamptz" json:"closed_at"`                                   // 关闭时间
	CloseReason         *string            `gorm:"column:close_reason;type:varchar(64)" json:"close_reason"`                             // 取消、拒单、超时等原因
	DisputeStatus       OrderDisputeStatus `gorm:"column:dispute_status;type:varchar(16);not null;default:'NONE'" json:"dispute_status"` // 默认 NONE；NONE/OPEN/RESOLVED
	Version             int32              `gorm:"column:version;type:integer;not null;default:1" json:"version"`                        // 默认 1
}

func (Order) TableName() string { return "orders" }

// OrderEvent 订单事件时间轴。
type OrderEvent struct {
	EventBase
	OrderID    int64               `gorm:"column:order_id;type:bigint;not null" json:"order_id,string"`     // →orders
	FromStatus *OrderStatus        `gorm:"column:from_status;type:varchar(24)" json:"from_status"`          // 原履约态
	ToStatus   OrderStatus         `gorm:"column:to_status;type:varchar(24);not null" json:"to_status"`     // 目标履约态，可与原态相同
	EventType  string              `gorm:"column:event_type;type:varchar(48);not null" json:"event_type"`   // 事件代码
	ActorType  OrderEventActorType `gorm:"column:actor_type;type:varchar(16);not null" json:"actor_type"`   // USER/COACH/ADMIN/SYSTEM
	ActorID    *int64              `gorm:"column:actor_id;type:bigint" json:"actor_id,string"`              // 对应用户或管理员；SYSTEM 可空
	RequestID  string              `gorm:"column:request_id;type:varchar(64);not null" json:"request_id"`   // 链路追踪号，不单独唯一
	EventKey   string              `gorm:"column:event_key;type:varchar(160);not null" json:"event_key"`    // 全局事件幂等键
	Payload    JSONB               `gorm:"column:payload;type:jsonb;not null;default:'{}'::jsonb" json:"-"` // 默认 {}；取消原因、退款引用、履约证据 ID 等
}

func (OrderEvent) TableName() string { return "order_events" }
