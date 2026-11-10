package service

import (
	"context"
	businessv1 "dd_web/back/api/business/v1"
	commonv1 "dd_web/back/api/common/v1"
	kerrors "github.com/go-kratos/kratos/v3/errors"
)

// CallbackIdentity 接收实名认证渠道的结果回调。
func (s *BusinessService) CallbackIdentity(context.Context, *businessv1.CallbackIdentityRequest) (*commonv1.ApiResponse, error) {
	return nil, kerrors.ServiceUnavailable("CHANNEL_NOT_CONFIGURED", "渠道验签与业务处理尚未接入")
}

// CallbackContent 接收内容审核渠道的结果回调。
func (s *BusinessService) CallbackContent(context.Context, *businessv1.CallbackContentRequest) (*commonv1.ApiResponse, error) {
	return nil, kerrors.ServiceUnavailable("CHANNEL_NOT_CONFIGURED", "渠道验签与业务处理尚未接入")
}
