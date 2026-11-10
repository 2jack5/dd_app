package service

import (
	"context"
	businessv1 "dd_web/back/api/business/v1"
	commonv1 "dd_web/back/api/common/v1"
	"dd_web/back/internal/ecode"
)

func (s *BusinessService) CallbackIdentity(context.Context, *businessv1.CallbackIdentityRequest) (*commonv1.ApiResponse, error) {
	return nil, ecode.New(ecode.ChannelNotConfigured)
}
func (s *BusinessService) CallbackContent(context.Context, *businessv1.CallbackContentRequest) (*commonv1.ApiResponse, error) {
	return nil, ecode.New(ecode.ChannelNotConfigured)
}
