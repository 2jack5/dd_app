package app

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"dd_web/back/internal/biz"
	"dd_web/back/internal/config"
	"dd_web/back/internal/data"
	"dd_web/back/internal/infra"
	"dd_web/back/internal/model"
	"dd_web/back/internal/server"
	"dd_web/back/internal/service"

	"github.com/go-kratos/kratos/v3"
	khttp "github.com/go-kratos/kratos/v3/transport/http"
)

//func RunBusiness() error {
//	return run("dongdong-business-api", func(cfg *config.Config, dataSet *data.Data) (*khttp.Server, time.Duration, error) {
//		catalog := biz.NewCatalogUsecase(data.NewCatalogRepo(dataSet))
//		svc := service.NewBusinessService(biz.NewSystemUsecase(dataSet), catalog)
//		httpServer, err := server.NewBusinessHTTPServer(cfg, svc, nil)
//		return httpServer, cfg.HTTP.ShutdownTimeout, err
//	})
//}

func RunBusiness() error {
	return run("dongdong-business-api",
		func(cfg *config.Config, dataSet *data.Data) (*khttp.Server, time.Duration, error) {
			catalog := biz.NewCatalogUsecase(data.NewCatalogRepo(dataSet))
			// 第一步：创建管理员数据库访问对象。
			adminRepo := data.NewAdminRepo(dataSet)
			// 第二步：把数据库访问对象交给业务用例。
			adminUsecase := biz.NewAdminUsecase(adminRepo)
			// 第三步：把业务用例交给接口层。
			svc := service.NewBusinessService(biz.NewSystemUsecase(dataSet), catalog, adminUsecase)

			httpServer, err := server.NewBusinessHTTPServer(
				cfg, svc, nil,
			)
			return httpServer, cfg.HTTP.ShutdownTimeout, err
		},
	)
}

func RunOrder() error {
	return run("dongdong-order-api", func(cfg *config.Config, dataSet *data.Data) (*khttp.Server, time.Duration, error) {
		svc := service.NewOrderService(biz.NewSystemUsecase(dataSet))
		httpServer, err := server.NewOrderHTTPServer(cfg, svc, nil)
		return httpServer, cfg.OrderHTTP.ShutdownTimeout, err
	})
}

type serverBuilder func(*config.Config, *data.Data) (*khttp.Server, time.Duration, error)

func run(name string, build serverBuilder) error {
	cfg, err := config.Load("")
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}
	infra.Init(cfg.Log.Level, cfg.Log.Format)
	pgCfg, err := PostgresConfig(cfg.PG)
	if err != nil {
		return err
	}
	pg, err := infra.NewPG(context.Background(), &pgCfg)
	if err != nil {
		return fmt.Errorf("initialize postgres: %w", err)
	}
	if err := pg.AutoMigrate(model.AllModels()...); err != nil {
		return fmt.Errorf("auto migrate postgres: %w", err)
	}
	defer func() {
		if err := infra.ClosePG(pg); err != nil {
			infra.L().Error("close postgres failed", "error", err)
		}
	}()
	httpServer, stopTimeout, err := build(cfg, data.NewData(pg))
	if err != nil {
		return err
	}
	opts := []kratos.Option{kratos.Name(name), kratos.Version("0.3.0"), kratos.Logger(infra.L()), kratos.Server(httpServer)}
	if stopTimeout > 0 {
		opts = append(opts, kratos.StopTimeout(stopTimeout))
	}
	return kratos.New(opts...).Run()
}

func PostgresConfig(c config.PG) (config.PG, error) {
	if v := os.Getenv("DD_DB_HOST"); v != "" {
		c.Host = v
	}
	if v := os.Getenv("DD_DB_NAME"); v != "" {
		c.Database = v
	}
	if v := os.Getenv("DD_DB_USER"); v != "" {
		c.User = v
	}
	if v := os.Getenv("DD_DB_PASSWORD"); v != "" {
		c.Password = v
	}
	if v := os.Getenv("DD_DB_SSL_MODE"); v != "" {
		c.SSLMode = v
	}
	if v := os.Getenv("DD_DB_TIME_ZONE"); v != "" {
		c.TimeZone = v
	}
	if v := os.Getenv("DD_DB_PORT"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			return c, fmt.Errorf("DD_DB_PORT must be a valid port number")
		}
		c.Port = n
	}
	c.Host, c.Database, c.User = strings.TrimSpace(c.Host), strings.TrimSpace(c.Database), strings.TrimSpace(c.User)
	if c.Host == "" || c.Database == "" || c.User == "" || c.Password == "" {
		return c, fmt.Errorf("PostgreSQL host, database, user, and password are required")
	}
	if c.Port == 0 {
		c.Port = 5432
	}
	if c.Port < 1 || c.Port > 65535 {
		return c, fmt.Errorf("PostgreSQL port must be between 1 and 65535")
	}
	if strings.TrimSpace(c.SSLMode) == "" {
		c.SSLMode = "disable"
	}
	if strings.TrimSpace(c.TimeZone) == "" {
		c.TimeZone = "UTC"
	}
	if c.MaxOpenConns == 0 {
		c.MaxOpenConns = 20
	}
	if c.ConnMaxLifetime == 0 {
		c.ConnMaxLifetime = 30 * time.Minute
	}
	for key, value := range map[string]*int{"DD_DB_MAX_OPEN": &c.MaxOpenConns, "DD_DB_MAX_IDLE": &c.MaxIdleConns} {
		if raw := os.Getenv(key); raw != "" {
			n, err := strconv.Atoi(raw)
			if err != nil || n < 0 {
				return c, fmt.Errorf("%s must be nonnegative", key)
			}
			*value = n
		}
	}
	if c.MaxOpenConns < 1 || c.MaxIdleConns < 0 || c.MaxIdleConns > c.MaxOpenConns || c.ConnMaxLifetime < 0 {
		return c, fmt.Errorf("invalid PostgreSQL pool limits or lifetime")
	}
	return c, nil
}
