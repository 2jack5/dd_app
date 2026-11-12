package service

import (
	"context"
	"errors"
	"strings"
	"time"

	commonv1 "dd_web/back/api/common/v1"
	orderv1 "dd_web/back/api/order/v1"
	"dd_web/back/internal/biz"
	"dd_web/back/internal/ecode"
	"github.com/go-kratos/kratos/v3/transport"
)

func (s *OrderService) WalletCurrent(ctx context.Context, _ *orderv1.WalletCurrentRequest) (*commonv1.ApiResponse, error) {
	principal, err := walletPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	if s.wallet == nil {
		return nil, ecode.New(ecode.WalletUnavailable)
	}
	result, err := s.wallet.Wallet(ctx, principal.UserID)
	if err != nil {
		return nil, walletError(err)
	}
	return ok(result)
}

func (s *OrderService) RechargePackages(ctx context.Context, req *orderv1.RechargePackagesRequest) (*commonv1.ApiResponse, error) {
	if _, err := walletPrincipal(ctx); err != nil {
		return nil, err
	}
	if s.wallet == nil {
		return nil, ecode.New(ecode.WalletUnavailable)
	}
	page, err := biz.ParsePage(req.GetCursor(), req.GetPageSize())
	if err != nil {
		return nil, ecode.New(ecode.InvalidArgument)
	}
	result, err := s.wallet.RechargePackages(ctx, page)
	if err != nil {
		return nil, walletError(err)
	}
	return ok(map[string]any{"items": result.Items, "next_cursor": result.NextCursor})
}

func (s *OrderService) RechargeCreate(ctx context.Context, req *orderv1.RechargeCreateRequest) (*commonv1.ApiResponse, error) {
	principal, err := walletPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	if s.wallet == nil {
		return nil, ecode.New(ecode.WalletUnavailable)
	}
	if req.GetBody() == nil {
		return nil, ecode.New(ecode.InvalidArgument)
	}
	packageID, err := biz.ParseID(req.GetBody().GetPackageId())
	if err != nil {
		return nil, ecode.New(ecode.InvalidArgument)
	}
	result, err := s.wallet.CreateRecharge(ctx, biz.RechargeCreateInput{UserID: principal.UserID, PackageID: packageID, Provider: req.GetBody().GetProvider(), ClientType: req.GetBody().GetClientType(), IdempotencyKey: idempotencyKey(ctx)})
	if err != nil {
		return nil, walletError(err)
	}
	return ok(result)
}

func (s *OrderService) RechargeDetail(ctx context.Context, req *orderv1.RechargeDetailRequest) (*commonv1.ApiResponse, error) {
	principal, err := walletPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	if s.wallet == nil {
		return nil, ecode.New(ecode.WalletUnavailable)
	}
	rechargeID, err := biz.ParseID(req.GetRechargeId())
	if err != nil {
		return nil, ecode.New(ecode.InvalidArgument)
	}
	result, err := s.wallet.Recharge(ctx, principal.UserID, rechargeID)
	if err != nil {
		return nil, walletError(err)
	}
	return ok(result)
}

func (s *OrderService) RechargeList(ctx context.Context, req *orderv1.RechargeListRequest) (*commonv1.ApiResponse, error) {
	principal, err := walletPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	if s.wallet == nil {
		return nil, ecode.New(ecode.WalletUnavailable)
	}
	page, err := biz.ParsePage(req.GetCursor(), req.GetPageSize())
	if err != nil {
		return nil, ecode.New(ecode.InvalidArgument)
	}
	result, err := s.wallet.Recharges(ctx, principal.UserID, page, req.GetStatus())
	if err != nil {
		return nil, walletError(err)
	}
	return ok(map[string]any{"items": result.Items, "next_cursor": result.NextCursor, "server_time": time.Now().UTC().Format(time.RFC3339)})
}

func (s *OrderService) WalletTransactions(ctx context.Context, req *orderv1.WalletTransactionsRequest) (*commonv1.ApiResponse, error) {
	principal, err := walletPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	if s.wallet == nil {
		return nil, ecode.New(ecode.WalletUnavailable)
	}
	page, err := biz.ParsePage(req.GetCursor(), req.GetPageSize())
	if err != nil {
		return nil, ecode.New(ecode.InvalidArgument)
	}
	result, err := s.wallet.Transactions(ctx, principal.UserID, page, req.GetType())
	if err != nil {
		return nil, walletError(err)
	}
	return ok(map[string]any{"items": result.Items, "next_cursor": result.NextCursor})
}

func (s *OrderService) WalletCallback(ctx context.Context, req *orderv1.WalletCallbackRequest) (*commonv1.ApiResponse, error) {
	if s.wallet == nil {
		return nil, ecode.New(ecode.ChannelNotConfigured)
	}
	payload := map[string]any{}
	if req.GetBody() != nil {
		payload = req.GetBody().AsMap()
	}
	result, err := s.wallet.Callback(ctx, req.GetProvider(), payload)
	if err != nil {
		return nil, walletError(err)
	}
	return ok(result)
}

func (s *OrderService) AdminRechargePackages(ctx context.Context, req *orderv1.AdminRechargePackagesRequest) (*commonv1.ApiResponse, error) {
	if _, err := walletAdminPrincipal(ctx); err != nil {
		return nil, err
	}
	if s.wallet == nil {
		return nil, ecode.New(ecode.WalletUnavailable)
	}
	page, err := biz.ParsePage(req.GetCursor(), req.GetPageSize())
	if err != nil {
		return nil, ecode.New(ecode.InvalidArgument)
	}
	result, err := s.wallet.AdminRechargePackages(ctx, page, req.GetStatus())
	if err != nil {
		return nil, walletError(err)
	}
	return ok(map[string]any{"items": result.Items, "next_cursor": result.NextCursor})
}

func (s *OrderService) AdminCreateRechargePackage(ctx context.Context, req *orderv1.AdminCreateRechargePackageRequest) (*commonv1.ApiResponse, error) {
	principal, err := walletAdminPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	if s.wallet == nil {
		return nil, ecode.New(ecode.WalletUnavailable)
	}
	if req.GetBody() == nil {
		return nil, ecode.New(ecode.InvalidArgument)
	}
	result, err := s.wallet.AdminCreateRechargePackage(ctx, principal.AdminID, rechargePackageInput(req.GetBody()), idempotencyKey(ctx))
	if err != nil {
		return nil, walletError(err)
	}
	return ok(result)
}

func (s *OrderService) AdminUpdateRechargePackage(ctx context.Context, req *orderv1.AdminUpdateRechargePackageRequest) (*commonv1.ApiResponse, error) {
	principal, err := walletAdminPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	if s.wallet == nil {
		return nil, ecode.New(ecode.WalletUnavailable)
	}
	if req.GetBody() == nil {
		return nil, ecode.New(ecode.InvalidArgument)
	}
	packageID, err := biz.ParseID(req.GetPackageId())
	if err != nil {
		return nil, ecode.New(ecode.InvalidArgument)
	}
	result, err := s.wallet.AdminUpdateRechargePackage(ctx, principal.AdminID, packageID, rechargePackageInput(req.GetBody()), idempotencyKey(ctx))
	if err != nil {
		return nil, walletError(err)
	}
	return ok(result)
}

func walletPrincipal(ctx context.Context) (biz.Principal, error) {
	principal, found := biz.PrincipalFromContext(ctx)
	if !found || principal.UserID <= 0 {
		return biz.Principal{}, ecode.New(ecode.AuthRequired)
	}
	return principal, nil
}

func walletAdminPrincipal(ctx context.Context) (biz.Principal, error) {
	principal, found := biz.PrincipalFromContext(ctx)
	if !found || principal.AdminID <= 0 {
		return biz.Principal{}, ecode.New(ecode.AuthRequired)
	}
	return principal, nil
}

func rechargePackageInput(body *orderv1.RechargePackageBody) biz.RechargePackageInput {
	return biz.RechargePackageInput{Name: body.GetName(), PayAmount: body.GetPayAmount(), PackageMultiplier: body.GetPackageMultiplier(), Currency: body.GetCurrency(), AssetCode: body.GetAssetCode(), Status: body.GetStatus(), EffectiveAt: body.GetEffectiveAt(), ExpiresAt: body.GetExpiresAt(), Version: body.GetVersion()}
}

func idempotencyKey(ctx context.Context) string {
	if tr, ok := transport.FromServerContext(ctx); ok {
		return strings.TrimSpace(tr.RequestHeader().Get("Idempotency-Key"))
	}
	return ""
}

func walletError(err error) error {
	switch {
	case errors.Is(err, biz.ErrInvalidArgument):
		return ecode.New(ecode.InvalidArgument)
	case errors.Is(err, biz.ErrWalletNotFound):
		return ecode.New(ecode.NotFound)
	case errors.Is(err, biz.ErrWalletUnavailable):
		return ecode.New(ecode.WalletUnavailable)
	case errors.Is(err, biz.ErrIdempotencyConflict):
		return ecode.New(ecode.IdempotencyConflict)
	case errors.Is(err, biz.ErrRechargeStateConflict):
		return ecode.New(ecode.RechargeStateConflict)
	case errors.Is(err, biz.ErrRechargePackageConflict):
		return ecode.New(ecode.RechargePackageConflict)
	case errors.Is(err, biz.ErrRechargeChannel):
		return ecode.New(ecode.ChannelNotConfigured)
	default:
		return ecode.Wrap(ecode.InternalError, err)
	}
}
