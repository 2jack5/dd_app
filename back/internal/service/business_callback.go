package service

import (
	"context"
	businessv1 "dd_web/back/api/business/v1"
	commonv1 "dd_web/back/api/common/v1"
	"dd_web/back/internal/ecode"
)

// CallbackIdentity 接收实名认证渠道的结果回调。
func (s *BusinessService) CallbackIdentity(context.Context, *businessv1.CallbackIdentityRequest) (*commonv1.ApiResponse, error) {
	return nil, ecode.New(ecode.ChannelNotConfigured)
}

// CallbackContent 接收内容审核渠道的结果回调。
func (s *BusinessService) CallbackContent(context.Context, *businessv1.CallbackContentRequest) (*commonv1.ApiResponse, error) {
	return nil, ecode.New(ecode.ChannelNotConfigured)
}
