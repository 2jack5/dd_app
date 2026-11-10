package service

import (
	"context"
	businessv1 "dd_web/back/api/business/v1"
	commonv1 "dd_web/back/api/common/v1"
	"dd_web/back/internal/biz"

	kerrors "github.com/go-kratos/kratos/v3/errors"
)

// SystemHealth 检查业务服务进程是否存活。
func (s *BusinessService) SystemHealth(context.Context, *businessv1.SystemHealthRequest) (*commonv1.ApiResponse, error) {
	return ok(map[string]any{"status": "alive"})
}

// SystemReady 检查业务服务及数据库连接是否就绪。
func (s *BusinessService) SystemReady(ctx context.Context, _ *businessv1.SystemReadyRequest) (*commonv1.ApiResponse, error) {
	if err := s.system.Ready(ctx); err != nil {
		return nil, kerrors.ServiceUnavailable("NOT_READY", "数据库连接未就绪")
	}
	return ok(map[string]any{"status": "ready"})
}

// AuthWechatLogin 处理普通用户的微信登录。
func (s *BusinessService) AuthWechatLogin(context.Context, *businessv1.AuthWechatLoginRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("auth.wechatLogin")
}

// AuthRefreshToken 刷新普通用户的登录凭证。
func (s *BusinessService) AuthRefreshToken(context.Context, *businessv1.AuthRefreshTokenRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("auth.refreshToken")
}

// AgreementCurrent 获取指定协议的当前版本。
func (s *BusinessService) AgreementCurrent(context.Context, *businessv1.AgreementCurrentRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("agreement.current")
}

// HomeOverview 获取用户首页概览信息。
func (s *BusinessService) HomeOverview(context.Context, *businessv1.HomeOverviewRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("home.overview")
}

// HomeBanners 获取用户首页展示的运营横幅。
func (s *BusinessService) HomeBanners(context.Context, *businessv1.HomeBannersRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("home.banners")
}

// HomeRecommendations 获取用户首页推荐的陪练师。
func (s *BusinessService) HomeRecommendations(context.Context, *businessv1.HomeRecommendationsRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("home.recommendations")
}

// ActivityList 分页查询用户可查看的活动列表。
func (s *BusinessService) ActivityList(context.Context, *businessv1.ActivityListRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("activity.list")
}

// ActivityDetail 查询指定活动的详细信息。
func (s *BusinessService) ActivityDetail(context.Context, *businessv1.ActivityDetailRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("activity.detail")
}

// CatalogGames 分页查询已启用的游戏列表。
func (s *BusinessService) CatalogGames(ctx context.Context, req *businessv1.CatalogGamesRequest) (*commonv1.ApiResponse, error) {
	p, err := biz.ParsePage(req.Cursor, req.PageSize)
	if err != nil {
		return nil, catalogError(err)
	}
	v, err := s.catalog.Games(ctx, p)
	if err != nil {
		return nil, catalogError(err)
	}
	return ok(map[string]any{"items": v.Items, "next_cursor": v.NextCursor})
}

// CatalogRegions 分页查询指定游戏已启用的区服列表。
func (s *BusinessService) CatalogRegions(ctx context.Context, req *businessv1.CatalogRegionsRequest) (*commonv1.ApiResponse, error) {
	id, err := biz.ParseID(req.GameId)
	if err != nil {
		return nil, catalogError(err)
	}
	p, err := biz.ParsePage(req.Cursor, req.PageSize)
	if err != nil {
		return nil, catalogError(err)
	}
	v, err := s.catalog.Regions(ctx, id, p)
	if err != nil {
		return nil, catalogError(err)
	}
	return ok(map[string]any{"items": v.Items, "next_cursor": v.NextCursor})
}

// CatalogRanks 分页查询指定游戏已启用的段位列表。
func (s *BusinessService) CatalogRanks(ctx context.Context, req *businessv1.CatalogRanksRequest) (*commonv1.ApiResponse, error) {
	id, err := biz.ParseID(req.GameId)
	if err != nil {
		return nil, catalogError(err)
	}
	p, err := biz.ParsePage(req.Cursor, req.PageSize)
	if err != nil {
		return nil, catalogError(err)
	}
	v, err := s.catalog.Ranks(ctx, id, p)
	if err != nil {
		return nil, catalogError(err)
	}
	return ok(map[string]any{"items": v.Items, "next_cursor": v.NextCursor})
}

// CatalogCategories 分页查询指定游戏已启用的服务分类。
func (s *BusinessService) CatalogCategories(ctx context.Context, req *businessv1.CatalogCategoriesRequest) (*commonv1.ApiResponse, error) {
	id, err := biz.ParseID(req.GameId)
	if err != nil {
		return nil, catalogError(err)
	}
	p, err := biz.ParsePage(req.Cursor, req.PageSize)
	if err != nil {
		return nil, catalogError(err)
	}
	v, err := s.catalog.Categories(ctx, id, p)
	if err != nil {
		return nil, catalogError(err)
	}
	return ok(map[string]any{"items": v.Items, "next_cursor": v.NextCursor})
}

// SearchSuggestions 获取搜索建议。
func (s *BusinessService) SearchSuggestions(context.Context, *businessv1.SearchSuggestionsRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("search.suggestions")
}

// DiscoveryCoaches 查询供用户浏览和筛选的陪练师列表。
func (s *BusinessService) DiscoveryCoaches(context.Context, *businessv1.DiscoveryCoachesRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("discovery.coaches")
}

// DiscoveryCoachDetail 查询指定陪练师的公开资料。
func (s *BusinessService) DiscoveryCoachDetail(context.Context, *businessv1.DiscoveryCoachDetailRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("discovery.coachDetail")
}

// DiscoveryServices 查询指定陪练师对外展示的服务列表。
func (s *BusinessService) DiscoveryServices(context.Context, *businessv1.DiscoveryServicesRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("discovery.services")
}

// DiscoveryAvailability 查询指定陪练师的可预约时段。
func (s *BusinessService) DiscoveryAvailability(context.Context, *businessv1.DiscoveryAvailabilityRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("discovery.availability")
}

// DiscoveryReviews 查询指定陪练师的公开评价。
func (s *BusinessService) DiscoveryReviews(context.Context, *businessv1.DiscoveryReviewsRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("discovery.reviews")
}

// DiscoveryServiceDetail 查询指定陪练服务的公开详情。
func (s *BusinessService) DiscoveryServiceDetail(context.Context, *businessv1.DiscoveryServiceDetailRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("discovery.serviceDetail")
}

// SupportConfig 获取客服联系方式及相关配置。
func (s *BusinessService) SupportConfig(context.Context, *businessv1.SupportConfigRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("support.config")
}
