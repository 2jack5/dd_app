package infra

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"testing"
	"time"

	"dd_web/back/internal/biz"
	"dd_web/back/internal/config"
	"github.com/shopspring/decimal"
)

func TestLocalRechargeProviderPrepayAndVerifiedCallback(t *testing.T) {
	const secret = "0123456789abcdef0123456789abcdef"
	provider, err := NewLocalRechargeProvider(config.WalletPayment{Provider: "LOCAL", MerchantID: "merchant", AppID: "app", CallbackSecret: secret, PrepayExpire: 30 * time.Minute})
	if err != nil {
		t.Fatal(err)
	}
	provider.now = func() time.Time { return time.Date(2026, 11, 12, 12, 0, 0, 0, time.UTC) }
	prepay, err := provider.CreatePrepay(context.Background(), biz.RechargePrepayRequest{OutTradeNo: "WRT1", PayAmount: decimal.RequireFromString("30.00"), Currency: "CNY", ClientType: "APP"})
	if err != nil {
		t.Fatal(err)
	}
	if got := prepay.PaymentParams["pay_amount"]; got != "30.00" {
		t.Fatalf("pay amount = %v", got)
	}

	paidAt := "2026-11-12T12:01:00Z"
	canonical := "merchant\nevent-1\nWRT1\ntx-1\n30.00\nCNY\n" + paidAt
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(canonical))
	payload := map[string]any{"merchant_id": "merchant", "event_key": "event-1", "out_trade_no": "WRT1", "transaction_id": "tx-1", "pay_amount": "30.00", "currency": "CNY", "paid_at": paidAt, "signature": hex.EncodeToString(mac.Sum(nil))}
	event, err := provider.VerifyCallback(context.Background(), payload)
	if err != nil {
		t.Fatal(err)
	}
	if event.OutTradeNo != "WRT1" || event.PayAmount.StringFixed(2) != "30.00" {
		t.Fatalf("event = %#v", event)
	}
	query, err := provider.QueryPayment(context.Background(), biz.RechargeQueryRequest{OutTradeNo: "WRT1", MerchantID: "merchant"})
	if err != nil || query.State != biz.RechargePaymentStatePaid || query.ProviderTransactionID != "tx-1" {
		t.Fatalf("query=%#v err=%v", query, err)
	}
	if err := provider.ClosePayment(context.Background(), biz.RechargeQueryRequest{OutTradeNo: "WRT1", MerchantID: "merchant"}); err == nil {
		t.Fatal("paid recharge must not be closed")
	}
}

func TestLocalRechargeProviderClosesUnpaidRecharge(t *testing.T) {
	provider, err := NewLocalRechargeProvider(config.WalletPayment{Provider: "LOCAL", MerchantID: "merchant", AppID: "app", CallbackSecret: "0123456789abcdef0123456789abcdef"})
	if err != nil {
		t.Fatal(err)
	}
	request := biz.RechargeQueryRequest{OutTradeNo: "WRT-unpaid", MerchantID: "merchant"}
	if err := provider.ClosePayment(context.Background(), request); err != nil {
		t.Fatal(err)
	}
	query, err := provider.QueryPayment(context.Background(), request)
	if err != nil || query.State != biz.RechargePaymentStateClosed {
		t.Fatalf("query=%#v err=%v", query, err)
	}
}
