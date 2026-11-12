package service

import (
	"context"

	businessv1 "dd_web/back/api/business/v1"
	commonv1 "dd_web/back/api/common/v1"
)

// CoachApplicationCurrent 查询当前用户的陪练师入驻申请。
func (s *BusinessService) CoachApplicationCurrent(context.Context, *businessv1.CoachApplicationCurrentRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("coach.applicationCurrent")
}

// CoachApply 提交陪练师入驻申请。
func (s *BusinessService) CoachApply(context.Context, *businessv1.CoachApplyRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("coach.apply")
}

// CoachWithdraw 撤回指定的陪练师入驻申请。
func (s *BusinessService) CoachWithdraw(context.Context, *businessv1.CoachWithdrawRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("coach.withdraw")
}

// CoachSubmitSkill 提交陪练师技能认证申请。
func (s *BusinessService) CoachSubmitSkill(context.Context, *businessv1.CoachSubmitSkillRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("coach.submitSkill")
}

// CoachSkills 分页查询当前用户的技能认证记录。
func (s *BusinessService) CoachSkills(context.Context, *businessv1.CoachSkillsRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("coach.skills")
}

// CoachSkillDetail 查询指定技能认证的详细信息。
func (s *BusinessService) CoachSkillDetail(context.Context, *businessv1.CoachSkillDetailRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("coach.skillDetail")
}

// CoachProfile 查询当前陪练师的个人资料。
func (s *BusinessService) CoachProfile(context.Context, *businessv1.CoachProfileRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("coach.profile")
}

// CoachUpdateProfile 修改当前陪练师的个人资料。
func (s *BusinessService) CoachUpdateProfile(context.Context, *businessv1.CoachUpdateProfileRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("coach.updateProfile")
}

// CoachWorkbench 获取当前陪练师的工作台信息。
func (s *BusinessService) CoachWorkbench(context.Context, *businessv1.CoachWorkbenchRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("coach.workbench")
}

// CoachAccepting 设置当前陪练师是否接单。
func (s *BusinessService) CoachAccepting(context.Context, *businessv1.CoachAcceptingRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("coach.accepting")
}

// CoachHeartbeat 上报当前陪练师的在线心跳。
func (s *BusinessService) CoachHeartbeat(context.Context, *businessv1.CoachHeartbeatRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("coach.heartbeat")
}

// CoachServices 分页查询当前陪练师的服务列表。
func (s *BusinessService) CoachServices(context.Context, *businessv1.CoachServicesRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("coach.services")
}

// CoachCreateService 创建当前陪练师的服务。
func (s *BusinessService) CoachCreateService(context.Context, *businessv1.CoachCreateServiceRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("coach.createService")
}

// CoachUpdateService 修改当前陪练师的指定服务。
func (s *BusinessService) CoachUpdateService(context.Context, *businessv1.CoachUpdateServiceRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("coach.updateService")
}

// CoachPublishService 发布当前陪练师的指定服务。
func (s *BusinessService) CoachPublishService(context.Context, *businessv1.CoachPublishServiceRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("coach.publishService")
}

// CoachUnpublishService 下架当前陪练师的指定服务。
func (s *BusinessService) CoachUnpublishService(context.Context, *businessv1.CoachUnpublishServiceRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("coach.unpublishService")
}

// CoachAvailability 查询当前陪练师指定日期的排期。
func (s *BusinessService) CoachAvailability(context.Context, *businessv1.CoachAvailabilityRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("coach.availability")
}

// CoachUpdateAvailability 更新当前陪练师的可预约排期。
func (s *BusinessService) CoachUpdateAvailability(context.Context, *businessv1.CoachUpdateAvailabilityRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("coach.updateAvailability")
}
