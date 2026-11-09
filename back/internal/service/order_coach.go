package service

import (
	"context"
	commonv1 "dd_web/back/api/common/v1"
	orderv1 "dd_web/back/api/order/v1"
)

func (s *OrderService) CoachOrders(context.Context, *orderv1.CoachOrdersRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("coach.orders")
}
func (s *OrderService) CoachOrderDetail(context.Context, *orderv1.CoachOrderDetailRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("coach.orderDetail")
}
func (s *OrderService) CoachOrderEvents(context.Context, *orderv1.CoachOrderEventsRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("coach.orderEvents")
}
func (s *OrderService) CoachAcceptOrder(context.Context, *orderv1.CoachAcceptOrderRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("coach.acceptOrder")
}
func (s *OrderService) CoachRejectOrder(context.Context, *orderv1.CoachRejectOrderRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("coach.rejectOrder")
}
func (s *OrderService) CoachStartOrder(context.Context, *orderv1.CoachStartOrderRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("coach.startOrder")
}
func (s *OrderService) CoachCompleteOrder(context.Context, *orderv1.CoachCompleteOrderRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("coach.completeOrder")
}
func (s *OrderService) CoachEarningsSummary(context.Context, *orderv1.CoachEarningsSummaryRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("coach.earningsSummary")
}
func (s *OrderService) CoachEarningsItems(context.Context, *orderv1.CoachEarningsItemsRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("coach.earningsItems")
}
func (s *OrderService) CoachSettlements(context.Context, *orderv1.CoachSettlementsRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("coach.settlements")
}
func (s *OrderService) CoachSettlementDetail(context.Context, *orderv1.CoachSettlementDetailRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("coach.settlementDetail")
}
func (s *OrderService) CoachSupplements(context.Context, *orderv1.CoachSupplementsRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("coach.supplements")
}
func (s *OrderService) CoachSettlementAccount(context.Context, *orderv1.CoachSettlementAccountRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("coach.settlementAccount")
}
func (s *OrderService) CoachOnboarding(context.Context, *orderv1.CoachOnboardingRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("coach.onboarding")
}
