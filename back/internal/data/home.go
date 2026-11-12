package data

import (
	"context"
	"dd_web/back/internal/biz"
	"dd_web/back/internal/model"
	"time"
)

type HomeRepo struct {
	data *Data
}

func NewHomeRepo(data *Data) *HomeRepo {
	return &HomeRepo{data: data}
}

var _ biz.HomeRepository = (*HomeRepo)(nil)

func (r *HomeRepo) Games(ctx context.Context, limit int) ([]model.Game, error) {
	var rows []model.Game
	err := r.data.db.WithContext(ctx).
		Select("id", "code", "name").
		Where("status = ?", model.GameStatusEnabled).
		Order("sort_order ASC").
		Order("id ASC").
		Limit(limit).
		Find(&rows).Error

	return rows, err
}
func (r *HomeRepo) Banners(ctx context.Context, now time.Time, limit int) ([]model.Banner, error) {
	var rows []model.Banner
	err := r.data.db.WithContext(ctx).
		Select("id", "title", "target_type", "target_id").
		Where("placement = ?", "HOME").
		Where("status = ?", model.BannerStatusPublished).
		Where("starts_at <= ? AND ends_at > ?", now, now).
		Order("sort_order ASC").
		Order("id ASC").
		Limit(limit).
		Find(&rows).Error
	return rows, err
}
func (r *HomeRepo) RecommendedCoaches(ctx context.Context, now time.Time, limit int) ([]biz.HomeCoach, error) {
	var rows []biz.HomeCoach

	err := r.data.db.WithContext(ctx).
		Table("coach_recommendations AS r").
		Select("c.id, c.display_name, c.accepting_orders").
		Joins("JOIN coaches AS c ON c.id = r.coach_id").
		Joins("JOIN users AS u ON u.id = c.user_id").
		Where("r.placement = ?", "HOME").
		Where("r.game_id IS NULL").
		Where("r.status = ?", model.CoachRecommendationStatusEnabled).
		Where("r.starts_at <= ? AND r.ends_at > ?", now, now).
		Where("c.status = ?", model.CoachStatusActive).
		Where("u.status = ?", model.UserStatusActive).
		Group("c.id, c.display_name, c.accepting_orders").
		Order("MIN(r.sort_order) ASC").
		Order("c.id ASC").
		Limit(limit).
		Scan(&rows).Error

	return rows, err
}
func (r *HomeRepo) Activities(ctx context.Context, now time.Time, limit int) ([]model.Activity, error) {
	var rows []model.Activity
	err := r.data.db.WithContext(ctx).
		Select("id", "title").
		Where("status = ?", model.ActivityStatusPublished).
		Where("starts_at <= ? AND ends_at > ?", now, now).
		Order("starts_at DESC").
		Order("id DESC").
		Limit(limit).
		Find(&rows).Error
	return rows, err
}
