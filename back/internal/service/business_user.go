package service

import (
	"context"
	businessv1 "dd_web/back/api/business/v1"
	commonv1 "dd_web/back/api/common/v1"
)

// UserProfile 查询当前登录用户的个人资料。
func (s *BusinessService) UserProfile(context.Context, *businessv1.UserProfileRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("user.profile")
}

// UserUpdateProfile 修改当前登录用户的个人资料。
func (s *BusinessService) UserUpdateProfile(context.Context, *businessv1.UserUpdateProfileRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("user.updateProfile")
}

// AuthLogout 退出当前用户登录。
func (s *BusinessService) AuthLogout(context.Context, *businessv1.AuthLogoutRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("auth.logout")
}

// UserCancelAccount 提交当前用户的账号注销申请。
func (s *BusinessService) UserCancelAccount(context.Context, *businessv1.UserCancelAccountRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("user.cancelAccount")
}

// NotificationList 分页查询当前用户的通知列表。
func (s *BusinessService) NotificationList(context.Context, *businessv1.NotificationListRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("notification.list")
}

// NotificationRead 将指定用户通知标记为已读。
func (s *BusinessService) NotificationRead(context.Context, *businessv1.NotificationReadRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("notification.read")
}

// AgreementAccept 记录当前用户对指定版本协议的同意。
func (s *BusinessService) AgreementAccept(context.Context, *businessv1.AgreementAcceptRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("agreement.accept")
}

// MediaPrepareUpload 为用户媒体资源创建上传准备信息。
func (s *BusinessService) MediaPrepareUpload(context.Context, *businessv1.MediaPrepareUploadRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("media.prepareUpload")
}

// MediaCompleteUpload 确认用户媒体资源上传完成。
func (s *BusinessService) MediaCompleteUpload(context.Context, *businessv1.MediaCompleteUploadRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("media.completeUpload")
}

// MediaAccess 获取用户媒体资源的访问凭证或地址。
func (s *BusinessService) MediaAccess(context.Context, *businessv1.MediaAccessRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("media.access")
}

// AnalyticsActive 上报用户活跃事件。
func (s *BusinessService) AnalyticsActive(context.Context, *businessv1.AnalyticsActiveRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("analytics.active")
}
