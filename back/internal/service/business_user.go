package service

import (
	"context"
	businessv1 "dd_web/back/api/business/v1"
	commonv1 "dd_web/back/api/common/v1"
)

func (s *BusinessService) UserProfile(context.Context, *businessv1.UserProfileRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("user.profile")
}
func (s *BusinessService) UserUpdateProfile(context.Context, *businessv1.UserUpdateProfileRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("user.updateProfile")
}
func (s *BusinessService) AuthLogout(context.Context, *businessv1.AuthLogoutRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("auth.logout")
}
func (s *BusinessService) UserCancelAccount(context.Context, *businessv1.UserCancelAccountRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("user.cancelAccount")
}
func (s *BusinessService) NotificationList(context.Context, *businessv1.NotificationListRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("notification.list")
}
func (s *BusinessService) NotificationRead(context.Context, *businessv1.NotificationReadRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("notification.read")
}
func (s *BusinessService) AgreementAccept(context.Context, *businessv1.AgreementAcceptRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("agreement.accept")
}
func (s *BusinessService) MediaPrepareUpload(context.Context, *businessv1.MediaPrepareUploadRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("media.prepareUpload")
}
func (s *BusinessService) MediaCompleteUpload(context.Context, *businessv1.MediaCompleteUploadRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("media.completeUpload")
}
func (s *BusinessService) MediaAccess(context.Context, *businessv1.MediaAccessRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("media.access")
}
func (s *BusinessService) AnalyticsActive(context.Context, *businessv1.AnalyticsActiveRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("analytics.active")
}
