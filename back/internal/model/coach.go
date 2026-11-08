package model

import (
	"time"
)

// Coach 陪练师档案。
type Coach struct {
	Base
	UserID          int64       `gorm:"column:user_id;type:bigint;not null;comment:→users" json:"user_id,string"`                                                               // →users
	DisplayName     string      `gorm:"column:display_name;type:varchar(64);not null;comment:已审核展示名" json:"display_name"`                                                       // 已审核展示名
	Bio             string      `gorm:"column:bio;type:text;not null;default:'';comment:默认空；已审核简介" json:"bio"`                                                                  // 默认空；已审核简介
	Status          CoachStatus `gorm:"column:status;type:varchar(24);not null;default:'APPLYING';comment:默认 APPLYING；APPLYING/ACTIVE/SUSPENDED/REJECTED/CLOSED" json:"status"` // 默认 APPLYING；APPLYING/ACTIVE/SUSPENDED/REJECTED/CLOSED
	AcceptingOrders bool        `gorm:"column:accepting_orders;type:boolean;not null;default:false;comment:默认 false；接单开关" json:"accepting_orders"`                              // 默认 false；接单开关
	LastActiveAt    *time.Time  `gorm:"column:last_active_at;type:timestamptz;comment:最近活跃；在线展示结合心跳 TTL" json:"last_active_at"`                                                 // 最近活跃；在线展示结合心跳 TTL
	RatingSum       int64       `gorm:"column:rating_sum;type:bigint;not null;default:0;comment:默认 0；已公开主评星级和" json:"rating_sum"`                                               // 默认 0；已公开主评星级和
	RatingCount     int64       `gorm:"column:rating_count;type:bigint;not null;default:0;comment:默认 0；已公开主评数" json:"rating_count"`                                             // 默认 0；已公开主评数
	CompletedOrders int64       `gorm:"column:completed_orders;type:bigint;not null;default:0;comment:默认 0；完成订单数缓存" json:"completed_orders"`                                    // 默认 0；完成订单数缓存
	Version         int32       `gorm:"column:version;type:integer;not null;default:1;comment:默认 1" json:"version"`                                                             // 默认 1
	ContentVersion  int32       `gorm:"column:content_version;type:integer;not null;default:1;comment:默认 1；内容审核版本，提交新版本时递增；独立于登录或价格版本" json:"content_version"`                  // 默认 1；内容审核版本，提交新版本时递增；独立于登录或价格版本
}

func (Coach) TableName() string { return "coaches" }

// CoachApplication 入驻申请历史。
type CoachApplication struct {
	Base
	CoachID          int64                  `gorm:"column:coach_id;type:bigint;not null;comment:→coaches" json:"coach_id,string"`                                        // →coaches
	Revision         int32                  `gorm:"column:revision;type:integer;not null;comment:申请版本，从 1 开始" json:"revision"`                                           // 申请版本，从 1 开始
	IdentityCheckID  *int64                 `gorm:"column:identity_check_id;type:bigint;comment:→identity_checks" json:"identity_check_id,string"`                       // →identity_checks
	SubmittedPayload JSONB                  `gorm:"column:submitted_payload;type:jsonb;not null;comment:不可变申请资料及附件引用；敏感值不明文保存" json:"-"`                                 // 不可变申请资料及附件引用；敏感值不明文保存
	Status           CoachApplicationStatus `gorm:"column:status;type:varchar(24);not null;comment:DRAFT/SUBMITTED/REVIEWING/APPROVED/REJECTED/WITHDRAWN" json:"status"` // DRAFT/SUBMITTED/REVIEWING/APPROVED/REJECTED/WITHDRAWN
	SubmittedAt      *time.Time             `gorm:"column:submitted_at;type:timestamptz;comment:提交时间" json:"submitted_at"`                                               // 提交时间
	ReviewerID       *int64                 `gorm:"column:reviewer_id;type:bigint;comment:→admin_users" json:"reviewer_id,string"`                                       // →admin_users
	ReviewedAt       *time.Time             `gorm:"column:reviewed_at;type:timestamptz;comment:审核时间" json:"reviewed_at"`                                                 // 审核时间
	RejectReason     *string                `gorm:"column:reject_reason;type:text;comment:拒绝原因" json:"reject_reason"`                                                    // 拒绝原因
}

func (CoachApplication) TableName() string { return "coach_applications" }

// IdentityCheck 实名核验记录。
type IdentityCheck struct {
	Base
	UserID            int64                       `gorm:"column:user_id;type:bigint;not null;comment:→users" json:"user_id,string"`                                              // →users
	Provider          string                      `gorm:"column:provider;type:varchar(32);not null;comment:核验渠道" json:"provider"`                                                // 核验渠道
	ProviderRequestID string                      `gorm:"column:provider_request_id;type:varchar(128);not null;comment:渠道请求号" json:"provider_request_id"`                        // 渠道请求号
	Status            IdentityCheckStatus         `gorm:"column:status;type:varchar(24);not null;comment:PENDING/PASSED/FAILED/EXPIRED/REVOKED" json:"status"`                   // PENDING/PASSED/FAILED/EXPIRED/REVOKED
	NameMasked        *string                     `gorm:"column:name_masked;type:varchar(64);comment:脱敏姓名" json:"name_masked"`                                                   // 脱敏姓名
	AgeCheckResult    IdentityCheckAgeCheckResult `gorm:"column:age_check_result;type:varchar(24);not null;comment:UNKNOWN/ELIGIBLE/INELIGIBLE；年龄阈值按规则" json:"age_check_result"` // UNKNOWN/ELIGIBLE/INELIGIBLE；年龄阈值按规则
	ResultRef         *string                     `gorm:"column:result_ref;type:varchar(256);comment:渠道结果引用" json:"result_ref"`                                                  // 渠道结果引用
	VerifiedAt        *time.Time                  `gorm:"column:verified_at;type:timestamptz;comment:成功时间" json:"verified_at"`                                                   // 成功时间
	ExpiresAt         *time.Time                  `gorm:"column:expires_at;type:timestamptz;comment:有效期" json:"expires_at"`                                                      // 有效期
	ReviewerID        *int64                      `gorm:"column:reviewer_id;type:bigint;comment:→admin_users" json:"reviewer_id,string"`                                         // →admin_users
	ReviewNote        *string                     `gorm:"column:review_note;type:text;comment:脱敏复核意见" json:"review_note"`                                                        // 脱敏复核意见
	RevokedAt         *time.Time                  `gorm:"column:revoked_at;type:timestamptz;comment:核验明确撤销时间；普通失败重试不自动撤销旧通过" json:"revoked_at"`                                  // 核验明确撤销时间；普通失败重试不自动撤销旧通过
	RevokeReason      *string                     `gorm:"column:revoke_reason;type:text;comment:撤销原因" json:"revoke_reason"`                                                      // 撤销原因
}

func (IdentityCheck) TableName() string { return "identity_checks" }

// CoachSkill 游戏技能认证版本。
type CoachSkill struct {
	Base
	CoachID          int64            `gorm:"column:coach_id;type:bigint;not null;comment:→coaches" json:"coach_id,string"`                                      // →coaches
	GameID           int64            `gorm:"column:game_id;type:bigint;not null;comment:→games" json:"game_id,string"`                                          // →games
	RankID           int64            `gorm:"column:rank_id;type:bigint;not null;comment:→game_ranks；须属于 game_id" json:"rank_id,string"`                         // →game_ranks；须属于 game_id
	Revision         int32            `gorm:"column:revision;type:integer;not null;comment:同陪练同游戏递增" json:"revision"`                                            // 同陪练同游戏递增
	Status           CoachSkillStatus `gorm:"column:status;type:varchar(24);not null;comment:DRAFT/SUBMITTED/REVIEWING/APPROVED/REJECTED/EXPIRED" json:"status"` // DRAFT/SUBMITTED/REVIEWING/APPROVED/REJECTED/EXPIRED
	IsCurrent        bool             `gorm:"column:is_current;type:boolean;not null;default:false;comment:默认 false；当前有效通过版本" json:"is_current"`                 // 默认 false；当前有效通过版本
	EvidenceAssetIds JSONB            `gorm:"column:evidence_asset_ids;type:jsonb;not null;default:'[]'::jsonb;comment:默认 []；私有证明资源 ID 数组" json:"-"`             // 默认 []；私有证明资源 ID 数组
	ReviewerID       *int64           `gorm:"column:reviewer_id;type:bigint;comment:→admin_users" json:"reviewer_id,string"`                                     // →admin_users
	ReviewNote       *string          `gorm:"column:review_note;type:text;comment:审核意见" json:"review_note"`                                                      // 审核意见
	VerifiedAt       *time.Time       `gorm:"column:verified_at;type:timestamptz;comment:认证通过时间" json:"verified_at"`                                             // 认证通过时间
	ExpiresAt        *time.Time       `gorm:"column:expires_at;type:timestamptz;comment:认证有效期" json:"expires_at"`                                                // 认证有效期
}

func (CoachSkill) TableName() string { return "coach_skills" }
