package api

import (
	khttp "github.com/go-kratos/kratos/v3/transport/http"
)

// CallbackPayment 微信支付结果。
func (a *API) CallbackPayment(ctx khttp.Context) error {
	return reply(ctx, 503, "CHANNEL_NOT_CONFIGURED", "渠道验签与业务处理尚未接入", nil)
}

// CallbackRefund 微信退款结果。
func (a *API) CallbackRefund(ctx khttp.Context) error {
	return reply(ctx, 503, "CHANNEL_NOT_CONFIGURED", "渠道验签与业务处理尚未接入", nil)
}

// CallbackSettlement 结算及补款渠道结果。
func (a *API) CallbackSettlement(ctx khttp.Context) error {
	return reply(ctx, 503, "CHANNEL_NOT_CONFIGURED", "渠道验签与业务处理尚未接入", nil)
}

// CallbackIdentity 实名核验结果。
func (a *API) CallbackIdentity(ctx khttp.Context) error {
	return reply(ctx, 503, "CHANNEL_NOT_CONFIGURED", "渠道验签与业务处理尚未接入", nil)
}

// CallbackContent 内容审核结果。
func (a *API) CallbackContent(ctx khttp.Context) error {
	return reply(ctx, 503, "CHANNEL_NOT_CONFIGURED", "渠道验签与业务处理尚未接入", nil)
}
