package service

import (
	"context"
	"testing"
	"time"

	businessv1 "dd_web/back/api/business/v1"
	commonv1 "dd_web/back/api/common/v1"
	"dd_web/back/internal/biz"
)

type serviceIdentityVerifier struct{}

func (serviceIdentityVerifier) Verify(context.Context, biz.IdentityVerifyInput) (biz.IdentityProviderResult, error) {
	return biz.IdentityProviderResult{RequestID: "request-1", Matched: true}, nil
}

type serviceIdentityRepo struct{ record biz.IdentityRecord }

func (r *serviceIdentityRepo) CreateIdentityCheck(_ context.Context, record biz.IdentityRecord) (biz.IdentityRecord, error) {
	record.ID = 12
	r.record = record
	return record, nil
}
func (r *serviceIdentityRepo) CurrentIdentityCheck(context.Context, int64) (biz.IdentityRecord, error) {
	return r.record, nil
}
func (r *serviceIdentityRepo) ValidAdultIdentityCheck(context.Context, int64, time.Time) (biz.IdentityRecord, error) {
	return r.record, nil
}

func TestUserVerifyIdentityReturnsSanitizedResult(t *testing.T) {
	repo := &serviceIdentityRepo{}
	identity := biz.NewIdentityUsecase(serviceIdentityVerifier{}, repo)
	service := NewBusinessService(nil, nil, nil, nil, identity)
	ctx := biz.WithPrincipal(context.Background(), biz.Principal{UserID: 7})
	response, err := service.UserVerifyIdentity(ctx, &businessv1.UserVerifyIdentityRequest{Body: &commonv1.IdentityCheck{
		Provider: "GENERIC_HTTP", ConsentVersion: "v1", RealName: "张三", IdentityNumber: "11010519491231002X",
	}})
	if err != nil {
		t.Fatal(err)
	}
	data := response.Data.AsInterface().(map[string]any)
	if data["status"] != biz.IdentityStatusPassed || data["name_masked"] != "张*" {
		t.Fatalf("unexpected response: %#v", data)
	}
	if _, exists := data["identity_number"]; exists {
		t.Fatal("identity number leaked in response")
	}
}
