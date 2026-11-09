package service

import (
	"context"
	commonv1 "dd_web/back/api/common/v1"
	orderv1 "dd_web/back/api/order/v1"
	kerrors "github.com/go-kratos/kratos/v3/errors"
)

func (s *OrderService) CallbackPayment(context.Context, *orderv1.CallbackPaymentRequest) (*commonv1.ApiResponse, error) {
	return nil, kerrors.ServiceUnavailable("CHANNEL_NOT_CONFIGURED", "渠道验签与业务处理尚未接入")
}
func (s *OrderService) CallbackRefund(context.Context, *orderv1.CallbackRefundRequest) (*commonv1.ApiResponse, error) {
	return nil, kerrors.ServiceUnavailable("CHANNEL_NOT_CONFIGURED", "渠道验签与业务处理尚未接入")
}
func (s *OrderService) CallbackSettlement(context.Context, *orderv1.CallbackSettlementRequest) (*commonv1.ApiResponse, error) {
	return nil, kerrors.ServiceUnavailable("CHANNEL_NOT_CONFIGURED", "渠道验签与业务处理尚未接入")
}
