package infra

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"dd_web/back/internal/biz"
	"dd_web/back/internal/config"
	"github.com/shopspring/decimal"
)

// LocalRechargeProvider 是开发和自动化测试使用的确定性渠道适配器。
// 它不模拟生产支付能力，但完整执行预支付、回调验签和入账契约。
type LocalRechargeProvider struct {
	merchantID string
	appID      string
	secret     []byte
	expire     time.Duration
	now        func() time.Time
	mu         sync.Mutex
	payments   map[string]biz.RechargeQueryResult
	closed     map[string]bool
}

func NewLocalRechargeProvider(cfg config.WalletPayment) (*LocalRechargeProvider, error) {
	if !strings.EqualFold(strings.TrimSpace(cfg.Provider), "LOCAL") {
		return nil, errors.New("local recharge provider is disabled")
	}
	if strings.TrimSpace(cfg.MerchantID) == "" || strings.TrimSpace(cfg.AppID) == "" {
		return nil, errors.New("local recharge merchant_id and app_id are required")
	}
	if len(cfg.CallbackSecret) < 32 {
		return nil, errors.New("local recharge callback_secret must contain at least 32 bytes")
	}
	expire := cfg.PrepayExpire
	if expire <= 0 {
		expire = 30 * time.Minute
	}
	return &LocalRechargeProvider{merchantID: strings.TrimSpace(cfg.MerchantID), appID: strings.TrimSpace(cfg.AppID), secret: []byte(cfg.CallbackSecret), expire: expire, now: func() time.Time { return time.Now().UTC() }, payments: make(map[string]biz.RechargeQueryResult), closed: make(map[string]bool)}, nil
}

func (*LocalRechargeProvider) Code() string { return "LOCAL" }

func (p *LocalRechargeProvider) MerchantAndApp(clientType string) (string, string, error) {
	if strings.TrimSpace(clientType) == "" {
		return "", "", errors.New("client_type is required")
	}
	return p.merchantID, p.appID, nil
}

func (p *LocalRechargeProvider) CreatePrepay(_ context.Context, req biz.RechargePrepayRequest) (biz.RechargePrepayResult, error) {
	if req.PayAmount.LessThanOrEqual(decimal.Zero) || req.OutTradeNo == "" {
		return biz.RechargePrepayResult{}, errors.New("invalid local prepay request")
	}
	expiresAt := req.ExpiresAt
	if expiresAt.IsZero() {
		expiresAt = p.now().Add(p.expire)
	}
	return biz.RechargePrepayResult{ExpiresAt: expiresAt, PaymentParams: map[string]any{"provider": "LOCAL", "out_trade_no": req.OutTradeNo, "pay_amount": req.PayAmount.StringFixed(2), "currency": req.Currency, "expires_at": expiresAt.Format(time.RFC3339)}}, nil
}

func (p *LocalRechargeProvider) QueryPayment(_ context.Context, req biz.RechargeQueryRequest) (biz.RechargeQueryResult, error) {
	if req.MerchantID != p.merchantID || strings.TrimSpace(req.OutTradeNo) == "" {
		return biz.RechargeQueryResult{}, errors.New("invalid local query identity")
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	if paid, ok := p.payments[req.OutTradeNo]; ok {
		return paid, nil
	}
	if p.closed[req.OutTradeNo] {
		return biz.RechargeQueryResult{State: biz.RechargePaymentStateClosed}, nil
	}
	return biz.RechargeQueryResult{State: biz.RechargePaymentStateUnpaid}, nil
}

func (p *LocalRechargeProvider) ClosePayment(_ context.Context, req biz.RechargeQueryRequest) error {
	if req.MerchantID != p.merchantID || strings.TrimSpace(req.OutTradeNo) == "" {
		return errors.New("invalid local close identity")
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	if _, paid := p.payments[req.OutTradeNo]; paid {
		return errors.New("local recharge is already paid")
	}
	p.closed[req.OutTradeNo] = true
	return nil
}

func (p *LocalRechargeProvider) VerifyCallback(_ context.Context, payload map[string]any) (biz.VerifiedRechargeEvent, error) {
	value := func(key string) string { raw, _ := payload[key].(string); return strings.TrimSpace(raw) }
	event := biz.VerifiedRechargeEvent{Provider: "LOCAL", MerchantID: value("merchant_id"), EventKey: value("event_key"), EventType: "PAY_NOTIFY", OutTradeNo: value("out_trade_no"), ProviderTransactionID: value("transaction_id"), Currency: value("currency")}
	if event.MerchantID != p.merchantID || event.EventKey == "" || event.OutTradeNo == "" || event.ProviderTransactionID == "" || event.Currency == "" {
		return event, errors.New("invalid local callback identity")
	}
	amount, err := decimal.NewFromString(value("pay_amount"))
	if err != nil || amount.Exponent() < -2 || amount.LessThanOrEqual(decimal.Zero) {
		return event, errors.New("invalid local callback amount")
	}
	event.PayAmount = amount.Round(2)
	paidAt, err := time.Parse(time.RFC3339, value("paid_at"))
	if err != nil {
		return event, errors.New("invalid local callback paid_at")
	}
	event.PaidAt = paidAt.UTC()
	canonical := fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n%s", event.MerchantID, event.EventKey, event.OutTradeNo, event.ProviderTransactionID, event.PayAmount.StringFixed(2), event.Currency, event.PaidAt.Format(time.RFC3339))
	mac := hmac.New(sha256.New, p.secret)
	_, _ = mac.Write([]byte(canonical))
	expected := hex.EncodeToString(mac.Sum(nil))
	provided, err := hex.DecodeString(value("signature"))
	if err != nil || !hmac.Equal(provided, mustDecodeHex(expected)) {
		return event, errors.New("invalid local callback signature")
	}
	digest := sha256.Sum256([]byte(canonical))
	event.PayloadDigest = hex.EncodeToString(digest[:])
	p.mu.Lock()
	if !p.closed[event.OutTradeNo] {
		p.payments[event.OutTradeNo] = biz.RechargeQueryResult{State: biz.RechargePaymentStatePaid, EventKey: event.EventKey, PayloadDigest: event.PayloadDigest, ProviderTransactionID: event.ProviderTransactionID, PayAmount: event.PayAmount, Currency: event.Currency, PaidAt: event.PaidAt}
	}
	p.mu.Unlock()
	return event, nil
}

func mustDecodeHex(value string) []byte {
	decoded, _ := hex.DecodeString(value)
	return decoded
}
