package service

import (
	commonv1 "dd_web/back/api/common/v1"
	"dd_web/back/internal/biz"
	"dd_web/back/internal/data"
	"dd_web/back/internal/ecode"
	"errors"

	"google.golang.org/protobuf/types/known/structpb"
)

type BusinessService struct {
	system   *biz.SystemUsecase
	catalog  *biz.CatalogUsecase
	auth     *biz.AuthUsecase
	admin    *biz.AdminUsecase
	identity *biz.IdentityUsecase
	home     *biz.HomeUsecase
}

func NewBusinessService(system *biz.SystemUsecase, catalog *biz.CatalogUsecase, auth *biz.AuthUsecase, admin *biz.AdminUsecase, home *biz.HomeUsecase, identity ...*biz.IdentityUsecase) *BusinessService {
	svc := &BusinessService{system: system, catalog: catalog, auth: auth, admin: admin, home: home}
	if len(identity) > 0 {
		svc.identity = identity[0]
	}

	return svc
}
func ok(data any) (*commonv1.ApiResponse, error) {
	value, err := structpb.NewValue(data)
	if err != nil {
		return nil, ecode.Wrap(ecode.InternalError, err)
	}
	return &commonv1.ApiResponse{Code: int32(ecode.Success), ErrorMsg: "", Data: value}, nil
}
func notImplemented(op string) error {
	return ecode.New(ecode.NotImplemented).WithData(map[string]string{"operation": op})
}
func catalogError(err error) error {
	switch {
	case errors.Is(err, biz.ErrInvalidArgument):
		return ecode.New(ecode.InvalidArgument)
	case errors.Is(err, data.ErrNotFound):
		return ecode.New(ecode.NotFound)
	default:
		return ecode.Wrap(ecode.InternalError, err)
	}
}
