package model

import (
	"time"
)

// Activity 展示活动。
type Activity struct {
	Base
	Title          string         `gorm:"column:title;type:varchar(128);not null;comment:标题" json:"title"`                                                       // 标题
	Content        JSONB          `gorm:"column:content;type:jsonb;not null;comment:结构化展示内容" json:"content"`                                                     // 结构化展示内容
	CoverAssetID   *int64         `gorm:"column:cover_asset_id;type:bigint;comment:→media_assets" json:"cover_asset_id,string"`                                  // →media_assets
	StartsAt       time.Time      `gorm:"column:starts_at;type:timestamptz;not null;comment:开始" json:"starts_at"`                                                // 开始
	EndsAt         time.Time      `gorm:"column:ends_at;type:timestamptz;not null;comment:结束" json:"ends_at"`                                                    // 结束
	Status         ActivityStatus `gorm:"column:status;type:varchar(16);not null;comment:DRAFT/PUBLISHED/OFFLINE" json:"status"`                                 // DRAFT/PUBLISHED/OFFLINE
	CreatedBy      int64          `gorm:"column:created_by;type:bigint;not null;comment:→admin_users" json:"created_by"`                                         // →admin_users
	ContentVersion int32          `gorm:"column:content_version;type:integer;not null;default:1;comment:默认 1；内容审核版本，提交新版本时递增；独立于登录或价格版本" json:"content_version"` // 默认 1；内容审核版本，提交新版本时递增；独立于登录或价格版本
	Version        int32          `gorm:"column:version;type:integer;not null;default:1;comment:默认 1；更新时通过 WHERE version 校验并显式递增" json:"version"`                // 默认 1；更新时通过 WHERE version 校验并显式递增
}

func (Activity) TableName() string { return "activities" }

// Banner 首页 Banner 与运营位。
type Banner struct {
	Base
	Placement    string           `gorm:"column:placement;type:varchar(48);not null;comment:首页位置编码" json:"placement"`                          // 首页位置编码
	Title        string           `gorm:"column:title;type:varchar(128);not null;comment:标题" json:"title"`                                     // 标题
	ImageAssetID int64            `gorm:"column:image_asset_id;type:bigint;not null;comment:→media_assets" json:"image_asset_id,string"`       // →media_assets
	TargetType   BannerTargetType `gorm:"column:target_type;type:varchar(24);not null;comment:GAME/COACH/SERVICE/ACTIVITY" json:"target_type"` // GAME/COACH/SERVICE/ACTIVITY
	TargetID     int64            `gorm:"column:target_id;type:bigint;not null;comment:对应目标，逻辑关联" json:"target_id,string"`                     // 对应目标，逻辑关联
	StartsAt     time.Time        `gorm:"column:starts_at;type:timestamptz;not null;comment:开始" json:"starts_at"`                              // 开始
	EndsAt       time.Time        `gorm:"column:ends_at;type:timestamptz;not null;comment:结束" json:"ends_at"`                                  // 结束
	Status       BannerStatus     `gorm:"column:status;type:varchar(16);not null;comment:DRAFT/PUBLISHED/OFFLINE" json:"status"`               // DRAFT/PUBLISHED/OFFLINE
	SortOrder    int32            `gorm:"column:sort_order;type:integer;not null;default:0;comment:默认 0" json:"sort_order"`                    // 默认 0
}

func (Banner) TableName() string { return "banners" }

// CoachRecommendation 首页推荐陪练。
type CoachRecommendation struct {
	Base
	Placement string                    `gorm:"column:placement;type:varchar(48);not null;comment:推荐位" json:"placement"`          // 推荐位
	CoachID   int64                     `gorm:"column:coach_id;type:bigint;not null;comment:→coaches" json:"coach_id,string"`     // →coaches
	GameID    *int64                    `gorm:"column:game_id;type:bigint;comment:→games；空表示全站" json:"game_id,string"`            // →games；空表示全站
	StartsAt  time.Time                 `gorm:"column:starts_at;type:timestamptz;not null;comment:开始" json:"starts_at"`           // 开始
	EndsAt    time.Time                 `gorm:"column:ends_at;type:timestamptz;not null;comment:结束" json:"ends_at"`               // 结束
	Status    CoachRecommendationStatus `gorm:"column:status;type:varchar(16);not null;comment:ENABLED/DISABLED" json:"status"`   // ENABLED/DISABLED
	SortOrder int32                     `gorm:"column:sort_order;type:integer;not null;default:0;comment:默认 0" json:"sort_order"` // 默认 0
}

func (CoachRecommendation) TableName() string { return "coach_recommendations" }

// ContentAudit 内容审核任务与提交版本。
type ContentAudit struct {
	Base
	SubjectType        ContentAuditSubjectType `gorm:"column:subject_type;type:varchar(32);not null;comment:USER/COACH/SERVICE/SKILL/REVIEW_CONTENT/ASSET/ACTIVITY" json:"-"`        // USER/COACH/SERVICE/SKILL/REVIEW_CONTENT/ASSET/ACTIVITY
	SubjectID          int64                   `gorm:"column:subject_id;type:bigint;not null;comment:业务对象 ID，逻辑关联" json:"-"`                                                         // 业务对象 ID，逻辑关联
	Revision           int32                   `gorm:"column:revision;type:integer;not null;comment:该对象提交版本" json:"revision"`                                                        // 该对象提交版本
	SubmittedByUserID  *int64                  `gorm:"column:submitted_by_user_id;type:bigint;comment:→users" json:"submitted_by_user_id,string"`                                    // →users
	SubmittedByAdminID *int64                  `gorm:"column:submitted_by_admin_id;type:bigint;comment:→admin_users" json:"submitted_by_admin_id,string"`                            // →admin_users
	ContentSnapshot    JSONB                   `gorm:"column:content_snapshot;type:jsonb;not null;comment:本次待审文本及资源引用" json:"-"`                                                     // 本次待审文本及资源引用
	SourceVersion      int32                   `gorm:"column:source_version;type:integer;not null;comment:提交时业务版本，防止旧审核覆盖新内容" json:"source_version"`                                 // 提交时业务版本，防止旧审核覆盖新内容
	ProviderTaskRef    *string                 `gorm:"column:provider_task_ref;type:varchar(128);comment:自动审核任务引用" json:"provider_task_ref"`                                         // 自动审核任务引用
	Status             ContentAuditStatus      `gorm:"column:status;type:varchar(24);not null;comment:PENDING/AUTO_PASSED/MANUAL_REVIEW/APPROVED/REJECTED/SUPERSEDED" json:"status"` // PENDING/AUTO_PASSED/MANUAL_REVIEW/APPROVED/REJECTED/SUPERSEDED
	RiskCodes          JSONB                   `gorm:"column:risk_codes;type:jsonb;not null;default:'[]'::jsonb;comment:默认 []" json:"risk_codes"`                                    // 默认 []
	ReviewerID         *int64                  `gorm:"column:reviewer_id;type:bigint;comment:→admin_users" json:"reviewer_id,string"`                                                // →admin_users
	Reason             *string                 `gorm:"column:reason;type:text;comment:审核依据" json:"reason"`                                                                           // 审核依据
	DecidedAt          *time.Time              `gorm:"column:decided_at;type:timestamptz;comment:决策时间" json:"decided_at"`                                                            // 决策时间
}

func (ContentAudit) TableName() string { return "content_audits" }

// UserPenalty 违规处罚与恢复。
type UserPenalty struct {
	Base
	UserID    int64             `gorm:"column:user_id;type:bigint;not null;comment:→users" json:"user_id,string"`                           // →users
	CoachID   *int64            `gorm:"column:coach_id;type:bigint;comment:→coaches" json:"coach_id,string"`                                // →coaches
	ServiceID *int64            `gorm:"column:service_id;type:bigint;comment:→coach_services" json:"service_id,string"`                     // →coach_services
	CaseID    *int64            `gorm:"column:case_id;type:bigint;comment:→after_sales_cases" json:"case_id,string"`                        // →after_sales_cases
	Action    UserPenaltyAction `gorm:"column:action;type:varchar(24);not null;comment:BAN_USER/SUSPEND_COACH/BLOCK_SERVICE" json:"action"` // BAN_USER/SUSPEND_COACH/BLOCK_SERVICE
	Reason    string            `gorm:"column:reason;type:text;not null;comment:处罚依据" json:"reason"`                                        // 处罚依据
	StartsAt  time.Time         `gorm:"column:starts_at;type:timestamptz;not null;comment:生效" json:"starts_at"`                             // 生效
	EndsAt    *time.Time        `gorm:"column:ends_at;type:timestamptz;comment:到期；空为长期" json:"ends_at"`                                     // 到期；空为长期
	Status    UserPenaltyStatus `gorm:"column:status;type:varchar(16);not null;comment:ACTIVE/REVOKED/EXPIRED" json:"status"`               // ACTIVE/REVOKED/EXPIRED
	AdminID   int64             `gorm:"column:admin_id;type:bigint;not null;comment:→admin_users" json:"admin_id,string"`                   // →admin_users
	RevokedAt *time.Time        `gorm:"column:revoked_at;type:timestamptz;comment:撤销时间" json:"revoked_at"`                                  // 撤销时间
}

func (UserPenalty) TableName() string { return "user_penalties" }
