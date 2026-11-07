package model

import (
	"fmt"
	"time"
)

const OrderEventServiceCompletionSubmitted = "SERVICE_COMPLETION_SUBMITTED"

// ServiceCompletionPayload 是 OrderEvent.Payload 在提交完成事件中的固定结构，不是额外一张表。
// 数量单位读取订单快照；实际服务真实性仍依赖买家确认及争议取证。
type ServiceCompletionPayload struct {
	SchemaVersion     int32     `json:"schema_version"`
	CompletedQuantity int32     `json:"completed_quantity"`
	ActualEndAt       time.Time `json:"actual_end_at"`
	EvidenceAssetIDs  []int64   `json:"evidence_asset_ids"`
	CompletionNote    string    `json:"completion_note"`
}

// Validate 仅校验格式。调用方还须校验订单状态、单位、实际开始/结束时间、数量范围及证据归属。
func (p ServiceCompletionPayload) Validate() error {
	if p.SchemaVersion != 1 {
		return fmt.Errorf("unsupported completion schema version: %d", p.SchemaVersion)
	}
	if p.CompletedQuantity <= 0 {
		return fmt.Errorf("completed quantity must be positive")
	}
	if p.ActualEndAt.IsZero() {
		return fmt.Errorf("actual end time is required")
	}
	for _, id := range p.EvidenceAssetIDs {
		if id <= 0 {
			return fmt.Errorf("evidence asset ID must be positive")
		}
	}
	return nil
}
