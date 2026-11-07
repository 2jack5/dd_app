package model

// Game 游戏目录。
type Game struct {
	Base
	Code        string     `gorm:"column:code;type:varchar(48);not null" json:"code"`                   // 稳定业务编码
	Name        string     `gorm:"column:name;type:varchar(64);not null" json:"name"`                   // 游戏名
	IconAssetID *int64     `gorm:"column:icon_asset_id;type:bigint" json:"icon_asset_id,string"`        // →media_assets
	Keywords    string     `gorm:"column:keywords;type:text;not null;default:''" json:"keywords"`       // 默认空；搜索别名
	Status      GameStatus `gorm:"column:status;type:varchar(16);not null" json:"status"`               // ENABLED/DISABLED
	SortOrder   int32      `gorm:"column:sort_order;type:integer;not null;default:0" json:"sort_order"` // 默认 0
}

func (Game) TableName() string { return "games" }

// GameRegion 游戏区服。
type GameRegion struct {
	Base
	GameID    int64            `gorm:"column:game_id;type:bigint;not null" json:"game_id,string"`           // →games
	Code      string           `gorm:"column:code;type:varchar(48);not null" json:"code"`                   // 区服编码
	Name      string           `gorm:"column:name;type:varchar(64);not null" json:"name"`                   // 展示名
	Status    GameRegionStatus `gorm:"column:status;type:varchar(16);not null" json:"status"`               // ENABLED/DISABLED
	SortOrder int32            `gorm:"column:sort_order;type:integer;not null;default:0" json:"sort_order"` // 默认 0
}

func (GameRegion) TableName() string { return "game_regions" }

// GameRank 游戏段位字典。
type GameRank struct {
	Base
	GameID    int64          `gorm:"column:game_id;type:bigint;not null" json:"game_id,string"` // →games
	Code      string         `gorm:"column:code;type:varchar(48);not null" json:"code"`         // 段位编码
	Name      string         `gorm:"column:name;type:varchar(64);not null" json:"name"`         // 如王者/钻石
	RankLevel int32          `gorm:"column:rank_level;type:integer;not null" json:"rank_level"` // 排序等级
	Status    GameRankStatus `gorm:"column:status;type:varchar(16);not null" json:"status"`     // ENABLED/DISABLED
}

func (GameRank) TableName() string { return "game_ranks" }

// ServiceCategory 服务分类与价格边界。
type ServiceCategory struct {
	Base
	GameID         int64                   `gorm:"column:game_id;type:bigint;not null" json:"game_id,string"`           // →games
	Code           string                  `gorm:"column:code;type:varchar(48);not null" json:"code"`                   // 分类编码
	Name           string                  `gorm:"column:name;type:varchar(64);not null" json:"name"`                   // 服务名
	UnitType       ServiceCategoryUnitType `gorm:"column:unit_type;type:varchar(16);not null" json:"unit_type"`         // HOUR/ROUND
	BillingMinutes *int32                  `gorm:"column:billing_minutes;type:integer" json:"billing_minutes"`          // 小时类一个计价单位的分钟数，如 30/60；局类 NULL
	MinPriceCent   int64                   `gorm:"column:min_price_cent;type:bigint;not null" json:"min_price_cent"`    // 最低单价，>0
	MaxPriceCent   int64                   `gorm:"column:max_price_cent;type:bigint;not null" json:"max_price_cent"`    // 最高单价，>=最低单价
	Status         ServiceCategoryStatus   `gorm:"column:status;type:varchar(16);not null" json:"status"`               // ENABLED/DISABLED
	SortOrder      int32                   `gorm:"column:sort_order;type:integer;not null;default:0" json:"sort_order"` // 默认 0
}

func (ServiceCategory) TableName() string { return "service_categories" }
