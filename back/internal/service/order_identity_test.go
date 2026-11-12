package service

import (
	"context"
	"testing"

	orderv1 "dd_web/back/api/order/v1"
	"dd_web/back/internal/biz"
	"dd_web/back/internal/ecode"
)

func TestTradeCreateOrderRequiresAdultIdentity(t *testing.T) {
	identity := biz.NewIdentityUsecase(nil, &serviceIdentityRepo{})
	service := NewOrderService(nil, identity)
	ctx := biz.WithPrincipal(context.Background(), biz.Principal{UserID: 7})
	_, err := service.TradeCreateOrder(ctx, &orderv1.TradeCreateOrderRequest{})
	businessError, ok := ecode.FromError(err)
	if !ok || businessError.Code() != ecode.AdultIdentityRequired {
		t.Fatalf("order without adult identity must be rejected: %v", err)
	}
}
