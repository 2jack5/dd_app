package service

import (
	commonv1 "dd_web/back/api/common/v1"
	"dd_web/back/internal/biz"
	"dd_web/back/internal/data"
	"errors"

	kerrors "github.com/go-kratos/kratos/v3/errors"
	"google.golang.org/protobuf/types/known/structpb"
)

type BusinessService struct {
	system  *biz.SystemUsecase
	catalog *biz.CatalogUsecase
	admin   *biz.AdminUsecase // 新增
}

func NewBusinessService(system *biz.SystemUsecase, catalog *biz.CatalogUsecase, admin *biz.AdminUsecase) *BusinessService {
	return &BusinessService{system: system, catalog: catalog, admin: admin}
}
func ok(data any) (*commonv1.ApiResponse, error) {
	value, err := structpb.NewValue(data)
	if err != nil {
		return nil, kerrors.InternalServer("INTERNAL_ERROR", "服务暂时不可用")
	}
	return &commonv1.ApiResponse{Code: "OK", Data: value}, nil
}
func notImplemented(op string) error {
	return kerrors.New(501, "NOT_IMPLEMENTED", "接口已设计，业务逻辑尚未实现").WithMetadata(map[string]string{"operation": op})
}
func catalogError(err error) error {
	switch {
	case errors.Is(err, biz.ErrInvalidArgument):
		return kerrors.BadRequest("INVALID_ARGUMENT", "ID、日期或分页参数无效")
	case errors.Is(err, data.ErrNotFound):
		return kerrors.NotFound("NOT_FOUND", "资源不存在或未启用")
	default:
		return kerrors.InternalServer("INTERNAL_ERROR", "服务暂时不可用")
	}
}
