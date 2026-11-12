package data

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"dd_web/back/internal/biz"
	"dd_web/back/internal/model"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type WalletRepo struct{ data *Data }

func NewWalletRepo(data *Data) *WalletRepo { return &WalletRepo{data: data} }

func (r *WalletRepo) ActiveRechargePackages(ctx context.Context, after int64, limit int, now time.Time) ([]model.RechargePackage, error) {
	var rows []model.RechargePackage
	err := r.data.db.WithContext(ctx).
		Where("id > ? AND status = ? AND (effective_at IS NULL OR effective_at <= ?) AND (expires_at IS NULL OR expires_at > ?)", after, model.RechargePackageStatusEnabled, now, now).
		Order("id ASC").Limit(limit + 1).Find(&rows).Error
	return rows, err
}

func (r *WalletRepo) RechargePackage(ctx context.Context, packageID int64, now time.Time) (model.RechargePackage, error) {
	var row model.RechargePackage
	err := r.data.db.WithContext(ctx).
		Where("id = ? AND status = ? AND (effective_at IS NULL OR effective_at <= ?) AND (expires_at IS NULL OR expires_at > ?)", packageID, model.RechargePackageStatusEnabled, now, now).
		First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return row, biz.ErrWalletNotFound
	}
	return row, err
}

func (r *WalletRepo) AdminRechargePackages(ctx context.Context, after int64, limit int, status string) ([]model.RechargePackage, error) {
	query := r.data.db.WithContext(ctx).Where("id > ?", after)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	var rows []model.RechargePackage
	err := query.Order("id ASC").Limit(limit + 1).Find(&rows).Error
	return rows, err
}

func (r *WalletRepo) CreateRechargePackage(ctx context.Context, adminID int64, candidate model.RechargePackage, key, requestHash string) (result model.RechargePackage, replayed bool, err error) {
	err = r.data.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		existing, found, err := lockIdempotency(tx, model.IdempotencyKeyActorTypeAdmin, adminID, "wallet.adminCreateRechargePackage", key)
		if err != nil {
			return err
		}
		if found {
			return replayPackage(tx, existing, requestHash, &result, &replayed)
		}
		idem, err := createIdempotency(tx, model.IdempotencyKeyActorTypeAdmin, adminID, "wallet.adminCreateRechargePackage", key, requestHash)
		if err != nil {
			return err
		}
		candidate.Version = 1
		if err := tx.Create(&candidate).Error; err != nil {
			return err
		}
		if err := finishIdempotency(tx, &idem, candidate.ID, "RECHARGE_PACKAGE"); err != nil {
			return err
		}
		result = candidate
		return nil
	})
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return r.replayAdminPackage(ctx, adminID, "wallet.adminCreateRechargePackage", key, requestHash)
	}
	return result, replayed, err
}

func (r *WalletRepo) UpdateRechargePackage(ctx context.Context, adminID, packageID int64, candidate model.RechargePackage, key, requestHash string) (result model.RechargePackage, replayed bool, err error) {
	err = r.data.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		existing, found, err := lockIdempotency(tx, model.IdempotencyKeyActorTypeAdmin, adminID, "wallet.adminUpdateRechargePackage", key)
		if err != nil {
			return err
		}
		if found {
			return replayPackage(tx, existing, requestHash, &result, &replayed)
		}
		var current model.RechargePackage
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&current, packageID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return biz.ErrWalletNotFound
			}
			return err
		}
		if current.Version != candidate.Version {
			return biz.ErrRechargePackageConflict
		}
		idem, err := createIdempotency(tx, model.IdempotencyKeyActorTypeAdmin, adminID, "wallet.adminUpdateRechargePackage", key, requestHash)
		if err != nil {
			return err
		}
		now := time.Now().UTC()
		updates := map[string]any{"name": candidate.Name, "pay_amount": candidate.PayAmount, "currency": candidate.Currency, "package_multiplier": candidate.PackageMultiplier, "credited_amount": candidate.CreditedAmount, "asset_code": candidate.AssetCode, "status": candidate.Status, "effective_at": candidate.EffectiveAt, "expires_at": candidate.ExpiresAt, "version": current.Version + 1, "updated_at": now}
		updated := tx.Model(&model.RechargePackage{}).Where("id = ? AND version = ?", current.ID, current.Version).Updates(updates)
		if updated.Error != nil {
			return updated.Error
		}
		if updated.RowsAffected != 1 {
			return biz.ErrRechargePackageConflict
		}
		candidate.Base, candidate.Version = current.Base, current.Version+1
		candidate.UpdatedAt = now
		if err := finishIdempotency(tx, &idem, current.ID, "RECHARGE_PACKAGE"); err != nil {
			return err
		}
		result = candidate
		return nil
	})
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return r.replayAdminPackage(ctx, adminID, "wallet.adminUpdateRechargePackage", key, requestHash)
	}
	return result, replayed, err
}

func lockIdempotency(tx *gorm.DB, actorType model.IdempotencyKeyActorType, actorID int64, operation, key string) (model.IdempotencyKey, bool, error) {
	var row model.IdempotencyKey
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("actor_type = ? AND actor_id = ? AND operation = ? AND biz_key = ?", actorType, actorID, operation, key).First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return row, false, nil
	}
	return row, err == nil, err
}

func createIdempotency(tx *gorm.DB, actorType model.IdempotencyKeyActorType, actorID int64, operation, key, requestHash string) (model.IdempotencyKey, error) {
	row := model.IdempotencyKey{ActorType: actorType, ActorID: actorID, Operation: operation, BizKey: key, RequestHash: requestHash, Status: model.IdempotencyKeyStatusProcessing, ResponseSummary: model.JSONB(`{}`), ExpiresAt: time.Now().UTC().Add(24 * time.Hour)}
	err := tx.Create(&row).Error
	return row, err
}

func finishIdempotency(tx *gorm.DB, row *model.IdempotencyKey, resourceID int64, resourceType string) error {
	return tx.Model(row).Updates(map[string]any{"status": model.IdempotencyKeyStatusSucceeded, "resource_type": resourceType, "resource_id": resourceID, "updated_at": time.Now().UTC()}).Error
}

func replayPackage(tx *gorm.DB, idem model.IdempotencyKey, requestHash string, result *model.RechargePackage, replayed *bool) error {
	if idem.RequestHash != requestHash {
		return biz.ErrIdempotencyConflict
	}
	if idem.ResourceID == nil {
		return biz.ErrRechargeStateConflict
	}
	if err := tx.First(result, *idem.ResourceID).Error; err != nil {
		return err
	}
	*replayed = true
	return nil
}

func (r *WalletRepo) replayAdminPackage(ctx context.Context, adminID int64, operation, key, requestHash string) (model.RechargePackage, bool, error) {
	var idem model.IdempotencyKey
	err := r.data.db.WithContext(ctx).Where("actor_type = ? AND actor_id = ? AND operation = ? AND biz_key = ?", model.IdempotencyKeyActorTypeAdmin, adminID, operation, key).First(&idem).Error
	if err != nil {
		return model.RechargePackage{}, false, err
	}
	var result model.RechargePackage
	replayed := false
	err = replayPackage(r.data.db.WithContext(ctx), idem, requestHash, &result, &replayed)
	return result, replayed, err
}

func (r *WalletRepo) Wallet(ctx context.Context, userID int64, assetCode string) (*model.WalletAccount, error) {
	var row model.WalletAccount
	err := r.data.db.WithContext(ctx).Where("user_id = ? AND asset_code = ?", userID, assetCode).First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *WalletRepo) CreateRecharge(ctx context.Context, packageRow model.RechargePackage, candidate model.WalletRecharge, key, requestHash string) (result model.WalletRecharge, replayed bool, err error) {
	err = r.data.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var existing model.IdempotencyKey
		query := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("actor_type = ? AND actor_id = ? AND operation = ? AND biz_key = ?", model.IdempotencyKeyActorTypeUser, candidate.UserID, "wallet.createRecharge", key)
		findErr := query.First(&existing).Error
		if findErr == nil {
			if existing.RequestHash != requestHash {
				return biz.ErrIdempotencyConflict
			}
			if existing.ResourceID == nil {
				return biz.ErrRechargeStateConflict
			}
			if err := tx.First(&result, *existing.ResourceID).Error; err != nil {
				return err
			}
			replayed = true
			return nil
		}
		if !errors.Is(findErr, gorm.ErrRecordNotFound) {
			return findErr
		}

		var current model.RechargePackage
		now := time.Now().UTC()
		if err := tx.Clauses(clause.Locking{Strength: "SHARE"}).Where("id = ? AND version = ? AND status = ? AND (effective_at IS NULL OR effective_at <= ?) AND (expires_at IS NULL OR expires_at > ?)", packageRow.ID, packageRow.Version, model.RechargePackageStatusEnabled, now, now).First(&current).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return biz.ErrWalletNotFound
			}
			return err
		}
		idem := model.IdempotencyKey{ActorType: model.IdempotencyKeyActorTypeUser, ActorID: candidate.UserID, Operation: "wallet.createRecharge", BizKey: key, RequestHash: requestHash, Status: model.IdempotencyKeyStatusProcessing, ResponseSummary: model.JSONB(`{}`), ExpiresAt: time.Now().UTC().Add(24 * time.Hour)}
		if err := tx.Create(&idem).Error; err != nil {
			return err
		}
		if err := tx.Create(&candidate).Error; err != nil {
			return err
		}
		resourceType := "WALLET_RECHARGE"
		if err := tx.Model(&idem).Updates(map[string]any{"status": model.IdempotencyKeyStatusSucceeded, "resource_type": resourceType, "resource_id": candidate.ID, "updated_at": time.Now().UTC()}).Error; err != nil {
			return err
		}
		result = candidate
		return nil
	})
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		var idem model.IdempotencyKey
		lookup := r.data.db.WithContext(ctx).Where("actor_type = ? AND actor_id = ? AND operation = ? AND biz_key = ?", model.IdempotencyKeyActorTypeUser, candidate.UserID, "wallet.createRecharge", key).First(&idem).Error
		if lookup == nil && idem.RequestHash != requestHash {
			return result, false, biz.ErrIdempotencyConflict
		}
		if lookup == nil && idem.ResourceID != nil {
			lookup = r.data.db.WithContext(ctx).First(&result, *idem.ResourceID).Error
			return result, true, lookup
		}
	}
	return result, replayed, err
}

func (r *WalletRepo) MarkRechargePrepayCreated(ctx context.Context, rechargeID int64, prepay biz.RechargePrepayResult) (model.WalletRecharge, error) {
	var row model.WalletRecharge
	err := r.data.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&row, rechargeID).Error; err != nil {
			return err
		}
		if row.Status == model.WalletRechargeStatusPrepayCreated {
			return nil
		}
		if row.Status != model.WalletRechargeStatusPending {
			return biz.ErrRechargeStateConflict
		}
		expiresAt := prepay.ExpiresAt
		if expiresAt.IsZero() {
			expiresAt = row.ExpiresAt
		}
		updates := map[string]any{"status": model.WalletRechargeStatusPrepayCreated, "prepay_ref_cipher": prepay.PrepayRefCipher, "expires_at": expiresAt, "version": row.Version + 1, "updated_at": time.Now().UTC()}
		result := tx.Model(&model.WalletRecharge{}).Where("id = ? AND version = ?", row.ID, row.Version).Updates(updates)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected != 1 {
			return biz.ErrRechargeStateConflict
		}
		row.Status, row.PrepayRefCipher, row.ExpiresAt, row.Version = model.WalletRechargeStatusPrepayCreated, prepay.PrepayRefCipher, expiresAt, row.Version+1
		return nil
	})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return row, biz.ErrWalletNotFound
	}
	return row, err
}

func (r *WalletRepo) Recharge(ctx context.Context, userID, rechargeID int64) (model.WalletRecharge, error) {
	var row model.WalletRecharge
	err := r.data.db.WithContext(ctx).Where("id = ? AND user_id = ?", rechargeID, userID).First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return row, biz.ErrWalletNotFound
	}
	return row, err
}

func (r *WalletRepo) Recharges(ctx context.Context, userID, before int64, limit int, status string) ([]model.WalletRecharge, error) {
	query := r.data.db.WithContext(ctx).Where("user_id = ?", userID)
	if before > 0 {
		query = query.Where("id < ?", before)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	var rows []model.WalletRecharge
	err := query.Order("id DESC").Limit(limit + 1).Find(&rows).Error
	return rows, err
}

func (r *WalletRepo) ExpiredRecharges(ctx context.Context, limit int, now, retryBefore time.Time) ([]model.WalletRecharge, error) {
	var rows []model.WalletRecharge
	err := r.data.db.WithContext(ctx).
		Where("(expires_at <= ? AND status IN ?) OR (status = ? AND updated_at <= ?)", now, []model.WalletRechargeStatus{model.WalletRechargeStatusPending, model.WalletRechargeStatusPrepayCreated}, model.WalletRechargeStatusClosing, retryBefore).
		Order("id ASC").Limit(limit).Find(&rows).Error
	return rows, err
}

func (r *WalletRepo) MarkRechargeClosing(ctx context.Context, rechargeID int64, version int32) (model.WalletRecharge, error) {
	now := time.Now().UTC()
	result := r.data.db.WithContext(ctx).Model(&model.WalletRecharge{}).
		Where("id = ? AND version = ? AND status IN ?", rechargeID, version, []model.WalletRechargeStatus{model.WalletRechargeStatusPending, model.WalletRechargeStatusPrepayCreated, model.WalletRechargeStatusClosing}).
		Updates(map[string]any{"status": model.WalletRechargeStatusClosing, "version": version + 1, "updated_at": now})
	if result.Error != nil {
		return model.WalletRecharge{}, result.Error
	}
	if result.RowsAffected != 1 {
		return model.WalletRecharge{}, biz.ErrRechargeStateConflict
	}
	var row model.WalletRecharge
	if err := r.data.db.WithContext(ctx).First(&row, rechargeID).Error; err != nil {
		return row, err
	}
	if row.Status != model.WalletRechargeStatusClosing || row.Version != version+1 {
		return row, biz.ErrRechargeStateConflict
	}
	return row, nil
}

func (r *WalletRepo) MarkRechargeClosed(ctx context.Context, rechargeID int64, version int32) (model.WalletRecharge, error) {
	now := time.Now().UTC()
	result := r.data.db.WithContext(ctx).Model(&model.WalletRecharge{}).
		Where("id = ? AND version = ? AND status = ?", rechargeID, version, model.WalletRechargeStatusClosing).
		Updates(map[string]any{"status": model.WalletRechargeStatusClosed, "failure_code": "PAYMENT_EXPIRED", "version": version + 1, "updated_at": now})
	if result.Error != nil {
		return model.WalletRecharge{}, result.Error
	}
	if result.RowsAffected != 1 {
		return model.WalletRecharge{}, biz.ErrRechargeStateConflict
	}
	var row model.WalletRecharge
	if err := r.data.db.WithContext(ctx).First(&row, rechargeID).Error; err != nil {
		return row, err
	}
	return row, nil
}

func (r *WalletRepo) LedgerEntries(ctx context.Context, userID, before int64, limit int, entryType string) ([]model.WalletLedgerEntry, error) {
	query := r.data.db.WithContext(ctx).Where("user_id = ?", userID)
	if before > 0 {
		query = query.Where("id < ?", before)
	}
	if entryType != "" {
		query = query.Where("entry_type = ?", entryType)
	}
	var rows []model.WalletLedgerEntry
	err := query.Order("id DESC").Limit(limit + 1).Find(&rows).Error
	return rows, err
}

func (r *WalletRepo) CreditRecharge(ctx context.Context, verified biz.VerifiedRechargeEvent) (recharge model.WalletRecharge, wallet model.WalletAccount, replayed bool, err error) {
	var businessErr error
	err = r.data.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now().UTC()
		digest := strings.TrimSpace(verified.PayloadDigest)
		if digest == "" {
			sum := sha256.Sum256([]byte(verified.EventKey))
			digest = hex.EncodeToString(sum[:])
		}
		event := model.WalletRechargeEvent{Provider: verified.Provider, MerchantID: verified.MerchantID, EventKey: verified.EventKey, EventType: verified.EventType, PayloadDigest: digest, Status: model.WalletRechargeEventStatusReceived}
		created := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&event)
		if created.Error != nil {
			return created.Error
		}
		if created.RowsAffected == 0 {
			if err := tx.Where("provider = ? AND merchant_id = ? AND event_key = ?", verified.Provider, verified.MerchantID, verified.EventKey).First(&event).Error; err != nil {
				return err
			}
			if event.Status != model.WalletRechargeEventStatusProcessed || event.RechargeID == nil {
				return biz.ErrRechargeStateConflict
			}
			if err := tx.First(&recharge, *event.RechargeID).Error; err != nil {
				return err
			}
			if recharge.WalletID == nil {
				return biz.ErrRechargeStateConflict
			}
			if err := tx.First(&wallet, *recharge.WalletID).Error; err != nil {
				return err
			}
			replayed = true
			return nil
		}

		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("provider = ? AND merchant_id = ? AND out_trade_no = ?", verified.Provider, verified.MerchantID, verified.OutTradeNo).First(&recharge).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return biz.ErrWalletNotFound
			}
			return err
		}
		if verified.ProviderTransactionID == "" || verified.EventKey == "" || !recharge.PayAmount.Equal(verified.PayAmount.Round(2)) || recharge.Currency != verified.Currency {
			code := "CALLBACK_MISMATCH"
			if err := tx.Model(&event).Updates(map[string]any{"recharge_id": recharge.ID, "status": model.WalletRechargeEventStatusRejected, "error_code": code, "handled_at": now}).Error; err != nil {
				return err
			}
			businessErr = biz.ErrInvalidArgument
			return nil
		}
		if recharge.Status == model.WalletRechargeStatusSucceeded {
			if recharge.ProviderTransactionID == nil || *recharge.ProviderTransactionID != verified.ProviderTransactionID || recharge.WalletID == nil {
				return biz.ErrRechargeStateConflict
			}
			if err := tx.First(&wallet, *recharge.WalletID).Error; err != nil {
				return err
			}
			if err := tx.Model(&event).Updates(map[string]any{"recharge_id": recharge.ID, "status": model.WalletRechargeEventStatusProcessed, "handled_at": now}).Error; err != nil {
				return err
			}
			replayed = true
			return nil
		}
		if recharge.Status != model.WalletRechargeStatusPending && recharge.Status != model.WalletRechargeStatusPrepayCreated && recharge.Status != model.WalletRechargeStatusClosing {
			return biz.ErrRechargeStateConflict
		}

		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("user_id = ? AND asset_code = ?", recharge.UserID, recharge.AssetCode).First(&wallet).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			wallet = model.WalletAccount{UserID: recharge.UserID, AssetCode: recharge.AssetCode, AvailableBalance: decimal.Zero, FrozenBalance: decimal.Zero, Status: model.WalletAccountStatusActive, Version: 1}
			if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&wallet).Error; err != nil {
				return err
			}
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("user_id = ? AND asset_code = ?", recharge.UserID, recharge.AssetCode).First(&wallet).Error; err != nil {
				return err
			}
		} else if err != nil {
			return err
		}
		if wallet.Status != model.WalletAccountStatusActive {
			return biz.ErrWalletUnavailable
		}
		newBalance := wallet.AvailableBalance.Add(recharge.CreditedAmount).Round(2)
		limit := decimal.RequireFromString("999999999999999999.99")
		if newBalance.IsNegative() || newBalance.GreaterThan(limit) {
			return biz.ErrInvalidArgument
		}
		result := tx.Model(&model.WalletAccount{}).Where("id = ? AND version = ?", wallet.ID, wallet.Version).Updates(map[string]any{"available_balance": newBalance, "version": wallet.Version + 1, "updated_at": now})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected != 1 {
			return biz.ErrRechargeStateConflict
		}
		wallet.AvailableBalance, wallet.Version = newBalance, wallet.Version+1
		ledger := model.WalletLedgerEntry{WalletID: wallet.ID, UserID: recharge.UserID, AssetCode: recharge.AssetCode, EntryType: "RECHARGE_CREDIT", BizType: "RECHARGE", BizID: recharge.ID, AvailableDelta: recharge.CreditedAmount, FrozenDelta: decimal.Zero, AvailableAfter: wallet.AvailableBalance, FrozenAfter: wallet.FrozenBalance, OccurredAt: now}
		if err := tx.Create(&ledger).Error; err != nil {
			return err
		}
		paidAt := verified.PaidAt
		if paidAt.IsZero() {
			paidAt = now
		}
		result = tx.Model(&model.WalletRecharge{}).Where("id = ? AND version = ?", recharge.ID, recharge.Version).Updates(map[string]any{"wallet_id": wallet.ID, "provider_transaction_id": verified.ProviderTransactionID, "status": model.WalletRechargeStatusSucceeded, "paid_at": paidAt, "credited_at": now, "version": recharge.Version + 1, "updated_at": now})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected != 1 {
			return biz.ErrRechargeStateConflict
		}
		recharge.WalletID, recharge.ProviderTransactionID, recharge.Status, recharge.PaidAt, recharge.CreditedAt, recharge.Version = &wallet.ID, &verified.ProviderTransactionID, model.WalletRechargeStatusSucceeded, &paidAt, &now, recharge.Version+1
		if err := tx.Model(&event).Updates(map[string]any{"recharge_id": recharge.ID, "status": model.WalletRechargeEventStatusProcessed, "handled_at": now}).Error; err != nil {
			return err
		}
		return nil
	})
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return recharge, wallet, false, fmt.Errorf("wallet credit uniqueness conflict: %w", err)
	}
	if err == nil && businessErr != nil {
		return recharge, wallet, false, businessErr
	}
	return recharge, wallet, replayed, err
}
