package model

import (
	"time"
)

// Payment 逻辑支付单。
type Payment struct {
	Base
	OrderID            int64         `gorm:"column:order_id;type:bigint;not null" json:"order_id,string"`                            // →orders
	Provider           string        `gorm:"column:provider;type:varchar(32);not null" json:"provider"`                              // 渠道适配器代码
	MerchantID         string        `gorm:"column:merchant_id;type:varchar(64);not null" json:"merchant_id"`                        // 商户标识
	AppID              string        `gorm:"column:app_id;type:varchar(64);not null" json:"app_id"`                                  // 发起支付应用
	OutTradeNo         string        `gorm:"column:out_trade_no;type:varchar(64);not null" json:"out_trade_no"`                      // 商户单号
	TransactionID      *string       `gorm:"column:transaction_id;type:varchar(128)" json:"transaction_id"`                          // 渠道交易号；未支付为 NULL
	AmountCent         int64         `gorm:"column:amount_cent;type:bigint;not null" json:"amount_cent"`                             // 应付总额，>0
	Currency           string        `gorm:"column:currency;type:char(3);not null;default:'CNY'" json:"currency"`                    // 默认 CNY
	Status             PaymentStatus `gorm:"column:status;type:varchar(24);not null" json:"status"`                                  // UNPAID/PREPAY_CREATED/PAID/CLOSING/CLOSED/FAILED
	PrepayRefCipher    []byte        `gorm:"column:prepay_ref_cipher;type:bytea" json:"-"`                                           // 预支付引用密文
	PaidAt             *time.Time    `gorm:"column:paid_at;type:timestamptz" json:"paid_at"`                                         // 渠道真实支付时间
	RefundReservedCent int64         `gorm:"column:refund_reserved_cent;type:bigint;not null;default:0" json:"refund_reserved_cent"` // 默认 0；已预留但未成功退款额
	RefundedCent       int64         `gorm:"column:refunded_cent;type:bigint;not null;default:0" json:"refunded_cent"`               // 默认 0；成功退款累计
	Version            int32         `gorm:"column:version;type:integer;not null;default:1" json:"version"`                          // 默认 1
}

func (Payment) TableName() string { return "payments" }

// PaymentEvent 支付通知与查单事实。
// 原始事实列不可变；处理账号仅更新 PaymentID、Status、ErrorCode、HandledAt 等处理字段。
type PaymentEvent struct {
	Base
	PaymentID      *int64                `gorm:"column:payment_id;type:bigint" json:"payment_id,string"`          // →payments；未匹配回调可空
	Provider       string                `gorm:"column:provider;type:varchar(32);not null" json:"provider"`       // 渠道
	MerchantID     string                `gorm:"column:merchant_id;type:varchar(64);not null" json:"merchant_id"` // 商户
	EventKey       string                `gorm:"column:event_key;type:varchar(160);not null" json:"event_key"`    // 回调通知号或主动查询唯一标识
	EventType      PaymentEventEventType `gorm:"column:event_type;type:varchar(48);not null" json:"event_type"`   // PAY_NOTIFY/QUERY/CLOSE
	PayloadDigest  string                `gorm:"column:payload_digest;type:varchar(64);not null" json:"-"`        // 原始内容摘要
	PayloadAssetID *int64                `gorm:"column:payload_asset_id;type:bigint" json:"-"`                    // →media_assets；受控加密原始报文引用
	Status         PaymentEventStatus    `gorm:"column:status;type:varchar(16);not null" json:"status"`           // RECEIVED/PROCESSED/RETRY/REJECTED
	ErrorCode      *string               `gorm:"column:error_code;type:varchar(64)" json:"error_code"`            // 异常分类
	HandledAt      *time.Time            `gorm:"column:handled_at;type:timestamptz" json:"handled_at"`            // 处理时间
}

func (PaymentEvent) TableName() string { return "payment_events" }

// Refund 退款单。
type Refund struct {
	Base
	RefundNo         string            `gorm:"column:refund_no;type:varchar(40);not null" json:"refund_no"`                // 业务退款号
	OrderID          int64             `gorm:"column:order_id;type:bigint;not null" json:"order_id,string"`                // →orders
	PaymentID        int64             `gorm:"column:payment_id;type:bigint;not null" json:"payment_id,string"`            // →payments
	ApplicantUserID  *int64            `gorm:"column:applicant_user_id;type:bigint" json:"applicant_user_id,string"`       // →users；系统发起可空
	CaseID           *int64            `gorm:"column:case_id;type:bigint" json:"case_id,string"`                           // →after_sales_cases
	AmountCent       int64             `gorm:"column:amount_cent;type:bigint;not null" json:"amount_cent"`                 // 申请及批准金额，>0；批准后不可改
	ReasonCode       string            `gorm:"column:reason_code;type:varchar(48);not null" json:"reason_code"`            // 原因
	ReasonDetail     string            `gorm:"column:reason_detail;type:text;not null;default:''" json:"reason_detail"`    // 默认空
	EvidenceAssetIds JSONB             `gorm:"column:evidence_asset_ids;type:jsonb;not null;default:'[]'::jsonb" json:"-"` // 默认 []
	Status           RefundStatus      `gorm:"column:status;type:varchar(24);not null" json:"status"`                      // APPLIED/REVIEWING/APPROVED/SUBMITTED/PROCESSING/SUCCESS/FAILED/REJECTED/CANCELLED
	FundsStatus      RefundFundsStatus `gorm:"column:funds_status;type:varchar(16);not null" json:"funds_status"`          // NONE/RESERVED/CONSUMED/RELEASED
	ReviewerID       *int64            `gorm:"column:reviewer_id;type:bigint" json:"reviewer_id,string"`                   // →admin_users
	ApprovalID       *int64            `gorm:"column:approval_id;type:bigint" json:"approval_id,string"`                   // →admin_approvals
	ProviderRefundID *string           `gorm:"column:provider_refund_id;type:varchar(128)" json:"provider_refund_id"`      // 渠道退款号
	SubmittedAt      *time.Time        `gorm:"column:submitted_at;type:timestamptz" json:"submitted_at"`                   // 首次提交渠道时间
	SuccessAt        *time.Time        `gorm:"column:success_at;type:timestamptz" json:"success_at"`                       // 渠道确认成功
	FailureCode      *string           `gorm:"column:failure_code;type:varchar(64)" json:"failure_code"`                   // 失败原因
	Version          int32             `gorm:"column:version;type:integer;not null;default:1" json:"version"`              // 默认 1
}

func (Refund) TableName() string { return "refunds" }

// RefundEvent 退款事件。
type RefundEvent struct {
	EventBase
	RefundID   int64                `gorm:"column:refund_id;type:bigint;not null" json:"refund_id,string"`   // →refunds
	EventKey   string               `gorm:"column:event_key;type:varchar(160);not null" json:"event_key"`    // 唯一事件键
	EventType  string               `gorm:"column:event_type;type:varchar(48);not null" json:"event_type"`   // 事件类型
	FromStatus *RefundStatus        `gorm:"column:from_status;type:varchar(24)" json:"from_status"`          // 原状态
	ToStatus   RefundStatus         `gorm:"column:to_status;type:varchar(24);not null" json:"to_status"`     // 新状态
	ActorType  RefundEventActorType `gorm:"column:actor_type;type:varchar(16);not null" json:"actor_type"`   // USER/ADMIN/PROVIDER/SYSTEM
	ActorID    *int64               `gorm:"column:actor_id;type:bigint" json:"actor_id,string"`              // 操作者
	Payload    JSONB                `gorm:"column:payload;type:jsonb;not null;default:'{}'::jsonb" json:"-"` // 默认 {}；脱敏结果和证据引用
}

func (RefundEvent) TableName() string { return "refund_events" }
