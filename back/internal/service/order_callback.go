package service

import (
	"context"
	commonv1 "dd_web/back/api/common/v1"
	orderv1 "dd_web/back/api/order/v1"
	"dd_web/back/internal/ecode"
)

func (s *OrderService) CallbackPayment(context.Context, *orderv1.CallbackPaymentRequest) (*commonv1.ApiResponse, error) {
	return nil, ecode.New(ecode.ChannelNotConfigured)
}
func (s *OrderService) CallbackRefund(context.Context, *orderv1.CallbackRefundRequest) (*commonv1.ApiResponse, error) {
	return nil, ecode.New(ecode.ChannelNotConfigured)
}
func (s *OrderService) CallbackSettlement(context.Context, *orderv1.CallbackSettlementRequest) (*commonv1.ApiResponse, error) {
	return nil, ecode.New(ecode.ChannelNotConfigured)
}
