package model

import (
	"time"
)

// ReconciliationRun 对账任务与重跑版本。
type ReconciliationRun struct {
	Base
	Provider      string                    `gorm:"column:provider;type:varchar(32);not null;comment:渠道" json:"provider"`                                     // 渠道
	MerchantID    string                    `gorm:"column:merchant_id;type:varchar(64);not null;comment:商户" json:"merchant_id"`                               // 商户
	BillType      ReconciliationRunBillType `gorm:"column:bill_type;type:varchar(24);not null;comment:PAYMENT/REFUND/SETTLEMENT" json:"bill_type"`            // PAYMENT/REFUND/SETTLEMENT
	BillDate      Date                      `gorm:"column:bill_date;type:date;not null;comment:渠道账单日期，按账单时区" json:"bill_date"`                                // 渠道账单日期，按账单时区
	RunNo         int32                     `gorm:"column:run_no;type:integer;not null;comment:重跑版本" json:"run_no"`                                           // 重跑版本
	Status        ReconciliationRunStatus   `gorm:"column:status;type:varchar(24);not null;comment:PENDING/DOWNLOADED/PARSED/COMPLETED/FAILED" json:"status"` // PENDING/DOWNLOADED/PARSED/COMPLETED/FAILED
	BillAssetID   *int64                    `gorm:"column:bill_asset_id;type:bigint;comment:→media_assets；原账单" json:"bill_asset_id,string"`                   // →media_assets；原账单
	BillTimezone  string                    `gorm:"column:bill_timezone;type:varchar(64);not null;comment:渠道账单时区" json:"bill_timezone"`                       // 渠道账单时区
	MismatchCount int32                     `gorm:"column:mismatch_count;type:integer;not null;default:0;comment:默认 0" json:"mismatch_count"`                 // 默认 0
	StartedAt     *time.Time                `gorm:"column:started_at;type:timestamptz;comment:开始" json:"started_at"`                                          // 开始
	FinishedAt    *time.Time                `gorm:"column:finished_at;type:timestamptz;comment:结束" json:"finished_at"`                                        // 结束
}

func (ReconciliationRun) TableName() string { return "reconciliation_runs" }

// ReconciliationLine 渠道账单标准行。
type ReconciliationLine struct {
	Base
	RunID         int64     `gorm:"column:run_id;type:bigint;not null;comment:→reconciliation_runs" json:"run_id,string"` // →reconciliation_runs
	LineNo        int32     `gorm:"column:line_no;type:integer;not null;comment:原文件行号" json:"line_no"`                    // 原文件行号
	BusinessNo    *string   `gorm:"column:business_no;type:varchar(128);comment:商户支付/退款/结算单号" json:"business_no"`         // 商户支付/退款/结算单号
	ProviderRef   string    `gorm:"column:provider_ref;type:varchar(128);not null;comment:渠道流水" json:"provider_ref"`      // 渠道流水
	AmountCent    int64     `gorm:"column:amount_cent;type:bigint;not null;comment:账单金额" json:"amount_cent"`              // 账单金额
	FeeCent       int64     `gorm:"column:fee_cent;type:bigint;not null;default:0;comment:默认 0；手续费" json:"fee_cent"`      // 默认 0；手续费
	Currency      string    `gorm:"column:currency;type:char(3);not null;default:'CNY';comment:默认 CNY" json:"currency"`   // 默认 CNY
	ChannelStatus string    `gorm:"column:channel_status;type:varchar(32);not null;comment:原渠道状态" json:"channel_status"`  // 原渠道状态
	OccurredAt    time.Time `gorm:"column:occurred_at;type:timestamptz;not null;comment:交易时间" json:"occurred_at"`         // 交易时间
	RowDigest     string    `gorm:"column:row_digest;type:varchar(64);not null;comment:原行摘要" json:"-"`                    // 原行摘要
}

func (ReconciliationLine) TableName() string { return "reconciliation_lines" }

// ReconciliationDiff 对账差异工单。
type ReconciliationDiff struct {
	Base
	RunID             int64                       `gorm:"column:run_id;type:bigint;not null;comment:→reconciliation_runs" json:"run_id,string"`                          // →reconciliation_runs
	LineID            *int64                      `gorm:"column:line_id;type:bigint;comment:→reconciliation_lines；本地单边时可空" json:"line_id,string"`                        // →reconciliation_lines；本地单边时可空
	LocalType         ReconciliationDiffLocalType `gorm:"column:local_type;type:varchar(24);not null;comment:PAYMENT/REFUND/SETTLEMENT" json:"local_type"`               // PAYMENT/REFUND/SETTLEMENT/SETTLEMENT_SUPPLEMENT
	LocalID           *int64                      `gorm:"column:local_id;type:bigint;comment:本地业务记录 ID，逻辑关联" json:"local_id,string"`                                     // 本地业务记录 ID，逻辑关联
	DiffKey           string                      `gorm:"column:diff_key;type:varchar(160);not null;comment:本次差异稳定键" json:"diff_key"`                                    // 本次差异稳定键
	DiffType          ReconciliationDiffDiffType  `gorm:"column:diff_type;type:varchar(32);not null;comment:LOCAL_ONLY/CHANNEL_ONLY/AMOUNT/STATUS/FEE" json:"diff_type"` // LOCAL_ONLY/CHANNEL_ONLY/AMOUNT/STATUS/FEE
	LocalAmountCent   *int64                      `gorm:"column:local_amount_cent;type:bigint;comment:本地金额" json:"local_amount_cent"`                                    // 本地金额
	ChannelAmountCent *int64                      `gorm:"column:channel_amount_cent;type:bigint;comment:渠道金额" json:"channel_amount_cent"`                                // 渠道金额
	Status            ReconciliationDiffStatus    `gorm:"column:status;type:varchar(24);not null;comment:OPEN/INVESTIGATING/RESOLVED/IGNORED" json:"status"`             // OPEN/INVESTIGATING/RESOLVED/IGNORED
	AssigneeID        *int64                      `gorm:"column:assignee_id;type:bigint;comment:→admin_users" json:"assignee_id,string"`                                 // →admin_users
	Resolution        *string                     `gorm:"column:resolution;type:text;comment:处理依据" json:"resolution"`                                                    // 处理依据
	ResolvedAt        *time.Time                  `gorm:"column:resolved_at;type:timestamptz;comment:处理时间" json:"resolved_at"`                                           // 处理时间
}

func (ReconciliationDiff) TableName() string { return "reconciliation_diffs" }
