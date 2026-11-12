package biz

import (
	"context"
	"errors"
	"testing"
	"time"

	"dd_web/back/internal/model"
	"github.com/shopspring/decimal"
)

type walletRepoFake struct {
	packageRow model.RechargePackage
	created    model.WalletRecharge
	expired    []model.WalletRecharge
	closed     model.WalletRecharge
}

func (f *walletRepoFake) ActiveRechargePackages(context.Context, int64, int, time.Time) ([]model.RechargePackage, error) {
	return []model.RechargePackage{f.packageRow}, nil
}
func (f *walletRepoFake) RechargePackage(context.Context, int64, time.Time) (model.RechargePackage, error) {
	return f.packageRow, nil
}
func (f *walletRepoFake) Wallet(context.Context, int64, string) (*model.WalletAccount, error) {
	return nil, nil
}
func (f *walletRepoFake) CreateRecharge(_ context.Context, _ model.RechargePackage, row model.WalletRecharge, _, _ string) (model.WalletRecharge, bool, error) {
	row.ID = 9
	f.created = row
	return row, false, nil
}
func (f *walletRepoFake) MarkRechargePrepayCreated(_ context.Context, _ int64, result RechargePrepayResult) (model.WalletRecharge, error) {
	row := f.created
	row.Status = model.WalletRechargeStatusPrepayCreated
	row.ExpiresAt = result.ExpiresAt
	return row, nil
}
func (f *walletRepoFake) Recharge(context.Context, int64, int64) (model.WalletRecharge, error) {
	return f.created, nil
}
func (f *walletRepoFake) Recharges(context.Context, int64, int64, int, string) ([]model.WalletRecharge, error) {
	return f.expired, nil
}
func (f *walletRepoFake) ExpiredRecharges(context.Context, int, time.Time, time.Time) ([]model.WalletRecharge, error) {
	return f.expired, nil
}
func (f *walletRepoFake) MarkRechargeClosing(_ context.Context, id int64, version int32) (model.WalletRecharge, error) {
	for _, row := range f.expired {
		if row.ID == id && row.Version == version {
			row.Status, row.Version = model.WalletRechargeStatusClosing, row.Version+1
			return row, nil
		}
	}
	return model.WalletRecharge{}, ErrRechargeStateConflict
}
func (f *walletRepoFake) MarkRechargeClosed(_ context.Context, id int64, version int32) (model.WalletRecharge, error) {
	f.closed = model.WalletRecharge{Base: model.Base{ID: id}, Status: model.WalletRechargeStatusClosed, Version: version + 1}
	return f.closed, nil
}
func (f *walletRepoFake) LedgerEntries(context.Context, int64, int64, int, string) ([]model.WalletLedgerEntry, error) {
	return nil, nil
}
func (f *walletRepoFake) CreditRecharge(context.Context, VerifiedRechargeEvent) (model.WalletRecharge, model.WalletAccount, bool, error) {
	return model.WalletRecharge{}, model.WalletAccount{}, false, nil
}
func (f *walletRepoFake) AdminRechargePackages(context.Context, int64, int, string) ([]model.RechargePackage, error) {
	return []model.RechargePackage{f.packageRow}, nil
}
func (f *walletRepoFake) CreateRechargePackage(_ context.Context, _ int64, row model.RechargePackage, _, _ string) (model.RechargePackage, bool, error) {
	row.ID = 10
	f.packageRow = row
	return row, false, nil
}
func (f *walletRepoFake) UpdateRechargePackage(_ context.Context, _ int64, id int64, row model.RechargePackage, _, _ string) (model.RechargePackage, bool, error) {
	row.ID, row.Version = id, row.Version+1
	f.packageRow = row
	return row, false, nil
}

type rechargeProviderFake struct{}

func (rechargeProviderFake) Code() string { return "FAKE" }
func (rechargeProviderFake) MerchantAndApp(string) (string, string, error) {
	return "merchant", "app", nil
}
func (rechargeProviderFake) CreatePrepay(context.Context, RechargePrepayRequest) (RechargePrepayResult, error) {
	return RechargePrepayResult{PaymentParams: map[string]any{"token": "test"}, ExpiresAt: time.Date(2026, 11, 12, 12, 30, 0, 0, time.UTC)}, nil
}
func (rechargeProviderFake) QueryPayment(context.Context, RechargeQueryRequest) (RechargeQueryResult, error) {
	return RechargeQueryResult{State: RechargePaymentStateUnpaid}, nil
}
func (rechargeProviderFake) ClosePayment(context.Context, RechargeQueryRequest) error { return nil }
func (rechargeProviderFake) VerifyCallback(context.Context, map[string]any) (VerifiedRechargeEvent, error) {
	return VerifiedRechargeEvent{}, nil
}

func TestRechargeListShowsPaymentActionsAndExpiredConfirmation(t *testing.T) {
	now := time.Date(2026, 11, 12, 12, 0, 0, 0, time.UTC)
	repo := &walletRepoFake{expired: []model.WalletRecharge{
		{Base: model.Base{ID: 2, CreatedAt: now.Add(-time.Minute)}, PackageID: 3, PayAmount: decimal.NewFromInt(30), PackageMultiplier: decimal.NewFromInt(1), CreditedAmount: decimal.NewFromInt(30), Currency: "CNY", AssetCode: "DD_COIN", Provider: "FAKE", Status: model.WalletRechargeStatusPrepayCreated, ExpiresAt: now.Add(time.Minute)},
		{Base: model.Base{ID: 1, CreatedAt: now.Add(-time.Hour)}, PackageID: 3, PayAmount: decimal.NewFromInt(30), PackageMultiplier: decimal.NewFromInt(1), CreditedAmount: decimal.NewFromInt(30), Currency: "CNY", AssetCode: "DD_COIN", Provider: "FAKE", Status: model.WalletRechargeStatusPrepayCreated, ExpiresAt: now.Add(-time.Minute)},
	}}
	usecase := NewWalletUsecase(repo, rechargeProviderFake{})
	usecase.now = func() time.Time { return now }
	list, err := usecase.Recharges(context.Background(), 7, Page{Limit: 20}, "")
	if err != nil {
		t.Fatal(err)
	}
	if list.Items[0]["display_status"] != "PENDING_PAYMENT" || list.Items[0]["can_pay"] != true {
		t.Fatalf("payable item = %#v", list.Items[0])
	}
	if list.Items[1]["display_status"] != "CONFIRMING" || list.Items[1]["can_pay"] != false {
		t.Fatalf("expired item = %#v", list.Items[1])
	}
}

func TestCloseExpiredRechargeQueriesProviderBeforeClosing(t *testing.T) {
	now := time.Date(2026, 11, 12, 12, 0, 0, 0, time.UTC)
	repo := &walletRepoFake{expired: []model.WalletRecharge{{Base: model.Base{ID: 9}, RechargeNo: "WR1", OutTradeNo: "WRT1", Provider: "FAKE", MerchantID: "merchant", AppID: "app", Status: model.WalletRechargeStatusPrepayCreated, ExpiresAt: now.Add(-time.Minute), Version: 2}}}
	usecase := NewWalletUsecase(repo, rechargeProviderFake{})
	usecase.now = func() time.Time { return now }
	result, err := usecase.CloseExpiredRecharges(context.Background(), 100, time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	if result.Closed != 1 || repo.closed.ID != 9 || repo.closed.Status != model.WalletRechargeStatusClosed {
		t.Fatalf("result=%#v closed=%#v", result, repo.closed)
	}
}

func TestWalletCreateRechargeUsesDecimalPackageMultiplier(t *testing.T) {
	repo := &walletRepoFake{packageRow: model.RechargePackage{
		Base: model.Base{ID: 3}, Name: "赠送5%", PayAmount: decimal.RequireFromString("30.00"), Currency: "CNY",
		PackageMultiplier: decimal.RequireFromString("1.0500"), CreditedAmount: decimal.RequireFromString("31.50"),
		AssetCode: "DD_COIN", Status: model.RechargePackageStatusEnabled, Version: 2,
	}}
	usecase := NewWalletUsecase(repo, rechargeProviderFake{})
	usecase.now = func() time.Time { return time.Date(2026, 11, 12, 12, 0, 0, 0, time.UTC) }

	result, err := usecase.CreateRecharge(context.Background(), RechargeCreateInput{UserID: 7, PackageID: 3, Provider: "fake", ClientType: "APP", IdempotencyKey: "key-1"})
	if err != nil {
		t.Fatal(err)
	}
	if got := repo.created.CreditedAmount.StringFixed(2); got != "31.50" {
		t.Fatalf("credited amount = %s", got)
	}
	if got := result["package_multiplier"]; got != "1.0500" {
		t.Fatalf("multiplier = %v", got)
	}
	if got := result["credited_amount"]; got != "31.50" {
		t.Fatalf("response credited amount = %v", got)
	}
	if got := result["payment_params"].(map[string]any)["token"]; got != "test" {
		t.Fatalf("payment token = %v", got)
	}
}

func TestWalletWithoutAccountReturnsTwoDecimalZero(t *testing.T) {
	usecase := NewWalletUsecaseWithOptions(&walletRepoFake{}, WalletOptions{AssetCode: "game_credit", RechargeExpire: 45 * time.Minute})
	result, err := usecase.Wallet(context.Background(), 7)
	if err != nil {
		t.Fatal(err)
	}
	if result["available_balance"] != "0.00" || result["frozen_balance"] != "0.00" {
		t.Fatalf("wallet = %#v", result)
	}
	if result["asset_code"] != "GAME_CREDIT" {
		t.Fatalf("asset code = %v", result["asset_code"])
	}
}

func TestAdminCreateRechargePackageCalculatesCreditedAmount(t *testing.T) {
	repo := &walletRepoFake{}
	usecase := NewWalletUsecase(repo)
	result, err := usecase.AdminCreateRechargePackage(context.Background(), 3, RechargePackageInput{Name: "赠送7.5%", PayAmount: "100.00", PackageMultiplier: "1.0750", Currency: "CNY", AssetCode: "DD_COIN", Status: "ENABLED"}, "admin-key")
	if err != nil {
		t.Fatal(err)
	}
	if repo.packageRow.CreditedAmount.StringFixed(2) != "107.50" {
		t.Fatalf("credited amount = %s", repo.packageRow.CreditedAmount)
	}
	if result["package_multiplier"] != "1.0750" || result["credited_amount"] != "107.50" {
		t.Fatalf("result = %#v", result)
	}
}

func TestAdminRechargePackageRejectsExcessPrecision(t *testing.T) {
	usecase := NewWalletUsecase(&walletRepoFake{})
	_, err := usecase.AdminCreateRechargePackage(context.Background(), 3, RechargePackageInput{Name: "bad", PayAmount: "10.001", PackageMultiplier: "1.0000", Currency: "CNY", AssetCode: "DD_COIN", Status: "ENABLED"}, "admin-key")
	if !errors.Is(err, ErrInvalidArgument) {
		t.Fatalf("error = %v", err)
	}
}
