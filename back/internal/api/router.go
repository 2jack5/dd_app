package api

import (
	"context"
	"dd_web/back/internal/config"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	khttp "github.com/go-kratos/kratos/v3/transport/http"
	"gorm.io/gorm"
)

// API 统一保存接口需要的数据库与身份适配器；路由只在本文件注册，处理函数按板块分文件。
type API struct {
	db   *gorm.DB
	auth Authenticator
}
type routerConfig struct {
	address string
	timeout time.Duration
	db      *gorm.DB
	auth    Authenticator
}

// NewRouter 只接收main.go初始化的连接并注册HTTP接口，不创建或关闭数据库。
func NewRouter(cfg *config.Config, db *gorm.DB) (*khttp.Server, error) {
	if cfg == nil {
		return nil, fmt.Errorf("configuration is nil")
	}
	c := routerConfig{address: "127.0.0.1:8000", timeout: 5 * time.Second, db: db}
	if v := strings.TrimSpace(cfg.HTTP.Port); v != "" {
		if _, _, err := net.SplitHostPort(v); err == nil {
			c.address = v
		} else {
			n, err := strconv.Atoi(v)
			if err != nil || n < 1 || n > 65535 {
				return nil, fmt.Errorf("http.port must be a port number or host:port")
			}
			c.address = net.JoinHostPort("127.0.0.1", v)
		}
	}
	if cfg.HTTP.ShutdownTimeout < 0 {
		return nil, fmt.Errorf("http.shutdown_timeout must be nonnegative")
	}
	if v, ok := os.LookupEnv("DD_HTTP_ADDR"); ok {
		c.address = v
	}
	if c.address == "" {
		c.address = "127.0.0.1:8000"
	}
	if v := os.Getenv("DD_HTTP_TIMEOUT"); v != "" {
		d, err := time.ParseDuration(v)
		if err != nil || d <= 0 {
			return nil, fmt.Errorf("DD_HTTP_TIMEOUT must be a positive duration")
		}
		c.timeout = d
	}
	return newRouter(c)
}

func newRouter(c routerConfig) (*khttp.Server, error) {
	a := &API{db: c.db, auth: c.auth}
	srv := khttp.NewServer(khttp.Address(c.address), khttp.Timeout(c.timeout), khttp.Filter(requestFilter()),
		khttp.NotFoundHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { writeError(w, r, 404, "NOT_FOUND", "接口不存在") })),
		khttp.MethodNotAllowedHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			writeError(w, r, 405, "METHOD_NOT_ALLOWED", "请求方法不支持")
		})),
	)
	// 健康检查不带业务前缀，处理函数与公共接口放在public.go。
	health := srv.Route("/")
	health.GET("/healthz", a.SystemHealth)
	health.GET("/readyz", a.SystemReady)
	// 全部业务路由集中在此注册；其他api文件只包含处理函数、请求结构与相关辅助。
	r := srv.Route("/api/v1")

	// ==================== 公共基础与发现 ====================
	// 登录与协议
	r.POST("/auth/wechat/login", a.protect(policy{scope: "public"}, a.AuthWechatLogin))         // 微信App或小程序登录
	r.POST("/auth/token/refresh", a.protect(policy{scope: "public"}, a.AuthRefreshToken))       // 使用刷新凭据换取新登录态
	r.GET("/agreements/{code}/current", a.protect(policy{scope: "public"}, a.AgreementCurrent)) // 获取当前有效协议版本

	// 首页与发现
	r.GET("/home", a.protect(policy{scope: "public"}, a.HomeOverview))                        // 首页组合数据：Banner、推荐、活动、游戏
	r.GET("/home/banners", a.protect(policy{scope: "public"}, a.HomeBanners))                 // 首页Banner列表
	r.GET("/home/recommendations", a.protect(policy{scope: "public"}, a.HomeRecommendations)) // 推荐陪练列表
	r.GET("/activities", a.protect(policy{scope: "public"}, a.ActivityList))                  // 展示活动列表
	r.GET("/activities/{activity_id}", a.protect(policy{scope: "public"}, a.ActivityDetail))  // 展示活动详情
	r.GET("/games", a.protect(policy{scope: "public"}, a.CatalogGames))                       // 游戏列表
	r.GET("/games/{game_id}/regions", a.protect(policy{scope: "public"}, a.CatalogRegions))   // 游戏区服
	r.GET("/games/{game_id}/ranks", a.protect(policy{scope: "public"}, a.CatalogRanks))       // 游戏段位
	r.GET("/service-categories", a.protect(policy{scope: "public"}, a.CatalogCategories))     // 服务分类及计价边界
	r.GET("/search/suggestions", a.protect(policy{scope: "public"}, a.SearchSuggestions))     // 游戏、陪练和服务关键词联想

	// 找陪练
	r.GET("/coaches", a.protect(policy{scope: "public"}, a.DiscoveryCoaches))                              // 按游戏、区服、价格、段位、服务方式及时间筛选陪练
	r.GET("/coaches/{coach_id}", a.protect(policy{scope: "public"}, a.DiscoveryCoachDetail))               // 陪练详情及认证标识
	r.GET("/coaches/{coach_id}/services", a.protect(policy{scope: "public"}, a.DiscoveryServices))         // 陪练可售服务
	r.GET("/coaches/{coach_id}/availability", a.protect(policy{scope: "public"}, a.DiscoveryAvailability)) // 陪练可预约时段
	r.GET("/coaches/{coach_id}/reviews", a.protect(policy{scope: "public"}, a.DiscoveryReviews))           // 已公开评价
	r.GET("/services/{service_id}", a.protect(policy{scope: "public"}, a.DiscoveryServiceDetail))          // 可售服务详情
	r.GET("/support/config", a.protect(policy{scope: "public"}, a.SupportConfig))                          // 客服入口、联系方式与帮助配置

	// ==================== 用户端 ====================
	// 账户与通知
	r.GET("/users/me", a.protect(policy{scope: "user"}, a.UserProfile))                                                              // 我的资料
	r.PATCH("/users/me", a.protect(policy{scope: "user", idempotent: true}, a.UserUpdateProfile))                                    // 提交资料修改审核
	r.POST("/auth/logout", a.protect(policy{scope: "user", idempotent: true}, a.AuthLogout))                                         // 退出并吊销当前会话
	r.POST("/users/me/cancellation", a.protect(policy{scope: "user", idempotent: true}, a.UserCancelAccount))                        // 申请注销并校验未结交易
	r.GET("/users/me/notifications", a.protect(policy{scope: "user"}, a.NotificationList))                                           // 业务通知列表
	r.POST("/users/me/notifications/{notification_id}/read", a.protect(policy{scope: "user", idempotent: true}, a.NotificationRead)) // 标记本人通知已读
	r.POST("/agreements/acceptances", a.protect(policy{scope: "user", idempotent: true}, a.AgreementAccept))                         // 记录协议同意
	r.POST("/media/uploads", a.protect(policy{scope: "user", idempotent: true}, a.MediaPrepareUpload))                               // 申请私有资源上传凭证
	r.POST("/media/uploads/{asset_id}/complete", a.protect(policy{scope: "user", idempotent: true}, a.MediaCompleteUpload))          // 校验上传完成并进入审核
	r.GET("/media/{asset_id}/access", a.protect(policy{scope: "user"}, a.MediaAccess))                                               // 按归属、用途和审核状态获取资源访问地址
	r.POST("/activity-events", a.protect(policy{scope: "user", idempotent: true}, a.AnalyticsActive))                                // 上报本人活跃事件

	// 报价与订单
	r.POST("/orders/quotes", a.protect(policy{scope: "user", idempotent: true}, a.TradeQuote))                    // 服务端报价、不锁库存
	r.POST("/orders", a.protect(policy{scope: "user", idempotent: true}, a.TradeCreateOrder))                     // 消费报价并原子预占时段
	r.GET("/orders", a.protect(policy{scope: "user"}, a.TradeOrders))                                             // 本人订单列表
	r.GET("/orders/{order_id}", a.protect(policy{scope: "user"}, a.TradeOrderDetail))                             // 本人订单详情、履约与资金状态
	r.GET("/orders/{order_id}/events", a.protect(policy{scope: "user"}, a.TradeOrderEvents))                      // 本人订单时间轴
	r.POST("/orders/{order_id}/cancel", a.protect(policy{scope: "user", idempotent: true}, a.TradeCancelOrder))   // 按阶段取消并协调关单或退款
	r.POST("/orders/{order_id}/confirm", a.protect(policy{scope: "user", idempotent: true}, a.TradeConfirmOrder)) // 确认完成并激活结算账期

	// 支付与退款
	r.POST("/orders/{order_id}/payment", a.protect(policy{scope: "user", idempotent: true}, a.PaymentPrepay)) // 生成或复用支付会话，不接收客户端金额
	r.GET("/orders/{order_id}/payment", a.protect(policy{scope: "user"}, a.PaymentResult))                    // 查询真实支付结果
	r.POST("/orders/{order_id}/refunds", a.protect(policy{scope: "user", idempotent: true}, a.RefundApply))   // 申请部分或全额退款
	r.GET("/refunds", a.protect(policy{scope: "user"}, a.RefundList))                                         // 本人退款列表
	r.GET("/refunds/{refund_id}", a.protect(policy{scope: "user"}, a.RefundDetail))                           // 本人退款详情与通道进度
	r.POST("/refunds/{refund_id}/cancel", a.protect(policy{scope: "user", idempotent: true}, a.RefundCancel)) // 仅撤销尚未批准且政策允许的退款申请

	// 评价与售后
	r.POST("/orders/{order_id}/reviews", a.protect(policy{scope: "user", idempotent: true}, a.ReviewCreate))        // 已完成订单提交主评
	r.POST("/reviews/{review_id}/follow-ups", a.protect(policy{scope: "user", idempotent: true}, a.ReviewFollowUp)) // 提交一次追加评价
	r.GET("/users/me/reviews", a.protect(policy{scope: "user"}, a.ReviewMine))                                      // 本人评价及审核进度
	r.POST("/cases", a.protect(policy{scope: "user", idempotent: true}, a.CaseCreate))                              // 投诉、客服、评价申诉统一工单
	r.GET("/cases", a.protect(policy{scope: "user"}, a.CaseList))                                                   // 本人发起工单列表
	r.GET("/cases/{case_id}", a.protect(policy{scope: "user"}, a.CaseDetail))                                       // 工单详情与对用户可见进度
	r.POST("/cases/{case_id}/messages", a.protect(policy{scope: "user", idempotent: true}, a.CaseMessage))          // 补充说明及证据

	// ==================== 陪练师端 ====================
	// 入驻与认证
	r.GET("/coach/applications/current", a.protect(policy{scope: "coach_applicant"}, a.CoachApplicationCurrent))                            // 当前入驻状态及驳回原因
	r.POST("/coach/applications", a.protect(policy{scope: "coach_applicant", idempotent: true}, a.CoachApply))                              // 提交入驻申请
	r.POST("/coach/applications/{application_id}/withdraw", a.protect(policy{scope: "coach_applicant", idempotent: true}, a.CoachWithdraw)) // 撤回本人待审入驻申请
	r.POST("/coach/identity-checks", a.protect(policy{scope: "coach_applicant", idempotent: true}, a.CoachVerifyIdentity))                  // 发起第三方实名与年龄核验
	r.GET("/coach/identity-checks/current", a.protect(policy{scope: "coach_applicant"}, a.CoachIdentityCurrent))                            // 当前有效实名及年龄资格
	r.POST("/coach/skills", a.protect(policy{scope: "coach_applicant", idempotent: true}, a.CoachSubmitSkill))                              // 提交游戏技能段位证明
	r.GET("/coach/skills", a.protect(policy{scope: "coach_applicant"}, a.CoachSkills))                                                      // 本人技能版本与审核状态
	r.GET("/coach/skills/{skill_id}", a.protect(policy{scope: "coach_applicant"}, a.CoachSkillDetail))                                      // 本人技能审核详情

	// 服务与排期
	r.GET("/coach/me", a.protect(policy{scope: "coach"}, a.CoachProfile))                                                          // 陪练档案
	r.PATCH("/coach/me", a.protect(policy{scope: "coach", idempotent: true}, a.CoachUpdateProfile))                                // 提交陪练资料审核
	r.GET("/coach/workbench", a.protect(policy{scope: "coach"}, a.CoachWorkbench))                                                 // 工作台订单、审核、收益摘要
	r.PUT("/coach/accepting-orders", a.protect(policy{scope: "coach", idempotent: true}, a.CoachAccepting))                        // 开启或关闭接新单
	r.POST("/coach/heartbeat", a.protect(policy{scope: "coach", idempotent: true}, a.CoachHeartbeat))                              // 更新本人在线心跳
	r.GET("/coach/services", a.protect(policy{scope: "coach"}, a.CoachServices))                                                   // 本人全部服务
	r.POST("/coach/services", a.protect(policy{scope: "coach", idempotent: true}, a.CoachCreateService))                           // 新建服务草稿
	r.PATCH("/coach/services/{service_id}", a.protect(policy{scope: "coach", idempotent: true}, a.CoachUpdateService))             // 修改价格、文案与服务区服
	r.POST("/coach/services/{service_id}/publish", a.protect(policy{scope: "coach", idempotent: true}, a.CoachPublishService))     // 校验认证、价格和审核后上架
	r.POST("/coach/services/{service_id}/unpublish", a.protect(policy{scope: "coach", idempotent: true}, a.CoachUnpublishService)) // 下架本人服务
	r.GET("/coach/availability", a.protect(policy{scope: "coach"}, a.CoachAvailability))                                           // 本人排期与占用情况
	r.PUT("/coach/availability", a.protect(policy{scope: "coach", idempotent: true}, a.CoachUpdateAvailability))                   // 带版本批量修改指定日期排期

	// 接单与履约
	r.GET("/coach/orders", a.protect(policy{scope: "coach"}, a.CoachOrders))                                               // 本人待接、进行中及历史订单
	r.GET("/coach/orders/{order_id}", a.protect(policy{scope: "coach"}, a.CoachOrderDetail))                               // 本人服务订单详情
	r.GET("/coach/orders/{order_id}/events", a.protect(policy{scope: "coach"}, a.CoachOrderEvents))                        // 本人订单履约时间轴
	r.POST("/coach/orders/{order_id}/accept", a.protect(policy{scope: "coach", idempotent: true}, a.CoachAcceptOrder))     // 期限内接单并重新校验资格
	r.POST("/coach/orders/{order_id}/reject", a.protect(policy{scope: "coach", idempotent: true}, a.CoachRejectOrder))     // 拒单、释放时段并申请退款
	r.POST("/coach/orders/{order_id}/start", a.protect(policy{scope: "coach", idempotent: true}, a.CoachStartOrder))       // 允许时间范围内开始服务
	r.POST("/coach/orders/{order_id}/complete", a.protect(policy{scope: "coach", idempotent: true}, a.CoachCompleteOrder)) // 提交数量、实际结束、证据与说明

	// 收益与结算
	r.GET("/coach/earnings/summary", a.protect(policy{scope: "coach"}, a.CoachEarningsSummary))                            // 应得、佣金、待结算、追偿汇总
	r.GET("/coach/earnings/items", a.protect(policy{scope: "coach"}, a.CoachEarningsItems))                                // 逐单收入与佣金调整
	r.GET("/coach/settlements", a.protect(policy{scope: "coach"}, a.CoachSettlements))                                     // 本人渠道结算批次
	r.GET("/coach/settlements/{batch_id}", a.protect(policy{scope: "coach"}, a.CoachSettlementDetail))                     // 本人结算明细及失败原因
	r.GET("/coach/settlement-supplements", a.protect(policy{scope: "coach"}, a.CoachSupplements))                          // 本人独立补款记录
	r.GET("/coach/settlement-account", a.protect(policy{scope: "coach"}, a.CoachSettlementAccount))                        // 本人渠道收款账户及进件状态
	r.POST("/coach/settlement-account/onboarding", a.protect(policy{scope: "coach", idempotent: true}, a.CoachOnboarding)) // 发起渠道收款人进件，不开平台钱包

	// ==================== 运营后台 ====================
	// 登录与权限
	r.POST("/admin/auth/login", a.protect(policy{scope: "public"}, a.AdminLogin))                // 后台账号密码及二次验证登录
	r.POST("/admin/auth/token/refresh", a.protect(policy{scope: "public"}, a.AdminRefreshToken)) // 后台刷新凭据换取新登录态

	// 权限与后台设置
	r.GET("/admin/me", a.protect(policy{scope: "admin", permission: "admin.self.read"}, a.AdminMe))                                                                    // 管理员资料、角色、权限
	r.POST("/admin/auth/logout", a.protect(policy{scope: "admin", permission: "admin.self.read", idempotent: true}, a.AdminLogout))                                    // 注销当前后台会话
	r.GET("/admin/accounts", a.protect(policy{scope: "admin", permission: "admin.account.read"}, a.AdminAccounts))                                                     // 管理员账号列表
	r.POST("/admin/accounts", a.protect(policy{scope: "admin", permission: "admin.account.write", idempotent: true}, a.AdminCreateAccount))                            // 创建管理员
	r.PATCH("/admin/accounts/{admin_id}", a.protect(policy{scope: "admin", permission: "admin.account.write", idempotent: true}, a.AdminUpdateAccount))                // 停用账号或调整角色
	r.GET("/admin/roles", a.protect(policy{scope: "admin", permission: "admin.role.read"}, a.AdminRoles))                                                              // 角色列表及数据范围
	r.POST("/admin/roles", a.protect(policy{scope: "admin", permission: "admin.role.write", idempotent: true}, a.AdminCreateRole))                                     // 创建角色
	r.PATCH("/admin/roles/{role_id}", a.protect(policy{scope: "admin", permission: "admin.role.write", idempotent: true}, a.AdminUpdateRole))                          // 修改角色权限及数据范围
	r.GET("/admin/permissions", a.protect(policy{scope: "admin", permission: "admin.role.read"}, a.AdminPermissions))                                                  // 接口动作权限字典
	r.GET("/admin/audit-logs", a.protect(policy{scope: "admin", permission: "admin.audit.read"}, a.AdminAuditLogs))                                                    // 敏感操作审计查询
	r.GET("/admin/approvals", a.protect(policy{scope: "admin", permission: "admin.approval.read"}, a.AdminApprovals))                                                  // 本人可处理敏感操作复核队列
	r.POST("/admin/approvals", a.protect(policy{scope: "admin", permission: "admin.approval.request", idempotent: true}, a.AdminRequestApproval))                      // 提交带资源版本和参数哈希的敏感操作
	r.POST("/admin/approvals/{approval_id}/decision", a.protect(policy{scope: "admin", permission: "admin.approval.review", idempotent: true}, a.AdminDecideApproval)) // 双人复核决定，禁止自审

	// 人员与审核
	r.GET("/admin/users", a.protect(policy{scope: "admin", permission: "user.read"}, a.AdminUsers))                                                                                      // 用户列表及状态
	r.GET("/admin/users/{user_id}", a.protect(policy{scope: "admin", permission: "user.read"}, a.AdminUserDetail))                                                                       // 用户详情，敏感信息脱敏
	r.POST("/admin/users/{user_id}/penalties", a.protect(policy{scope: "admin", permission: "user.penalize", idempotent: true}, a.AdminPenalizeUser))                                    // 记录用户、陪练或服务违规处罚
	r.POST("/admin/penalties/{penalty_id}/revoke", a.protect(policy{scope: "admin", permission: "user.penalize", idempotent: true}, a.AdminRevokePenalty))                               // 撤销处罚并重算其他有效处罚
	r.GET("/admin/coaches", a.protect(policy{scope: "admin", permission: "coach.read"}, a.AdminCoaches))                                                                                 // 陪练档案及资格列表
	r.GET("/admin/coaches/{coach_id}", a.protect(policy{scope: "admin", permission: "coach.read"}, a.AdminCoachDetail))                                                                  // 陪练服务、审核和处罚历史
	r.GET("/admin/coach-applications", a.protect(policy{scope: "admin", permission: "coach.application.read"}, a.AdminApplications))                                                     // 入驻审核队列
	r.POST("/admin/coach-applications/{application_id}/decision", a.protect(policy{scope: "admin", permission: "coach.application.review", idempotent: true}, a.AdminDecideApplication)) // 审核指定版本入驻资料
	r.GET("/admin/identity-checks", a.protect(policy{scope: "admin", permission: "coach.identity.read"}, a.AdminIdentityChecks))                                                         // 实名及年龄核验结果队列
	r.POST("/admin/identity-checks/{check_id}/review", a.protect(policy{scope: "admin", permission: "coach.identity.review", idempotent: true}, a.AdminReviewIdentity))                  // 人工复核，不伪造第三方实名通过
	r.GET("/admin/coach-skills", a.protect(policy{scope: "admin", permission: "coach.skill.read"}, a.AdminSkills))                                                                       // 技能段位审核队列
	r.POST("/admin/coach-skills/{skill_id}/decision", a.protect(policy{scope: "admin", permission: "coach.skill.review", idempotent: true}, a.AdminDecideSkill))                         // 审核指定技能证明版本

	// 商品与运营
	r.GET("/admin/games", a.protect(policy{scope: "admin", permission: "catalog.read"}, a.AdminGames))                                                                           // 游戏目录管理
	r.POST("/admin/games", a.protect(policy{scope: "admin", permission: "catalog.write", idempotent: true}, a.AdminCreateGame))                                                  // 新增游戏
	r.PATCH("/admin/games/{game_id}", a.protect(policy{scope: "admin", permission: "catalog.write", idempotent: true}, a.AdminUpdateGame))                                       // 编辑或停用游戏
	r.GET("/admin/games/{game_id}/regions", a.protect(policy{scope: "admin", permission: "catalog.read"}, a.AdminRegions))                                                       // 区服管理列表
	r.POST("/admin/games/{game_id}/regions", a.protect(policy{scope: "admin", permission: "catalog.write", idempotent: true}, a.AdminCreateRegion))                              // 新增区服
	r.PATCH("/admin/game-regions/{region_id}", a.protect(policy{scope: "admin", permission: "catalog.write", idempotent: true}, a.AdminUpdateRegion))                            // 修改区服及状态
	r.GET("/admin/games/{game_id}/ranks", a.protect(policy{scope: "admin", permission: "catalog.read"}, a.AdminRanks))                                                           // 段位管理列表
	r.POST("/admin/games/{game_id}/ranks", a.protect(policy{scope: "admin", permission: "catalog.write", idempotent: true}, a.AdminCreateRank))                                  // 新增段位
	r.PATCH("/admin/game-ranks/{rank_id}", a.protect(policy{scope: "admin", permission: "catalog.write", idempotent: true}, a.AdminUpdateRank))                                  // 修改段位及状态
	r.GET("/admin/service-categories", a.protect(policy{scope: "admin", permission: "catalog.read"}, a.AdminCategories))                                                         // 分类和价格边界
	r.POST("/admin/service-categories", a.protect(policy{scope: "admin", permission: "catalog.write", idempotent: true}, a.AdminCreateCategory))                                 // 新增计价分类
	r.PATCH("/admin/service-categories/{category_id}", a.protect(policy{scope: "admin", permission: "catalog.write", idempotent: true}, a.AdminUpdateCategory))                  // 修改价格边界和状态，禁止改变已售计价语义
	r.GET("/admin/coach-services", a.protect(policy{scope: "admin", permission: "service.read"}, a.AdminServices))                                                               // 服务及价格异常查询
	r.POST("/admin/coach-services/{service_id}/moderate", a.protect(policy{scope: "admin", permission: "service.moderate", idempotent: true}, a.AdminModerateService))           // 禁售或解除运营禁售，不直接代陪练改价
	r.GET("/admin/banners", a.protect(policy{scope: "admin", permission: "operation.read"}, a.AdminBanners))                                                                     // Banner列表
	r.POST("/admin/banners", a.protect(policy{scope: "admin", permission: "operation.write", idempotent: true}, a.AdminCreateBanner))                                            // 新增运营位
	r.PATCH("/admin/banners/{banner_id}", a.protect(policy{scope: "admin", permission: "operation.write", idempotent: true}, a.AdminUpdateBanner))                               // 编辑、排序、上下线
	r.GET("/admin/activities", a.protect(policy{scope: "admin", permission: "operation.read"}, a.AdminActivities))                                                               // 展示活动列表
	r.POST("/admin/activities", a.protect(policy{scope: "admin", permission: "operation.write", idempotent: true}, a.AdminCreateActivity))                                       // 新增展示活动
	r.PATCH("/admin/activities/{activity_id}", a.protect(policy{scope: "admin", permission: "operation.write", idempotent: true}, a.AdminUpdateActivity))                        // 编辑活动并更新内容版本
	r.GET("/admin/coach-recommendations", a.protect(policy{scope: "admin", permission: "operation.read"}, a.AdminRecommendations))                                               // 陪练推荐配置
	r.POST("/admin/coach-recommendations", a.protect(policy{scope: "admin", permission: "operation.write", idempotent: true}, a.AdminCreateRecommendation))                      // 新增定时推荐
	r.PATCH("/admin/coach-recommendations/{recommendation_id}", a.protect(policy{scope: "admin", permission: "operation.write", idempotent: true}, a.AdminUpdateRecommendation)) // 编辑推荐及上下线
	r.GET("/admin/content-audits", a.protect(policy{scope: "admin", permission: "content.read"}, a.AdminContentAudits))                                                          // 内容审核队列
	r.POST("/admin/content-audits/{audit_id}/decision", a.protect(policy{scope: "admin", permission: "content.review", idempotent: true}, a.AdminDecideContent))                 // 审核指定版本内容，过期结果不可覆盖

	// 交易与售后
	r.GET("/admin/orders", a.protect(policy{scope: "admin", permission: "order.read"}, a.AdminOrders))                                                   // 按授权数据范围查询订单
	r.GET("/admin/orders/{order_id}", a.protect(policy{scope: "admin", permission: "order.read"}, a.AdminOrderDetail))                                   // 订单履约和资金详情
	r.GET("/admin/orders/{order_id}/events", a.protect(policy{scope: "admin", permission: "order.read"}, a.AdminOrderEvents))                            // 订单完整时间轴
	r.GET("/admin/refunds", a.protect(policy{scope: "admin", permission: "refund.read"}, a.AdminRefunds))                                                // 退款审核及处理列表
	r.GET("/admin/refunds/{refund_id}", a.protect(policy{scope: "admin", permission: "refund.read"}, a.AdminRefundDetail))                               // 退款详情与额度预留
	r.POST("/admin/refunds/{refund_id}/decision", a.protect(policy{scope: "admin", permission: "refund.review", idempotent: true}, a.AdminDecideRefund)) // 带批准和版本执行退款决定
	r.GET("/admin/cases", a.protect(policy{scope: "admin", permission: "case.read"}, a.AdminCases))                                                      // 投诉、客服与评价申诉工单
	r.GET("/admin/cases/{case_id}", a.protect(policy{scope: "admin", permission: "case.read"}, a.AdminCaseDetail))                                       // 工单证据、内部备注及流转
	r.POST("/admin/cases/{case_id}/assign", a.protect(policy{scope: "admin", permission: "case.assign", idempotent: true}, a.AdminAssignCase))           // 分配处理人
	r.POST("/admin/cases/{case_id}/messages", a.protect(policy{scope: "admin", permission: "case.handle", idempotent: true}, a.AdminCaseMessage))        // 回复用户或添加内部备注
	r.POST("/admin/cases/{case_id}/resolve", a.protect(policy{scope: "admin", permission: "case.handle", idempotent: true}, a.AdminResolveCase))         // 记录结论并重算全部阻断工单
	r.POST("/admin/cases/{case_id}/reopen", a.protect(policy{scope: "admin", permission: "case.handle", idempotent: true}, a.AdminReopenCase))           // 重新调查，重新评估结算阻断

	// 财务与统计
	r.GET("/admin/dashboard", a.protect(policy{scope: "admin", permission: "statistics.read"}, a.AdminDashboard))                                                                 // 用户、订单、GMV、佣金、退款与活跃统计
	r.GET("/admin/finance/payments", a.protect(policy{scope: "admin", permission: "finance.read"}, a.AdminPayments))                                                              // 真实支付记录
	r.GET("/admin/finance/settlement-items", a.protect(policy{scope: "admin", permission: "finance.read"}, a.AdminSettlementItems))                                               // 应结、冻结、已付与追偿明细
	r.GET("/admin/finance/settlement-batches", a.protect(policy{scope: "admin", permission: "finance.read"}, a.AdminBatches))                                                     // 出款批次及逐条结果
	r.POST("/admin/finance/settlement-batches", a.protect(policy{scope: "admin", permission: "settlement.submit", idempotent: true}, a.AdminCreateBatch))                         // 复核后创建出款批次
	r.GET("/admin/finance/adjustments", a.protect(policy{scope: "admin", permission: "finance.read"}, a.AdminAdjustments))                                                        // 退款、费用、修正及追偿记录
	r.POST("/admin/finance/adjustments", a.protect(policy{scope: "admin", permission: "settlement.adjust", idempotent: true}, a.AdminCreateAdjustment))                           // 复核后追加资金修正，不覆盖历史
	r.GET("/admin/finance/supplements", a.protect(policy{scope: "admin", permission: "finance.read"}, a.AdminSupplements))                                                        // 独立补款指令列表
	r.POST("/admin/finance/supplements", a.protect(policy{scope: "admin", permission: "settlement.supplement", idempotent: true}, a.AdminCreateSupplement))                       // 首次已付后按已批准正向调整补款
	r.POST("/admin/finance/supplements/{supplement_id}/retry", a.protect(policy{scope: "admin", permission: "settlement.supplement", idempotent: true}, a.AdminRetrySupplement))  // 查明失败后沿原请求号重试
	r.POST("/admin/finance/reconciliation-runs", a.protect(policy{scope: "admin", permission: "reconciliation.run", idempotent: true}, a.AdminCreateReconciliation))              // 按渠道账单日及类型创建对账任务
	r.GET("/admin/finance/reconciliation-runs", a.protect(policy{scope: "admin", permission: "finance.read"}, a.AdminReconciliationRuns))                                         // 对账任务历史与重跑版本
	r.GET("/admin/finance/reconciliation-runs/{run_id}", a.protect(policy{scope: "admin", permission: "finance.read"}, a.AdminReconciliationDetail))                              // 对账结果和差异汇总
	r.GET("/admin/finance/reconciliation-diffs", a.protect(policy{scope: "admin", permission: "finance.read"}, a.AdminReconciliationDiffs))                                       // 对账差异队列
	r.POST("/admin/finance/reconciliation-diffs/{diff_id}/resolve", a.protect(policy{scope: "admin", permission: "reconciliation.handle", idempotent: true}, a.AdminResolveDiff)) // 记录核验依据，不能直接改支付成功
	r.GET("/admin/commission-rules", a.protect(policy{scope: "admin", permission: "commission.read"}, a.AdminCommissionRules))                                                    // 佣金规则版本
	r.POST("/admin/commission-rules", a.protect(policy{scope: "admin", permission: "commission.write", idempotent: true}, a.AdminCreateCommission))                               // 新建佣金草稿
	r.POST("/admin/commission-rules/{rule_id}/publish", a.protect(policy{scope: "admin", permission: "commission.publish", idempotent: true}, a.AdminPublishCommission))          // 批准后发布规则，不影响存量单
	r.GET("/admin/business-rules", a.protect(policy{scope: "admin", permission: "rule.read"}, a.AdminBusinessRules))                                                              // 超时、退款、结算、客服规则版本
	r.POST("/admin/business-rules", a.protect(policy{scope: "admin", permission: "rule.write", idempotent: true}, a.AdminCreateBusinessRule))                                     // 创建规则草稿
	r.POST("/admin/business-rules/{rule_id}/publish", a.protect(policy{scope: "admin", permission: "rule.publish", idempotent: true}, a.AdminPublishBusinessRule))                // 批准后发布业务规则

	// 运营文件
	r.POST("/admin/media/uploads", a.protect(policy{scope: "admin", permission: "media.write", idempotent: true}, a.AdminPrepareUpload))                      // 申请后台运营图片上传凭证
	r.POST("/admin/media/uploads/{asset_id}/complete", a.protect(policy{scope: "admin", permission: "media.write", idempotent: true}, a.AdminCompleteUpload)) // 校验后台上传完成并进入审核
	r.GET("/admin/media/{asset_id}/access", a.protect(policy{scope: "admin", permission: "media.read"}, a.AdminMediaAccess))                                  // 按后台权限获取证据或运营资源访问地址

	// ==================== 渠道回调 ====================
	// 渠道通知
	r.POST("/callbacks/wechat/pay", a.protect(policy{scope: "channel"}, a.CallbackPayment))               // 微信支付结果
	r.POST("/callbacks/wechat/refund", a.protect(policy{scope: "channel"}, a.CallbackRefund))             // 微信退款结果
	r.POST("/callbacks/settlement/{provider}", a.protect(policy{scope: "channel"}, a.CallbackSettlement)) // 结算及补款渠道结果
	r.POST("/callbacks/identity/{provider}", a.protect(policy{scope: "channel"}, a.CallbackIdentity))     // 实名核验结果
	r.POST("/callbacks/content-audit/{provider}", a.protect(policy{scope: "channel"}, a.CallbackContent)) // 内容审核结果
	return srv, nil
}

// database 仅获取当前请求使用的连接；初始化统一由infra/pg.go负责。
func (a *API) database(ctx context.Context) (*gorm.DB, error) {
	if a.db == nil {
		return nil, errDatabaseUnavailable
	}
	return a.db.WithContext(ctx), nil
}
