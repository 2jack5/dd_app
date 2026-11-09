package server

import (
	businessv1 "dd_web/back/api/business/v1"
	"dd_web/back/internal/biz"
	"dd_web/back/internal/config"
	"dd_web/back/internal/service"
	khttp "github.com/go-kratos/kratos/v3/transport/http"
)

func NewBusinessHTTPServer(cfg *config.Config, svc *service.BusinessService, auth biz.Authenticator) (*khttp.Server, error) {
	srv, err := newHTTPServer(cfg.HTTP, "127.0.0.1:8000", "DD_BUSINESS_HTTP_ADDR", "DD_BUSINESS_HTTP_TIMEOUT", auth, businessPolicies)
	if err != nil {
		return nil, err
	}
	businessv1.RegisterPublicServiceHTTPServer(srv, svc)
	businessv1.RegisterAccountServiceHTTPServer(srv, svc)
	businessv1.RegisterCoachServiceHTTPServer(srv, svc)
	businessv1.RegisterAdminServiceHTTPServer(srv, svc)
	businessv1.RegisterCallbackServiceHTTPServer(srv, svc)
	return srv, nil
}
