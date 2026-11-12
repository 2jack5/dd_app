package service

import (
	"context"
	"errors"
	"strconv"
	"time"

	businessv1 "dd_web/back/api/business/v1"
	commonv1 "dd_web/back/api/common/v1"
	"dd_web/back/internal/biz"
	"dd_web/back/internal/data"
	"dd_web/back/internal/ecode"
)

// UserProfile 查询当前登录用户的个人资料。
func (s *BusinessService) UserProfile(context.Context, *businessv1.UserProfileRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("user.profile")
}

// UserUpdateProfile 修改当前登录用户的个人资料。
func (s *BusinessService) UserUpdateProfile(context.Context, *businessv1.UserUpdateProfileRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("user.updateProfile")
}

// UserVerifyIdentity 提交当前登录用户的实名认证请求。
func (s *BusinessService) UserVerifyIdentity(ctx context.Context, req *businessv1.UserVerifyIdentityRequest) (*commonv1.ApiResponse, error) {
	principal, found := biz.PrincipalFromContext(ctx)
	if !found || principal.UserID <= 0 {
		return nil, ecode.New(ecode.AuthRequired)
	}
	if s.identity == nil {
		return nil, ecode.New(ecode.IdentityNotConfigured)
	}
	if req == nil || req.Body == nil {
		return nil, ecode.New(ecode.InvalidArgument)
	}
	record, err := s.identity.Verify(ctx, principal.UserID, biz.IdentityVerifyInput{
		Provider: req.Body.Provider, ConsentVersion: req.Body.ConsentVersion,
		RealName: req.Body.RealName, IdentityNumber: req.Body.IdentityNumber,
	})
	if err != nil {
		return nil, identityError(err)
	}
	return ok(identityRecordData(record))
}

// UserIdentityCurrent 查询当前登录用户的实名认证结果。
func (s *BusinessService) UserIdentityCurrent(ctx context.Context, _ *businessv1.UserIdentityCurrentRequest) (*commonv1.ApiResponse, error) {
	principal, found := biz.PrincipalFromContext(ctx)
	if !found || principal.UserID <= 0 {
		return nil, ecode.New(ecode.AuthRequired)
	}
	if s.identity == nil {
		return nil, ecode.New(ecode.IdentityNotConfigured)
	}
	record, err := s.identity.Current(ctx, principal.UserID)
	if err != nil {
		return nil, identityError(err)
	}
	return ok(identityRecordData(record))
}

func identityError(err error) error {
	switch {
	case errors.Is(err, biz.ErrInvalidArgument):
		return ecode.New(ecode.InvalidArgument)
	case errors.Is(err, biz.ErrIdentityUnavailable):
		return ecode.New(ecode.IdentityNotConfigured)
	case errors.Is(err, data.ErrNotFound):
		return ecode.New(ecode.NotFound)
	default:
		return ecode.Wrap(ecode.InternalError, err)
	}
}

func identityRecordData(record biz.IdentityRecord) map[string]any {
	result := map[string]any{
		"id": strconv.FormatInt(record.ID, 10), "provider": record.Provider,
		"status": record.Status, "name_masked": record.NameMasked,
		"age_check_result": record.AgeCheckResult, "created_at": record.CreatedAt.UTC().Format(time.RFC3339),
	}
	if record.VerifiedAt != nil {
		result["verified_at"] = record.VerifiedAt.UTC().Format(time.RFC3339)
	}
	if record.ExpiresAt != nil {
		result["expires_at"] = record.ExpiresAt.UTC().Format(time.RFC3339)
	}
	return result
}

// AuthLogout 退出当前用户登录。
func (s *BusinessService) AuthLogout(ctx context.Context, _ *businessv1.AuthLogoutRequest) (*commonv1.ApiResponse, error) {
	principal, found := biz.PrincipalFromContext(ctx)
	if !found || principal.UserID <= 0 || principal.SessionID == "" {
		return nil, ecode.New(ecode.AuthRequired)
	}
	if s.auth == nil {
		return nil, ecode.New(ecode.AuthNotConfigured)
	}
	if err := s.auth.Logout(ctx, principal.UserID, principal.SessionID); err != nil {
		switch {
		case errors.Is(err, biz.ErrInvalidToken):
			return nil, ecode.New(ecode.InvalidToken)
		case errors.Is(err, biz.ErrAuthenticationUnavailable):
			return nil, ecode.New(ecode.AuthNotConfigured)
		default:
			return nil, ecode.Wrap(ecode.InternalError, err)
		}
	}
	return ok(map[string]any{"logged_out": true})
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
