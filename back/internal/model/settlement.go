package model

import (
	"time"
)

// SettlementAccount 陪练渠道结算账户映射。
type SettlementAccount struct {
	Base
	CoachID            int64                   `gorm:"column:coach_id;type:bigint;not null;comment:→coaches" json:"coach_id,string"`                             // →coaches
	Provider           string                  `gorm:"column:provider;type:varchar(32);not null;comment:结算渠道" json:"provider"`                                   // 结算渠道
	ProviderAccountRef *string                 `gorm:"column:provider_account_ref;type:varchar(128);comment:子商户或合法收款方引用" json:"-"`                               // 子商户或合法收款方引用
	Status             SettlementAccountStatus `gorm:"column:status;type:varchar(24);not null;comment:PENDING/REVIEWING/ACTIVE/REJECTED/DISABLED" json:"status"` // PENDING/REVIEWING/ACTIVE/REJECTED/DISABLED
	AccountMasked      *string                 `gorm:"column:account_masked;type:varchar(128);comment:脱敏收款信息" json:"account_masked"`                             // 脱敏收款信息
	OnboardingRef      *string                 `gorm:"column:onboarding_ref;type:varchar(128);comment:进件请求号" json:"onboarding_ref"`                              // 进件请求号
	VerifiedAt         *time.Time              `gorm:"column:verified_at;type:timestamptz;comment:渠道确认时间" json:"verified_at"`                                    // 渠道确认时间
	Version            int32                   `gorm:"column:version;type:integer;not null;default:1;comment:默认 1；更新时通过 WHERE version 校验并显式递增" json:"version"`   // 默认 1；更新时通过 WHERE version 校验并显式递增
}

func (SettlementAccount) TableName() string { return "settlement_accounts" }

// CommissionRule 佣金规则版本。
type CommissionRule struct {
	Base
	ScopeType   CommissionRuleScopeType `gorm:"column:scope_type;type:varchar(24);not null;comment:GLOBAL/GAME_CATEGORY/COACH" json:"scope_type"`       // GLOBAL/GAME_CATEGORY/COACH
	ScopeKey    string                  `gorm:"column:scope_key;type:varchar(96);not null;comment:GLOBAL 或标准化范围键" json:"scope_key"`                     // GLOBAL 或标准化范围键
	CoachID     *int64                  `gorm:"column:coach_id;type:bigint;comment:→coaches；COACH 时必填" json:"coach_id,string"`                          // →coaches；COACH 时必填
	GameID      *int64                  `gorm:"column:game_id;type:bigint;comment:→games；GAME_CATEGORY 时必填" json:"game_id,string"`                      // →games；GAME_CATEGORY 时必填
	CategoryID  *int64                  `gorm:"column:category_id;type:bigint;comment:→service_categories；GAME_CATEGORY 时必填" json:"category_id,string"` // →service_categories；GAME_CATEGORY 时必填
	Revision    int32                   `gorm:"column:revision;type:integer;not null;comment:范围内版本" json:"revision"`                                    // 范围内版本
	BPS         int32                   `gorm:"column:bps;type:integer;not null;comment:费率 0..10000" json:"bps"`                                        // 费率 0..10000
	EffectiveAt time.Time               `gorm:"column:effective_at;type:timestamptz;not null;comment:生效开始" json:"effective_at"`                         // 生效开始
	ExpiresAt   *time.Time              `gorm:"column:expires_at;type:timestamptz;comment:结束；NULL 无上限" json:"expires_at"`                               // 结束；NULL 无上限
	Status      CommissionRuleStatus    `gorm:"column:status;type:varchar(16);not null;comment:DRAFT/PUBLISHED/RETIRED" json:"status"`                  // DRAFT/PUBLISHED/RETIRED
	CreatedBy   int64                   `gorm:"column:created_by;type:bigint;not null;comment:→admin_users" json:"created_by"`                          // →admin_users
	ApprovalID  *int64                  `gorm:"column:approval_id;type:bigint;comment:→admin_approvals" json:"approval_id,string"`                      // →admin_approvals
	Version     int32                   `gorm:"column:version;type:integer;not null;default:1;comment:默认 1；更新时通过 WHERE version 校验并显式递增" json:"version"` // 默认 1；更新时通过 WHERE version 校验并显式递增
}

func (CommissionRule) TableName() string { return "commission_rules" }

// SettlementItem 订单应结算明细。
// 支付成功即创建 NOT_ELIGIBLE 基准，确认完成只激活账期；不能等待完成后才记退款调整。
type SettlementItem struct {
	Base
	OrderID                 int64                `gorm:"column:order_id;type:bigint;not null;comment:→orders" json:"order_id,string"`                                                    // →orders
	CoachID                 int64                `gorm:"column:coach_id;type:bigint;not null;comment:→coaches" json:"coach_id,string"`                                                   // →coaches
	AccountID               *int64               `gorm:"column:account_id;type:bigint;comment:→settlement_accounts" json:"account_id,string"`                                            // →settlement_accounts
	GrossCent               int64                `gorm:"column:gross_cent;type:bigint;not null;comment:支付原始服务额" json:"gross_cent"`                                                       // 支付原始服务额
	CommissionCent          int64                `gorm:"column:commission_cent;type:bigint;not null;comment:原始佣金" json:"commission_cent"`                                                // 原始佣金
	CoachNetCent            int64                `gorm:"column:coach_net_cent;type:bigint;not null;comment:原始理论陪练收入" json:"coach_net_cent"`                                              // 原始理论陪练收入
	RefundReductionCent     int64                `gorm:"column:refund_reduction_cent;type:bigint;not null;default:0;comment:默认 0；成功退款扣减服务额汇总缓存" json:"refund_reduction_cent"`            // 默认 0；成功退款扣减服务额汇总缓存
	CommissionReductionCent int64                `gorm:"column:commission_reduction_cent;type:bigint;not null;default:0;comment:默认 0；佣金冲减汇总缓存" json:"commission_reduction_cent"`         // 默认 0；佣金冲减汇总缓存
	NetAdjustmentCent       int64                `gorm:"column:net_adjustment_cent;type:bigint;not null;default:0;comment:默认 0；陪练应得净调整，可为负" json:"net_adjustment_cent"`                  // 默认 0；陪练应得净调整，可为负
	ChannelFeeCent          int64                `gorm:"column:channel_fee_cent;type:bigint;not null;default:0;comment:默认 0；陪练承担通道费，实际规则决定" json:"channel_fee_cent"`                     // 默认 0；陪练承担通道费，实际规则决定
	PayableCent             int64                `gorm:"column:payable_cent;type:bigint;not null;comment:当前应付=max(原始应得+净调整-通道费,0)，超出部分单独追偿" json:"payable_cent"`                         // 当前应付=max(原始应得+净调整-通道费,0)，超出部分单独追偿
	PaidOutCent             int64                `gorm:"column:paid_out_cent;type:bigint;not null;default:0;comment:首次出款与独立补款成功金额累计，幂等更新" json:"paid_out_cent"`                          // 默认 0；首次出款与独立补款成功金额的累计，幂等更新
	Status                  SettlementItemStatus `gorm:"column:status;type:varchar(24);not null;comment:NOT_ELIGIBLE/PENDING/HOLDING/SUBMITTED/SETTLED/FAILED/NO_PAYABLE" json:"status"` // NOT_ELIGIBLE/PENDING/HOLDING/SUBMITTED/SETTLED/FAILED/NO_PAYABLE
	HoldReason              *string              `gorm:"column:hold_reason;type:varchar(64);comment:冻结原因" json:"hold_reason"`                                                            // 冻结原因
	EligibleAt              *time.Time           `gorm:"column:eligible_at;type:timestamptz;comment:满足账期时间" json:"eligible_at"`                                                          // 满足账期时间
	Version                 int32                `gorm:"column:version;type:integer;not null;default:1;comment:默认 1" json:"version"`                                                     // 默认 1
}

func (SettlementItem) TableName() string { return "settlement_items" }

// SettlementAdjustment 结算调整与已结算追偿记录。
// 金额事实仅追加；追偿进度可更新并追加审计，不能覆盖历史出款金额。
type SettlementAdjustment struct {
	Base
	SettlementItemID    int64                              `gorm:"column:settlement_item_id;type:bigint;not null;comment:→settlement_items" json:"settlement_item_id,string"`                // →settlement_items
	RefundID            *int64                             `gorm:"column:refund_id;type:bigint;comment:→refunds" json:"refund_id,string"`                                                    // →refunds
	Kind                SettlementAdjustmentKind           `gorm:"column:kind;type:varchar(24);not null;comment:REFUND/FEE/REVERSAL/CORRECTION" json:"kind"`                                 // REFUND/FEE/REVERSAL/CORRECTION
	BusinessKey         string                             `gorm:"column:business_key;type:varchar(160);not null;comment:唯一业务调整键" json:"business_key"`                                       // 唯一业务调整键
	GrossDeltaCent      int64                              `gorm:"column:gross_delta_cent;type:bigint;not null;comment:服务总额调整，带符号" json:"gross_delta_cent"`                                  // 服务总额调整，带符号
	CommissionDeltaCent int64                              `gorm:"column:commission_delta_cent;type:bigint;not null;comment:佣金调整，带符号" json:"commission_delta_cent"`                          // 佣金调整，带符号
	CoachDeltaCent      int64                              `gorm:"column:coach_delta_cent;type:bigint;not null;comment:陪练应得调整，带符号" json:"coach_delta_cent"`                                  // 陪练应得调整，带符号
	ChannelFeeDeltaCent int64                              `gorm:"column:channel_fee_delta_cent;type:bigint;not null;comment:陪练承担手续费调整，带符号" json:"channel_fee_delta_cent"`                   // 陪练承担手续费调整，带符号
	RecoveryCent        int64                              `gorm:"column:recovery_cent;type:bigint;not null;default:0;comment:默认 0；已结算后需追偿额" json:"recovery_cent"`                           // 默认 0；已结算后需追偿额
	RecoveryStatus      SettlementAdjustmentRecoveryStatus `gorm:"column:recovery_status;type:varchar(24);not null;comment:NONE/PENDING/PROCESSING/RECOVERED/FAILED" json:"recovery_status"` // NONE/PENDING/PROCESSING/RECOVERED/FAILED
	ProviderRef         *string                            `gorm:"column:provider_ref;type:varchar(128);comment:渠道追偿或冲正引用" json:"provider_ref"`                                              // 渠道追偿或冲正引用
	Reason              string                             `gorm:"column:reason;type:text;not null;comment:调整依据" json:"reason"`                                                              // 调整依据
	ApprovalID          *int64                             `gorm:"column:approval_id;type:bigint;comment:→admin_approvals" json:"approval_id,string"`                                        // →admin_approvals
	RecoveredCent       int64                              `gorm:"column:recovered_cent;type:bigint;not null;default:0;comment:默认 0；累计实际已追回金额，可信渠道结果确认" json:"recovered_cent"`               // 默认 0；累计实际已追回金额，可信渠道结果确认
}

func (SettlementAdjustment) TableName() string { return "settlement_adjustments" }

// SettlementBatch 结算批次。
type SettlementBatch struct {
	Base
	BatchNo          string                `gorm:"column:batch_no;type:varchar(40);not null;comment:业务批次号" json:"batch_no"`                                                     // 业务批次号
	CoachID          int64                 `gorm:"column:coach_id;type:bigint;not null;comment:→coaches" json:"coach_id,string"`                                                // →coaches
	AccountID        int64                 `gorm:"column:account_id;type:bigint;not null;comment:→settlement_accounts" json:"account_id,string"`                                // →settlement_accounts
	Provider         string                `gorm:"column:provider;type:varchar(32);not null;comment:渠道" json:"provider"`                                                        // 渠道
	AccountSnapshot  JSONB                 `gorm:"column:account_snapshot;type:jsonb;not null;comment:渠道账户引用与脱敏收款方快照" json:"-"`                                                 // 渠道账户引用与脱敏收款方快照
	TotalCent        int64                 `gorm:"column:total_cent;type:bigint;not null;comment:本批应付，>0" json:"total_cent"`                                                    // 本批应付，>0
	ItemCount        int32                 `gorm:"column:item_count;type:integer;not null;comment:明细数量" json:"item_count"`                                                      // 明细数量
	Status           SettlementBatchStatus `gorm:"column:status;type:varchar(24);not null;comment:CREATED/SUBMITTED/PROCESSING/PARTIAL/SETTLED/FAILED/CANCELLED" json:"status"` // CREATED/SUBMITTED/PROCESSING/PARTIAL/SETTLED/FAILED/CANCELLED
	ProviderBatchRef *string               `gorm:"column:provider_batch_ref;type:varchar(128);comment:渠道批号" json:"provider_batch_ref"`                                          // 渠道批号
	SubmittedAt      *time.Time            `gorm:"column:submitted_at;type:timestamptz;comment:提交时间" json:"submitted_at"`                                                       // 提交时间
	FinishedAt       *time.Time            `gorm:"column:finished_at;type:timestamptz;comment:终态时间" json:"finished_at"`                                                         // 终态时间
}

func (SettlementBatch) TableName() string { return "settlement_batches" }

// SettlementBatchItem 批次出款尝试明细。
type SettlementBatchItem struct {
	Base
	BatchID          int64                     `gorm:"column:batch_id;type:bigint;not null;comment:→settlement_batches" json:"batch_id,string"`                             // →settlement_batches
	SettlementItemID int64                     `gorm:"column:settlement_item_id;type:bigint;not null;comment:→settlement_items" json:"settlement_item_id,string"`           // →settlement_items
	AttemptNo        int32                     `gorm:"column:attempt_no;type:integer;not null;comment:同一应结明细尝试序号" json:"attempt_no"`                                        // 同一应结明细尝试序号
	AmountCent       int64                     `gorm:"column:amount_cent;type:bigint;not null;comment:本次应付金额，>0" json:"amount_cent"`                                        // 本次应付金额，>0
	RequestNo        string                    `gorm:"column:request_no;type:varchar(64);not null;comment:渠道幂等请求号" json:"request_no"`                                       // 渠道幂等请求号
	Status           SettlementBatchItemStatus `gorm:"column:status;type:varchar(24);not null;comment:CREATED/SUBMITTED/PROCESSING/SUCCESS/FAILED/CANCELLED" json:"status"` // CREATED/SUBMITTED/PROCESSING/SUCCESS/FAILED/CANCELLED
	ProviderItemRef  *string                   `gorm:"column:provider_item_ref;type:varchar(128);comment:渠道逐条流水" json:"provider_item_ref"`                                  // 渠道逐条流水
	FailureCode      *string                   `gorm:"column:failure_code;type:varchar(64);comment:失败原因" json:"failure_code"`                                               // 失败原因
	SucceededAt      *time.Time                `gorm:"column:succeeded_at;type:timestamptz;comment:成功时间" json:"succeeded_at"`                                               // 成功时间
	CoachID          int64                     `gorm:"column:coach_id;type:bigint;not null;comment:→coaches；与结算批次及应结明细一致" json:"coach_id,string"`                           // →coaches；与结算批次及应结明细一致
}

func (SettlementBatchItem) TableName() string { return "settlement_batch_items" }

// SettlementSupplement 已结算订单的独立补款指令，不占用首次出款明细。
// 一个已批准的正向调整只生成一条指令；失败重试仍沿用原 RequestNo，未知结果先查单。
// 创建/提交时锁应结记录，校验首次出款成功、调整归属及可补金额；成功与累计已付同事务。
type SettlementSupplement struct {
	Base
	SupplementNo     string                     `gorm:"column:supplement_no;type:varchar(40);not null;comment:独立补款业务号" json:"supplement_no"`                      // 独立补款业务号
	SettlementItemID int64                      `gorm:"column:settlement_item_id;type:bigint;not null;comment:应结基准" json:"settlement_item_id,string"`             // 应结基准
	AdjustmentID     int64                      `gorm:"column:adjustment_id;type:bigint;not null;comment:唯一正向调整；补款金额不得超出调整净增额及应结剩余额" json:"adjustment_id,string"` // 唯一正向调整，必须属于同一应结记录
	CoachID          int64                      `gorm:"column:coach_id;type:bigint;not null;comment:与应结记录及账户一致" json:"coach_id,string"`                           // 与应结记录及账户一致
	AccountID        int64                      `gorm:"column:account_id;type:bigint;not null;comment:渠道收款账户" json:"account_id,string"`                           // 渠道收款账户
	Provider         string                     `gorm:"column:provider;type:varchar(32);not null;comment:真实结算渠道" json:"provider"`                                 // 真实结算渠道
	AccountSnapshot  JSONB                      `gorm:"column:account_snapshot;type:jsonb;not null;comment:提交前冻结的渠道收款信息" json:"-"`                                // 提交前冻结的渠道收款信息
	AmountCent       int64                      `gorm:"column:amount_cent;type:bigint;not null;comment:本指令补款额，>0，提交后不可变" json:"amount_cent"`                      // 本指令补款额，>0，提交后不可变
	RequestNo        string                     `gorm:"column:request_no;type:varchar(64);not null;comment:渠道幂等请求号；失败重试复用，未知结果先查单" json:"request_no"`             // 渠道幂等请求号，同一指令重试复用
	Status           SettlementSupplementStatus `gorm:"column:status;type:varchar(24);not null;default:'PENDING';comment:PENDING/SUBMITTED/PROCESSING/SUCCESS/FAILED" json:"status"`
	ApprovalID       int64                      `gorm:"column:approval_id;type:bigint;not null;comment:经参数/金额核验的双人复核批准" json:"approval_id,string"` // 经参数/金额核验的双人复核批准
	ProviderRef      *string                    `gorm:"column:provider_ref;type:varchar(128);comment:渠道真实出款引用" json:"provider_ref"`                // 渠道真实出款引用
	SubmittedAt      *time.Time                 `gorm:"column:submitted_at;type:timestamptz;comment:提交时间" json:"submitted_at"`
	SucceededAt      *time.Time                 `gorm:"column:succeeded_at;type:timestamptz;comment:成功时间" json:"succeeded_at"`
	FailureCode      *string                    `gorm:"column:failure_code;type:varchar(64);comment:失败代码" json:"failure_code"`
	Version          int32                      `gorm:"column:version;type:integer;not null;default:1;comment:乐观锁版本" json:"version"`
}

func (SettlementSupplement) TableName() string { return "settlement_supplements" }
