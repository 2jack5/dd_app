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
	"github.com/go-kratos/kratos/v3/transport"
	"github.com/redis/go-redis/v9"
)

func RunBusiness() error {
	return run("dongdong-business-api", func(cfg *config.Config, dataSet *data.Data, redisClient *redis.Client) ([]transport.Server, time.Duration, error) {
		catalog := biz.NewCatalogUsecase(data.NewCatalogRepo(dataSet))
		auth := newAuthUsecase(cfg, dataSet, redisClient)
		adminUsecase := newAdminUsecase(cfg, dataSet, redisClient)
		identity := biz.NewIdentityUsecase(infra.NewIdentityClient(cfg.Identity), data.NewIdentityRepo(dataSet))
		home := biz.NewHomeUsecase(data.NewHomeRepo(dataSet))
		svc := service.NewBusinessService(biz.NewSystemUsecase(data.NewDependencyHealth(dataSet, redisClient)), catalog, auth, adminUsecase, home, identity)
		authenticator := biz.ScopedAuthenticator{
			User:  auth,
			Admin: adminUsecase,
		}
		httpServer, err := server.NewBusinessHTTPServer(cfg, svc, authenticator)
		if err != nil {
			return nil, cfg.HTTP.ShutdownTimeout, err
		}
		return []transport.Server{httpServer}, cfg.HTTP.ShutdownTimeout, nil
	})
}

func RunOrder() error {
	return run("dongdong-order-api", func(cfg *config.Config, dataSet *data.Data, redisClient *redis.Client) ([]transport.Server, time.Duration, error) {
		identity := biz.NewIdentityUsecase(nil, data.NewIdentityRepo(dataSet))
		var rechargeProviders []biz.RechargeProvider
		if strings.EqualFold(strings.TrimSpace(cfg.WalletPayment.Provider), "LOCAL") {
			if strings.EqualFold(cfg.Env, "prod") {
				return nil, 0, fmt.Errorf("LOCAL wallet payment provider is not allowed in prod")
			}
			provider, err := infra.NewLocalRechargeProvider(cfg.WalletPayment)
			if err != nil {
				return nil, 0, fmt.Errorf("initialize wallet payment provider: %w", err)
			}
			rechargeProviders = append(rechargeProviders, provider)
		}
		wallet := biz.NewWalletUsecaseWithOptions(data.NewWalletRepo(dataSet), biz.WalletOptions{AssetCode: cfg.Wallet.AssetCode, RechargeExpire: cfg.Wallet.RechargeExpire}, rechargeProviders...)
		svc := service.NewOrderServiceWithWallet(biz.NewSystemUsecase(data.NewDependencyHealth(dataSet, redisClient)), identity, wallet)
		authenticator := biz.ScopedAuthenticator{
			User:  newAuthUsecase(cfg, dataSet, redisClient),
			Admin: newAdminUsecase(cfg, dataSet, redisClient),
		}
		httpServer, err := server.NewOrderHTTPServer(cfg, svc, authenticator)
		if err != nil {
			return nil, cfg.OrderHTTP.ShutdownTimeout, err
		}
		closer := server.NewWalletRechargeCloser(wallet, infra.L(), cfg.Wallet.CloseInterval, cfg.Wallet.CloseRetryDelay, cfg.Wallet.CloseBatchSize)
		return []transport.Server{httpServer, closer}, cfg.OrderHTTP.ShutdownTimeout, nil
	})
}

func newAuthUsecase(cfg *config.Config, dataSet *data.Data, redisClient *redis.Client) *biz.AuthUsecase {
	return biz.NewAuthUsecase(infra.NewWechatClient(cfg.Wechat), data.NewAuthRepo(dataSet), infra.NewTokenManager(cfg.Auth), data.NewRedisSessionStore(redisClient))
}
func newAdminUsecase(cfg *config.Config, dataSet *data.Data, redisClient *redis.Client) *biz.AdminUsecase {
	tokenConfig := cfg.Auth

	issuer := strings.TrimSpace(tokenConfig.Issuer)
	if issuer == "" {
		issuer = "dongdong-business-api"
	}
	tokenConfig.Issuer = issuer + "/admin"

	return biz.NewAdminUsecase(
		data.NewAdminRepo(dataSet),
		infra.NewTokenManager(tokenConfig),
		data.NewRedisSessionStore(redisClient, "admin"),
	)
}

type serverBuilder func(*config.Config, *data.Data, *redis.Client) ([]transport.Server, time.Duration, error)

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
	redisClient, err := infra.NewRedis(cfg.Redis)
	if err != nil {
		return fmt.Errorf("initialize redis: %w", err)
	}
	defer func() {
		if err := infra.CloseRedis(redisClient); err != nil {
			infra.L().Error("close redis failed", "error", err)
		}
	}()
	servers, stopTimeout, err := build(cfg, data.NewData(pg), redisClient)
	if err != nil {
		return err
	}
	opts := []kratos.Option{kratos.Name(name), kratos.Version("0.3.0"), kratos.Logger(infra.L()), kratos.Server(servers...)}
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
