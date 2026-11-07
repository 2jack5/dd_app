package api

import (
	"dd_web/back/internal/model"
	khttp "github.com/go-kratos/kratos/v3/transport/http"
	"gorm.io/gorm"
	"strconv"
)

// AuthWechatLoginRequest 微信App或小程序登录的请求契约；字段校验与业务用例尚未实现。
type AuthWechatLoginRequest struct {
	ClientType string `json:"client_type"` // 必填。 取值：WECHAT_APP/WECHAT_MINI
	Code       string `json:"code"`        // 必填。
}

// AuthWechatLogin 微信App或小程序登录。
func (a *API) AuthWechatLogin(ctx khttp.Context) error {
	return notImplemented(ctx, "auth.wechatLogin")
}

// AuthRefreshTokenRequest 使用刷新凭据换取新登录态的请求契约；字段校验与业务用例尚未实现。
type AuthRefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"` // 必填。
}

// AuthRefreshToken 使用刷新凭据换取新登录态。
func (a *API) AuthRefreshToken(ctx khttp.Context) error {
	return notImplemented(ctx, "auth.refreshToken")
}

// AgreementCurrent 获取当前有效协议版本。
func (a *API) AgreementCurrent(ctx khttp.Context) error {
	return notImplemented(ctx, "agreement.current")
}

// HomeOverview 首页组合数据：Banner、推荐、活动、游戏。
func (a *API) HomeOverview(ctx khttp.Context) error {
	return notImplemented(ctx, "home.overview")
}

// HomeBanners 首页Banner列表。
func (a *API) HomeBanners(ctx khttp.Context) error {
	return notImplemented(ctx, "home.banners")
}

// HomeRecommendations 推荐陪练列表。
func (a *API) HomeRecommendations(ctx khttp.Context) error {
	return notImplemented(ctx, "home.recommendations")
}

// ActivityList 展示活动列表。
func (a *API) ActivityList(ctx khttp.Context) error {
	return notImplemented(ctx, "activity.list")
}

// ActivityDetail 展示活动详情。
func (a *API) ActivityDetail(ctx khttp.Context) error {
	return notImplemented(ctx, "activity.detail")
}

// SearchSuggestions 游戏、陪练和服务关键词联想。
func (a *API) SearchSuggestions(ctx khttp.Context) error {
	return notImplemented(ctx, "search.suggestions")
}

// DiscoveryCoaches 按游戏、区服、价格、段位、服务方式及时间筛选陪练。
func (a *API) DiscoveryCoaches(ctx khttp.Context) error {
	return notImplemented(ctx, "discovery.coaches")
}

// DiscoveryCoachDetail 陪练详情及认证标识。
func (a *API) DiscoveryCoachDetail(ctx khttp.Context) error {
	return notImplemented(ctx, "discovery.coachDetail")
}

// DiscoveryServices 陪练可售服务。
func (a *API) DiscoveryServices(ctx khttp.Context) error {
	return notImplemented(ctx, "discovery.services")
}

// DiscoveryAvailability 陪练可预约时段。
func (a *API) DiscoveryAvailability(ctx khttp.Context) error {
	return notImplemented(ctx, "discovery.availability")
}

// DiscoveryReviews 已公开评价。
func (a *API) DiscoveryReviews(ctx khttp.Context) error {
	return notImplemented(ctx, "discovery.reviews")
}

// DiscoveryServiceDetail 可售服务详情。
func (a *API) DiscoveryServiceDetail(ctx khttp.Context) error {
	return notImplemented(ctx, "discovery.serviceDetail")
}

// SupportConfig 客服入口、联系方式与帮助配置。
func (a *API) SupportConfig(ctx khttp.Context) error {
	return notImplemented(ctx, "support.config")
}

// 以下为已接通GORM的基础查询。其余本文件处理函数保持明确的501占位。
func (a *API) SystemHealth(ctx khttp.Context) error {
	return reply(ctx, 200, "OK", "", map[string]string{"status": "alive"})
}
func (a *API) SystemReady(ctx khttp.Context) error {
	db, err := a.database(ctx)
	if err != nil {
		return reply(ctx, 503, "NOT_READY", "数据库连接未就绪", nil)
	}
	pool, err := db.DB()
	if err != nil || pool.PingContext(ctx) != nil {
		return reply(ctx, 503, "NOT_READY", "数据库连接未就绪", nil)
	}
	return reply(ctx, 200, "OK", "", map[string]string{"status": "ready"})
}

type catalogPage struct {
	afterID int64
	limit   int
}
type catalogList struct {
	Items      any    `json:"items"`
	NextCursor string `json:"next_cursor,omitempty"`
}
type gameItem struct {
	ID          string  `json:"id"`
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	IconAssetID *string `json:"icon_asset_id"`
}
type regionItem struct {
	ID     string `json:"id"`
	GameID string `json:"game_id"`
	Code   string `json:"code"`
	Name   string `json:"name"`
}
type rankItem struct {
	ID        string `json:"id"`
	GameID    string `json:"game_id"`
	Code      string `json:"code"`
	Name      string `json:"name"`
	RankLevel int32  `json:"rank_level"`
}
type categoryItem struct {
	ID             string `json:"id"`
	GameID         string `json:"game_id"`
	Code           string `json:"code"`
	Name           string `json:"name"`
	UnitType       string `json:"unit_type"`
	BillingMinutes *int32 `json:"billing_minutes"`
	MinPriceCent   int64  `json:"min_price_cent"`
	MaxPriceCent   int64  `json:"max_price_cent"`
}

func parseCatalogPage(ctx khttp.Context) (catalogPage, error) {
	p := catalogPage{limit: 20}
	var err error
	if v := ctx.Query().Get("page_size"); v != "" {
		p.limit, err = strconv.Atoi(v)
		if err != nil || p.limit < 1 || p.limit > 100 {
			return p, errInvalidArgument
		}
	}
	if v := ctx.Query().Get("cursor"); v != "" {
		p.afterID, err = strconv.ParseInt(v, 10, 64)
		if err != nil || p.afterID < 0 {
			return p, errInvalidArgument
		}
	}
	return p, nil
}
func parsePositiveID(s string) (int64, error) {
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil || n < 1 {
		return 0, errInvalidArgument
	}
	return n, nil
}
func formatID(n int64) string { return strconv.FormatInt(n, 10) }
func catalogQuery(db *gorm.DB, p catalogPage) *gorm.DB {
	return db.Where("id > ?", p.afterID).Order("id ASC").Limit(p.limit + 1)
}
func trimCatalog[T any](rows []T, p catalogPage, getID func(T) int64) ([]T, string) {
	if len(rows) > p.limit {
		return rows[:p.limit], formatID(getID(rows[p.limit-1]))
	}
	return rows, ""
}
func enabledGame(db *gorm.DB, id int64) error {
	var n int64
	if err := db.Model(&model.Game{}).Where("id = ? AND status = ?", id, model.GameStatusEnabled).Count(&n).Error; err != nil {
		return err
	}
	if n == 0 {
		return errNotFound
	}
	return nil
}
func (a *API) CatalogGames(ctx khttp.Context) error {
	p, err := parseCatalogPage(ctx)
	if err != nil {
		return dataError(ctx, err)
	}
	db, err := a.database(ctx)
	if err != nil {
		return dataError(ctx, err)
	}
	var rows []model.Game
	if err = catalogQuery(db, p).Select("id", "code", "name", "icon_asset_id").Where("status = ?", model.GameStatusEnabled).Find(&rows).Error; err != nil {
		return dataError(ctx, err)
	}
	rows, next := trimCatalog(rows, p, func(v model.Game) int64 { return v.ID })
	out := make([]gameItem, 0, len(rows))
	for _, v := range rows {
		var asset *string
		if v.IconAssetID != nil {
			x := formatID(*v.IconAssetID)
			asset = &x
		}
		out = append(out, gameItem{formatID(v.ID), v.Code, v.Name, asset})
	}
	return reply(ctx, 200, "OK", "", catalogList{out, next})
}
func (a *API) CatalogRegions(ctx khttp.Context) error {
	id, err := parsePositiveID(ctx.Vars().Get("game_id"))
	if err != nil {
		return dataError(ctx, err)
	}
	p, err := parseCatalogPage(ctx)
	if err != nil {
		return dataError(ctx, err)
	}
	db, err := a.database(ctx)
	if err != nil {
		return dataError(ctx, err)
	}
	if err = enabledGame(db, id); err != nil {
		return dataError(ctx, err)
	}
	var rows []model.GameRegion
	if err = catalogQuery(db, p).Select("id", "game_id", "code", "name").Where("game_id = ? AND status = ?", id, model.GameRegionStatusEnabled).Find(&rows).Error; err != nil {
		return dataError(ctx, err)
	}
	rows, next := trimCatalog(rows, p, func(v model.GameRegion) int64 { return v.ID })
	out := make([]regionItem, 0, len(rows))
	for _, v := range rows {
		out = append(out, regionItem{formatID(v.ID), formatID(v.GameID), v.Code, v.Name})
	}
	return reply(ctx, 200, "OK", "", catalogList{out, next})
}
func (a *API) CatalogRanks(ctx khttp.Context) error {
	id, err := parsePositiveID(ctx.Vars().Get("game_id"))
	if err != nil {
		return dataError(ctx, err)
	}
	p, err := parseCatalogPage(ctx)
	if err != nil {
		return dataError(ctx, err)
	}
	db, err := a.database(ctx)
	if err != nil {
		return dataError(ctx, err)
	}
	if err = enabledGame(db, id); err != nil {
		return dataError(ctx, err)
	}
	var rows []model.GameRank
	if err = catalogQuery(db, p).Select("id", "game_id", "code", "name", "rank_level").Where("game_id = ? AND status = ?", id, model.GameRankStatusEnabled).Find(&rows).Error; err != nil {
		return dataError(ctx, err)
	}
	rows, next := trimCatalog(rows, p, func(v model.GameRank) int64 { return v.ID })
	out := make([]rankItem, 0, len(rows))
	for _, v := range rows {
		out = append(out, rankItem{formatID(v.ID), formatID(v.GameID), v.Code, v.Name, v.RankLevel})
	}
	return reply(ctx, 200, "OK", "", catalogList{out, next})
}
func (a *API) CatalogCategories(ctx khttp.Context) error {
	id, err := parsePositiveID(ctx.Query().Get("game_id"))
	if err != nil {
		return dataError(ctx, err)
	}
	p, err := parseCatalogPage(ctx)
	if err != nil {
		return dataError(ctx, err)
	}
	db, err := a.database(ctx)
	if err != nil {
		return dataError(ctx, err)
	}
	if err = enabledGame(db, id); err != nil {
		return dataError(ctx, err)
	}
	var rows []model.ServiceCategory
	if err = catalogQuery(db, p).Select("id", "game_id", "code", "name", "unit_type", "billing_minutes", "min_price_cent", "max_price_cent").Where("game_id = ? AND status = ?", id, model.ServiceCategoryStatusEnabled).Find(&rows).Error; err != nil {
		return dataError(ctx, err)
	}
	rows, next := trimCatalog(rows, p, func(v model.ServiceCategory) int64 { return v.ID })
	out := make([]categoryItem, 0, len(rows))
	for _, v := range rows {
		out = append(out, categoryItem{formatID(v.ID), formatID(v.GameID), v.Code, v.Name, string(v.UnitType), v.BillingMinutes, v.MinPriceCent, v.MaxPriceCent})
	}
	return reply(ctx, 200, "OK", "", catalogList{out, next})
}
