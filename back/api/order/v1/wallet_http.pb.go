// Code generated from api/order/v1/wallet.proto. DO NOT EDIT.

package orderv1

import (
	"context"

	commonv1 "dd_web/back/api/common/v1"
	khttp "github.com/go-kratos/kratos/v3/transport/http"
)

const (
	OperationWalletServiceWalletCurrent      = "/api.order.v1.WalletService/WalletCurrent"
	OperationWalletServiceRechargePackages   = "/api.order.v1.WalletService/RechargePackages"
	OperationWalletServiceRechargeCreate     = "/api.order.v1.WalletService/RechargeCreate"
	OperationWalletServiceRechargeDetail     = "/api.order.v1.WalletService/RechargeDetail"
	OperationWalletServiceRechargeList       = "/api.order.v1.WalletService/RechargeList"
	OperationWalletServiceWalletTransactions = "/api.order.v1.WalletService/WalletTransactions"
	OperationWalletCallbackServiceCallback   = "/api.order.v1.WalletCallbackService/WalletCallback"
	OperationAdminWalletServicePackages      = "/api.order.v1.AdminWalletService/AdminRechargePackages"
	OperationAdminWalletServiceCreatePackage = "/api.order.v1.AdminWalletService/AdminCreateRechargePackage"
	OperationAdminWalletServiceUpdatePackage = "/api.order.v1.AdminWalletService/AdminUpdateRechargePackage"
)

type WalletServiceHTTPServer interface {
	WalletCurrent(context.Context, *WalletCurrentRequest) (*commonv1.ApiResponse, error)
	RechargePackages(context.Context, *RechargePackagesRequest) (*commonv1.ApiResponse, error)
	RechargeCreate(context.Context, *RechargeCreateRequest) (*commonv1.ApiResponse, error)
	RechargeDetail(context.Context, *RechargeDetailRequest) (*commonv1.ApiResponse, error)
	RechargeList(context.Context, *RechargeListRequest) (*commonv1.ApiResponse, error)
	WalletTransactions(context.Context, *WalletTransactionsRequest) (*commonv1.ApiResponse, error)
}

type WalletCallbackServiceHTTPServer interface {
	WalletCallback(context.Context, *WalletCallbackRequest) (*commonv1.ApiResponse, error)
}

type AdminWalletServiceHTTPServer interface {
	AdminRechargePackages(context.Context, *AdminRechargePackagesRequest) (*commonv1.ApiResponse, error)
	AdminCreateRechargePackage(context.Context, *AdminCreateRechargePackageRequest) (*commonv1.ApiResponse, error)
	AdminUpdateRechargePackage(context.Context, *AdminUpdateRechargePackageRequest) (*commonv1.ApiResponse, error)
}

func RegisterWalletServiceHTTPServer(s *khttp.Server, srv WalletServiceHTTPServer) {
	r := s.Route("/")
	r.Handle("GET", "/api/v1/wallet", walletHandler(OperationWalletServiceWalletCurrent, func(ctx context.Context, in any) (*commonv1.ApiResponse, error) {
		return srv.WalletCurrent(ctx, in.(*WalletCurrentRequest))
	}, func(ctx khttp.Context) (any, error) { var in WalletCurrentRequest; return &in, ctx.BindQuery(&in) }))
	r.Handle("GET", "/api/v1/wallet/recharge-packages", walletHandler(OperationWalletServiceRechargePackages, func(ctx context.Context, in any) (*commonv1.ApiResponse, error) {
		return srv.RechargePackages(ctx, in.(*RechargePackagesRequest))
	}, func(ctx khttp.Context) (any, error) { var in RechargePackagesRequest; return &in, ctx.BindQuery(&in) }))
	r.Handle("POST", "/api/v1/wallet/recharges", walletHandler(OperationWalletServiceRechargeCreate, func(ctx context.Context, in any) (*commonv1.ApiResponse, error) {
		return srv.RechargeCreate(ctx, in.(*RechargeCreateRequest))
	}, func(ctx khttp.Context) (any, error) {
		var in RechargeCreateRequest
		if err := ctx.Bind(&in.Body); err != nil {
			return nil, err
		}
		return &in, ctx.BindQuery(&in)
	}))
	r.Handle("GET", "/api/v1/wallet/recharges/{recharge_id}", walletHandler(OperationWalletServiceRechargeDetail, func(ctx context.Context, in any) (*commonv1.ApiResponse, error) {
		return srv.RechargeDetail(ctx, in.(*RechargeDetailRequest))
	}, func(ctx khttp.Context) (any, error) {
		var in RechargeDetailRequest
		if err := ctx.BindQuery(&in); err != nil {
			return nil, err
		}
		return &in, ctx.BindVars(&in)
	}))
	r.Handle("GET", "/api/v1/wallet/recharges", walletHandler(OperationWalletServiceRechargeList, func(ctx context.Context, in any) (*commonv1.ApiResponse, error) {
		return srv.RechargeList(ctx, in.(*RechargeListRequest))
	}, func(ctx khttp.Context) (any, error) { var in RechargeListRequest; return &in, ctx.BindQuery(&in) }))
	r.Handle("GET", "/api/v1/wallet/transactions", walletHandler(OperationWalletServiceWalletTransactions, func(ctx context.Context, in any) (*commonv1.ApiResponse, error) {
		return srv.WalletTransactions(ctx, in.(*WalletTransactionsRequest))
	}, func(ctx khttp.Context) (any, error) { var in WalletTransactionsRequest; return &in, ctx.BindQuery(&in) }))
}

func RegisterWalletCallbackServiceHTTPServer(s *khttp.Server, srv WalletCallbackServiceHTTPServer) {
	r := s.Route("/")
	r.Handle("POST", "/api/v1/callbacks/wallet/{provider}", walletHandler(OperationWalletCallbackServiceCallback, func(ctx context.Context, in any) (*commonv1.ApiResponse, error) {
		return srv.WalletCallback(ctx, in.(*WalletCallbackRequest))
	}, func(ctx khttp.Context) (any, error) {
		var in WalletCallbackRequest
		if err := ctx.Bind(&in.Body); err != nil {
			return nil, err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return nil, err
		}
		return &in, ctx.BindVars(&in)
	}))
}

func RegisterAdminWalletServiceHTTPServer(s *khttp.Server, srv AdminWalletServiceHTTPServer) {
	r := s.Route("/")
	r.Handle("GET", "/api/v1/admin/wallet/recharge-packages", walletHandler(OperationAdminWalletServicePackages, func(ctx context.Context, in any) (*commonv1.ApiResponse, error) {
		return srv.AdminRechargePackages(ctx, in.(*AdminRechargePackagesRequest))
	}, func(ctx khttp.Context) (any, error) {
		var in AdminRechargePackagesRequest
		return &in, ctx.BindQuery(&in)
	}))
	r.Handle("POST", "/api/v1/admin/wallet/recharge-packages", walletHandler(OperationAdminWalletServiceCreatePackage, func(ctx context.Context, in any) (*commonv1.ApiResponse, error) {
		return srv.AdminCreateRechargePackage(ctx, in.(*AdminCreateRechargePackageRequest))
	}, func(ctx khttp.Context) (any, error) {
		var in AdminCreateRechargePackageRequest
		if err := ctx.Bind(&in.Body); err != nil {
			return nil, err
		}
		return &in, ctx.BindQuery(&in)
	}))
	r.Handle("PUT", "/api/v1/admin/wallet/recharge-packages/{package_id}", walletHandler(OperationAdminWalletServiceUpdatePackage, func(ctx context.Context, in any) (*commonv1.ApiResponse, error) {
		return srv.AdminUpdateRechargePackage(ctx, in.(*AdminUpdateRechargePackageRequest))
	}, func(ctx khttp.Context) (any, error) {
		var in AdminUpdateRechargePackageRequest
		if err := ctx.Bind(&in.Body); err != nil {
			return nil, err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return nil, err
		}
		return &in, ctx.BindVars(&in)
	}))
}

type walletHTTPCall func(context.Context, any) (*commonv1.ApiResponse, error)
type walletHTTPBind func(khttp.Context) (any, error)

func walletHandler(operation string, call walletHTTPCall, bind walletHTTPBind) func(khttp.Context) error {
	return func(ctx khttp.Context) error {
		in, err := bind(ctx)
		if err != nil {
			return err
		}
		khttp.SetOperation(ctx, operation)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) { return call(ctx, req) })
		out, err := h(ctx, in)
		if err != nil {
			return err
		}
		return ctx.Result(200, out.(*commonv1.ApiResponse))
	}
}
