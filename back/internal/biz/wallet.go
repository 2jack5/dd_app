package biz

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"dd_web/back/internal/model"
	"github.com/shopspring/decimal"
)

var (
	ErrWalletNotFound          = errors.New("wallet resource not found")
	ErrWalletUnavailable       = errors.New("wallet unavailable")
	ErrRechargeStateConflict   = errors.New("recharge state conflict")
	ErrRechargePackageConflict = errors.New("recharge package version conflict")
	ErrIdempotencyConflict     = errors.New("idempotency key conflict")
	ErrRechargeChannel         = errors.New("recharge channel unavailable")
)

type RechargeCreateInput struct {
	UserID         int64
	PackageID      int64
	Provider       string
	ClientType     string
	IdempotencyKey string
}

type RechargePackageInput struct {
	Name              string
	PayAmount         string
	PackageMultiplier string
	Currency          string
	AssetCode         string
	Status            string
	EffectiveAt       string
	ExpiresAt         string
	Version           int32
}

type RechargePrepayRequest struct {
	RechargeNo string
	OutTradeNo string
	PayAmount  decimal.Decimal
	Currency   string
	ClientType string
	ExpiresAt  time.Time
}

type RechargePrepayResult struct {
	PaymentParams   map[string]any
	PrepayRefCipher []byte
	ExpiresAt       time.Time
}

type RechargePaymentState string

const (
	RechargePaymentStateUnknown RechargePaymentState = "UNKNOWN"
	RechargePaymentStateUnpaid  RechargePaymentState = "UNPAID"
	RechargePaymentStatePaid    RechargePaymentState = "PAID"
	RechargePaymentStateClosed  RechargePaymentState = "CLOSED"
)

type RechargeQueryRequest struct {
	RechargeNo string
	OutTradeNo string
	MerchantID string
	AppID      string
}

type RechargeQueryResult struct {
	State                 RechargePaymentState
	EventKey              string
	PayloadDigest         string
	ProviderTransactionID string
	PayAmount             decimal.Decimal
	Currency              string
	PaidAt                time.Time
}

type RechargeCloseResult struct {
	Scanned   int
	Closed    int
	Credited  int
	Retrying  int
	Conflicts int
}

type VerifiedRechargeEvent struct {
	Provider              string
	MerchantID            string
	EventKey              string
	EventType             string
	PayloadDigest         string
	OutTradeNo            string
	ProviderTransactionID string
	PayAmount             decimal.Decimal
	Currency              string
	PaidAt                time.Time
}

type RechargeProvider interface {
	Code() string
	MerchantAndApp(clientType string) (merchantID, appID string, err error)
	CreatePrepay(context.Context, RechargePrepayRequest) (RechargePrepayResult, error)
	QueryPayment(context.Context, RechargeQueryRequest) (RechargeQueryResult, error)
	ClosePayment(context.Context, RechargeQueryRequest) error
	VerifyCallback(context.Context, map[string]any) (VerifiedRechargeEvent, error)
}

type WalletRepository interface {
	ActiveRechargePackages(context.Context, int64, int, time.Time) ([]model.RechargePackage, error)
	RechargePackage(context.Context, int64, time.Time) (model.RechargePackage, error)
	Wallet(context.Context, int64, string) (*model.WalletAccount, error)
	CreateRecharge(context.Context, model.RechargePackage, model.WalletRecharge, string, string) (model.WalletRecharge, bool, error)
	MarkRechargePrepayCreated(context.Context, int64, RechargePrepayResult) (model.WalletRecharge, error)
	Recharge(context.Context, int64, int64) (model.WalletRecharge, error)
	Recharges(context.Context, int64, int64, int, string) ([]model.WalletRecharge, error)
	ExpiredRecharges(context.Context, int, time.Time, time.Time) ([]model.WalletRecharge, error)
	MarkRechargeClosing(context.Context, int64, int32) (model.WalletRecharge, error)
	MarkRechargeClosed(context.Context, int64, int32) (model.WalletRecharge, error)
	LedgerEntries(context.Context, int64, int64, int, string) ([]model.WalletLedgerEntry, error)
	CreditRecharge(context.Context, VerifiedRechargeEvent) (model.WalletRecharge, model.WalletAccount, bool, error)
	AdminRechargePackages(context.Context, int64, int, string) ([]model.RechargePackage, error)
	CreateRechargePackage(context.Context, int64, model.RechargePackage, string, string) (model.RechargePackage, bool, error)
	UpdateRechargePackage(context.Context, int64, int64, model.RechargePackage, string, string) (model.RechargePackage, bool, error)
}

type WalletUsecase struct {
	repo           WalletRepository
	providers      map[string]RechargeProvider
	assetCode      string
	rechargeExpire time.Duration
	now            func() time.Time
}

func NewWalletUsecase(repo WalletRepository, providers ...RechargeProvider) *WalletUsecase {
	return NewWalletUsecaseWithOptions(repo, WalletOptions{}, providers...)
}

type WalletOptions struct {
	AssetCode      string
	RechargeExpire time.Duration
}

func NewWalletUsecaseWithOptions(repo WalletRepository, options WalletOptions, providers ...RechargeProvider) *WalletUsecase {
	registry := make(map[string]RechargeProvider, len(providers))
	for _, provider := range providers {
		if provider != nil && strings.TrimSpace(provider.Code()) != "" {
			registry[strings.ToUpper(strings.TrimSpace(provider.Code()))] = provider
		}
	}
	assetCode := strings.ToUpper(strings.TrimSpace(options.AssetCode))
	if assetCode == "" {
		assetCode = "DD_COIN"
	}
	rechargeExpire := options.RechargeExpire
	if rechargeExpire <= 0 {
		rechargeExpire = 30 * time.Minute
	}
	return &WalletUsecase{repo: repo, providers: registry, assetCode: assetCode, rechargeExpire: rechargeExpire, now: func() time.Time { return time.Now().UTC() }}
}

func (u *WalletUsecase) RechargePackages(ctx context.Context, page Page) (List, error) {
	rows, err := u.repo.ActiveRechargePackages(ctx, page.After, page.Limit, u.now())
	if err != nil {
		return List{}, err
	}
	next := trim(&rows, page.Limit, func(v model.RechargePackage) int64 { return v.ID })
	items := make([]map[string]any, 0, len(rows))
	for _, row := range rows {
		items = append(items, rechargePackageMap(row))
	}
	return List{Items: items, NextCursor: next}, nil
}

func (u *WalletUsecase) AdminRechargePackages(ctx context.Context, page Page, status string) (List, error) {
	status = strings.ToUpper(strings.TrimSpace(status))
	if status != "" && status != string(model.RechargePackageStatusEnabled) && status != string(model.RechargePackageStatusDisabled) {
		return List{}, ErrInvalidArgument
	}
	rows, err := u.repo.AdminRechargePackages(ctx, page.After, page.Limit, status)
	if err != nil {
		return List{}, err
	}
	next := trim(&rows, page.Limit, func(v model.RechargePackage) int64 { return v.ID })
	items := make([]map[string]any, 0, len(rows))
	for _, row := range rows {
		items = append(items, rechargePackageMap(row))
	}
	return List{Items: items, NextCursor: next}, nil
}

func (u *WalletUsecase) AdminCreateRechargePackage(ctx context.Context, adminID int64, input RechargePackageInput, idempotencyKey string) (map[string]any, error) {
	if adminID <= 0 || strings.TrimSpace(idempotencyKey) == "" {
		return nil, ErrInvalidArgument
	}
	row, err := u.validateRechargePackage(input, false)
	if err != nil {
		return nil, err
	}
	hash := rechargePackageRequestHash("CREATE", 0, input)
	row, _, err = u.repo.CreateRechargePackage(ctx, adminID, row, idempotencyKey, hash)
	if err != nil {
		return nil, err
	}
	return rechargePackageMap(row), nil
}

func (u *WalletUsecase) AdminUpdateRechargePackage(ctx context.Context, adminID, packageID int64, input RechargePackageInput, idempotencyKey string) (map[string]any, error) {
	if adminID <= 0 || packageID <= 0 || input.Version <= 0 || strings.TrimSpace(idempotencyKey) == "" {
		return nil, ErrInvalidArgument
	}
	row, err := u.validateRechargePackage(input, true)
	if err != nil {
		return nil, err
	}
	hash := rechargePackageRequestHash("UPDATE", packageID, input)
	row, _, err = u.repo.UpdateRechargePackage(ctx, adminID, packageID, row, idempotencyKey, hash)
	if err != nil {
		return nil, err
	}
	return rechargePackageMap(row), nil
}

func (u *WalletUsecase) validateRechargePackage(input RechargePackageInput, updating bool) (model.RechargePackage, error) {
	name := strings.TrimSpace(input.Name)
	if name == "" || utf8.RuneCountInString(name) > 64 {
		return model.RechargePackage{}, ErrInvalidArgument
	}
	payAmount, err := decimal.NewFromString(strings.TrimSpace(input.PayAmount))
	if err != nil || payAmount.Exponent() < -2 || payAmount.LessThanOrEqual(decimal.Zero) {
		return model.RechargePackage{}, ErrInvalidArgument
	}
	multiplier, err := decimal.NewFromString(strings.TrimSpace(input.PackageMultiplier))
	if err != nil || multiplier.Exponent() < -4 || multiplier.LessThan(decimal.NewFromInt(1)) || multiplier.GreaterThan(decimal.NewFromInt(10)) {
		return model.RechargePackage{}, ErrInvalidArgument
	}
	currency := strings.ToUpper(strings.TrimSpace(input.Currency))
	if currency != "CNY" {
		return model.RechargePackage{}, ErrInvalidArgument
	}
	assetCode := strings.ToUpper(strings.TrimSpace(input.AssetCode))
	if assetCode != u.assetCode {
		return model.RechargePackage{}, ErrInvalidArgument
	}
	status := model.RechargePackageStatus(strings.ToUpper(strings.TrimSpace(input.Status)))
	if status != model.RechargePackageStatusEnabled && status != model.RechargePackageStatusDisabled {
		return model.RechargePackage{}, ErrInvalidArgument
	}
	effectiveAt, err := optionalRFC3339(input.EffectiveAt)
	if err != nil {
		return model.RechargePackage{}, ErrInvalidArgument
	}
	expiresAt, err := optionalRFC3339(input.ExpiresAt)
	if err != nil || (effectiveAt != nil && expiresAt != nil && !expiresAt.After(*effectiveAt)) {
		return model.RechargePackage{}, ErrInvalidArgument
	}
	version := int32(1)
	if updating {
		version = input.Version
	}
	return model.RechargePackage{Name: name, PayAmount: payAmount.Round(2), Currency: currency, PackageMultiplier: multiplier, CreditedAmount: payAmount.Mul(multiplier).Round(2), AssetCode: assetCode, Status: status, EffectiveAt: effectiveAt, ExpiresAt: expiresAt, Version: version}, nil
}

func (u *WalletUsecase) Wallet(ctx context.Context, userID int64) (map[string]any, error) {
	if userID <= 0 {
		return nil, ErrInvalidArgument
	}
	row, err := u.repo.Wallet(ctx, userID, u.assetCode)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return map[string]any{"asset_code": u.assetCode, "available_balance": "0.00", "frozen_balance": "0.00", "status": string(model.WalletAccountStatusActive)}, nil
	}
	return walletMap(*row), nil
}

func (u *WalletUsecase) CreateRecharge(ctx context.Context, input RechargeCreateInput) (map[string]any, error) {
	if input.UserID <= 0 || input.PackageID <= 0 || strings.TrimSpace(input.IdempotencyKey) == "" || strings.TrimSpace(input.ClientType) == "" {
		return nil, ErrInvalidArgument
	}
	providerCode := strings.ToUpper(strings.TrimSpace(input.Provider))
	provider := u.providers[providerCode]
	if provider == nil {
		return nil, ErrRechargeChannel
	}
	merchantID, appID, err := provider.MerchantAndApp(strings.TrimSpace(input.ClientType))
	if err != nil || strings.TrimSpace(merchantID) == "" || strings.TrimSpace(appID) == "" {
		return nil, ErrRechargeChannel
	}
	packageRow, err := u.repo.RechargePackage(ctx, input.PackageID, u.now())
	if err != nil {
		return nil, err
	}
	credited := packageRow.PayAmount.Mul(packageRow.PackageMultiplier).Round(2)
	if !credited.Equal(packageRow.CreditedAmount.Round(2)) || credited.LessThanOrEqual(decimal.Zero) {
		return nil, ErrInvalidArgument
	}
	if strings.ToUpper(strings.TrimSpace(packageRow.AssetCode)) != u.assetCode {
		return nil, ErrInvalidArgument
	}
	rechargeNo, err := walletNumber("WR")
	if err != nil {
		return nil, err
	}
	outTradeNo, err := walletNumber("WRT")
	if err != nil {
		return nil, err
	}
	now := u.now()
	row := model.WalletRecharge{
		RechargeNo: rechargeNo, UserID: input.UserID, PackageID: packageRow.ID, PackageVersion: packageRow.Version,
		PayAmount: packageRow.PayAmount.Round(2), Currency: packageRow.Currency, PackageMultiplier: packageRow.PackageMultiplier,
		CreditedAmount: credited, AssetCode: packageRow.AssetCode, Provider: providerCode, MerchantID: merchantID, AppID: appID,
		ClientType: strings.TrimSpace(input.ClientType), OutTradeNo: outTradeNo, Status: model.WalletRechargeStatusPending,
		ExpiresAt: now.Add(u.rechargeExpire), Version: 1,
	}
	hash := rechargeRequestHash(input.PackageID, providerCode, input.ClientType)
	row, replayed, err := u.repo.CreateRecharge(ctx, packageRow, row, input.IdempotencyKey, hash)
	if err != nil {
		return nil, err
	}
	responseTime := u.now()
	if replayed && (row.Status != model.WalletRechargeStatusPending && row.Status != model.WalletRechargeStatusPrepayCreated || !row.ExpiresAt.After(responseTime)) {
		result := rechargeMap(row, responseTime)
		result["payment_params"] = map[string]any{}
		return result, nil
	}
	prepay, err := provider.CreatePrepay(ctx, RechargePrepayRequest{RechargeNo: row.RechargeNo, OutTradeNo: row.OutTradeNo, PayAmount: row.PayAmount, Currency: row.Currency, ClientType: row.ClientType, ExpiresAt: row.ExpiresAt})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRechargeChannel, err)
	}
	row, err = u.repo.MarkRechargePrepayCreated(ctx, row.ID, prepay)
	if err != nil {
		return nil, err
	}
	result := rechargeMap(row, u.now())
	result["payment_params"] = prepay.PaymentParams
	return result, nil
}

func (u *WalletUsecase) Recharge(ctx context.Context, userID, rechargeID int64) (map[string]any, error) {
	if userID <= 0 || rechargeID <= 0 {
		return nil, ErrInvalidArgument
	}
	row, err := u.repo.Recharge(ctx, userID, rechargeID)
	if err != nil {
		return nil, err
	}
	return rechargeMap(row, u.now()), nil
}

func (u *WalletUsecase) Recharges(ctx context.Context, userID int64, page Page, status string) (List, error) {
	if userID <= 0 {
		return List{}, ErrInvalidArgument
	}
	status = strings.ToUpper(strings.TrimSpace(status))
	if status != "" && !validRechargeStatus(status) {
		return List{}, ErrInvalidArgument
	}
	rows, err := u.repo.Recharges(ctx, userID, page.After, page.Limit, status)
	if err != nil {
		return List{}, err
	}
	next := trim(&rows, page.Limit, func(v model.WalletRecharge) int64 { return v.ID })
	now := u.now()
	items := make([]map[string]any, 0, len(rows))
	for _, row := range rows {
		items = append(items, rechargeMap(row, now))
	}
	return List{Items: items, NextCursor: next}, nil
}

// CloseExpiredRecharges 处理已经超过支付截止时间的充值单。先进入 CLOSING，
// 再以渠道查单结果为准：已支付则入账，未支付则渠道关单后关闭；不确定时保留
// CLOSING，等待下一轮重试，避免仅依据本地时间误关已付款订单。
func (u *WalletUsecase) CloseExpiredRecharges(ctx context.Context, limit int, retryDelay time.Duration) (RechargeCloseResult, error) {
	if limit <= 0 {
		limit = 100
	}
	if retryDelay <= 0 {
		retryDelay = time.Minute
	}
	now := u.now()
	rows, err := u.repo.ExpiredRecharges(ctx, limit, now, now.Add(-retryDelay))
	if err != nil {
		return RechargeCloseResult{}, err
	}
	result := RechargeCloseResult{Scanned: len(rows)}
	for _, candidate := range rows {
		if err := ctx.Err(); err != nil {
			return result, err
		}
		row, err := u.repo.MarkRechargeClosing(ctx, candidate.ID, candidate.Version)
		if errors.Is(err, ErrRechargeStateConflict) {
			result.Conflicts++
			continue
		}
		if err != nil {
			return result, err
		}
		provider := u.providers[strings.ToUpper(strings.TrimSpace(row.Provider))]
		if provider == nil {
			result.Retrying++
			continue
		}
		request := RechargeQueryRequest{RechargeNo: row.RechargeNo, OutTradeNo: row.OutTradeNo, MerchantID: row.MerchantID, AppID: row.AppID}
		query, err := provider.QueryPayment(ctx, request)
		if err != nil || query.State == RechargePaymentStateUnknown {
			result.Retrying++
			continue
		}
		if query.State == RechargePaymentStatePaid {
			eventKey := strings.TrimSpace(query.EventKey)
			if eventKey == "" {
				eventKey = "QUERY:" + strings.TrimSpace(query.ProviderTransactionID)
			}
			event := VerifiedRechargeEvent{Provider: strings.ToUpper(provider.Code()), MerchantID: row.MerchantID, EventKey: eventKey, EventType: "QUERY", PayloadDigest: query.PayloadDigest, OutTradeNo: row.OutTradeNo, ProviderTransactionID: query.ProviderTransactionID, PayAmount: query.PayAmount, Currency: query.Currency, PaidAt: query.PaidAt}
			if _, _, _, err := u.repo.CreditRecharge(ctx, event); err != nil {
				if errors.Is(err, ErrRechargeStateConflict) {
					result.Conflicts++
					continue
				}
				return result, err
			}
			result.Credited++
			continue
		}
		if query.State == RechargePaymentStateUnpaid {
			if err := provider.ClosePayment(ctx, request); err != nil {
				result.Retrying++
				continue
			}
		}
		if query.State != RechargePaymentStateUnpaid && query.State != RechargePaymentStateClosed {
			result.Retrying++
			continue
		}
		if _, err := u.repo.MarkRechargeClosed(ctx, row.ID, row.Version); err != nil {
			if errors.Is(err, ErrRechargeStateConflict) {
				result.Conflicts++
				continue
			}
			return result, err
		}
		result.Closed++
	}
	return result, nil
}

func (u *WalletUsecase) Transactions(ctx context.Context, userID int64, page Page, entryType string) (List, error) {
	if userID <= 0 {
		return List{}, ErrInvalidArgument
	}
	entryType = strings.ToUpper(strings.TrimSpace(entryType))
	if entryType != "" && entryType != "RECHARGE_CREDIT" {
		return List{}, ErrInvalidArgument
	}
	rows, err := u.repo.LedgerEntries(ctx, userID, page.After, page.Limit, entryType)
	if err != nil {
		return List{}, err
	}
	next := trim(&rows, page.Limit, func(v model.WalletLedgerEntry) int64 { return v.ID })
	items := make([]map[string]any, 0, len(rows))
	for _, row := range rows {
		items = append(items, ledgerMap(row))
	}
	return List{Items: items, NextCursor: next}, nil
}

func (u *WalletUsecase) Callback(ctx context.Context, providerCode string, payload map[string]any) (map[string]any, error) {
	provider := u.providers[strings.ToUpper(strings.TrimSpace(providerCode))]
	if provider == nil {
		return nil, ErrRechargeChannel
	}
	event, err := provider.VerifyCallback(ctx, payload)
	if err != nil {
		return nil, ErrInvalidArgument
	}
	event.Provider = strings.ToUpper(strings.TrimSpace(provider.Code()))
	recharge, wallet, replayed, err := u.repo.CreditRecharge(ctx, event)
	if err != nil {
		return nil, err
	}
	return map[string]any{"recharge_no": recharge.RechargeNo, "status": string(recharge.Status), "available_balance": wallet.AvailableBalance.StringFixed(2), "replayed": replayed}, nil
}

func rechargePackageMap(row model.RechargePackage) map[string]any {
	result := map[string]any{"id": id(row.ID), "name": row.Name, "pay_amount": row.PayAmount.StringFixed(2), "currency": row.Currency, "package_multiplier": row.PackageMultiplier.StringFixed(4), "credited_amount": row.CreditedAmount.StringFixed(2), "asset_code": row.AssetCode, "status": string(row.Status), "version": row.Version}
	if row.EffectiveAt != nil {
		result["effective_at"] = row.EffectiveAt.Format(time.RFC3339)
	}
	if row.ExpiresAt != nil {
		result["expires_at"] = row.ExpiresAt.Format(time.RFC3339)
	}
	return result
}

func walletMap(row model.WalletAccount) map[string]any {
	return map[string]any{"asset_code": row.AssetCode, "available_balance": row.AvailableBalance.StringFixed(2), "frozen_balance": row.FrozenBalance.StringFixed(2), "status": string(row.Status), "version": row.Version}
}

func rechargeMap(row model.WalletRecharge, now time.Time) map[string]any {
	displayStatus, canPay, canRetry := rechargeDisplay(row, now)
	result := map[string]any{"id": id(row.ID), "recharge_no": row.RechargeNo, "package_id": id(row.PackageID), "pay_amount": row.PayAmount.StringFixed(2), "currency": row.Currency, "package_multiplier": row.PackageMultiplier.StringFixed(4), "credited_amount": row.CreditedAmount.StringFixed(2), "asset_code": row.AssetCode, "provider": row.Provider, "status": string(row.Status), "display_status": displayStatus, "can_pay": canPay, "can_retry": canRetry, "created_at": row.CreatedAt.Format(time.RFC3339), "expires_at": row.ExpiresAt.Format(time.RFC3339)}
	if row.PaidAt != nil {
		result["paid_at"] = row.PaidAt.Format(time.RFC3339)
	}
	if row.CreditedAt != nil {
		result["credited_at"] = row.CreditedAt.Format(time.RFC3339)
	}
	if row.FailureCode != nil {
		result["failure_code"] = *row.FailureCode
	}
	return result
}

func rechargeDisplay(row model.WalletRecharge, now time.Time) (string, bool, bool) {
	switch row.Status {
	case model.WalletRechargeStatusPending:
		if !row.ExpiresAt.After(now) {
			return "CONFIRMING", false, false
		}
		return "CREATING", false, false
	case model.WalletRechargeStatusPrepayCreated:
		if !row.ExpiresAt.After(now) {
			return "CONFIRMING", false, false
		}
		return "PENDING_PAYMENT", true, false
	case model.WalletRechargeStatusClosing:
		return "CONFIRMING", false, false
	case model.WalletRechargeStatusSucceeded:
		return "SUCCEEDED", false, false
	case model.WalletRechargeStatusClosed:
		return "EXPIRED", false, true
	case model.WalletRechargeStatusFailed:
		return "FAILED", false, true
	default:
		return "UNKNOWN", false, false
	}
}

func ledgerMap(row model.WalletLedgerEntry) map[string]any {
	return map[string]any{"id": id(row.ID), "asset_code": row.AssetCode, "entry_type": row.EntryType, "biz_type": row.BizType, "biz_id": id(row.BizID), "available_delta": row.AvailableDelta.StringFixed(2), "frozen_delta": row.FrozenDelta.StringFixed(2), "available_after": row.AvailableAfter.StringFixed(2), "frozen_after": row.FrozenAfter.StringFixed(2), "occurred_at": row.OccurredAt.Format(time.RFC3339)}
}

func rechargeRequestHash(packageID int64, provider, clientType string) string {
	raw, _ := json.Marshal([]any{packageID, provider, strings.TrimSpace(clientType)})
	sum := sha256.Sum256(raw)
	return hex.EncodeToString(sum[:])
}

func rechargePackageRequestHash(operation string, packageID int64, input RechargePackageInput) string {
	raw, _ := json.Marshal([]any{operation, packageID, strings.TrimSpace(input.Name), strings.TrimSpace(input.PayAmount), strings.TrimSpace(input.PackageMultiplier), strings.ToUpper(strings.TrimSpace(input.Currency)), strings.ToUpper(strings.TrimSpace(input.AssetCode)), strings.ToUpper(strings.TrimSpace(input.Status)), strings.TrimSpace(input.EffectiveAt), strings.TrimSpace(input.ExpiresAt), input.Version})
	sum := sha256.Sum256(raw)
	return hex.EncodeToString(sum[:])
}

func optionalRFC3339(value string) (*time.Time, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, nil
	}
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return nil, err
	}
	parsed = parsed.UTC()
	return &parsed, nil
}

func walletNumber(prefix string) (string, error) {
	var raw [12]byte
	if _, err := rand.Read(raw[:]); err != nil {
		return "", err
	}
	return prefix + hex.EncodeToString(raw[:]), nil
}

func validRechargeStatus(status string) bool {
	switch model.WalletRechargeStatus(status) {
	case model.WalletRechargeStatusPending, model.WalletRechargeStatusPrepayCreated, model.WalletRechargeStatusSucceeded, model.WalletRechargeStatusClosing, model.WalletRechargeStatusClosed, model.WalletRechargeStatusFailed:
		return true
	default:
		return false
	}
}
