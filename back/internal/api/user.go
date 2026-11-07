package api

import (
	"encoding/json"
	khttp "github.com/go-kratos/kratos/v3/transport/http"
)

// UserProfile 我的资料。
func (a *API) UserProfile(ctx khttp.Context) error {
	return notImplemented(ctx, "user.profile")
}

// UserUpdateProfileRequest 提交资料修改审核的请求契约；字段校验与业务用例尚未实现。
type UserUpdateProfileRequest struct {
	Nickname       *string `json:"nickname,omitempty"`        // 可选。
	AvatarAssetID  *string `json:"avatar_asset_id,omitempty"` // 可选。十进制int64 ID，以字符串传输，服务端还应校验int64上限
	ContentVersion int32   `json:"content_version"`           // 必填。
}

// UserUpdateProfile 提交资料修改审核。
func (a *API) UserUpdateProfile(ctx khttp.Context) error {
	return notImplemented(ctx, "user.updateProfile")
}

// AuthLogout 退出并吊销当前会话。
func (a *API) AuthLogout(ctx khttp.Context) error {
	return notImplemented(ctx, "auth.logout")
}

// UserCancelAccountRequest 申请注销并校验未结交易的请求契约；字段校验与业务用例尚未实现。
type UserCancelAccountRequest struct {
	ReasonCode   string  `json:"reason_code"`             // 必填。
	ReasonDetail *string `json:"reason_detail,omitempty"` // 可选。
	Version      int32   `json:"version"`                 // 必填。
}

// UserCancelAccount 申请注销并校验未结交易。
func (a *API) UserCancelAccount(ctx khttp.Context) error {
	return notImplemented(ctx, "user.cancelAccount")
}

// NotificationList 业务通知列表。
func (a *API) NotificationList(ctx khttp.Context) error {
	return notImplemented(ctx, "notification.list")
}

// NotificationRead 标记本人通知已读。
func (a *API) NotificationRead(ctx khttp.Context) error {
	return notImplemented(ctx, "notification.read")
}

// AgreementAcceptRequest 记录协议同意的请求契约；字段校验与业务用例尚未实现。
type AgreementAcceptRequest struct {
	AgreementCode    string  `json:"agreement_code"`     // 必填。 取值：USER/PRIVACY/COACH/ORDER
	AgreementVersion string  `json:"agreement_version"`  // 必填。
	ContentHash      string  `json:"content_hash"`       // 必填。
	OrderID          *string `json:"order_id,omitempty"` // 可选。十进制int64 ID，以字符串传输，服务端还应校验int64上限
}

// AgreementAccept 记录协议同意。
func (a *API) AgreementAccept(ctx khttp.Context) error {
	return notImplemented(ctx, "agreement.accept")
}

// MediaPrepareUploadRequest 申请私有资源上传凭证的请求契约；字段校验与业务用例尚未实现。
type MediaPrepareUploadRequest struct {
	Purpose   string `json:"purpose"`    // 必填。 取值：AVATAR/SKILL/REVIEW/CASE
	MIMEType  string `json:"mime_type"`  // 必填。
	SizeBytes int64  `json:"size_bytes"` // 必填。
	SHA256    string `json:"sha256"`     // 必填。
}

// MediaPrepareUpload 申请私有资源上传凭证。
func (a *API) MediaPrepareUpload(ctx khttp.Context) error {
	return notImplemented(ctx, "media.prepareUpload")
}

// MediaCompleteUpload 校验上传完成并进入审核。
func (a *API) MediaCompleteUpload(ctx khttp.Context) error {
	return notImplemented(ctx, "media.completeUpload")
}

// MediaAccess 按归属、用途和审核状态获取资源访问地址。
func (a *API) MediaAccess(ctx khttp.Context) error {
	return notImplemented(ctx, "media.access")
}

// AnalyticsActiveRequest 上报本人活跃事件的请求契约；字段校验与业务用例尚未实现。
type AnalyticsActiveRequest struct {
	EventKey   string          `json:"event_key"`         // 必填。
	EventType  string          `json:"event_type"`        // 必填。 取值：APP_ACTIVE/VIEW_SERVICE/CREATE_ORDER
	OccurredAt string          `json:"occurred_at"`       // 必填。
	Context    json.RawMessage `json:"context,omitempty"` // 可选。
}

// AnalyticsActive 上报本人活跃事件。
func (a *API) AnalyticsActive(ctx khttp.Context) error {
	return notImplemented(ctx, "analytics.active")
}

// TradeQuoteRequest 服务端报价、不锁库存的请求契约；字段校验与业务用例尚未实现。
type TradeQuoteRequest struct {
	ServiceID          string `json:"service_id"`           // 必填。十进制int64 ID，以字符串传输，服务端还应校验int64上限
	RegionID           string `json:"region_id"`            // 必填。十进制int64 ID，以字符串传输，服务端还应校验int64上限
	Quantity           int64  `json:"quantity"`             // 必填。
	AppointmentStartAt string `json:"appointment_start_at"` // 必填。
}

// TradeQuote 服务端报价、不锁库存。
func (a *API) TradeQuote(ctx khttp.Context) error {
	return notImplemented(ctx, "trade.quote")
}

// TradeCreateOrderRequest 消费报价并原子预占时段的请求契约；字段校验与业务用例尚未实现。
type TradeCreateOrderRequest struct {
	QuoteNo               string  `json:"quote_no"`                // 必填。
	GamePlayerID          string  `json:"game_player_id"`          // 必填。
	Remark                *string `json:"remark,omitempty"`        // 可选。
	AgreementAcceptanceID string  `json:"agreement_acceptance_id"` // 必填。十进制int64 ID，以字符串传输，服务端还应校验int64上限
}

// TradeCreateOrder 消费报价并原子预占时段。
func (a *API) TradeCreateOrder(ctx khttp.Context) error {
	return notImplemented(ctx, "trade.createOrder")
}

// TradeOrders 本人订单列表。
func (a *API) TradeOrders(ctx khttp.Context) error {
	return notImplemented(ctx, "trade.orders")
}

// TradeOrderDetail 本人订单详情、履约与资金状态。
func (a *API) TradeOrderDetail(ctx khttp.Context) error {
	return notImplemented(ctx, "trade.orderDetail")
}

// TradeOrderEvents 本人订单时间轴。
func (a *API) TradeOrderEvents(ctx khttp.Context) error {
	return notImplemented(ctx, "trade.orderEvents")
}

// TradeCancelOrderRequest 按阶段取消并协调关单或退款的请求契约；字段校验与业务用例尚未实现。
type TradeCancelOrderRequest struct {
	ReasonCode   string  `json:"reason_code"`             // 必填。
	ReasonDetail *string `json:"reason_detail,omitempty"` // 可选。
	Version      int32   `json:"version"`                 // 必填。
}

// TradeCancelOrder 按阶段取消并协调关单或退款。
func (a *API) TradeCancelOrder(ctx khttp.Context) error {
	return notImplemented(ctx, "trade.cancelOrder")
}

// TradeConfirmOrder 确认完成并激活结算账期。
func (a *API) TradeConfirmOrder(ctx khttp.Context) error {
	return notImplemented(ctx, "trade.confirmOrder")
}

// PaymentPrepayRequest 生成或复用支付会话，不接收客户端金额的请求契约；字段校验与业务用例尚未实现。
type PaymentPrepayRequest struct {
	ClientType string `json:"client_type"` // 必填。 取值：WECHAT_APP/WECHAT_MINI
}

// PaymentPrepay 生成或复用支付会话，不接收客户端金额。
func (a *API) PaymentPrepay(ctx khttp.Context) error {
	return notImplemented(ctx, "payment.prepay")
}

// PaymentResult 查询真实支付结果。
func (a *API) PaymentResult(ctx khttp.Context) error {
	return notImplemented(ctx, "payment.result")
}

// RefundApplyRequest 申请部分或全额退款的请求契约；字段校验与业务用例尚未实现。
type RefundApplyRequest struct {
	AmountCent       int64    `json:"amount_cent"`                  // 必填。
	ReasonCode       string   `json:"reason_code"`                  // 必填。
	ReasonDetail     *string  `json:"reason_detail,omitempty"`      // 可选。
	EvidenceAssetIDs []string `json:"evidence_asset_ids,omitempty"` // 可选。需校验存在、用途、归属与审核状态
	CaseID           *string  `json:"case_id,omitempty"`            // 可选。十进制int64 ID，以字符串传输，服务端还应校验int64上限
}

// RefundApply 申请部分或全额退款。
func (a *API) RefundApply(ctx khttp.Context) error {
	return notImplemented(ctx, "refund.apply")
}

// RefundList 本人退款列表。
func (a *API) RefundList(ctx khttp.Context) error {
	return notImplemented(ctx, "refund.list")
}

// RefundDetail 本人退款详情与通道进度。
func (a *API) RefundDetail(ctx khttp.Context) error {
	return notImplemented(ctx, "refund.detail")
}

// RefundCancel 仅撤销尚未批准且政策允许的退款申请。
func (a *API) RefundCancel(ctx khttp.Context) error {
	return notImplemented(ctx, "refund.cancel")
}

// ReviewCreateRequest 已完成订单提交主评的请求契约；字段校验与业务用例尚未实现。
type ReviewCreateRequest struct {
	Rating        int64    `json:"rating"`                    // 必填。
	Content       string   `json:"content"`                   // 必填。
	ImageAssetIDs []string `json:"image_asset_ids,omitempty"` // 可选。需校验存在、用途、归属与审核状态
}

// ReviewCreate 已完成订单提交主评。
func (a *API) ReviewCreate(ctx khttp.Context) error {
	return notImplemented(ctx, "review.create")
}

// ReviewFollowUpRequest 提交一次追加评价的请求契约；字段校验与业务用例尚未实现。
type ReviewFollowUpRequest struct {
	Content       string   `json:"content"`                   // 必填。
	ImageAssetIDs []string `json:"image_asset_ids,omitempty"` // 可选。需校验存在、用途、归属与审核状态
}

// ReviewFollowUp 提交一次追加评价。
func (a *API) ReviewFollowUp(ctx khttp.Context) error {
	return notImplemented(ctx, "review.followUp")
}

// ReviewMine 本人评价及审核进度。
func (a *API) ReviewMine(ctx khttp.Context) error {
	return notImplemented(ctx, "review.mine")
}

// CaseCreateRequest 投诉、客服、评价申诉统一工单的请求契约；字段校验与业务用例尚未实现。
type CaseCreateRequest struct {
	Kind             string   `json:"kind"`                         // 必填。 取值：COMPLAINT/SUPPORT/REVIEW_APPEAL
	OrderID          *string  `json:"order_id,omitempty"`           // 可选。十进制int64 ID，以字符串传输，服务端还应校验int64上限
	TargetUserID     *string  `json:"target_user_id,omitempty"`     // 可选。十进制int64 ID，以字符串传输，服务端还应校验int64上限
	ReviewID         *string  `json:"review_id,omitempty"`          // 可选。十进制int64 ID，以字符串传输，服务端还应校验int64上限
	CategoryCode     string   `json:"category_code"`                // 必填。
	Content          string   `json:"content"`                      // 必填。
	EvidenceAssetIDs []string `json:"evidence_asset_ids,omitempty"` // 可选。需校验存在、用途、归属与审核状态
}

// CaseCreate 投诉、客服、评价申诉统一工单。
func (a *API) CaseCreate(ctx khttp.Context) error {
	return notImplemented(ctx, "case.create")
}

// CaseList 本人发起工单列表。
func (a *API) CaseList(ctx khttp.Context) error {
	return notImplemented(ctx, "case.list")
}

// CaseDetail 工单详情与对用户可见进度。
func (a *API) CaseDetail(ctx khttp.Context) error {
	return notImplemented(ctx, "case.detail")
}

// CaseMessageRequest 补充说明及证据的请求契约；字段校验与业务用例尚未实现。
type CaseMessageRequest struct {
	Content          string   `json:"content"`                      // 必填。
	EvidenceAssetIDs []string `json:"evidence_asset_ids,omitempty"` // 可选。需校验存在、用途、归属与审核状态
}

// CaseMessage 补充说明及证据。
func (a *API) CaseMessage(ctx khttp.Context) error {
	return notImplemented(ctx, "case.message")
}
