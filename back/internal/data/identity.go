package data

import (
	"context"
	"errors"
	"time"

	"dd_web/back/internal/biz"
	"dd_web/back/internal/model"
	"gorm.io/gorm"
)

type IdentityRepo struct{ data *Data }

func NewIdentityRepo(data *Data) *IdentityRepo { return &IdentityRepo{data: data} }

func (r *IdentityRepo) CreateIdentityCheck(ctx context.Context, record biz.IdentityRecord) (biz.IdentityRecord, error) {
	row := model.IdentityCheck{
		UserID: record.UserID, Provider: record.Provider, ProviderRequestID: record.ProviderRequestID,
		ConsentVersion: record.ConsentVersion, Status: model.IdentityCheckStatus(record.Status),
		AgeCheckResult: model.IdentityCheckAgeCheckResult(record.AgeCheckResult),
		VerifiedAt:     record.VerifiedAt, ExpiresAt: record.ExpiresAt,
	}
	if record.NameMasked != "" {
		row.NameMasked = &record.NameMasked
	}
	if record.ResultRef != "" {
		row.ResultRef = &record.ResultRef
	}
	if err := r.data.db.WithContext(ctx).Create(&row).Error; err != nil {
		return biz.IdentityRecord{}, err
	}
	return identityRecord(row), nil
}

func (r *IdentityRepo) CurrentIdentityCheck(ctx context.Context, userID int64) (biz.IdentityRecord, error) {
	var row model.IdentityCheck
	err := r.data.db.WithContext(ctx).Where("user_id = ?", userID).Order("id DESC").First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return biz.IdentityRecord{}, ErrNotFound
	}
	if err != nil {
		return biz.IdentityRecord{}, err
	}
	return identityRecord(row), nil
}

func (r *IdentityRepo) ValidAdultIdentityCheck(ctx context.Context, userID int64, now time.Time) (biz.IdentityRecord, error) {
	var row model.IdentityCheck
	err := r.data.db.WithContext(ctx).
		Where("user_id = ? AND status = ? AND age_check_result = ? AND verified_at IS NOT NULL AND revoked_at IS NULL AND (expires_at IS NULL OR expires_at > ?)",
			userID, model.IdentityCheckStatusPassed, model.IdentityCheckAgeCheckResultEligible, now).
		Order("verified_at DESC, id DESC").First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return biz.IdentityRecord{}, biz.ErrAdultIdentityRequired
	}
	if err != nil {
		return biz.IdentityRecord{}, err
	}
	return identityRecord(row), nil
}

func identityRecord(row model.IdentityCheck) biz.IdentityRecord {
	record := biz.IdentityRecord{
		ID: row.ID, UserID: row.UserID, Provider: row.Provider, ProviderRequestID: row.ProviderRequestID,
		ConsentVersion: row.ConsentVersion, Status: string(row.Status), AgeCheckResult: string(row.AgeCheckResult),
		VerifiedAt: row.VerifiedAt, ExpiresAt: row.ExpiresAt, CreatedAt: row.CreatedAt,
	}
	if row.NameMasked != nil {
		record.NameMasked = *row.NameMasked
	}
	if row.ResultRef != nil {
		record.ResultRef = *row.ResultRef
	}
	return record
}
