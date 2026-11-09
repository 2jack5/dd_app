package service

import (
	"dd_web/back/internal/biz"
)

type OrderService struct{ system *biz.SystemUsecase }

func NewOrderService(system *biz.SystemUsecase) *OrderService { return &OrderService{system: system} }
