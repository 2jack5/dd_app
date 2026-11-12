package service

import (
	"context"
	"errors"

	commonv1 "dd_web/back/api/common/v1"
	orderv1 "dd_web/back/api/order/v1"
	"dd_web/back/internal/biz"
	"dd_web/back/internal/ecode"
)

func (s *OrderService) TradeQuote(context.Context, *orderv1.TradeQuoteRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("trade.quote")
}
func (s *OrderService) TradeCreateOrder(ctx context.Context, _ *orderv1.TradeCreateOrderRequest) (*commonv1.ApiResponse, error) {
	principal, found := biz.PrincipalFromContext(ctx)
	if !found || principal.UserID <= 0 {
		return nil, ecode.New(ecode.AuthRequired)
	}
	if s.identity == nil {
		return nil, ecode.New(ecode.IdentityNotConfigured)
	}
	if err := s.identity.RequireAdult(ctx, principal.UserID); err != nil {
		if errors.Is(err, biz.ErrAdultIdentityRequired) {
			return nil, ecode.New(ecode.AdultIdentityRequired)
		}
		if errors.Is(err, biz.ErrIdentityUnavailable) {
			return nil, ecode.New(ecode.IdentityNotConfigured)
		}
		return nil, ecode.Wrap(ecode.InternalError, err)
	}
	return nil, notImplemented("trade.createOrder")
}
func (s *OrderService) TradeOrders(context.Context, *orderv1.TradeOrdersRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("trade.orders")
}
func (s *OrderService) TradeOrderDetail(context.Context, *orderv1.TradeOrderDetailRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("trade.orderDetail")
}
func (s *OrderService) TradeOrderEvents(context.Context, *orderv1.TradeOrderEventsRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("trade.orderEvents")
}
func (s *OrderService) TradeCancelOrder(context.Context, *orderv1.TradeCancelOrderRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("trade.cancelOrder")
}
func (s *OrderService) TradeConfirmOrder(context.Context, *orderv1.TradeConfirmOrderRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("trade.confirmOrder")
}
func (s *OrderService) PaymentPrepay(context.Context, *orderv1.PaymentPrepayRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("payment.prepay")
}
func (s *OrderService) PaymentResult(context.Context, *orderv1.PaymentResultRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("payment.result")
}
func (s *OrderService) RefundApply(context.Context, *orderv1.RefundApplyRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("refund.apply")
}
func (s *OrderService) RefundList(context.Context, *orderv1.RefundListRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("refund.list")
}
func (s *OrderService) RefundDetail(context.Context, *orderv1.RefundDetailRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("refund.detail")
}
func (s *OrderService) RefundCancel(context.Context, *orderv1.RefundCancelRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("refund.cancel")
}
func (s *OrderService) ReviewCreate(context.Context, *orderv1.ReviewCreateRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("review.create")
}
func (s *OrderService) ReviewFollowUp(context.Context, *orderv1.ReviewFollowUpRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("review.followUp")
}
func (s *OrderService) ReviewMine(context.Context, *orderv1.ReviewMineRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("review.mine")
}
func (s *OrderService) CaseCreate(context.Context, *orderv1.CaseCreateRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("case.create")
}
func (s *OrderService) CaseList(context.Context, *orderv1.CaseListRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("case.list")
}
func (s *OrderService) CaseDetail(context.Context, *orderv1.CaseDetailRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("case.detail")
}
func (s *OrderService) CaseMessage(context.Context, *orderv1.CaseMessageRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("case.message")
}
