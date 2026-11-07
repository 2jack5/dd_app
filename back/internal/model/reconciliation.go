package model

import (
	"time"
)

// ReconciliationRun 对账任务与重跑版本。
type ReconciliationRun struct {
	Base
	Provider      string                    `gorm:"column:provider;type:varchar(32);not null" json:"provider"`                   // 渠道
	MerchantID    string                    `gorm:"column:merchant_id;type:varchar(64);not null" json:"merchant_id"`             // 商户
	BillType      ReconciliationRunBillType `gorm:"column:bill_type;type:varchar(24);not null" json:"bill_type"`                 // PAYMENT/REFUND/SETTLEMENT
	BillDate      Date                      `gorm:"column:bill_date;type:date;not null" json:"bill_date"`                        // 渠道账单日期，按账单时区
	RunNo         int32                     `gorm:"column:run_no;type:integer;not null" json:"run_no"`                           // 重跑版本
	Status        ReconciliationRunStatus   `gorm:"column:status;type:varchar(24);not null" json:"status"`                       // PENDING/DOWNLOADED/PARSED/COMPLETED/FAILED
	BillAssetID   *int64                    `gorm:"column:bill_asset_id;type:bigint" json:"bill_asset_id,string"`                // →media_assets；原账单
	BillTimezone  string                    `gorm:"column:bill_timezone;type:varchar(64);not null" json:"bill_timezone"`         // 渠道账单时区
	MismatchCount int32                     `gorm:"column:mismatch_count;type:integer;not null;default:0" json:"mismatch_count"` // 默认 0
	StartedAt     *time.Time                `gorm:"column:started_at;type:timestamptz" json:"started_at"`                        // 开始
	FinishedAt    *time.Time                `gorm:"column:finished_at;type:timestamptz" json:"finished_at"`                      // 结束
}

func (ReconciliationRun) TableName() string { return "reconciliation_runs" }

// ReconciliationLine 渠道账单标准行。
type ReconciliationLine struct {
	Base
	RunID         int64     `gorm:"column:run_id;type:bigint;not null" json:"run_id,string"`               // →reconciliation_runs
	LineNo        int32     `gorm:"column:line_no;type:integer;not null" json:"line_no"`                   // 原文件行号
	BusinessNo    *string   `gorm:"column:business_no;type:varchar(128)" json:"business_no"`               // 商户支付/退款/结算单号
	ProviderRef   string    `gorm:"column:provider_ref;type:varchar(128);not null" json:"provider_ref"`    // 渠道流水
	AmountCent    int64     `gorm:"column:amount_cent;type:bigint;not null" json:"amount_cent"`            // 账单金额
	FeeCent       int64     `gorm:"column:fee_cent;type:bigint;not null;default:0" json:"fee_cent"`        // 默认 0；手续费
	Currency      string    `gorm:"column:currency;type:char(3);not null;default:'CNY'" json:"currency"`   // 默认 CNY
	ChannelStatus string    `gorm:"column:channel_status;type:varchar(32);not null" json:"channel_status"` // 原渠道状态
	OccurredAt    time.Time `gorm:"column:occurred_at;type:timestamptz;not null" json:"occurred_at"`       // 交易时间
	RowDigest     string    `gorm:"column:row_digest;type:varchar(64);not null" json:"-"`                  // 原行摘要
}

func (ReconciliationLine) TableName() string { return "reconciliation_lines" }

// ReconciliationDiff 对账差异工单。
type ReconciliationDiff struct {
	Base
	RunID             int64                       `gorm:"column:run_id;type:bigint;not null" json:"run_id,string"`           // →reconciliation_runs
	LineID            *int64                      `gorm:"column:line_id;type:bigint" json:"line_id,string"`                  // →reconciliation_lines；本地单边时可空
	LocalType         ReconciliationDiffLocalType `gorm:"column:local_type;type:varchar(24);not null" json:"local_type"`     // PAYMENT/REFUND/SETTLEMENT/SETTLEMENT_SUPPLEMENT
	LocalID           *int64                      `gorm:"column:local_id;type:bigint" json:"local_id,string"`                // 本地业务记录 ID，逻辑关联
	DiffKey           string                      `gorm:"column:diff_key;type:varchar(160);not null" json:"diff_key"`        // 本次差异稳定键
	DiffType          ReconciliationDiffDiffType  `gorm:"column:diff_type;type:varchar(32);not null" json:"diff_type"`       // LOCAL_ONLY/CHANNEL_ONLY/AMOUNT/STATUS/FEE
	LocalAmountCent   *int64                      `gorm:"column:local_amount_cent;type:bigint" json:"local_amount_cent"`     // 本地金额
	ChannelAmountCent *int64                      `gorm:"column:channel_amount_cent;type:bigint" json:"channel_amount_cent"` // 渠道金额
	Status            ReconciliationDiffStatus    `gorm:"column:status;type:varchar(24);not null" json:"status"`             // OPEN/INVESTIGATING/RESOLVED/IGNORED
	AssigneeID        *int64                      `gorm:"column:assignee_id;type:bigint" json:"assignee_id,string"`          // →admin_users
	Resolution        *string                     `gorm:"column:resolution;type:text" json:"resolution"`                     // 处理依据
	ResolvedAt        *time.Time                  `gorm:"column:resolved_at;type:timestamptz" json:"resolved_at"`            // 处理时间
}

func (ReconciliationDiff) TableName() string { return "reconciliation_diffs" }
