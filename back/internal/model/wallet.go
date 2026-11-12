package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type RechargePackageStatus string
type WalletAccountStatus string
type WalletRechargeStatus string
type WalletRechargeEventStatus string

const (
	RechargePackageStatusEnabled  RechargePackageStatus = "ENABLED"
	RechargePackageStatusDisabled RechargePackageStatus = "DISABLED"

	WalletAccountStatusActive WalletAccountStatus = "ACTIVE"
	WalletAccountStatusFrozen WalletAccountStatus = "FROZEN"
	WalletAccountStatusClosed WalletAccountStatus = "CLOSED"

	WalletRechargeStatusPending       WalletRechargeStatus = "PENDING"
	WalletRechargeStatusPrepayCreated WalletRechargeStatus = "PREPAY_CREATED"
	WalletRechargeStatusSucceeded     WalletRechargeStatus = "SUCCEEDED"
	WalletRechargeStatusClosing       WalletRechargeStatus = "CLOSING"
	WalletRechargeStatusClosed        WalletRechargeStatus = "CLOSED"
	WalletRechargeStatusFailed        WalletRechargeStatus = "FAILED"

	WalletRechargeEventStatusReceived  WalletRechargeEventStatus = "RECEIVED"
	WalletRechargeEventStatusProcessed WalletRechargeEventStatus = "PROCESSED"
	WalletRechargeEventStatusRejected  WalletRechargeEventStatus = "REJECTED"
	WalletRechargeEventStatusRetry     WalletRechargeEventStatus = "RETRY"
)

// RechargePackage 是运营配置的固定充值套餐。
type RechargePackage struct {
	Base
	Name              string                `gorm:"column:name;type:varchar(64);not null;comment:套餐名称" json:"name"`
	PayAmount         decimal.Decimal       `gorm:"column:pay_amount;type:decimal(20,2);not null;comment:人民币支付金额" json:"pay_amount"`
	Currency          string                `gorm:"column:currency;type:char(3);not null;default:'CNY';comment:支付币种" json:"currency"`
	PackageMultiplier decimal.Decimal       `gorm:"column:package_multiplier;type:decimal(10,4);not null;comment:套餐到账倍率" json:"package_multiplier"`
	CreditedAmount    decimal.Decimal       `gorm:"column:credited_amount;type:decimal(20,2);not null;comment:预计到账钱包金额" json:"credited_amount"`
	AssetCode         string                `gorm:"column:asset_code;type:varchar(32);not null;default:'DD_COIN';comment:钱包资产代码" json:"asset_code"`
	Status            RechargePackageStatus `gorm:"column:status;type:varchar(16);not null;comment:ENABLED/DISABLED" json:"status"`
	EffectiveAt       *time.Time            `gorm:"column:effective_at;type:timestamptz;comment:生效时间" json:"effective_at"`
	ExpiresAt         *time.Time            `gorm:"column:expires_at;type:timestamptz;comment:失效时间" json:"expires_at"`
	Version           int32                 `gorm:"column:version;type:integer;not null;default:1;comment:套餐版本" json:"version"`
}

func (RechargePackage) TableName() string { return "recharge_packages" }

// WalletAccount 保存用户钱包余额快照。
type WalletAccount struct {
	Base
	UserID           int64               `gorm:"column:user_id;type:bigint;not null;uniqueIndex:uk_wallet_user_asset,priority:1;comment:用户ID" json:"user_id,string"`
	AssetCode        string              `gorm:"column:asset_code;type:varchar(32);not null;uniqueIndex:uk_wallet_user_asset,priority:2;comment:钱包资产代码" json:"asset_code"`
	AvailableBalance decimal.Decimal     `gorm:"column:available_balance;type:decimal(20,2);not null;default:0;comment:可用余额" json:"available_balance"`
	FrozenBalance    decimal.Decimal     `gorm:"column:frozen_balance;type:decimal(20,2);not null;default:0;comment:冻结余额" json:"frozen_balance"`
	Status           WalletAccountStatus `gorm:"column:status;type:varchar(16);not null;default:'ACTIVE';comment:ACTIVE/FROZEN/CLOSED" json:"status"`
	Version          int32               `gorm:"column:version;type:integer;not null;default:1;comment:并发版本" json:"version"`
}

func (WalletAccount) TableName() string { return "wallet_accounts" }

// WalletRecharge 保存充值请求、套餐快照和渠道支付状态。
type WalletRecharge struct {
	Base
	RechargeNo            string               `gorm:"column:recharge_no;type:varchar(40);not null;uniqueIndex:uk_wallet_recharge_no;comment:充值编号" json:"recharge_no"`
	UserID                int64                `gorm:"column:user_id;type:bigint;not null;index:idx_wallet_recharges_user_id;comment:发起用户" json:"user_id,string"`
	WalletID              *int64               `gorm:"column:wallet_id;type:bigint;comment:成功入账钱包ID" json:"wallet_id,string"`
	PackageID             int64                `gorm:"column:package_id;type:bigint;not null;comment:充值套餐ID" json:"package_id,string"`
	PackageVersion        int32                `gorm:"column:package_version;type:integer;not null;comment:套餐版本快照" json:"package_version"`
	PayAmount             decimal.Decimal      `gorm:"column:pay_amount;type:decimal(20,2);not null;comment:人民币支付金额快照" json:"pay_amount"`
	Currency              string               `gorm:"column:currency;type:char(3);not null;comment:支付币种" json:"currency"`
	PackageMultiplier     decimal.Decimal      `gorm:"column:package_multiplier;type:decimal(10,4);not null;comment:套餐倍率快照" json:"package_multiplier"`
	CreditedAmount        decimal.Decimal      `gorm:"column:credited_amount;type:decimal(20,2);not null;comment:实际到账钱包金额" json:"credited_amount"`
	AssetCode             string               `gorm:"column:asset_code;type:varchar(32);not null;comment:钱包资产代码" json:"asset_code"`
	Provider              string               `gorm:"column:provider;type:varchar(32);not null;uniqueIndex:uk_wallet_recharge_trade,priority:1;uniqueIndex:uk_wallet_provider_transaction,priority:1;comment:支付渠道" json:"provider"`
	MerchantID            string               `gorm:"column:merchant_id;type:varchar(64);not null;uniqueIndex:uk_wallet_recharge_trade,priority:2;uniqueIndex:uk_wallet_provider_transaction,priority:2;comment:渠道商户号" json:"merchant_id"`
	AppID                 string               `gorm:"column:app_id;type:varchar(64);not null;comment:渠道应用ID" json:"app_id"`
	ClientType            string               `gorm:"column:client_type;type:varchar(32);not null;comment:发起支付客户端类型" json:"client_type"`
	OutTradeNo            string               `gorm:"column:out_trade_no;type:varchar(64);not null;uniqueIndex:uk_wallet_recharge_trade,priority:3;comment:渠道商户单号" json:"out_trade_no"`
	ProviderTransactionID *string              `gorm:"column:provider_transaction_id;type:varchar(128);uniqueIndex:uk_wallet_provider_transaction,priority:3;comment:渠道交易号" json:"provider_transaction_id"`
	Status                WalletRechargeStatus `gorm:"column:status;type:varchar(24);not null;index:idx_wallet_recharges_expire,priority:1;comment:充值状态" json:"status"`
	PrepayRefCipher       []byte               `gorm:"column:prepay_ref_cipher;type:bytea;comment:预支付引用密文" json:"-"`
	FailureCode           *string              `gorm:"column:failure_code;type:varchar(64);comment:稳定失败分类" json:"failure_code"`
	ExpiresAt             time.Time            `gorm:"column:expires_at;type:timestamptz;not null;index:idx_wallet_recharges_expire,priority:2;comment:支付截止时间" json:"expires_at"`
	PaidAt                *time.Time           `gorm:"column:paid_at;type:timestamptz;comment:渠道支付时间" json:"paid_at"`
	CreditedAt            *time.Time           `gorm:"column:credited_at;type:timestamptz;comment:钱包入账时间" json:"credited_at"`
	Version               int32                `gorm:"column:version;type:integer;not null;default:1;comment:状态版本" json:"version"`
}

func (WalletRecharge) TableName() string { return "wallet_recharges" }

// WalletLedgerEntry 是只追加的钱包资金事实。
type WalletLedgerEntry struct {
	EventBase
	WalletID       int64           `gorm:"column:wallet_id;type:bigint;not null;uniqueIndex:uk_wallet_ledger_biz,priority:1;comment:钱包ID" json:"wallet_id,string"`
	UserID         int64           `gorm:"column:user_id;type:bigint;not null;index:idx_wallet_ledger_user_id;comment:用户ID" json:"user_id,string"`
	AssetCode      string          `gorm:"column:asset_code;type:varchar(32);not null;comment:钱包资产代码" json:"asset_code"`
	EntryType      string          `gorm:"column:entry_type;type:varchar(32);not null;uniqueIndex:uk_wallet_ledger_biz,priority:2;comment:流水类型" json:"entry_type"`
	BizType        string          `gorm:"column:biz_type;type:varchar(32);not null;uniqueIndex:uk_wallet_ledger_biz,priority:3;comment:业务类型" json:"biz_type"`
	BizID          int64           `gorm:"column:biz_id;type:bigint;not null;uniqueIndex:uk_wallet_ledger_biz,priority:4;comment:业务ID" json:"biz_id,string"`
	AvailableDelta decimal.Decimal `gorm:"column:available_delta;type:decimal(20,2);not null;comment:可用余额变化" json:"available_delta"`
	FrozenDelta    decimal.Decimal `gorm:"column:frozen_delta;type:decimal(20,2);not null;comment:冻结余额变化" json:"frozen_delta"`
	AvailableAfter decimal.Decimal `gorm:"column:available_after;type:decimal(20,2);not null;comment:变化后可用余额" json:"available_after"`
	FrozenAfter    decimal.Decimal `gorm:"column:frozen_after;type:decimal(20,2);not null;comment:变化后冻结余额" json:"frozen_after"`
	OccurredAt     time.Time       `gorm:"column:occurred_at;type:timestamptz;not null;comment:资金事实时间" json:"occurred_at"`
}

func (WalletLedgerEntry) TableName() string { return "wallet_ledger_entries" }

// WalletRechargeEvent 保存通知或主动查单的处理事实。
type WalletRechargeEvent struct {
	Base
	RechargeID     *int64                    `gorm:"column:recharge_id;type:bigint;comment:匹配的充值单ID" json:"recharge_id,string"`
	Provider       string                    `gorm:"column:provider;type:varchar(32);not null;uniqueIndex:uk_wallet_recharge_event,priority:1;comment:支付渠道" json:"provider"`
	MerchantID     string                    `gorm:"column:merchant_id;type:varchar(64);not null;uniqueIndex:uk_wallet_recharge_event,priority:2;comment:商户号" json:"merchant_id"`
	EventKey       string                    `gorm:"column:event_key;type:varchar(160);not null;uniqueIndex:uk_wallet_recharge_event,priority:3;comment:渠道事件幂等键" json:"event_key"`
	EventType      string                    `gorm:"column:event_type;type:varchar(48);not null;comment:PAY_NOTIFY/QUERY" json:"event_type"`
	PayloadDigest  string                    `gorm:"column:payload_digest;type:varchar(64);not null;comment:原始内容摘要" json:"-"`
	PayloadAssetID *int64                    `gorm:"column:payload_asset_id;type:bigint;comment:加密原始报文资源ID" json:"-"`
	Status         WalletRechargeEventStatus `gorm:"column:status;type:varchar(16);not null;comment:RECEIVED/PROCESSED/REJECTED/RETRY" json:"status"`
	ErrorCode      *string                   `gorm:"column:error_code;type:varchar(64);comment:稳定处理错误码" json:"error_code"`
	HandledAt      *time.Time                `gorm:"column:handled_at;type:timestamptz;comment:处理完成时间" json:"handled_at"`
}

func (WalletRechargeEvent) TableName() string { return "wallet_recharge_events" }
