package infra

import (
	"context"
	"dd_web/back/internal/config"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// NewPG 初始化 GORM PostgreSQL 连接。
func NewPG(ctx context.Context, cfg *config.PG) (*gorm.DB, error) {
	logLevel := gormlogger.Warn
	if cfg == nil {
		return nil, fmt.Errorf("pg config is nil")
	}
	if cfg.DSN == "" {
		return nil, fmt.Errorf("pg dsn is empty")
	}

	db, err := gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{
		// 关联字段暂未在模型中声明完整约束，避免 GORM 推断出错误外键。
		DisableForeignKeyConstraintWhenMigrating: true,
		DisableAutomaticPing:                     true,
		TranslateError:                           true,
		Logger: gormlogger.New(
			gormWriter{},
			gormlogger.Config{
				SlowThreshold:             200 * time.Millisecond,
				LogLevel:                  logLevel,
				IgnoreRecordNotFoundError: true,
				ParameterizedQueries:      true,
			},
		),
	})
	if err != nil {
		return nil, fmt.Errorf("open postgres: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql.DB: %w", err)
	}
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	// 启动时探活
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(pingCtx); err != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	L().Info("postgres connected", TraceID("infra"))
	return db, nil
}

// gormWriter 使用统一slog实例；SQL参数已由ParameterizedQueries屏蔽。
type gormWriter struct{}

func (gormWriter) Printf(format string, args ...any) {
	L().Warn(fmt.Sprintf("[gorm] "+format, args...))
}

// ClosePG 在Kratos退出后释放连接池。
func ClosePG(db *gorm.DB) error {
	if db == nil {
		return nil
	}
	pool, err := db.DB()
	if err != nil {
		return err
	}
	return pool.Close()
}
