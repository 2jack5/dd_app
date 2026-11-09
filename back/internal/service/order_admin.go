package service

import (
	"context"
	commonv1 "dd_web/back/api/common/v1"
	orderv1 "dd_web/back/api/order/v1"
)

func (s *OrderService) AdminOrders(context.Context, *orderv1.AdminOrdersRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("admin.orders")
}
func (s *OrderService) AdminOrderDetail(context.Context, *orderv1.AdminOrderDetailRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("admin.orderDetail")
}
func (s *OrderService) AdminOrderEvents(context.Context, *orderv1.AdminOrderEventsRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("admin.orderEvents")
}
func (s *OrderService) AdminRefunds(context.Context, *orderv1.AdminRefundsRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("admin.refunds")
}
func (s *OrderService) AdminRefundDetail(context.Context, *orderv1.AdminRefundDetailRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("admin.refundDetail")
}
func (s *OrderService) AdminDecideRefund(context.Context, *orderv1.AdminDecideRefundRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("admin.decideRefund")
}
func (s *OrderService) AdminCases(context.Context, *orderv1.AdminCasesRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("admin.cases")
}
func (s *OrderService) AdminCaseDetail(context.Context, *orderv1.AdminCaseDetailRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("admin.caseDetail")
}
func (s *OrderService) AdminAssignCase(context.Context, *orderv1.AdminAssignCaseRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("admin.assignCase")
}
func (s *OrderService) AdminCaseMessage(context.Context, *orderv1.AdminCaseMessageRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("admin.caseMessage")
}
func (s *OrderService) AdminResolveCase(context.Context, *orderv1.AdminResolveCaseRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("admin.resolveCase")
}
func (s *OrderService) AdminReopenCase(context.Context, *orderv1.AdminReopenCaseRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("admin.reopenCase")
}
func (s *OrderService) AdminDashboard(context.Context, *orderv1.AdminDashboardRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("admin.dashboard")
}
func (s *OrderService) AdminPayments(context.Context, *orderv1.AdminPaymentsRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("admin.payments")
}
func (s *OrderService) AdminSettlementItems(context.Context, *orderv1.AdminSettlementItemsRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("admin.settlementItems")
}
func (s *OrderService) AdminBatches(context.Context, *orderv1.AdminBatchesRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("admin.batches")
}
func (s *OrderService) AdminCreateBatch(context.Context, *orderv1.AdminCreateBatchRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("admin.createBatch")
}
func (s *OrderService) AdminAdjustments(context.Context, *orderv1.AdminAdjustmentsRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("admin.adjustments")
}
func (s *OrderService) AdminCreateAdjustment(context.Context, *orderv1.AdminCreateAdjustmentRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("admin.createAdjustment")
}
func (s *OrderService) AdminSupplements(context.Context, *orderv1.AdminSupplementsRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("admin.supplements")
}
func (s *OrderService) AdminCreateSupplement(context.Context, *orderv1.AdminCreateSupplementRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("admin.createSupplement")
}
func (s *OrderService) AdminRetrySupplement(context.Context, *orderv1.AdminRetrySupplementRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("admin.retrySupplement")
}
func (s *OrderService) AdminCreateReconciliation(context.Context, *orderv1.AdminCreateReconciliationRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("admin.createReconciliation")
}
func (s *OrderService) AdminReconciliationRuns(context.Context, *orderv1.AdminReconciliationRunsRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("admin.reconciliationRuns")
}
func (s *OrderService) AdminReconciliationDetail(context.Context, *orderv1.AdminReconciliationDetailRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("admin.reconciliationDetail")
}
func (s *OrderService) AdminReconciliationDiffs(context.Context, *orderv1.AdminReconciliationDiffsRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("admin.reconciliationDiffs")
}
func (s *OrderService) AdminResolveDiff(context.Context, *orderv1.AdminResolveDiffRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("admin.resolveDiff")
}
func (s *OrderService) AdminCommissionRules(context.Context, *orderv1.AdminCommissionRulesRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("admin.commissionRules")
}
func (s *OrderService) AdminCreateCommission(context.Context, *orderv1.AdminCreateCommissionRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("admin.createCommission")
}
func (s *OrderService) AdminPublishCommission(context.Context, *orderv1.AdminPublishCommissionRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("admin.publishCommission")
}
func (s *OrderService) AdminBusinessRules(context.Context, *orderv1.AdminBusinessRulesRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("admin.businessRules")
}
func (s *OrderService) AdminCreateBusinessRule(context.Context, *orderv1.AdminCreateBusinessRuleRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("admin.createBusinessRule")
}
func (s *OrderService) AdminPublishBusinessRule(context.Context, *orderv1.AdminPublishBusinessRuleRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("admin.publishBusinessRule")
}
