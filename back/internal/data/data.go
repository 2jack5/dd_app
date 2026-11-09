package data

import (
	"context"
	"errors"

	"dd_web/back/internal/model"
	"gorm.io/gorm"
)

var ErrNotFound = errors.New("data not found")

type Data struct{ db *gorm.DB }

func NewData(db *gorm.DB) *Data { return &Data{db: db} }

func (d *Data) Ping(ctx context.Context) error {
	pool, err := d.db.DB()
	if err != nil {
		return err
	}
	return pool.PingContext(ctx)
}

type CatalogRepo struct{ data *Data }

func NewCatalogRepo(data *Data) *CatalogRepo { return &CatalogRepo{data: data} }

func (r *CatalogRepo) Games(ctx context.Context, after int64, limit int) ([]model.Game, error) {
	var rows []model.Game
	err := r.data.db.WithContext(ctx).Where("id > ? AND status = ?", after, model.GameStatusEnabled).
		Order("id ASC").Limit(limit + 1).Find(&rows).Error
	return rows, err
}

func (r *CatalogRepo) Regions(ctx context.Context, gameID, after int64, limit int) ([]model.GameRegion, error) {
	if err := r.requireGame(ctx, gameID); err != nil {
		return nil, err
	}
	var rows []model.GameRegion
	err := r.data.db.WithContext(ctx).Where("game_id = ? AND id > ? AND status = ?", gameID, after, model.GameRegionStatusEnabled).
		Order("id ASC").Limit(limit + 1).Find(&rows).Error
	return rows, err
}

func (r *CatalogRepo) Ranks(ctx context.Context, gameID, after int64, limit int) ([]model.GameRank, error) {
	if err := r.requireGame(ctx, gameID); err != nil {
		return nil, err
	}
	var rows []model.GameRank
	err := r.data.db.WithContext(ctx).Where("game_id = ? AND id > ? AND status = ?", gameID, after, model.GameRankStatusEnabled).
		Order("id ASC").Limit(limit + 1).Find(&rows).Error
	return rows, err
}

func (r *CatalogRepo) Categories(ctx context.Context, gameID, after int64, limit int) ([]model.ServiceCategory, error) {
	if err := r.requireGame(ctx, gameID); err != nil {
		return nil, err
	}
	var rows []model.ServiceCategory
	err := r.data.db.WithContext(ctx).Where("game_id = ? AND id > ? AND status = ?", gameID, after, model.ServiceCategoryStatusEnabled).
		Order("id ASC").Limit(limit + 1).Find(&rows).Error
	return rows, err
}

func (r *CatalogRepo) requireGame(ctx context.Context, id int64) error {
	var count int64
	err := r.data.db.WithContext(ctx).Model(&model.Game{}).Where("id = ? AND status = ?", id, model.GameStatusEnabled).Count(&count).Error
	if err != nil {
		return err
	}
	if count == 0 {
		return ErrNotFound
	}
	return nil
}
