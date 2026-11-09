package model_test

import (
	"dd_web/back/internal/model"
	"testing"
	"time"
)

func validCompletion() model.ServiceCompletionPayload {
	return model.ServiceCompletionPayload{
		SchemaVersion:     1,
		CompletedQuantity: 2,
		ActualEndAt:       time.Date(2026, time.November, 9, 12, 0, 0, 0, time.UTC),
		EvidenceAssetIDs:  []int64{10, 20},
	}
}

func TestServiceCompletionAcceptsValidPayload(t *testing.T) {
	if err := validCompletion().Validate(); err != nil {
		t.Fatalf("合法履约数据被拒绝: %v", err)
	}
}

func TestServiceCompletionRejectsUnsupportedSchema(t *testing.T) {
	payload := validCompletion()
	payload.SchemaVersion = 2
	if err := payload.Validate(); err == nil {
		t.Fatal("非法 schema 版本被接受")
	}
}

func TestServiceCompletionRejectsZeroQuantity(t *testing.T) {
	payload := validCompletion()
	payload.CompletedQuantity = 0
	if err := payload.Validate(); err == nil {
		t.Fatal("零完成数量被接受")
	}
}

func TestServiceCompletionRequiresEndTime(t *testing.T) {
	payload := validCompletion()
	payload.ActualEndAt = time.Time{}
	if err := payload.Validate(); err == nil {
		t.Fatal("缺少结束时间的数据被接受")
	}
}

func TestServiceCompletionRejectsInvalidEvidenceID(t *testing.T) {
	payload := validCompletion()
	payload.EvidenceAssetIDs = []int64{-1}
	if err := payload.Validate(); err == nil {
		t.Fatal("非法证据 ID 被接受")
	}
}
