package model

import (
	"time"
)

// AfterSalesCase 投诉、客服与申诉工单。
type AfterSalesCase struct {
	Base
	CaseNo           string               `gorm:"column:case_no;type:varchar(40);not null" json:"case_no"`                               // 工单编号
	Kind             AfterSalesCaseKind   `gorm:"column:kind;type:varchar(24);not null" json:"kind"`                                     // COMPLAINT/SUPPORT/REVIEW_APPEAL
	ReporterID       int64                `gorm:"column:reporter_id;type:bigint;not null" json:"reporter_id,string"`                     // →users
	OrderID          *int64               `gorm:"column:order_id;type:bigint" json:"order_id,string"`                                    // →orders；普通帮助可空
	TargetUserID     *int64               `gorm:"column:target_user_id;type:bigint" json:"target_user_id,string"`                        // →users
	ReviewID         *int64               `gorm:"column:review_id;type:bigint" json:"review_id,string"`                                  // →reviews
	CategoryCode     string               `gorm:"column:category_code;type:varchar(48);not null" json:"category_code"`                   // 投诉或问题分类
	Content          string               `gorm:"column:content;type:text;not null" json:"content"`                                      // 问题描述
	EvidenceAssetIds JSONB                `gorm:"column:evidence_asset_ids;type:jsonb;not null;default:'[]'::jsonb" json:"-"`            // 默认 []
	Status           AfterSalesCaseStatus `gorm:"column:status;type:varchar(24);not null" json:"status"`                                 // OPEN/ASSIGNED/INVESTIGATING/WAIT_USER/RESOLVED/CLOSED
	AssigneeID       *int64               `gorm:"column:assignee_id;type:bigint" json:"assignee_id,string"`                              // →admin_users
	BlocksSettlement bool                 `gorm:"column:blocks_settlement;type:boolean;not null;default:false" json:"blocks_settlement"` // 默认 false；是否冻结该订单结算
	ResolutionCode   *string              `gorm:"column:resolution_code;type:varchar(48)" json:"resolution_code"`                        // 处理结论
	ResolutionNote   *string              `gorm:"column:resolution_note;type:text" json:"resolution_note"`                               // 对用户可见结论
	ResolvedAt       *time.Time           `gorm:"column:resolved_at;type:timestamptz" json:"resolved_at"`                                // 结案时间
	Version          int32                `gorm:"column:version;type:integer;not null;default:1" json:"version"`                         // 默认 1；更新时通过 WHERE version 校验并显式递增
}

func (AfterSalesCase) TableName() string { return "after_sales_cases" }

// CaseEvent 工单流转与补充证据。
type CaseEvent struct {
	EventBase
	CaseID           int64               `gorm:"column:case_id;type:bigint;not null" json:"case_id,string"`                  // →after_sales_cases
	ActorType        CaseEventActorType  `gorm:"column:actor_type;type:varchar(16);not null" json:"actor_type"`              // USER/ADMIN/SYSTEM
	ActorID          *int64              `gorm:"column:actor_id;type:bigint" json:"actor_id,string"`                         // 操作者
	EventType        CaseEventEventType  `gorm:"column:event_type;type:varchar(32);not null" json:"event_type"`              // COMMENT/ASSIGN/REQUEST_EVIDENCE/RESOLVE/REOPEN
	Visibility       CaseEventVisibility `gorm:"column:visibility;type:varchar(16);not null" json:"visibility"`              // PUBLIC/INTERNAL
	Content          string              `gorm:"column:content;type:text;not null" json:"content"`                           // 补充说明
	EvidenceAssetIds JSONB               `gorm:"column:evidence_asset_ids;type:jsonb;not null;default:'[]'::jsonb" json:"-"` // 默认 []
	EventKey         string              `gorm:"column:event_key;type:varchar(160);not null" json:"event_key"`               // 唯一事件键
}

func (CaseEvent) TableName() string { return "case_events" }
