package biz

import (
	"context"
	"errors"
	"time"

	"dd_web/back/internal/model"
)

// HomeCoach 只包含首页推荐需要展示的陪练信息。
type HomeCoach struct {
	ID              int64
	DisplayName     string
	AcceptingOrders bool
}

type HomeRepository interface {
	Games(context.Context, int) ([]model.Game, error)
	Banners(context.Context, time.Time, int) ([]model.Banner, error)

	RecommendedCoaches(context.Context, time.Time, int) ([]HomeCoach, error)

	Activities(context.Context, time.Time, int) ([]model.Activity, error)
}

type HomeUsecase struct {
	repo HomeRepository
}

func NewHomeUsecase(repo HomeRepository) *HomeUsecase {
	return &HomeUsecase{repo: repo}
}

// HomeBannerPage 保存 Banner 的分页参数。
type HomeBannerPage struct {
	Limit          int
	AfterSortOrder int32
	AfterID        int64
}

// HomeBannerList 保存本页结果。
type HomeBannerList struct {
	Items      []any
	NextCursor string
}

// Overview 获取首页的组合数据。
func (u *HomeUsecase) Overview(ctx context.Context) (map[string]any, error) {
	if u.repo == nil {
		return nil, errors.New("home repository is not configured")
	}
	// 各板块使用同一个时间判断是否处于展示期间。
	now := time.Now().UTC()
	// 这些数量是本例的展示规则，可以以后按产品需求调整。
	gameRows, err := u.repo.Games(ctx, 12)
	if err != nil {
		return nil, err
	}
	bannerRows, err := u.repo.Banners(ctx, now, 5)
	if err != nil {
		return nil, err
	}
	coachRows, err := u.repo.RecommendedCoaches(ctx, now, 8)
	if err != nil {
		return nil, err
	}
	activityRows, err := u.repo.Activities(ctx, now, 4)
	if err != nil {
		return nil, err
	}
	// 只挑选允许返回给前端的字段。
	games := make([]any, 0, len(gameRows))
	for _, game := range gameRows {
		games = append(games, map[string]any{
			"id":   id(game.ID),
			"code": game.Code,
			"name": game.Name,
		})
	}
	banners := make([]any, 0, len(bannerRows))
	for _, banner := range bannerRows {
		banners = append(banners, map[string]any{
			"id":          id(banner.ID),
			"title":       banner.Title,
			"target_type": string(banner.TargetType),
			"target_id":   id(banner.TargetID),
		})
	}

	recommendations := make([]any, 0, len(coachRows))
	for _, coach := range coachRows {
		recommendations = append(recommendations, map[string]any{
			"id":               id(coach.ID),
			"display_name":     coach.DisplayName,
			"accepting_orders": coach.AcceptingOrders,
		})
	}

	activities := make([]any, 0, len(activityRows))
	for _, activity := range activityRows {
		activities = append(activities, map[string]any{
			"id":    id(activity.ID),
			"title": activity.Title,
		})
	}
	return map[string]any{
		"games":           games,
		"banners":         banners,
		"recommendations": recommendations,
		"activities":      activities,
	}, nil
}
