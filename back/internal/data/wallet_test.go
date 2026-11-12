package data

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"dd_web/back/internal/biz"
	"dd_web/back/internal/model"
	"github.com/shopspring/decimal"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func newWalletTestRepo(t *testing.T) (*WalletRepo, *gorm.DB, model.RechargePackage) {
	t.Helper()
	dsn := fmt.Sprintf("file:wallet-%d?mode=memory&cache=shared", time.Now().UnixNano())
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{TranslateError: true})
	if err != nil {
		t.Fatal(err)
	}
	schema := []string{
		`CREATE TABLE recharge_packages (
			id INTEGER PRIMARY KEY AUTOINCREMENT, created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			name VARCHAR(64) NOT NULL, pay_amount DECIMAL(20,2) NOT NULL, currency CHAR(3) NOT NULL,
			package_multiplier DECIMAL(10,4) NOT NULL, credited_amount DECIMAL(20,2) NOT NULL,
			asset_code VARCHAR(32) NOT NULL, status VARCHAR(16) NOT NULL, effective_at DATETIME, expires_at DATETIME, version INTEGER NOT NULL
		)`,
		`CREATE TABLE wallet_accounts (
			user_id BIGINT PRIMARY KEY, asset_code VARCHAR(32) NOT NULL, available_balance DECIMAL(20,2) NOT NULL,
			frozen_balance DECIMAL(20,2) NOT NULL, status VARCHAR(16) NOT NULL, version INTEGER NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE wallet_recharges (
			id INTEGER PRIMARY KEY AUTOINCREMENT, created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			recharge_no VARCHAR(40) NOT NULL UNIQUE, user_id BIGINT NOT NULL, package_id BIGINT NOT NULL,
			package_version INTEGER NOT NULL, pay_amount DECIMAL(20,2) NOT NULL, currency CHAR(3) NOT NULL,
			package_multiplier DECIMAL(10,4) NOT NULL, credited_amount DECIMAL(20,2) NOT NULL, asset_code VARCHAR(32) NOT NULL,
			provider VARCHAR(32) NOT NULL, merchant_id VARCHAR(64) NOT NULL, app_id VARCHAR(64) NOT NULL, client_type VARCHAR(32) NOT NULL,
			out_trade_no VARCHAR(64) NOT NULL, provider_transaction_id VARCHAR(128), status VARCHAR(24) NOT NULL,
			prepay_ref_cipher BLOB, failure_code VARCHAR(64), expires_at DATETIME NOT NULL, paid_at DATETIME, credited_at DATETIME, version INTEGER NOT NULL,
			UNIQUE(provider, merchant_id, out_trade_no), UNIQUE(provider, merchant_id, provider_transaction_id)
		)`,
		`CREATE TABLE wallet_ledger_entries (
			id INTEGER PRIMARY KEY AUTOINCREMENT, created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			user_id BIGINT NOT NULL, asset_code VARCHAR(32) NOT NULL, entry_type VARCHAR(32) NOT NULL,
			biz_type VARCHAR(32) NOT NULL, biz_id BIGINT NOT NULL, available_delta DECIMAL(20,2) NOT NULL,
			frozen_delta DECIMAL(20,2) NOT NULL, available_after DECIMAL(20,2) NOT NULL, frozen_after DECIMAL(20,2) NOT NULL,
			occurred_at DATETIME NOT NULL, UNIQUE(user_id, asset_code, entry_type, biz_type, biz_id)
		)`,
		`CREATE TABLE wallet_recharge_events (
			id INTEGER PRIMARY KEY AUTOINCREMENT, created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			recharge_id BIGINT, provider VARCHAR(32) NOT NULL, merchant_id VARCHAR(64) NOT NULL, event_key VARCHAR(160) NOT NULL,
			event_type VARCHAR(48) NOT NULL, payload_digest VARCHAR(64) NOT NULL, payload_asset_id BIGINT,
			status VARCHAR(16) NOT NULL, error_code VARCHAR(64), handled_at DATETIME, UNIQUE(provider, merchant_id, event_key)
		)`,
		`CREATE TABLE idempotency_keys (
		id INTEGER PRIMARY KEY AUTOINCREMENT, created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		actor_type VARCHAR(16) NOT NULL, actor_id BIGINT NOT NULL, operation VARCHAR(64) NOT NULL,
		biz_key VARCHAR(128) NOT NULL, request_hash VARCHAR(64) NOT NULL, status VARCHAR(16) NOT NULL,
		resource_type VARCHAR(32), resource_id BIGINT, response_summary TEXT NOT NULL, expires_at DATETIME NOT NULL,
		UNIQUE(actor_type, actor_id, operation, biz_key)
	)`,
	}
	for _, statement := range schema {
		if err := db.Exec(statement).Error; err != nil {
			t.Fatal(err)
		}
	}
	pkg := model.RechargePackage{Name: "赠送5%", PayAmount: decimal.RequireFromString("30.00"), Currency: "CNY", PackageMultiplier: decimal.RequireFromString("1.0500"), CreditedAmount: decimal.RequireFromString("31.50"), AssetCode: "DD_COIN", Status: model.RechargePackageStatusEnabled, Version: 1}
	if err := db.Create(&pkg).Error; err != nil {
		t.Fatal(err)
	}
	return NewWalletRepo(NewData(db)), db, pkg
}

func rechargeCandidate(pkg model.RechargePackage) model.WalletRecharge {
	return model.WalletRecharge{RechargeNo: "WR-test", UserID: 7, PackageID: pkg.ID, PackageVersion: pkg.Version, PayAmount: pkg.PayAmount, Currency: pkg.Currency, PackageMultiplier: pkg.PackageMultiplier, CreditedAmount: pkg.CreditedAmount, AssetCode: pkg.AssetCode, Provider: "LOCAL", MerchantID: "merchant", AppID: "app", ClientType: "APP", OutTradeNo: "WRT-test", Status: model.WalletRechargeStatusPending, ExpiresAt: time.Now().UTC().Add(30 * time.Minute), Version: 1}
}

func TestWalletRechargeIdempotency(t *testing.T) {
	repo, _, pkg := newWalletTestRepo(t)
	ctx := context.Background()
	first, replayed, err := repo.CreateRecharge(ctx, pkg, rechargeCandidate(pkg), "same-key", "hash-a")
	if err != nil {
		t.Fatal(err)
	}
	if replayed {
		t.Fatal("first request must not be replayed")
	}
	second, replayed, err := repo.CreateRecharge(ctx, pkg, rechargeCandidate(pkg), "same-key", "hash-a")
	if err != nil {
		t.Fatal(err)
	}
	if !replayed || second.ID != first.ID {
		t.Fatalf("replay = %v, ids = %d/%d", replayed, first.ID, second.ID)
	}
	_, _, err = repo.CreateRecharge(ctx, pkg, rechargeCandidate(pkg), "same-key", "hash-b")
	if !errors.Is(err, biz.ErrIdempotencyConflict) {
		t.Fatalf("conflict error = %v", err)
	}
}

func TestWalletRechargeCallbackCreditsExactlyOnce(t *testing.T) {
	repo, db, pkg := newWalletTestRepo(t)
	ctx := context.Background()
	recharge, _, err := repo.CreateRecharge(ctx, pkg, rechargeCandidate(pkg), "credit-key", "hash")
	if err != nil {
		t.Fatal(err)
	}
	recharge, err = repo.MarkRechargePrepayCreated(ctx, recharge.ID, biz.RechargePrepayResult{ExpiresAt: recharge.ExpiresAt})
	if err != nil {
		t.Fatal(err)
	}
	event := biz.VerifiedRechargeEvent{Provider: "LOCAL", MerchantID: "merchant", EventKey: "event-1", EventType: "PAY_NOTIFY", PayloadDigest: "digest", OutTradeNo: recharge.OutTradeNo, ProviderTransactionID: "tx-1", PayAmount: decimal.RequireFromString("30.00"), Currency: "CNY", PaidAt: time.Now().UTC()}
	credited, wallet, replayed, err := repo.CreditRecharge(ctx, event)
	if err != nil {
		t.Fatal(err)
	}
	if replayed || credited.Status != model.WalletRechargeStatusSucceeded || wallet.AvailableBalance.StringFixed(2) != "31.50" {
		t.Fatalf("first credit: replayed=%v recharge=%s balance=%s", replayed, credited.Status, wallet.AvailableBalance)
	}
	_, wallet, replayed, err = repo.CreditRecharge(ctx, event)
	if err != nil {
		t.Fatal(err)
	}
	if !replayed || wallet.AvailableBalance.StringFixed(2) != "31.50" {
		t.Fatalf("replay: replayed=%v balance=%s", replayed, wallet.AvailableBalance)
	}
	event.EventKey = "event-2"
	_, wallet, replayed, err = repo.CreditRecharge(ctx, event)
	if err != nil {
		t.Fatal(err)
	}
	if !replayed || wallet.AvailableBalance.StringFixed(2) != "31.50" {
		t.Fatalf("second event: replayed=%v balance=%s", replayed, wallet.AvailableBalance)
	}
	var ledgerCount int64
	if err := db.Model(&model.WalletLedgerEntry{}).Count(&ledgerCount).Error; err != nil {
		t.Fatal(err)
	}
	if ledgerCount != 1 {
		t.Fatalf("ledger count = %d", ledgerCount)
	}
}

func TestWalletRechargeClosingStillAcceptsVerifiedPayment(t *testing.T) {
	repo, _, pkg := newWalletTestRepo(t)
	ctx := context.Background()
	candidate := rechargeCandidate(pkg)
	candidate.ExpiresAt = time.Now().UTC().Add(-time.Minute)
	recharge, _, err := repo.CreateRecharge(ctx, pkg, candidate, "closing-credit-key", "hash")
	if err != nil {
		t.Fatal(err)
	}
	recharge, err = repo.MarkRechargePrepayCreated(ctx, recharge.ID, biz.RechargePrepayResult{ExpiresAt: candidate.ExpiresAt})
	if err != nil {
		t.Fatal(err)
	}
	recharge, err = repo.MarkRechargeClosing(ctx, recharge.ID, recharge.Version)
	if err != nil {
		t.Fatal(err)
	}
	event := biz.VerifiedRechargeEvent{Provider: "LOCAL", MerchantID: "merchant", EventKey: "closing-event", EventType: "PAY_NOTIFY", PayloadDigest: "digest", OutTradeNo: recharge.OutTradeNo, ProviderTransactionID: "tx-closing", PayAmount: decimal.RequireFromString("30.00"), Currency: "CNY", PaidAt: time.Now().UTC()}
	credited, wallet, _, err := repo.CreditRecharge(ctx, event)
	if err != nil {
		t.Fatal(err)
	}
	if credited.Status != model.WalletRechargeStatusSucceeded || wallet.AvailableBalance.StringFixed(2) != "31.50" {
		t.Fatalf("recharge=%s balance=%s", credited.Status, wallet.AvailableBalance)
	}
	if _, err := repo.MarkRechargeClosed(ctx, recharge.ID, recharge.Version); !errors.Is(err, biz.ErrRechargeStateConflict) {
		t.Fatalf("stale close error = %v", err)
	}
}

func TestWalletExpiredRechargeStateTransitionsUseVersionCAS(t *testing.T) {
	repo, _, pkg := newWalletTestRepo(t)
	ctx := context.Background()
	candidate := rechargeCandidate(pkg)
	candidate.ExpiresAt = time.Now().UTC().Add(-time.Minute)
	recharge, _, err := repo.CreateRecharge(ctx, pkg, candidate, "expire-key", "hash")
	if err != nil {
		t.Fatal(err)
	}
	rows, err := repo.ExpiredRecharges(ctx, 10, time.Now().UTC(), time.Now().UTC().Add(-time.Minute))
	if err != nil || len(rows) != 1 {
		t.Fatalf("expired rows=%d err=%v", len(rows), err)
	}
	closing, err := repo.MarkRechargeClosing(ctx, recharge.ID, recharge.Version)
	if err != nil || closing.Status != model.WalletRechargeStatusClosing {
		t.Fatalf("closing=%#v err=%v", closing, err)
	}
	if _, err := repo.MarkRechargeClosing(ctx, recharge.ID, recharge.Version); !errors.Is(err, biz.ErrRechargeStateConflict) {
		t.Fatalf("stale claim error = %v", err)
	}
	closed, err := repo.MarkRechargeClosed(ctx, closing.ID, closing.Version)
	if err != nil || closed.Status != model.WalletRechargeStatusClosed {
		t.Fatalf("closed=%#v err=%v", closed, err)
	}
	if closed.FailureCode == nil || *closed.FailureCode != "PAYMENT_EXPIRED" {
		t.Fatalf("failure code = %v", closed.FailureCode)
	}
}

func TestWalletRechargeMismatchIsRejectedWithoutBalanceChange(t *testing.T) {
	repo, db, pkg := newWalletTestRepo(t)
	ctx := context.Background()
	recharge, _, err := repo.CreateRecharge(ctx, pkg, rechargeCandidate(pkg), "mismatch-key", "hash")
	if err != nil {
		t.Fatal(err)
	}
	event := biz.VerifiedRechargeEvent{Provider: "LOCAL", MerchantID: "merchant", EventKey: "event-mismatch", EventType: "PAY_NOTIFY", PayloadDigest: "digest", OutTradeNo: recharge.OutTradeNo, ProviderTransactionID: "tx-mismatch", PayAmount: decimal.RequireFromString("29.99"), Currency: "CNY", PaidAt: time.Now().UTC()}
	_, _, _, err = repo.CreditRecharge(ctx, event)
	if !errors.Is(err, biz.ErrInvalidArgument) {
		t.Fatalf("callback error = %v", err)
	}
	var stored model.WalletRechargeEvent
	if err := db.Where("event_key = ?", event.EventKey).First(&stored).Error; err != nil {
		t.Fatal(err)
	}
	if stored.Status != model.WalletRechargeEventStatusRejected {
		t.Fatalf("event status = %s", stored.Status)
	}
	var wallets, ledgers int64
	if err := db.Model(&model.WalletAccount{}).Count(&wallets).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Model(&model.WalletLedgerEntry{}).Count(&ledgers).Error; err != nil {
		t.Fatal(err)
	}
	if wallets != 0 || ledgers != 0 {
		t.Fatalf("wallets=%d ledgers=%d", wallets, ledgers)
	}
}

func TestAdminRechargePackageCreateAndVersionedUpdate(t *testing.T) {
	repo, _, _ := newWalletTestRepo(t)
	ctx := context.Background()
	candidate := model.RechargePackage{Name: "100元套餐", PayAmount: decimal.RequireFromString("100.00"), Currency: "CNY", PackageMultiplier: decimal.RequireFromString("1.0500"), CreditedAmount: decimal.RequireFromString("105.00"), AssetCode: "DD_COIN", Status: model.RechargePackageStatusEnabled, Version: 1}
	created, replayed, err := repo.CreateRechargePackage(ctx, 11, candidate, "create-package", "hash-create")
	if err != nil {
		t.Fatal(err)
	}
	if replayed || created.ID <= 0 || created.Version != 1 {
		t.Fatalf("created=%#v replayed=%v", created, replayed)
	}
	replayedCreate, replayed, err := repo.CreateRechargePackage(ctx, 11, candidate, "create-package", "hash-create")
	if err != nil {
		t.Fatal(err)
	}
	if !replayed || replayedCreate.ID != created.ID {
		t.Fatalf("create replay=%v ids=%d/%d", replayed, created.ID, replayedCreate.ID)
	}

	update := created
	update.Name = "100元加赠套餐"
	update.PackageMultiplier = decimal.RequireFromString("1.0750")
	update.CreditedAmount = decimal.RequireFromString("107.50")
	updated, replayed, err := repo.UpdateRechargePackage(ctx, 11, created.ID, update, "update-package", "hash-update")
	if err != nil {
		t.Fatal(err)
	}
	if replayed || updated.Version != 2 || updated.CreditedAmount.StringFixed(2) != "107.50" {
		t.Fatalf("updated=%#v replayed=%v", updated, replayed)
	}
	_, _, err = repo.UpdateRechargePackage(ctx, 11, created.ID, update, "stale-update", "hash-stale")
	if !errors.Is(err, biz.ErrRechargePackageConflict) {
		t.Fatalf("stale update error = %v", err)
	}
}
