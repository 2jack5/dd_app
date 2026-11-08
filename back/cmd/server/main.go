package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"dd_web/back/internal/api"
	"dd_web/back/internal/config"
	"dd_web/back/internal/infra"
	"dd_web/back/internal/model"

	"github.com/go-kratos/kratos/v3"
	"gorm.io/gorm"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "server exited: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// 1.配置与日志
	cfg, err := config.Load("")
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}
	infra.Init(cfg.Log.Level, cfg.Log.Format)

	// 2.基础设施
	pgCfg, err := postgresConfig(cfg.PG)
	if err != nil {
		return err
	}
	var pg *gorm.DB
	// 配置留空时先运行接口框架；填写连接串后在这里初始化PostgreSQL。
	if strings.TrimSpace(pgCfg.DSN) != "" {
		pg, err = infra.NewPG(context.Background(), &pgCfg)
		if err != nil {
			return fmt.Errorf("initialize postgres: %w", err)
		}
		if err := pg.AutoMigrate(model.AllModels()...); err != nil {
			return fmt.Errorf("auto migrate postgres: %w", err)
		}
		// 路由创建失败、服务报错或正常停止，都会执行连接池关闭。
		defer func() {
			if err := infra.ClosePG(pg); err != nil {
				infra.L().Error("close postgres failed", "error", err)
			}
		}()
	}

	// 3.路由与服务
	httpServer, err := api.NewRouter(cfg, pg)
	if err != nil {
		return err
	}
	opts := []kratos.Option{kratos.Name("dongdong-esports-api"), kratos.Version("0.2.0"), kratos.Logger(infra.L()), kratos.Server(httpServer)}
	if cfg.HTTP.ShutdownTimeout > 0 {
		opts = append(opts, kratos.StopTimeout(cfg.HTTP.ShutdownTimeout))
	}
	return kratos.New(opts...).Run()
}

// postgresConfig 保留原有环境覆盖和空配置默认值，路由不再处理数据库配置。
func postgresConfig(c config.PG) (config.PG, error) {
	if v, ok := os.LookupEnv("DD_DATABASE_URL"); ok {
		c.DSN = v
	}
	if c.MaxOpenConns == 0 {
		c.MaxOpenConns = 20
	}
	if c.ConnMaxLifetime == 0 {
		c.ConnMaxLifetime = 30 * time.Minute
	}
	for key, p := range map[string]*int{"DD_DB_MAX_OPEN": &c.MaxOpenConns, "DD_DB_MAX_IDLE": &c.MaxIdleConns} {
		if v := os.Getenv(key); v != "" {
			n, err := strconv.Atoi(v)
			if err != nil || n < 0 {
				return c, fmt.Errorf("%s must be nonnegative", key)
			}
			*p = n
		}
	}
	if c.MaxOpenConns < 1 || c.MaxIdleConns < 0 || c.MaxIdleConns > c.MaxOpenConns || c.ConnMaxLifetime < 0 {
		return c, fmt.Errorf("invalid PostgreSQL pool limits or lifetime")
	}
	return c, nil
}
