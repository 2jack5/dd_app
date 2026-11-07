package model

import (
	"time"
)

// BusinessRuleVersion 业务规则版本。
type BusinessRuleVersion struct {
	Base
	Code        BusinessRuleVersionCode   `gorm:"column:code;type:varchar(48);not null" json:"code"`                 // ORDER/CANCELLATION/SETTLEMENT/SUPPORT 等
	Revision    int32                     `gorm:"column:revision;type:integer;not null" json:"revision"`             // 版本
	Config      JSONB                     `gorm:"column:config;type:jsonb;not null" json:"config"`                   // 具备校验 schema 的规则
	Status      BusinessRuleVersionStatus `gorm:"column:status;type:varchar(16);not null" json:"status"`             // DRAFT/PUBLISHED/RETIRED
	EffectiveAt time.Time                 `gorm:"column:effective_at;type:timestamptz;not null" json:"effective_at"` // 生效
	ExpiresAt   *time.Time                `gorm:"column:expires_at;type:timestamptz" json:"expires_at"`              // 结束
	CreatedBy   int64                     `gorm:"column:created_by;type:bigint;not null" json:"created_by"`          // →admin_users
	ApprovalID  *int64                    `gorm:"column:approval_id;type:bigint" json:"approval_id,string"`          // →admin_approvals
	Version     int32                     `gorm:"column:version;type:integer;not null;default:1" json:"version"`     // 默认 1；更新时通过 WHERE version 校验并显式递增
}

func (BusinessRuleVersion) TableName() string { return "business_rule_versions" }

// IdempotencyKey 请求幂等。
type IdempotencyKey struct {
	Base
	ActorType       IdempotencyKeyActorType `gorm:"column:actor_type;type:varchar(16);not null" json:"actor_type"`                           // USER/ADMIN/SYSTEM
	ActorID         int64                   `gorm:"column:actor_id;type:bigint;not null" json:"actor_id,string"`                             // 主体 ID；系统使用固定主体
	Operation       string                  `gorm:"column:operation;type:varchar(64);not null" json:"operation"`                             // 接口操作
	BizKey          string                  `gorm:"column:biz_key;type:varchar(128);not null" json:"biz_key"`                                // 客户端幂等键
	RequestHash     string                  `gorm:"column:request_hash;type:varchar(64);not null" json:"request_hash"`                       // 规范化请求摘要
	Status          IdempotencyKeyStatus    `gorm:"column:status;type:varchar(16);not null" json:"status"`                                   // PROCESSING/SUCCEEDED/FAILED
	ResourceType    *string                 `gorm:"column:resource_type;type:varchar(32)" json:"resource_type"`                              // 结果资源
	ResourceID      *int64                  `gorm:"column:resource_id;type:bigint" json:"resource_id,string"`                                // 结果 ID
	ResponseSummary JSONB                   `gorm:"column:response_summary;type:jsonb;not null;default:'{}'::jsonb" json:"response_summary"` // 默认 {}；脱敏可重放结果
	ExpiresAt       time.Time               `gorm:"column:expires_at;type:timestamptz;not null" json:"expires_at"`                           // 过期清理时间
}

func (IdempotencyKey) TableName() string { return "idempotency_keys" }

// OutboxEvent 事务消息与后台任务。
type OutboxEvent struct {
	Base
	EventKey      string                   `gorm:"column:event_key;type:varchar(160);not null" json:"event_key"`          // 稳定业务事件键
	AggregateType OutboxEventAggregateType `gorm:"column:aggregate_type;type:varchar(32);not null" json:"aggregate_type"` // 业务聚合类型，覆盖交易、审核、售后、通知与对账；见 OutboxEventAggregateType
	AggregateID   int64                    `gorm:"column:aggregate_id;type:bigint;not null" json:"aggregate_id,string"`   // 聚合 ID，逻辑关联
	EventType     string                   `gorm:"column:event_type;type:varchar(64);not null" json:"event_type"`         // 消息类型
	Payload       JSONB                    `gorm:"column:payload;type:jsonb;not null;default:'{}'::jsonb" json:"-"`       // 默认 {}；不存密钥及明文隐私
	Status        OutboxEventStatus        `gorm:"column:status;type:varchar(16);not null" json:"status"`                 // PENDING/PROCESSING/DONE/DEAD
	AvailableAt   time.Time                `gorm:"column:available_at;type:timestamptz;not null" json:"available_at"`     // 下次可执行时间
	LockedUntil   *time.Time               `gorm:"column:locked_until;type:timestamptz" json:"locked_until"`              // 任务租约截止
	LockToken     *string                  `gorm:"column:lock_token;type:varchar(64)" json:"lock_token"`                  // 本轮认领标识
	Attempts      int32                    `gorm:"column:attempts;type:integer;not null;default:0" json:"attempts"`       // 默认 0
	LastError     *string                  `gorm:"column:last_error;type:text" json:"last_error"`                         // 脱敏错误
}

func (OutboxEvent) TableName() string { return "outbox_events" }
