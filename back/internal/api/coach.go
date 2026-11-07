package api

import (
	khttp "github.com/go-kratos/kratos/v3/transport/http"
)

// CoachApplicationCurrent 当前入驻状态及驳回原因。
func (a *API) CoachApplicationCurrent(ctx khttp.Context) error {
	return notImplemented(ctx, "coach.applicationCurrent")
}

// CoachApplyRequest 提交入驻申请的请求契约；字段校验与业务用例尚未实现。
type CoachApplyRequest struct {
	IdentityCheckID       string   `json:"identity_check_id"`            // 必填。十进制int64 ID，以字符串传输，服务端还应校验int64上限
	AgreementAcceptanceID string   `json:"agreement_acceptance_id"`      // 必填。十进制int64 ID，以字符串传输，服务端还应校验int64上限
	DisplayName           string   `json:"display_name"`                 // 必填。
	Bio                   *string  `json:"bio,omitempty"`                // 可选。
	EvidenceAssetIDs      []string `json:"evidence_asset_ids,omitempty"` // 可选。需校验存在、用途、归属与审核状态
}

// CoachApply 提交入驻申请。
func (a *API) CoachApply(ctx khttp.Context) error {
	return notImplemented(ctx, "coach.apply")
}

// CoachWithdraw 撤回本人待审入驻申请。
func (a *API) CoachWithdraw(ctx khttp.Context) error {
	return notImplemented(ctx, "coach.withdraw")
}

// CoachVerifyIdentityRequest 发起第三方实名与年龄核验的请求契约；字段校验与业务用例尚未实现。
type CoachVerifyIdentityRequest struct {
	Provider       string `json:"provider"`        // 必填。
	ConsentVersion string `json:"consent_version"` // 必填。
}

// CoachVerifyIdentity 发起第三方实名与年龄核验。
func (a *API) CoachVerifyIdentity(ctx khttp.Context) error {
	return notImplemented(ctx, "coach.verifyIdentity")
}

// CoachIdentityCurrent 当前有效实名及年龄资格。
func (a *API) CoachIdentityCurrent(ctx khttp.Context) error {
	return notImplemented(ctx, "coach.identityCurrent")
}

// CoachSubmitSkillRequest 提交游戏技能段位证明的请求契约；字段校验与业务用例尚未实现。
type CoachSubmitSkillRequest struct {
	GameID           string   `json:"game_id"`            // 必填。十进制int64 ID，以字符串传输，服务端还应校验int64上限
	RankID           string   `json:"rank_id"`            // 必填。十进制int64 ID，以字符串传输，服务端还应校验int64上限
	EvidenceAssetIDs []string `json:"evidence_asset_ids"` // 必填。需校验存在、用途、归属与审核状态
}

// CoachSubmitSkill 提交游戏技能段位证明。
func (a *API) CoachSubmitSkill(ctx khttp.Context) error {
	return notImplemented(ctx, "coach.submitSkill")
}

// CoachSkills 本人技能版本与审核状态。
func (a *API) CoachSkills(ctx khttp.Context) error {
	return notImplemented(ctx, "coach.skills")
}

// CoachSkillDetail 本人技能审核详情。
func (a *API) CoachSkillDetail(ctx khttp.Context) error {
	return notImplemented(ctx, "coach.skillDetail")
}

// CoachProfile 陪练档案。
func (a *API) CoachProfile(ctx khttp.Context) error {
	return notImplemented(ctx, "coach.profile")
}

// CoachUpdateProfileRequest 提交陪练资料审核的请求契约；字段校验与业务用例尚未实现。
type CoachUpdateProfileRequest struct {
	DisplayName    *string `json:"display_name,omitempty"` // 可选。
	Bio            *string `json:"bio,omitempty"`          // 可选。
	ContentVersion int32   `json:"content_version"`        // 必填。
}

// CoachUpdateProfile 提交陪练资料审核。
func (a *API) CoachUpdateProfile(ctx khttp.Context) error {
	return notImplemented(ctx, "coach.updateProfile")
}

// CoachWorkbench 工作台订单、审核、收益摘要。
func (a *API) CoachWorkbench(ctx khttp.Context) error {
	return notImplemented(ctx, "coach.workbench")
}

// CoachAcceptingRequest 开启或关闭接新单的请求契约；字段校验与业务用例尚未实现。
type CoachAcceptingRequest struct {
	AcceptingOrders bool  `json:"accepting_orders"` // 必填。
	Version         int32 `json:"version"`          // 必填。
}

// CoachAccepting 开启或关闭接新单。
func (a *API) CoachAccepting(ctx khttp.Context) error {
	return notImplemented(ctx, "coach.accepting")
}

// CoachHeartbeat 更新本人在线心跳。
func (a *API) CoachHeartbeat(ctx khttp.Context) error {
	return notImplemented(ctx, "coach.heartbeat")
}

// CoachServices 本人全部服务。
func (a *API) CoachServices(ctx khttp.Context) error {
	return notImplemented(ctx, "coach.services")
}

// CoachCreateServiceRequest 新建服务草稿的请求契约；字段校验与业务用例尚未实现。
type CoachCreateServiceRequest struct {
	GameID                string   `json:"game_id"`                  // 必填。十进制int64 ID，以字符串传输，服务端还应校验int64上限
	CategoryID            string   `json:"category_id"`              // 必填。十进制int64 ID，以字符串传输，服务端还应校验int64上限
	Title                 string   `json:"title"`                    // 必填。
	Description           string   `json:"description"`              // 必填。
	PriceCent             int64    `json:"price_cent"`               // 必填。
	MinQty                int64    `json:"min_qty"`                  // 必填。
	MaxQty                int64    `json:"max_qty"`                  // 必填。
	BookingMinutesPerUnit int64    `json:"booking_minutes_per_unit"` // 必填。
	RegionIDs             []string `json:"region_ids"`               // 必填。
}

// CoachCreateService 新建服务草稿。
func (a *API) CoachCreateService(ctx khttp.Context) error {
	return notImplemented(ctx, "coach.createService")
}

// CoachUpdateServiceRequest 修改价格、文案与服务区服的请求契约；字段校验与业务用例尚未实现。
type CoachUpdateServiceRequest struct {
	Title                 *string  `json:"title,omitempty"`                    // 可选。
	Description           *string  `json:"description,omitempty"`              // 可选。
	PriceCent             *int64   `json:"price_cent,omitempty"`               // 可选。
	MinQty                *int64   `json:"min_qty,omitempty"`                  // 可选。
	MaxQty                *int64   `json:"max_qty,omitempty"`                  // 可选。
	BookingMinutesPerUnit *int64   `json:"booking_minutes_per_unit,omitempty"` // 可选。
	RegionIDs             []string `json:"region_ids,omitempty"`               // 可选。
	Version               int32    `json:"version"`                            // 必填。
	ContentVersion        *int32   `json:"content_version,omitempty"`          // 可选。
}

// CoachUpdateService 修改价格、文案与服务区服。
func (a *API) CoachUpdateService(ctx khttp.Context) error {
	return notImplemented(ctx, "coach.updateService")
}

// CoachPublishService 校验认证、价格和审核后上架。
func (a *API) CoachPublishService(ctx khttp.Context) error {
	return notImplemented(ctx, "coach.publishService")
}

// CoachUnpublishService 下架本人服务。
func (a *API) CoachUnpublishService(ctx khttp.Context) error {
	return notImplemented(ctx, "coach.unpublishService")
}

// CoachAvailability 本人排期与占用情况。
func (a *API) CoachAvailability(ctx khttp.Context) error {
	return notImplemented(ctx, "coach.availability")
}

// CoachUpdateAvailabilityRequest 带版本批量修改指定日期排期的请求契约；字段校验与业务用例尚未实现。
type CoachUpdateAvailabilityRequest struct {
	Date          string `json:"date"` // 必填。
	ExpectedSlots []struct {
		SlotID  string `json:"slot_id"` // 必填。十进制int64 ID，以字符串传输，服务端还应校验int64上限
		Version int64  `json:"version"` // 必填。
	} `json:"expected_slots"` // 必填。
	Slots []struct {
		StartAt string `json:"start_at"` // 必填。
		EndAt   string `json:"end_at"`   // 必填。
	} `json:"slots"` // 必填。
}

// CoachUpdateAvailability 带版本批量修改指定日期排期。
func (a *API) CoachUpdateAvailability(ctx khttp.Context) error {
	return notImplemented(ctx, "coach.updateAvailability")
}

// CoachOrders 本人待接、进行中及历史订单。
func (a *API) CoachOrders(ctx khttp.Context) error {
	return notImplemented(ctx, "coach.orders")
}

// CoachOrderDetail 本人服务订单详情。
func (a *API) CoachOrderDetail(ctx khttp.Context) error {
	return notImplemented(ctx, "coach.orderDetail")
}

// CoachOrderEvents 本人订单履约时间轴。
func (a *API) CoachOrderEvents(ctx khttp.Context) error {
	return notImplemented(ctx, "coach.orderEvents")
}

// CoachAcceptOrder 期限内接单并重新校验资格。
func (a *API) CoachAcceptOrder(ctx khttp.Context) error {
	return notImplemented(ctx, "coach.acceptOrder")
}

// CoachRejectOrderRequest 拒单、释放时段并申请退款的请求契约；字段校验与业务用例尚未实现。
type CoachRejectOrderRequest struct {
	ReasonCode   string  `json:"reason_code"`             // 必填。
	ReasonDetail *string `json:"reason_detail,omitempty"` // 可选。
	Version      int32   `json:"version"`                 // 必填。
}

// CoachRejectOrder 拒单、释放时段并申请退款。
func (a *API) CoachRejectOrder(ctx khttp.Context) error {
	return notImplemented(ctx, "coach.rejectOrder")
}

// CoachStartOrder 允许时间范围内开始服务。
func (a *API) CoachStartOrder(ctx khttp.Context) error {
	return notImplemented(ctx, "coach.startOrder")
}

// CoachCompleteOrderRequest 提交数量、实际结束、证据与说明的请求契约；字段校验与业务用例尚未实现。
type CoachCompleteOrderRequest struct {
	SchemaVersion     int64    `json:"schema_version"`               // 必填。 取值：1
	CompletedQuantity int64    `json:"completed_quantity"`           // 必填。
	ActualEndAt       string   `json:"actual_end_at"`                // 必填。
	EvidenceAssetIDs  []string `json:"evidence_asset_ids,omitempty"` // 可选。需校验存在、用途、归属与审核状态
	CompletionNote    *string  `json:"completion_note,omitempty"`    // 可选。
	Version           int32    `json:"version"`                      // 必填。
}

// CoachCompleteOrder 提交数量、实际结束、证据与说明。
func (a *API) CoachCompleteOrder(ctx khttp.Context) error {
	return notImplemented(ctx, "coach.completeOrder")
}

// CoachEarningsSummary 应得、佣金、待结算、追偿汇总。
func (a *API) CoachEarningsSummary(ctx khttp.Context) error {
	return notImplemented(ctx, "coach.earningsSummary")
}

// CoachEarningsItems 逐单收入与佣金调整。
func (a *API) CoachEarningsItems(ctx khttp.Context) error {
	return notImplemented(ctx, "coach.earningsItems")
}

// CoachSettlements 本人渠道结算批次。
func (a *API) CoachSettlements(ctx khttp.Context) error {
	return notImplemented(ctx, "coach.settlements")
}

// CoachSettlementDetail 本人结算明细及失败原因。
func (a *API) CoachSettlementDetail(ctx khttp.Context) error {
	return notImplemented(ctx, "coach.settlementDetail")
}

// CoachSupplements 本人独立补款记录。
func (a *API) CoachSupplements(ctx khttp.Context) error {
	return notImplemented(ctx, "coach.supplements")
}

// CoachSettlementAccount 本人渠道收款账户及进件状态。
func (a *API) CoachSettlementAccount(ctx khttp.Context) error {
	return notImplemented(ctx, "coach.settlementAccount")
}

// CoachOnboardingRequest 发起渠道收款人进件，不开平台钱包的请求契约；字段校验与业务用例尚未实现。
type CoachOnboardingRequest struct {
	Provider       string `json:"provider"`        // 必填。
	ConsentVersion string `json:"consent_version"` // 必填。
}

// CoachOnboarding 发起渠道收款人进件，不开平台钱包。
func (a *API) CoachOnboarding(ctx khttp.Context) error {
	return notImplemented(ctx, "coach.onboarding")
}
