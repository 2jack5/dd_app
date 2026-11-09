package biz

import (
	"context"
	"errors"
	"strconv"

	"dd_web/back/internal/model"
)

var ErrInvalidArgument = errors.New("invalid argument")

type Page struct {
	After int64
	Limit int
}
type List struct {
	Items      []map[string]any
	NextCursor string
}

type CatalogRepository interface {
	Games(context.Context, int64, int) ([]model.Game, error)
	Regions(context.Context, int64, int64, int) ([]model.GameRegion, error)
	Ranks(context.Context, int64, int64, int) ([]model.GameRank, error)
	Categories(context.Context, int64, int64, int) ([]model.ServiceCategory, error)
}

type CatalogUsecase struct{ repo CatalogRepository }

func NewCatalogUsecase(repo CatalogRepository) *CatalogUsecase { return &CatalogUsecase{repo: repo} }

func ParsePage(cursor string, size int32) (Page, error) {
	p := Page{Limit: 20}
	if size != 0 {
		p.Limit = int(size)
	}
	if p.Limit < 1 || p.Limit > 100 {
		return p, ErrInvalidArgument
	}
	if cursor != "" {
		n, err := strconv.ParseInt(cursor, 10, 64)
		if err != nil || n < 0 {
			return p, ErrInvalidArgument
		}
		p.After = n
	}
	return p, nil
}

func ParseID(value string) (int64, error) {
	n, err := strconv.ParseInt(value, 10, 64)
	if err != nil || n < 1 {
		return 0, ErrInvalidArgument
	}
	return n, nil
}

func (u *CatalogUsecase) Games(ctx context.Context, p Page) (List, error) {
	rows, err := u.repo.Games(ctx, p.After, p.Limit)
	if err != nil {
		return List{}, err
	}
	next := trim(&rows, p.Limit, func(v model.Game) int64 { return v.ID })
	items := make([]map[string]any, 0, len(rows))
	for _, v := range rows {
		items = append(items, map[string]any{"id": id(v.ID), "code": v.Code, "name": v.Name, "icon_asset_id": optionalID(v.IconAssetID)})
	}
	return List{items, next}, nil
}

func (u *CatalogUsecase) Regions(ctx context.Context, gameID int64, p Page) (List, error) {
	rows, err := u.repo.Regions(ctx, gameID, p.After, p.Limit)
	if err != nil {
		return List{}, err
	}
	next := trim(&rows, p.Limit, func(v model.GameRegion) int64 { return v.ID })
	items := make([]map[string]any, 0, len(rows))
	for _, v := range rows {
		items = append(items, map[string]any{"id": id(v.ID), "game_id": id(v.GameID), "code": v.Code, "name": v.Name})
	}
	return List{items, next}, nil
}

func (u *CatalogUsecase) Ranks(ctx context.Context, gameID int64, p Page) (List, error) {
	rows, err := u.repo.Ranks(ctx, gameID, p.After, p.Limit)
	if err != nil {
		return List{}, err
	}
	next := trim(&rows, p.Limit, func(v model.GameRank) int64 { return v.ID })
	items := make([]map[string]any, 0, len(rows))
	for _, v := range rows {
		items = append(items, map[string]any{"id": id(v.ID), "game_id": id(v.GameID), "code": v.Code, "name": v.Name, "rank_level": v.RankLevel})
	}
	return List{items, next}, nil
}

func (u *CatalogUsecase) Categories(ctx context.Context, gameID int64, p Page) (List, error) {
	rows, err := u.repo.Categories(ctx, gameID, p.After, p.Limit)
	if err != nil {
		return List{}, err
	}
	next := trim(&rows, p.Limit, func(v model.ServiceCategory) int64 { return v.ID })
	items := make([]map[string]any, 0, len(rows))
	for _, v := range rows {
		items = append(items, map[string]any{"id": id(v.ID), "game_id": id(v.GameID), "code": v.Code, "name": v.Name, "unit_type": string(v.UnitType), "billing_minutes": v.BillingMinutes, "min_price_cent": v.MinPriceCent, "max_price_cent": v.MaxPriceCent})
	}
	return List{items, next}, nil
}

func trim[T any](rows *[]T, limit int, getID func(T) int64) string {
	if len(*rows) <= limit {
		return ""
	}
	*rows = (*rows)[:limit]
	return id(getID((*rows)[limit-1]))
}
func id(v int64) string { return strconv.FormatInt(v, 10) }
func optionalID(v *int64) any {
	if v == nil {
		return nil
	}
	return id(*v)
}
