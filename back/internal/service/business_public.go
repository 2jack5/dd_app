package service

import (
	"context"
	businessv1 "dd_web/back/api/business/v1"
	commonv1 "dd_web/back/api/common/v1"
	"dd_web/back/internal/biz"
	kerrors "github.com/go-kratos/kratos/v3/errors"
)

func (s *BusinessService) SystemHealth(context.Context, *businessv1.SystemHealthRequest) (*commonv1.ApiResponse, error) {
	return ok(map[string]any{"status": "alive"})
}
func (s *BusinessService) SystemReady(ctx context.Context, _ *businessv1.SystemReadyRequest) (*commonv1.ApiResponse, error) {
	if err := s.system.Ready(ctx); err != nil {
		return nil, kerrors.ServiceUnavailable("NOT_READY", "数据库连接未就绪")
	}
	return ok(map[string]any{"status": "ready"})
}
func (s *BusinessService) AuthWechatLogin(context.Context, *businessv1.AuthWechatLoginRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("auth.wechatLogin")
}
func (s *BusinessService) AuthRefreshToken(context.Context, *businessv1.AuthRefreshTokenRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("auth.refreshToken")
}
func (s *BusinessService) AgreementCurrent(context.Context, *businessv1.AgreementCurrentRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("agreement.current")
}
func (s *BusinessService) HomeOverview(context.Context, *businessv1.HomeOverviewRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("home.overview")
}
func (s *BusinessService) HomeBanners(context.Context, *businessv1.HomeBannersRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("home.banners")
}
func (s *BusinessService) HomeRecommendations(context.Context, *businessv1.HomeRecommendationsRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("home.recommendations")
}
func (s *BusinessService) ActivityList(context.Context, *businessv1.ActivityListRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("activity.list")
}
func (s *BusinessService) ActivityDetail(context.Context, *businessv1.ActivityDetailRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("activity.detail")
}
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
func (s *BusinessService) SearchSuggestions(context.Context, *businessv1.SearchSuggestionsRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("search.suggestions")
}
func (s *BusinessService) DiscoveryCoaches(context.Context, *businessv1.DiscoveryCoachesRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("discovery.coaches")
}
func (s *BusinessService) DiscoveryCoachDetail(context.Context, *businessv1.DiscoveryCoachDetailRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("discovery.coachDetail")
}
func (s *BusinessService) DiscoveryServices(context.Context, *businessv1.DiscoveryServicesRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("discovery.services")
}
func (s *BusinessService) DiscoveryAvailability(context.Context, *businessv1.DiscoveryAvailabilityRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("discovery.availability")
}
func (s *BusinessService) DiscoveryReviews(context.Context, *businessv1.DiscoveryReviewsRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("discovery.reviews")
}
func (s *BusinessService) DiscoveryServiceDetail(context.Context, *businessv1.DiscoveryServiceDetailRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("discovery.serviceDetail")
}
func (s *BusinessService) SupportConfig(context.Context, *businessv1.SupportConfigRequest) (*commonv1.ApiResponse, error) {
	return nil, notImplemented("support.config")
}
