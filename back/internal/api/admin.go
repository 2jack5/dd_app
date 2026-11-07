package api

import (
	"encoding/json"
	khttp "github.com/go-kratos/kratos/v3/transport/http"
)

// AdminLoginRequest 后台账号密码及二次验证登录的请求契约；字段校验与业务用例尚未实现。
type AdminLoginRequest struct {
	Account  string  `json:"account"`            // 必填。
	Password string  `json:"password"`           // 必填。
	MFACode  *string `json:"mfa_code,omitempty"` // 可选。
}

// AdminLogin 后台账号密码及二次验证登录。
func (a *API) AdminLogin(ctx khttp.Context) error {
	return notImplemented(ctx, "admin.login")
}

// AdminRefreshTokenRequest 后台刷新凭据换取新登录态的请求契约；字段校验与业务用例尚未实现。
type AdminRefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"` // 必填。
}

// AdminRefreshToken 后台刷新凭据换取新登录态。
func (a *API) AdminRefreshToken(ctx khttp.Context) error {
	return notImplemented(ctx, "admin.refreshToken")
}

// AdminMe 管理员资料、角色、权限。
func (a *API) AdminMe(ctx khttp.Context) error {
	return notImplemented(ctx, "admin.me")
}

// AdminLogout 注销当前后台会话。
func (a *API) AdminLogout(ctx khttp.Context) error {
	return notImplemented(ctx, "admin.logout")
}

// AdminAccounts 管理员账号列表。
func (a *API) AdminAccounts(ctx khttp.Context) error {
	return notImplemented(ctx, "admin.accounts")
}

// AdminCreateAccountRequest 创建管理员的请求契约；字段校验与业务用例尚未实现。
type AdminCreateAccountRequest struct {
	Account         string   `json:"account"`          // 必填。
	InitialPassword string   `json:"initial_password"` // 必填。
	DisplayName     string   `json:"display_name"`     // 必填。
	RoleIDs         []string `json:"role_ids"`         // 必填。
}

// AdminCreateAccount 创建管理员。
func (a *API) AdminCreateAccount(ctx khttp.Context) error {
	return notImplemented(ctx, "admin.createAccount")
}

// AdminUpdateAccountRequest 停用账号或调整角色的请求契约；字段校验与业务用例尚未实现。
type AdminUpdateAccountRequest struct {
	DisplayName *string  `json:"display_name,omitempty"` // 可选。
	Status      *string  `json:"status,omitempty"`       // 可选。 取值：ACTIVE/DISABLED
	RoleIDs     []string `json:"role_ids,omitempty"`     // 可选。
}

// AdminUpdateAccount 停用账号或调整角色。
func (a *API) AdminUpdateAccount(ctx khttp.Context) error {
	return notImplemented(ctx, "admin.updateAccount")
}

// AdminRoles 角色列表及数据范围。
func (a *API) AdminRoles(ctx khttp.Context) error {
	return notImplemented(ctx, "admin.roles")
}

// AdminCreateRoleRequest 创建角色的请求契约；字段校验与业务用例尚未实现。
type AdminCreateRoleRequest struct {
	Code          string   `json:"code"`           // 必填。
	Name          string   `json:"name"`           // 必填。
	DataScope     string   `json:"data_scope"`     // 必填。 取值：ALL/ASSIGNED
	PermissionIDs []string `json:"permission_ids"` // 必填。
	Status        string   `json:"status"`         // 必填。 取值：ACTIVE/DISABLED
}

// AdminCreateRole 创建角色。
func (a *API) AdminCreateRole(ctx khttp.Context) error {
	return notImplemented(ctx, "admin.createRole")
}

// AdminUpdateRoleRequest 修改角色权限及数据范围的请求契约；字段校验与业务用例尚未实现。
type AdminUpdateRoleRequest struct {
	Name          *string  `json:"name,omitempty"`           // 可选。
	DataScope     *string  `json:"data_scope,omitempty"`     // 可选。 取值：ALL/ASSIGNED
	PermissionIDs []string `json:"permission_ids,omitempty"` // 可选。
	Status        *string  `json:"status,omitempty"`         // 可选。 取值：ACTIVE/DISABLED
}

// AdminUpdateRole 修改角色权限及数据范围。
func (a *API) AdminUpdateRole(ctx khttp.Context) error {
	return notImplemented(ctx, "admin.updateRole")
}

// AdminPermissions 接口动作权限字典。
func (a *API) AdminPermissions(ctx khttp.Context) error {
	return notImplemented(ctx, "admin.permissions")
}

// AdminAuditLogs 敏感操作审计查询。
func (a *API) AdminAuditLogs(ctx khttp.Context) error {
	return notImplemented(ctx, "admin.auditLogs")
}

// AdminApprovals 本人可处理敏感操作复核队列。
func (a *API) AdminApprovals(ctx khttp.Context) error {
	return notImplemented(ctx, "admin.approvals")
}

// AdminRequestApprovalRequest 提交带资源版本和参数哈希的敏感操作的请求契约；字段校验与业务用例尚未实现。
type AdminRequestApprovalRequest struct {
	Operation       string          `json:"operation"`        // 必填。operation白名单由服务端维护，例settlement.batch.create、settlement.adjustment.create、settlement.supplement.create
	ResourceType    string          `json:"resource_type"`    // 必填。 取值：REFUND/COMMISSION_RULE/BUSINESS_RULE/SETTLEMENT_ITEM/SETTLEMENT_SUPPLEMENT
	ResourceID      string          `json:"resource_id"`      // 必填。十进制int64 ID，以字符串传输，服务端还应校验int64上限
	ResourceVersion int32           `json:"resource_version"` // 必填。
	Parameters      json.RawMessage `json:"parameters"`       // 必填。随operation使用服务端白名单schema验证；创建类操作必须覆盖完整执行参数及服务端金额/账户/版本快照；不能仅哈希客户端可见ID。
}

// AdminRequestApproval 提交带资源版本和参数哈希的敏感操作。
func (a *API) AdminRequestApproval(ctx khttp.Context) error {
	return notImplemented(ctx, "admin.requestApproval")
}

// AdminDecideApprovalRequest 双人复核决定，禁止自审的请求契约；字段校验与业务用例尚未实现。
type AdminDecideApprovalRequest struct {
	Decision            string  `json:"decision"`              // 必填。 取值：APPROVE/REJECT
	ExpectedPayloadHash string  `json:"expected_payload_hash"` // 必填。
	Reason              *string `json:"reason,omitempty"`      // 可选。
}

// AdminDecideApproval 双人复核决定，禁止自审。
func (a *API) AdminDecideApproval(ctx khttp.Context) error {
	return notImplemented(ctx, "admin.decideApproval")
}

// AdminUsers 用户列表及状态。
func (a *API) AdminUsers(ctx khttp.Context) error {
	return notImplemented(ctx, "admin.users")
}

// AdminUserDetail 用户详情，敏感信息脱敏。
func (a *API) AdminUserDetail(ctx khttp.Context) error {
	return notImplemented(ctx, "admin.userDetail")
}

// AdminPenalizeUserRequest 记录用户、陪练或服务违规处罚的请求契约；字段校验与业务用例尚未实现。
type AdminPenalizeUserRequest struct {
	Action    string  `json:"action"`               // 必填。 取值：BAN_USER/SUSPEND_COACH/BLOCK_SERVICE
	CoachID   *string `json:"coach_id,omitempty"`   // 可选。十进制int64 ID，以字符串传输，服务端还应校验int64上限
	ServiceID *string `json:"service_id,omitempty"` // 可选。十进制int64 ID，以字符串传输，服务端还应校验int64上限
	CaseID    *string `json:"case_id,omitempty"`    // 可选。十进制int64 ID，以字符串传输，服务端还应校验int64上限
	Reason    string  `json:"reason"`               // 必填。
	StartsAt  string  `json:"starts_at"`            // 必填。
	EndsAt    *string `json:"ends_at,omitempty"`    // 可选。
}

// AdminPenalizeUser 记录用户、陪练或服务违规处罚。
func (a *API) AdminPenalizeUser(ctx khttp.Context) error {
	return notImplemented(ctx, "admin.penalizeUser")
}

// AdminRevokePenaltyRequest 撤销处罚并重算其他有效处罚的请求契约；字段校验与业务用例尚未实现。
type AdminRevokePenaltyRequest struct {
	ReasonCode   string  `json:"reason_code"`             // 必填。
	ReasonDetail *string `json:"reason_detail,omitempty"` // 可选。
}

// AdminRevokePenalty 撤销处罚并重算其他有效处罚。
func (a *API) AdminRevokePenalty(ctx khttp.Context) error {
	return notImplemented(ctx, "admin.revokePenalty")
}

// AdminCoaches 陪练档案及资格列表。
func (a *API) AdminCoaches(ctx khttp.Context) error {
	return notImplemented(ctx, "admin.coaches")
}

// AdminCoachDetail 陪练服务、审核和处罚历史。
func (a *API) AdminCoachDetail(ctx khttp.Context) error {
	return notImplemented(ctx, "admin.coachDetail")
}

// AdminApplications 入驻审核队列。
func (a *API) AdminApplications(ctx khttp.Context) error {
	return notImplemented(ctx, "admin.applications")
}

// AdminDecideApplicationRequest 审核指定版本入驻资料的请求契约；字段校验与业务用例尚未实现。
type AdminDecideApplicationRequest struct {
	Decision         string  `json:"decision"`          // 必填。 取值：APPROVE/REJECT
	ExpectedRevision int32   `json:"expected_revision"` // 必填。
	ExpectedStatus   string  `json:"expected_status"`   // 必填。 取值：SUBMITTED/REVIEWING
	Reason           *string `json:"reason,omitempty"`  // 可选。
}

// AdminDecideApplication 审核指定版本入驻资料。
func (a *API) AdminDecideApplication(ctx khttp.Context) error {
	return notImplemented(ctx, "admin.decideApplication")
}

// AdminIdentityChecks 实名及年龄核验结果队列。
func (a *API) AdminIdentityChecks(ctx khttp.Context) error {
	return notImplemented(ctx, "admin.identityChecks")
}

// AdminReviewIdentityRequest 人工复核，不伪造第三方实名通过的请求契约；字段校验与业务用例尚未实现。
type AdminReviewIdentityRequest struct {
	Decision       string `json:"decision"`        // 必填。 取值：APPROVE/REJECT
	ExpectedStatus string `json:"expected_status"` // 必填。 取值：PENDING/PASSED/FAILED/EXPIRED/REVOKED
	ReviewNote     string `json:"review_note"`     // 必填。
}

// AdminReviewIdentity 人工复核，不伪造第三方实名通过。
func (a *API) AdminReviewIdentity(ctx khttp.Context) error {
	return notImplemented(ctx, "admin.reviewIdentity")
}

// AdminSkills 技能段位审核队列。
func (a *API) AdminSkills(ctx khttp.Context) error {
	return notImplemented(ctx, "admin.skills")
}

// AdminDecideSkillRequest 审核指定技能证明版本的请求契约；字段校验与业务用例尚未实现。
type AdminDecideSkillRequest struct {
	Decision         string  `json:"decision"`          // 必填。 取值：APPROVE/REJECT
	ExpectedRevision int32   `json:"expected_revision"` // 必填。
	ExpectedStatus   string  `json:"expected_status"`   // 必填。 取值：SUBMITTED/REVIEWING
	Reason           *string `json:"reason,omitempty"`  // 可选。
}

// AdminDecideSkill 审核指定技能证明版本。
func (a *API) AdminDecideSkill(ctx khttp.Context) error {
	return notImplemented(ctx, "admin.decideSkill")
}

// AdminGames 游戏目录管理。
func (a *API) AdminGames(ctx khttp.Context) error {
	return notImplemented(ctx, "admin.games")
}

// AdminCreateGameRequest 新增游戏的请求契约；字段校验与业务用例尚未实现。
type AdminCreateGameRequest struct {
	Code        string  `json:"code"`                    // 必填。
	Name        string  `json:"name"`                    // 必填。
	IconAssetID *string `json:"icon_asset_id,omitempty"` // 可选。十进制int64 ID，以字符串传输，服务端还应校验int64上限
	Keywords    *string `json:"keywords,omitempty"`      // 可选。
	SortOrder   *int64  `json:"sort_order,omitempty"`    // 可选。
	Status      string  `json:"status"`                  // 必填。 取值：ENABLED/DISABLED
}

// AdminCreateGame 新增游戏。
func (a *API) AdminCreateGame(ctx khttp.Context) error {
	return notImplemented(ctx, "admin.createGame")
}

// AdminUpdateGameRequest 编辑或停用游戏的请求契约；字段校验与业务用例尚未实现。
type AdminUpdateGameRequest struct {
	Name        *string `json:"name,omitempty"`          // 可选。
	IconAssetID *string `json:"icon_asset_id,omitempty"` // 可选。十进制int64 ID，以字符串传输，服务端还应校验int64上限
	Keywords    *string `json:"keywords,omitempty"`      // 可选。
	SortOrder   *int64  `json:"sort_order,omitempty"`    // 可选。
	Status      *string `json:"status,omitempty"`        // 可选。 取值：ENABLED/DISABLED
}

// AdminUpdateGame 编辑或停用游戏。
func (a *API) AdminUpdateGame(ctx khttp.Context) error {
	return notImplemented(ctx, "admin.updateGame")
}

// AdminRegions 区服管理列表。
func (a *API) AdminRegions(ctx khttp.Context) error {
	return notImplemented(ctx, "admin.regions")
}

// AdminCreateRegionRequest 新增区服的请求契约；字段校验与业务用例尚未实现。
type AdminCreateRegionRequest struct {
	Code      string `json:"code"`                 // 必填。
	Name      string `json:"name"`                 // 必填。
	Status    string `json:"status"`               // 必填。 取值：ENABLED/DISABLED
	SortOrder *int64 `json:"sort_order,omitempty"` // 可选。
}

// AdminCreateRegion 新增区服。
func (a *API) AdminCreateRegion(ctx khttp.Context) error {
	return notImplemented(ctx, "admin.createRegion")
}

// AdminUpdateRegionRequest 修改区服及状态的请求契约；字段校验与业务用例尚未实现。
type AdminUpdateRegionRequest struct {
	Name      *string `json:"name,omitempty"`       // 可选。
	Status    *string `json:"status,omitempty"`     // 可选。 取值：ENABLED/DISABLED
	SortOrder *int64  `json:"sort_order,omitempty"` // 可选。
}

// AdminUpdateRegion 修改区服及状态。
func (a *API) AdminUpdateRegion(ctx khttp.Context) error {
	return notImplemented(ctx, "admin.updateRegion")
}

// AdminRanks 段位管理列表。
func (a *API) AdminRanks(ctx khttp.Context) error {
	return notImplemented(ctx, "admin.ranks")
}

// AdminCreateRankRequest 新增段位的请求契约；字段校验与业务用例尚未实现。
type AdminCreateRankRequest struct {
	Code      string `json:"code"`       // 必填。
	Name      string `json:"name"`       // 必填。
	RankLevel int64  `json:"rank_level"` // 必填。
	Status    string `json:"status"`     // 必填。 取值：ENABLED/DISABLED
}

// AdminCreateRank 新增段位。
func (a *API) AdminCreateRank(ctx khttp.Context) error {
	return notImplemented(ctx, "admin.createRank")
}

// AdminUpdateRankRequest 修改段位及状态的请求契约；字段校验与业务用例尚未实现。
type AdminUpdateRankRequest struct {
	Name      *string `json:"name,omitempty"`       // 可选。
	RankLevel *int64  `json:"rank_level,omitempty"` // 可选。
	Status    *string `json:"status,omitempty"`     // 可选。 取值：ENABLED/DISABLED
}

// AdminUpdateRank 修改段位及状态。
func (a *API) AdminUpdateRank(ctx khttp.Context) error {
	return notImplemented(ctx, "admin.updateRank")
}

// AdminCategories 分类和价格边界。
func (a *API) AdminCategories(ctx khttp.Context) error {
	return notImplemented(ctx, "admin.categories")
}

// AdminCreateCategoryRequest 新增计价分类的请求契约；字段校验与业务用例尚未实现。
type AdminCreateCategoryRequest struct {
	GameID         string `json:"game_id"`                   // 必填。十进制int64 ID，以字符串传输，服务端还应校验int64上限
	Code           string `json:"code"`                      // 必填。
	Name           string `json:"name"`                      // 必填。
	UnitType       string `json:"unit_type"`                 // 必填。 取值：HOUR/ROUND
	BillingMinutes *int64 `json:"billing_minutes,omitempty"` // 可选。
	MinPriceCent   int64  `json:"min_price_cent"`            // 必填。
	MaxPriceCent   int64  `json:"max_price_cent"`            // 必填。
	Status         string `json:"status"`                    // 必填。 取值：ENABLED/DISABLED
	SortOrder      *int64 `json:"sort_order,omitempty"`      // 可选。
}

// AdminCreateCategory 新增计价分类。
func (a *API) AdminCreateCategory(ctx khttp.Context) error {
	return notImplemented(ctx, "admin.createCategory")
}

// AdminUpdateCategoryRequest 修改价格边界和状态，禁止改变已售计价语义的请求契约；字段校验与业务用例尚未实现。
type AdminUpdateCategoryRequest struct {
	Name         *string `json:"name,omitempty"`           // 可选。
	MinPriceCent *int64  `json:"min_price_cent,omitempty"` // 可选。
	MaxPriceCent *int64  `json:"max_price_cent,omitempty"` // 可选。
	Status       *string `json:"status,omitempty"`         // 可选。 取值：ENABLED/DISABLED
	SortOrder    *int64  `json:"sort_order,omitempty"`     // 可选。
}

// AdminUpdateCategory 修改价格边界和状态，禁止改变已售计价语义。
func (a *API) AdminUpdateCategory(ctx khttp.Context) error {
	return notImplemented(ctx, "admin.updateCategory")
}

// AdminServices 服务及价格异常查询。
func (a *API) AdminServices(ctx khttp.Context) error {
	return notImplemented(ctx, "admin.services")
}

// AdminModerateServiceRequest 禁售或解除运营禁售，不直接代陪练改价的请求契约；字段校验与业务用例尚未实现。
type AdminModerateServiceRequest struct {
	Action  string `json:"action"`  // 必填。 取值：BLOCK/UNBLOCK
	Reason  string `json:"reason"`  // 必填。
	Version int32  `json:"version"` // 必填。
}

// AdminModerateService 禁售或解除运营禁售，不直接代陪练改价。
func (a *API) AdminModerateService(ctx khttp.Context) error {
	return notImplemented(ctx, "admin.moderateService")
}

// AdminBanners Banner列表。
func (a *API) AdminBanners(ctx khttp.Context) error {
	return notImplemented(ctx, "admin.banners")
}

// AdminCreateBannerRequest 新增运营位的请求契约；字段校验与业务用例尚未实现。
type AdminCreateBannerRequest struct {
	Placement    string `json:"placement"`            // 必填。
	Title        string `json:"title"`                // 必填。
	ImageAssetID string `json:"image_asset_id"`       // 必填。十进制int64 ID，以字符串传输，服务端还应校验int64上限
	TargetType   string `json:"target_type"`          // 必填。 取值：GAME/COACH/SERVICE/ACTIVITY
	TargetID     string `json:"target_id"`            // 必填。十进制int64 ID，以字符串传输，服务端还应校验int64上限
	StartsAt     string `json:"starts_at"`            // 必填。
	EndsAt       string `json:"ends_at"`              // 必填。
	Status       string `json:"status"`               // 必填。 取值：DRAFT/PUBLISHED/OFFLINE
	SortOrder    *int64 `json:"sort_order,omitempty"` // 可选。
}

// AdminCreateBanner 新增运营位。
func (a *API) AdminCreateBanner(ctx khttp.Context) error {
	return notImplemented(ctx, "admin.createBanner")
}

// AdminUpdateBannerRequest 编辑、排序、上下线的请求契约；字段校验与业务用例尚未实现。
type AdminUpdateBannerRequest struct {
	Placement    *string `json:"placement,omitempty"`      // 可选。
	Title        *string `json:"title,omitempty"`          // 可选。
	ImageAssetID *string `json:"image_asset_id,omitempty"` // 可选。十进制int64 ID，以字符串传输，服务端还应校验int64上限
	TargetType   *string `json:"target_type,omitempty"`    // 可选。 取值：GAME/COACH/SERVICE/ACTIVITY
	TargetID     *string `json:"target_id,omitempty"`      // 可选。十进制int64 ID，以字符串传输，服务端还应校验int64上限
	StartsAt     *string `json:"starts_at,omitempty"`      // 可选。
	EndsAt       *string `json:"ends_at,omitempty"`        // 可选。
	Status       *string `json:"status,omitempty"`         // 可选。 取值：DRAFT/PUBLISHED/OFFLINE
	SortOrder    *int64  `json:"sort_order,omitempty"`     // 可选。
}

// AdminUpdateBanner 编辑、排序、上下线。
func (a *API) AdminUpdateBanner(ctx khttp.Context) error {
	return notImplemented(ctx, "admin.updateBanner")
}

// AdminActivities 展示活动列表。
func (a *API) AdminActivities(ctx khttp.Context) error {
	return notImplemented(ctx, "admin.activities")
}

// AdminCreateActivityRequest 新增展示活动的请求契约；字段校验与业务用例尚未实现。
type AdminCreateActivityRequest struct {
	Title        string          `json:"title"`                    // 必填。
	Content      json.RawMessage `json:"content"`                  // 必填。
	CoverAssetID *string         `json:"cover_asset_id,omitempty"` // 可选。十进制int64 ID，以字符串传输，服务端还应校验int64上限
	StartsAt     string          `json:"starts_at"`                // 必填。
	EndsAt       string          `json:"ends_at"`                  // 必填。
	Status       string          `json:"status"`                   // 必填。 取值：DRAFT/PUBLISHED/OFFLINE
}

// AdminCreateActivity 新增展示活动。
func (a *API) AdminCreateActivity(ctx khttp.Context) error {
	return notImplemented(ctx, "admin.createActivity")
}

// AdminUpdateActivityRequest 编辑活动并更新内容版本的请求契约；字段校验与业务用例尚未实现。
type AdminUpdateActivityRequest struct {
	Title          *string         `json:"title,omitempty"`          // 可选。
	Content        json.RawMessage `json:"content,omitempty"`        // 可选。
	CoverAssetID   *string         `json:"cover_asset_id,omitempty"` // 可选。十进制int64 ID，以字符串传输，服务端还应校验int64上限
	StartsAt       *string         `json:"starts_at,omitempty"`      // 可选。
	EndsAt         *string         `json:"ends_at,omitempty"`        // 可选。
	Status         *string         `json:"status,omitempty"`         // 可选。 取值：DRAFT/PUBLISHED/OFFLINE
	ContentVersion int32           `json:"content_version"`          // 必填。
	Version        int32           `json:"version"`                  // 必填。
}

// AdminUpdateActivity 编辑活动并更新内容版本。
func (a *API) AdminUpdateActivity(ctx khttp.Context) error {
	return notImplemented(ctx, "admin.updateActivity")
}

// AdminRecommendations 陪练推荐配置。
func (a *API) AdminRecommendations(ctx khttp.Context) error {
	return notImplemented(ctx, "admin.recommendations")
}

// AdminCreateRecommendationRequest 新增定时推荐的请求契约；字段校验与业务用例尚未实现。
type AdminCreateRecommendationRequest struct {
	Placement string  `json:"placement"`            // 必填。
	CoachID   string  `json:"coach_id"`             // 必填。十进制int64 ID，以字符串传输，服务端还应校验int64上限
	GameID    *string `json:"game_id,omitempty"`    // 可选。十进制int64 ID，以字符串传输，服务端还应校验int64上限
	StartsAt  string  `json:"starts_at"`            // 必填。
	EndsAt    string  `json:"ends_at"`              // 必填。
	Status    string  `json:"status"`               // 必填。 取值：ENABLED/DISABLED
	SortOrder *int64  `json:"sort_order,omitempty"` // 可选。
}

// AdminCreateRecommendation 新增定时推荐。
func (a *API) AdminCreateRecommendation(ctx khttp.Context) error {
	return notImplemented(ctx, "admin.createRecommendation")
}

// AdminUpdateRecommendationRequest 编辑推荐及上下线的请求契约；字段校验与业务用例尚未实现。
type AdminUpdateRecommendationRequest struct {
	Placement *string `json:"placement,omitempty"`  // 可选。
	CoachID   *string `json:"coach_id,omitempty"`   // 可选。十进制int64 ID，以字符串传输，服务端还应校验int64上限
	GameID    *string `json:"game_id,omitempty"`    // 可选。十进制int64 ID，以字符串传输，服务端还应校验int64上限
	StartsAt  *string `json:"starts_at,omitempty"`  // 可选。
	EndsAt    *string `json:"ends_at,omitempty"`    // 可选。
	Status    *string `json:"status,omitempty"`     // 可选。 取值：ENABLED/DISABLED
	SortOrder *int64  `json:"sort_order,omitempty"` // 可选。
}

// AdminUpdateRecommendation 编辑推荐及上下线。
func (a *API) AdminUpdateRecommendation(ctx khttp.Context) error {
	return notImplemented(ctx, "admin.updateRecommendation")
}

// AdminContentAudits 内容审核队列。
func (a *API) AdminContentAudits(ctx khttp.Context) error {
	return notImplemented(ctx, "admin.contentAudits")
}

// AdminDecideContentRequest 审核指定版本内容，过期结果不可覆盖的请求契约；字段校验与业务用例尚未实现。
type AdminDecideContentRequest struct {
	Decision         string  `json:"decision"`          // 必填。 取值：APPROVE/REJECT
	ExpectedRevision int32   `json:"expected_revision"` // 必填。
	SourceVersion    int32   `json:"source_version"`    // 必填。
	ExpectedStatus   string  `json:"expected_status"`   // 必填。 取值：PENDING/AUTO_PASSED/MANUAL_REVIEW
	Reason           *string `json:"reason,omitempty"`  // 可选。
}

// AdminDecideContent 审核指定版本内容，过期结果不可覆盖。
func (a *API) AdminDecideContent(ctx khttp.Context) error {
	return notImplemented(ctx, "admin.decideContent")
}

// AdminOrders 按授权数据范围查询订单。
func (a *API) AdminOrders(ctx khttp.Context) error {
	return notImplemented(ctx, "admin.orders")
}

// AdminOrderDetail 订单履约和资金详情。
func (a *API) AdminOrderDetail(ctx khttp.Context) error {
	return notImplemented(ctx, "admin.orderDetail")
}

// AdminOrderEvents 订单完整时间轴。
func (a *API) AdminOrderEvents(ctx khttp.Context) error {
	return notImplemented(ctx, "admin.orderEvents")
}

// AdminRefunds 退款审核及处理列表。
func (a *API) AdminRefunds(ctx khttp.Context) error {
	return notImplemented(ctx, "admin.refunds")
}

// AdminRefundDetail 退款详情与额度预留。
func (a *API) AdminRefundDetail(ctx khttp.Context) error {
	return notImplemented(ctx, "admin.refundDetail")
}

// AdminDecideRefundRequest 带批准和版本执行退款决定的请求契约；字段校验与业务用例尚未实现。
type AdminDecideRefundRequest struct {
	Decision        string  `json:"decision"`              // 必填。 取值：APPROVE/REJECT
	ApprovalID      *string `json:"approval_id,omitempty"` // 可选。十进制int64 ID，以字符串传输，服务端还应校验int64上限
	ResourceVersion int32   `json:"resource_version"`      // 必填。
	Reason          *string `json:"reason,omitempty"`      // 可选。
}

// AdminDecideRefund 带批准和版本执行退款决定。
func (a *API) AdminDecideRefund(ctx khttp.Context) error {
	return notImplemented(ctx, "admin.decideRefund")
}

// AdminCases 投诉、客服与评价申诉工单。
func (a *API) AdminCases(ctx khttp.Context) error {
	return notImplemented(ctx, "admin.cases")
}

// AdminCaseDetail 工单证据、内部备注及流转。
func (a *API) AdminCaseDetail(ctx khttp.Context) error {
	return notImplemented(ctx, "admin.caseDetail")
}

// AdminAssignCaseRequest 分配处理人的请求契约；字段校验与业务用例尚未实现。
type AdminAssignCaseRequest struct {
	AssigneeID string `json:"assignee_id"` // 必填。十进制int64 ID，以字符串传输，服务端还应校验int64上限
	Version    int32  `json:"version"`     // 必填。
}

// AdminAssignCase 分配处理人。
func (a *API) AdminAssignCase(ctx khttp.Context) error {
	return notImplemented(ctx, "admin.assignCase")
}

// AdminCaseMessageRequest 回复用户或添加内部备注的请求契约；字段校验与业务用例尚未实现。
type AdminCaseMessageRequest struct {
	Content          string   `json:"content"`                      // 必填。
	EvidenceAssetIDs []string `json:"evidence_asset_ids,omitempty"` // 可选。需校验存在、用途、归属与审核状态
	Visibility       string   `json:"visibility"`                   // 必填。 取值：PUBLIC/INTERNAL
}

// AdminCaseMessage 回复用户或添加内部备注。
func (a *API) AdminCaseMessage(ctx khttp.Context) error {
	return notImplemented(ctx, "admin.caseMessage")
}

// AdminResolveCaseRequest 记录结论并重算全部阻断工单的请求契约；字段校验与业务用例尚未实现。
type AdminResolveCaseRequest struct {
	ResolutionCode string `json:"resolution_code"` // 必填。
	ResolutionNote string `json:"resolution_note"` // 必填。
	Version        int32  `json:"version"`         // 必填。
}

// AdminResolveCase 记录结论并重算全部阻断工单。
func (a *API) AdminResolveCase(ctx khttp.Context) error {
	return notImplemented(ctx, "admin.resolveCase")
}

// AdminReopenCaseRequest 重新调查，重新评估结算阻断的请求契约；字段校验与业务用例尚未实现。
type AdminReopenCaseRequest struct {
	ReasonCode   string  `json:"reason_code"`             // 必填。
	ReasonDetail *string `json:"reason_detail,omitempty"` // 可选。
	Version      int32   `json:"version"`                 // 必填。
}

// AdminReopenCase 重新调查，重新评估结算阻断。
func (a *API) AdminReopenCase(ctx khttp.Context) error {
	return notImplemented(ctx, "admin.reopenCase")
}

// AdminDashboard 用户、订单、GMV、佣金、退款与活跃统计。
func (a *API) AdminDashboard(ctx khttp.Context) error {
	return notImplemented(ctx, "admin.dashboard")
}

// AdminPayments 真实支付记录。
func (a *API) AdminPayments(ctx khttp.Context) error {
	return notImplemented(ctx, "admin.payments")
}

// AdminSettlementItems 应结、冻结、已付与追偿明细。
func (a *API) AdminSettlementItems(ctx khttp.Context) error {
	return notImplemented(ctx, "admin.settlementItems")
}

// AdminBatches 出款批次及逐条结果。
func (a *API) AdminBatches(ctx khttp.Context) error {
	return notImplemented(ctx, "admin.batches")
}

// AdminCreateBatchRequest 复核后创建出款批次的请求契约；字段校验与业务用例尚未实现。
type AdminCreateBatchRequest struct {
	SettlementItemIDs []string `json:"settlement_item_ids"` // 必填。
	ApprovalID        string   `json:"approval_id"`         // 必填。十进制int64 ID，以字符串传输，服务端还应校验int64上限
}

// AdminCreateBatch 复核后创建出款批次。
func (a *API) AdminCreateBatch(ctx khttp.Context) error {
	return notImplemented(ctx, "admin.createBatch")
}

// AdminAdjustments 退款、费用、修正及追偿记录。
func (a *API) AdminAdjustments(ctx khttp.Context) error {
	return notImplemented(ctx, "admin.adjustments")
}

// AdminCreateAdjustmentRequest 复核后追加资金修正，不覆盖历史的请求契约；字段校验与业务用例尚未实现。
type AdminCreateAdjustmentRequest struct {
	SettlementItemID    string `json:"settlement_item_id"`     // 必填。十进制int64 ID，以字符串传输，服务端还应校验int64上限
	Kind                string `json:"kind"`                   // 必填。 取值：FEE/REVERSAL/CORRECTION
	GrossDeltaCent      int64  `json:"gross_delta_cent"`       // 必填。
	CommissionDeltaCent int64  `json:"commission_delta_cent"`  // 必填。
	CoachDeltaCent      int64  `json:"coach_delta_cent"`       // 必填。
	ChannelFeeDeltaCent int64  `json:"channel_fee_delta_cent"` // 必填。
	Reason              string `json:"reason"`                 // 必填。
	ApprovalID          string `json:"approval_id"`            // 必填。十进制int64 ID，以字符串传输，服务端还应校验int64上限
}

// AdminCreateAdjustment 复核后追加资金修正，不覆盖历史。
func (a *API) AdminCreateAdjustment(ctx khttp.Context) error {
	return notImplemented(ctx, "admin.createAdjustment")
}

// AdminSupplements 独立补款指令列表。
func (a *API) AdminSupplements(ctx khttp.Context) error {
	return notImplemented(ctx, "admin.supplements")
}

// AdminCreateSupplementRequest 首次已付后按已批准正向调整补款的请求契约；字段校验与业务用例尚未实现。
type AdminCreateSupplementRequest struct {
	AdjustmentID string `json:"adjustment_id"` // 必填。十进制int64 ID，以字符串传输，服务端还应校验int64上限
	AccountID    string `json:"account_id"`    // 必填。十进制int64 ID，以字符串传输，服务端还应校验int64上限
	ApprovalID   string `json:"approval_id"`   // 必填。十进制int64 ID，以字符串传输，服务端还应校验int64上限
}

// AdminCreateSupplement 首次已付后按已批准正向调整补款。
func (a *API) AdminCreateSupplement(ctx khttp.Context) error {
	return notImplemented(ctx, "admin.createSupplement")
}

// AdminRetrySupplementRequest 查明失败后沿原请求号重试的请求契约；字段校验与业务用例尚未实现。
type AdminRetrySupplementRequest struct {
	ApprovalID      string `json:"approval_id"`      // 必填。十进制int64 ID，以字符串传输，服务端还应校验int64上限
	ResourceVersion int32  `json:"resource_version"` // 必填。
}

// AdminRetrySupplement 查明失败后沿原请求号重试。
func (a *API) AdminRetrySupplement(ctx khttp.Context) error {
	return notImplemented(ctx, "admin.retrySupplement")
}

// AdminCreateReconciliationRequest 按渠道账单日及类型创建对账任务的请求契约；字段校验与业务用例尚未实现。
type AdminCreateReconciliationRequest struct {
	Provider   string `json:"provider"`    // 必填。
	MerchantID string `json:"merchant_id"` // 必填。
	BillType   string `json:"bill_type"`   // 必填。 取值：PAYMENT/REFUND/SETTLEMENT
	BillDate   string `json:"bill_date"`   // 必填。
}

// AdminCreateReconciliation 按渠道账单日及类型创建对账任务。
func (a *API) AdminCreateReconciliation(ctx khttp.Context) error {
	return notImplemented(ctx, "admin.createReconciliation")
}

// AdminReconciliationRuns 对账任务历史与重跑版本。
func (a *API) AdminReconciliationRuns(ctx khttp.Context) error {
	return notImplemented(ctx, "admin.reconciliationRuns")
}

// AdminReconciliationDetail 对账结果和差异汇总。
func (a *API) AdminReconciliationDetail(ctx khttp.Context) error {
	return notImplemented(ctx, "admin.reconciliationDetail")
}

// AdminReconciliationDiffs 对账差异队列。
func (a *API) AdminReconciliationDiffs(ctx khttp.Context) error {
	return notImplemented(ctx, "admin.reconciliationDiffs")
}

// AdminResolveDiffRequest 记录核验依据，不能直接改支付成功的请求契约；字段校验与业务用例尚未实现。
type AdminResolveDiffRequest struct {
	Resolution string `json:"resolution"` // 必填。
	Result     string `json:"result"`     // 必填。 取值：RESOLVED/IGNORED
}

// AdminResolveDiff 记录核验依据，不能直接改支付成功。
func (a *API) AdminResolveDiff(ctx khttp.Context) error {
	return notImplemented(ctx, "admin.resolveDiff")
}

// AdminCommissionRules 佣金规则版本。
func (a *API) AdminCommissionRules(ctx khttp.Context) error {
	return notImplemented(ctx, "admin.commissionRules")
}

// AdminCreateCommissionRequest 新建佣金草稿的请求契约；字段校验与业务用例尚未实现。
type AdminCreateCommissionRequest struct {
	ScopeType   string  `json:"scope_type"`            // 必填。 取值：GLOBAL/GAME_CATEGORY/COACH
	CoachID     *string `json:"coach_id,omitempty"`    // 可选。十进制int64 ID，以字符串传输，服务端还应校验int64上限
	GameID      *string `json:"game_id,omitempty"`     // 可选。十进制int64 ID，以字符串传输，服务端还应校验int64上限
	CategoryID  *string `json:"category_id,omitempty"` // 可选。十进制int64 ID，以字符串传输，服务端还应校验int64上限
	BPS         int64   `json:"bps"`                   // 必填。
	EffectiveAt string  `json:"effective_at"`          // 必填。
	ExpiresAt   *string `json:"expires_at,omitempty"`  // 可选。null或不传表示无结束时间；若提供必须晚于effective_at
}

// AdminCreateCommission 新建佣金草稿。
func (a *API) AdminCreateCommission(ctx khttp.Context) error {
	return notImplemented(ctx, "admin.createCommission")
}

// AdminPublishCommissionRequest 批准后发布规则，不影响存量单的请求契约；字段校验与业务用例尚未实现。
type AdminPublishCommissionRequest struct {
	ApprovalID      string `json:"approval_id"`      // 必填。十进制int64 ID，以字符串传输，服务端还应校验int64上限
	ResourceVersion int32  `json:"resource_version"` // 必填。
}

// AdminPublishCommission 批准后发布规则，不影响存量单。
func (a *API) AdminPublishCommission(ctx khttp.Context) error {
	return notImplemented(ctx, "admin.publishCommission")
}

// AdminBusinessRules 超时、退款、结算、客服规则版本。
func (a *API) AdminBusinessRules(ctx khttp.Context) error {
	return notImplemented(ctx, "admin.businessRules")
}

// AdminCreateBusinessRuleRequest 创建规则草稿的请求契约；字段校验与业务用例尚未实现。
type AdminCreateBusinessRuleRequest struct {
	Code        string          `json:"code"`                 // 必填。 取值：ORDER/CANCELLATION/SETTLEMENT/SUPPORT
	Config      json.RawMessage `json:"config"`               // 必填。
	EffectiveAt string          `json:"effective_at"`         // 必填。
	ExpiresAt   *string         `json:"expires_at,omitempty"` // 可选。null或不传表示无结束时间；若提供必须晚于effective_at
}

// AdminCreateBusinessRule 创建规则草稿。
func (a *API) AdminCreateBusinessRule(ctx khttp.Context) error {
	return notImplemented(ctx, "admin.createBusinessRule")
}

// AdminPublishBusinessRuleRequest 批准后发布业务规则的请求契约；字段校验与业务用例尚未实现。
type AdminPublishBusinessRuleRequest struct {
	ApprovalID      string `json:"approval_id"`      // 必填。十进制int64 ID，以字符串传输，服务端还应校验int64上限
	ResourceVersion int32  `json:"resource_version"` // 必填。
}

// AdminPublishBusinessRule 批准后发布业务规则。
func (a *API) AdminPublishBusinessRule(ctx khttp.Context) error {
	return notImplemented(ctx, "admin.publishBusinessRule")
}

// AdminPrepareUploadRequest 申请后台运营图片上传凭证的请求契约；字段校验与业务用例尚未实现。
type AdminPrepareUploadRequest struct {
	Purpose   string `json:"purpose"`    // 必填。 取值：GAME_ICON/BANNER/ACTIVITY_COVER/SERVICE_IMAGE
	MIMEType  string `json:"mime_type"`  // 必填。
	SizeBytes int64  `json:"size_bytes"` // 必填。
	SHA256    string `json:"sha256"`     // 必填。
}

// AdminPrepareUpload 申请后台运营图片上传凭证。
func (a *API) AdminPrepareUpload(ctx khttp.Context) error {
	return notImplemented(ctx, "admin.prepareUpload")
}

// AdminCompleteUpload 校验后台上传完成并进入审核。
func (a *API) AdminCompleteUpload(ctx khttp.Context) error {
	return notImplemented(ctx, "admin.completeUpload")
}

// AdminMediaAccess 按后台权限获取证据或运营资源访问地址。
func (a *API) AdminMediaAccess(ctx khttp.Context) error {
	return notImplemented(ctx, "admin.mediaAccess")
}
