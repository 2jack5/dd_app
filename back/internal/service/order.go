package service

import (
	"dd_web/back/internal/biz"
)

type OrderService struct {
	system   *biz.SystemUsecase
	identity *biz.IdentityUsecase
}

func NewOrderService(system *biz.SystemUsecase, identity ...*biz.IdentityUsecase) *OrderService {
	service := &OrderService{system: system}
	if len(identity) > 0 {
		service.identity = identity[0]
	}
	return service
}
