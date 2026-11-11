package service

import (
	"context"
	businessv1 "dd_web/back/api/business/v1"
	commonv1 "dd_web/back/api/common/v1"
	"dd_web/back/internal/biz"
	"dd_web/back/internal/ecode"
	"errors"
	"strconv"
	"time"
)

// AdminLogin 处理后台管理员登录。
func (s *BusinessService) AdminLogin(ctx context.Context, req *businessv1.AdminLoginRequest) (*commonv1.ApiResponse, error) {
	if s.admin == nil {
		return nil, ecode.New(ecode.AuthNotConfigured)
	}
	if req == nil || req.Body == nil {
		return nil, ecode.New(ecode.InvalidArgument)
	}

	result, err := s.admin.Login(ctx, req.Body.Account, req.Body.Password)
	if err != nil {
		switch {
		case errors.Is(err, biz.ErrInvalidArgument):
			return nil, ecode.New(ecode.InvalidArgument)
		case errors.Is(err, biz.ErrInvalidCredentials):
			return nil, ecode.New(ecode.InvalidCredentials)
		case errors.Is(err, biz.ErrAccountDisabled):
			return nil, ecode.New(ecode.AccountDisabled)
		case errors.Is(err, biz.ErrAuthenticationUnavailable):
			return nil, ecode.New(ecode.AuthNotConfigured)
		default:
			return nil, ecode.Wrap(ecode.InternalError, err)
		}
	}

	return ok(map[string]any{
		"token_type":    "Bearer",
		"access_token":  result.Tokens.AccessToken,
		"refresh_token": result.Tokens.RefreshToken,
		"expires_in": float64(max(
			0,
			int(time.Until(result.Tokens.AccessExpiresAt).Seconds()),
		)),
		"refresh_expires_in": float64(max(
			0,
			int(time.Until(result.Tokens.RefreshExpiresAt).Seconds()),
		)),
		"admin": map[string]any{
			"id":           strconv.FormatInt(result.Admin.ID, 10),
			"account":      result.Admin.Account,
			"display_name": result.Admin.DisplayName,
		},
	})
}

// AdminRefreshToken 刷新后台管理员的登录凭证。
func (s *BusinessService) AdminRefreshToken(context.Context, *businessv1.AdminRefreshTokenRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("admin.refreshToken")
}

// AdminMe 查询当前登录管理员的信息。
func (s *BusinessService) AdminMe(context.Context, *businessv1.AdminMeRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("admin.me")
}

// AdminLogout 退出当前管理员登录。
func (s *BusinessService) AdminLogout(context.Context, *businessv1.AdminLogoutRequest) (*commonv1.ApiResponse, error) {
	if s.admin == nil {
		return nil, ecode.New(ecode.AuthNotConfigured)
	}
	return nil, notImplemented("admin.logout")
}

// AdminAccounts 分页查询后台管理员账号列表。
func (s *BusinessService) AdminAccounts(context.Context, *businessv1.AdminAccountsRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("admin.accounts")
}

// AdminCreateAccount 创建后台管理员账号。
func (s *BusinessService) AdminCreateAccount(ctx context.Context, req *businessv1.AdminCreateAccountRequest) (*commonv1.ApiResponse, error) {
	if s.admin == nil {
		return nil, ecode.New(ecode.AuthNotConfigured)
	}
	if req == nil || req.Body == nil {
		return nil, ecode.New(ecode.InvalidArgument)
	}

	// 创建操作必须由有权限的管理员执行。
	principal, found := biz.PrincipalFromContext(ctx)
	if !found ||
		principal.AdminID <= 0 ||
		!principal.Permissions["admin.account.write"] {
		return nil, ecode.New(ecode.Forbidden)
	}

	body := req.Body
	admin, err := s.admin.CreateAccount(ctx, body.Account, body.InitialPassword, body.DisplayName, body.RoleIds)
	if err != nil {
		switch {
		case errors.Is(err, biz.ErrInvalidArgument):
			return nil, ecode.New(ecode.InvalidArgument)

		case errors.Is(err, biz.ErrAdminAccountExists):
			return nil, ecode.New(ecode.AdminAccountExists)

		case errors.Is(err, biz.ErrAdminRoleInvalid):
			return nil, ecode.New(ecode.AdminRoleInvalid)

		case errors.Is(err, biz.ErrAuthenticationUnavailable):
			return nil, ecode.New(ecode.AuthNotConfigured)

		default:
			return nil, ecode.Wrap(ecode.InternalError, err)
		}
	}

	return ok(map[string]any{
		"id":           strconv.FormatInt(admin.ID, 10),
		"account":      admin.Account,
		"display_name": admin.DisplayName,
		"status":       string(admin.Status),
	})
}

// AdminUpdateAccount 修改后台管理员账号信息、状态和角色。
func (s *BusinessService) AdminUpdateAccount(context.Context, *businessv1.AdminUpdateAccountRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("admin.updateAccount")
}

// AdminRoles 分页查询后台角色列表。
func (s *BusinessService) AdminRoles(context.Context, *businessv1.AdminRolesRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("admin.roles")
}

// AdminCreateRole 创建后台角色并配置权限和数据范围。
func (s *BusinessService) AdminCreateRole(context.Context, *businessv1.AdminCreateRoleRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("admin.createRole")
}

// AdminUpdateRole 修改后台角色信息、权限和数据范围。
func (s *BusinessService) AdminUpdateRole(context.Context, *businessv1.AdminUpdateRoleRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("admin.updateRole")
}

// AdminPermissions 查询后台权限列表。
func (s *BusinessService) AdminPermissions(context.Context, *businessv1.AdminPermissionsRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("admin.permissions")
}

// AdminAuditLogs 分页查询后台操作审计日志。
func (s *BusinessService) AdminAuditLogs(context.Context, *businessv1.AdminAuditLogsRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("admin.auditLogs")
}

// AdminApprovals 分页查询敏感操作审批记录。
func (s *BusinessService) AdminApprovals(context.Context, *businessv1.AdminApprovalsRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("admin.approvals")
}

// AdminRequestApproval 提交敏感操作的审批申请。
func (s *BusinessService) AdminRequestApproval(context.Context, *businessv1.AdminRequestApprovalRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("admin.requestApproval")
}

// AdminDecideApproval 审核敏感操作审批申请。
func (s *BusinessService) AdminDecideApproval(context.Context, *businessv1.AdminDecideApprovalRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("admin.decideApproval")
}

// AdminUsers 分页查询用户列表。
func (s *BusinessService) AdminUsers(context.Context, *businessv1.AdminUsersRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("admin.users")
}

// AdminUserDetail 查询指定用户的详细信息。
func (s *BusinessService) AdminUserDetail(context.Context, *businessv1.AdminUserDetailRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("admin.userDetail")
}

// AdminPenalizeUser 对指定用户执行处罚。
func (s *BusinessService) AdminPenalizeUser(context.Context, *businessv1.AdminPenalizeUserRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("admin.penalizeUser")
}

// AdminRevokePenalty 撤销指定的用户处罚记录。
func (s *BusinessService) AdminRevokePenalty(context.Context, *businessv1.AdminRevokePenaltyRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("admin.revokePenalty")
}

// AdminCoaches 分页查询陪练师列表。
func (s *BusinessService) AdminCoaches(context.Context, *businessv1.AdminCoachesRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("admin.coaches")
}

// AdminCoachDetail 查询指定陪练师的详细信息。
func (s *BusinessService) AdminCoachDetail(context.Context, *businessv1.AdminCoachDetailRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("admin.coachDetail")
}

// AdminApplications 分页查询陪练师入驻申请。
func (s *BusinessService) AdminApplications(context.Context, *businessv1.AdminApplicationsRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("admin.applications")
}

// AdminDecideApplication 审核陪练师入驻申请。
func (s *BusinessService) AdminDecideApplication(context.Context, *businessv1.AdminDecideApplicationRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("admin.decideApplication")
}

// AdminIdentityChecks 分页查询实名认证记录。
func (s *BusinessService) AdminIdentityChecks(context.Context, *businessv1.AdminIdentityChecksRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("admin.identityChecks")
}

// AdminReviewIdentity 人工复核指定实名认证记录。
func (s *BusinessService) AdminReviewIdentity(context.Context, *businessv1.AdminReviewIdentityRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("admin.reviewIdentity")
}

// AdminSkills 分页查询陪练师技能认证记录。
func (s *BusinessService) AdminSkills(context.Context, *businessv1.AdminSkillsRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("admin.skills")
}

// AdminDecideSkill 审核陪练师技能认证申请。
func (s *BusinessService) AdminDecideSkill(context.Context, *businessv1.AdminDecideSkillRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("admin.decideSkill")
}

// AdminGames 分页查询后台游戏列表。
func (s *BusinessService) AdminGames(context.Context, *businessv1.AdminGamesRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("admin.games")
}

// AdminCreateGame 创建游戏目录条目。
func (s *BusinessService) AdminCreateGame(context.Context, *businessv1.AdminCreateGameRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("admin.createGame")
}

// AdminUpdateGame 修改游戏目录信息。
func (s *BusinessService) AdminUpdateGame(context.Context, *businessv1.AdminUpdateGameRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("admin.updateGame")
}

// AdminRegions 分页查询指定游戏的区服列表。
func (s *BusinessService) AdminRegions(context.Context, *businessv1.AdminRegionsRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("admin.regions")
}

// AdminCreateRegion 为指定游戏创建区服。
func (s *BusinessService) AdminCreateRegion(context.Context, *businessv1.AdminCreateRegionRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("admin.createRegion")
}

// AdminUpdateRegion 修改游戏区服信息。
func (s *BusinessService) AdminUpdateRegion(context.Context, *businessv1.AdminUpdateRegionRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("admin.updateRegion")
}

// AdminRanks 分页查询指定游戏的段位列表。
func (s *BusinessService) AdminRanks(context.Context, *businessv1.AdminRanksRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("admin.ranks")
}

// AdminCreateRank 为指定游戏创建段位。
func (s *BusinessService) AdminCreateRank(context.Context, *businessv1.AdminCreateRankRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("admin.createRank")
}

// AdminUpdateRank 修改游戏段位信息。
func (s *BusinessService) AdminUpdateRank(context.Context, *businessv1.AdminUpdateRankRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("admin.updateRank")
}

// AdminCategories 分页查询服务分类列表。
func (s *BusinessService) AdminCategories(context.Context, *businessv1.AdminCategoriesRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("admin.categories")
}

// AdminCreateCategory 创建服务分类。
func (s *BusinessService) AdminCreateCategory(context.Context, *businessv1.AdminCreateCategoryRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("admin.createCategory")
}

// AdminUpdateCategory 修改服务分类信息。
func (s *BusinessService) AdminUpdateCategory(context.Context, *businessv1.AdminUpdateCategoryRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("admin.updateCategory")
}

// AdminServices 分页查询陪练服务列表。
func (s *BusinessService) AdminServices(context.Context, *businessv1.AdminServicesRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("admin.services")
}

// AdminModerateService 对指定陪练服务执行运营管控操作。
func (s *BusinessService) AdminModerateService(context.Context, *businessv1.AdminModerateServiceRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("admin.moderateService")
}

// AdminBanners 分页查询运营横幅列表。
func (s *BusinessService) AdminBanners(context.Context, *businessv1.AdminBannersRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("admin.banners")
}

// AdminCreateBanner 创建运营横幅。
func (s *BusinessService) AdminCreateBanner(context.Context, *businessv1.AdminCreateBannerRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("admin.createBanner")
}

// AdminUpdateBanner 修改运营横幅信息。
func (s *BusinessService) AdminUpdateBanner(context.Context, *businessv1.AdminUpdateBannerRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("admin.updateBanner")
}

// AdminActivities 分页查询运营活动列表。
func (s *BusinessService) AdminActivities(context.Context, *businessv1.AdminActivitiesRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("admin.activities")
}

// AdminCreateActivity 创建运营活动。
func (s *BusinessService) AdminCreateActivity(context.Context, *businessv1.AdminCreateActivityRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("admin.createActivity")
}

// AdminUpdateActivity 修改运营活动信息。
func (s *BusinessService) AdminUpdateActivity(context.Context, *businessv1.AdminUpdateActivityRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("admin.updateActivity")
}

// AdminRecommendations 分页查询陪练推荐配置。
func (s *BusinessService) AdminRecommendations(context.Context, *businessv1.AdminRecommendationsRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("admin.recommendations")
}

// AdminCreateRecommendation 创建陪练推荐配置。
func (s *BusinessService) AdminCreateRecommendation(context.Context, *businessv1.AdminCreateRecommendationRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("admin.createRecommendation")
}

// AdminUpdateRecommendation 修改陪练推荐配置。
func (s *BusinessService) AdminUpdateRecommendation(context.Context, *businessv1.AdminUpdateRecommendationRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("admin.updateRecommendation")
}

// AdminContentAudits 分页查询内容审核记录。
func (s *BusinessService) AdminContentAudits(context.Context, *businessv1.AdminContentAuditsRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("admin.contentAudits")
}

// AdminDecideContent 审核指定内容并提交审核结果。
func (s *BusinessService) AdminDecideContent(context.Context, *businessv1.AdminDecideContentRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("admin.decideContent")
}

// AdminPrepareUpload 为后台媒体资源创建上传准备信息。
func (s *BusinessService) AdminPrepareUpload(context.Context, *businessv1.AdminPrepareUploadRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("admin.prepareUpload")
}

// AdminCompleteUpload 确认后台媒体资源上传完成。
func (s *BusinessService) AdminCompleteUpload(context.Context, *businessv1.AdminCompleteUploadRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("admin.completeUpload")
}

// AdminMediaAccess 获取后台媒体资源的访问凭证或地址。
func (s *BusinessService) AdminMediaAccess(context.Context, *businessv1.AdminMediaAccessRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("admin.mediaAccess")
}
