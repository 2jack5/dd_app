package service

import (
	"context"
	commonv1 "dd_web/back/api/common/v1"
	orderv1 "dd_web/back/api/order/v1"
	kerrors "github.com/go-kratos/kratos/v3/errors"
)

func (s *OrderService) SystemHealth(context.Context, *orderv1.SystemHealthRequest) (*commonv1.ApiResponse, error) {
	return ok(map[string]any{"status": "alive"})
}
func (s *OrderService) SystemReady(ctx context.Context, _ *orderv1.SystemReadyRequest) (*commonv1.ApiResponse, error) {
	if err := s.system.Ready(ctx); err != nil {
		return nil, kerrors.ServiceUnavailable("NOT_READY", "数据库连接未就绪")
	}
	return ok(map[string]any{"status": "ready"})
}
