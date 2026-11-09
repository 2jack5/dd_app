package server

import (
	orderv1 "dd_web/back/api/order/v1"
	"dd_web/back/internal/biz"
	"dd_web/back/internal/config"
	"dd_web/back/internal/service"
	khttp "github.com/go-kratos/kratos/v3/transport/http"
)

func NewOrderHTTPServer(cfg *config.Config, svc *service.OrderService, auth biz.Authenticator) (*khttp.Server, error) {
	srv, err := newHTTPServer(cfg.OrderHTTP, "127.0.0.1:8001", "DD_ORDER_HTTP_ADDR", "DD_ORDER_HTTP_TIMEOUT", auth, orderPolicies)
	if err != nil {
		return nil, err
	}
	orderv1.RegisterSystemServiceHTTPServer(srv, svc)
	orderv1.RegisterUserOrderServiceHTTPServer(srv, svc)
	orderv1.RegisterCoachOrderServiceHTTPServer(srv, svc)
	orderv1.RegisterAdminOrderServiceHTTPServer(srv, svc)
	orderv1.RegisterCallbackServiceHTTPServer(srv, svc)
	return srv, nil
}
