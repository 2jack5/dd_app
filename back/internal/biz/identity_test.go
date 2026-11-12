package biz

import (
	"context"
	"errors"
	"testing"
	"time"
)

type identityVerifierStub struct {
	result IdentityProviderResult
	err    error
}

func (s identityVerifierStub) Verify(context.Context, IdentityVerifyInput) (IdentityProviderResult, error) {
	return s.result, s.err
}

type identityRepoStub struct {
	record   IdentityRecord
	validErr error
}

func (s *identityRepoStub) CreateIdentityCheck(_ context.Context, record IdentityRecord) (IdentityRecord, error) {
	record.ID = 10
	record.CreatedAt = time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	s.record = record
	return record, nil
}
func (s *identityRepoStub) CurrentIdentityCheck(context.Context, int64) (IdentityRecord, error) {
	return s.record, nil
}
func (s *identityRepoStub) ValidAdultIdentityCheck(context.Context, int64, time.Time) (IdentityRecord, error) {
	return s.record, s.validErr
}

func TestIdentityUsecaseVerifiesAndStoresOnlyMaskedName(t *testing.T) {
	repo := &identityRepoStub{}
	usecase := NewIdentityUsecase(identityVerifierStub{result: IdentityProviderResult{RequestID: "provider-1", Matched: true}}, repo)
	usecase.now = func() time.Time { return time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC) }
	record, err := usecase.Verify(context.Background(), 7, IdentityVerifyInput{
		Provider: "generic_http", ConsentVersion: "v1", RealName: "张三丰", IdentityNumber: "11010519491231002X",
	})
	if err != nil {
		t.Fatal(err)
	}
	if record.Status != IdentityStatusPassed || record.AgeCheckResult != AgeEligible || record.NameMasked != "张*丰" || record.VerifiedAt == nil {
		t.Fatalf("unexpected identity result: %+v", record)
	}
	if repo.record.NameMasked == "张三丰" {
		t.Fatal("plain real name was persisted")
	}
}

func TestIdentityUsecaseRejectsInvalidIdentityNumberBeforeProvider(t *testing.T) {
	usecase := NewIdentityUsecase(identityVerifierStub{}, &identityRepoStub{})
	_, err := usecase.Verify(context.Background(), 7, IdentityVerifyInput{
		Provider: "GENERIC_HTTP", ConsentVersion: "v1", RealName: "张三", IdentityNumber: "110105194912310021",
	})
	if !errors.Is(err, ErrInvalidArgument) {
		t.Fatalf("invalid identity number: %v", err)
	}
}

func TestIdentityUsecaseRecordsProviderMismatchAsFailed(t *testing.T) {
	repo := &identityRepoStub{}
	usecase := NewIdentityUsecase(identityVerifierStub{result: IdentityProviderResult{RequestID: "provider-2", Matched: false}}, repo)
	record, err := usecase.Verify(context.Background(), 7, IdentityVerifyInput{
		Provider: "GENERIC_HTTP", ConsentVersion: "v1", RealName: "张三", IdentityNumber: "11010519491231002X",
	})
	if err != nil || record.Status != IdentityStatusFailed || record.AgeCheckResult != AgeUnknown || record.VerifiedAt != nil {
		t.Fatalf("unexpected failed result: %+v err=%v", record, err)
	}
}

func TestIdentityUsecaseRequiresPassedAdultIdentity(t *testing.T) {
	now := time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC)
	repo := &identityRepoStub{record: IdentityRecord{Status: IdentityStatusPassed, AgeCheckResult: AgeEligible, VerifiedAt: &now}}
	usecase := NewIdentityUsecase(nil, repo)
	usecase.now = func() time.Time { return now }
	if err := usecase.RequireAdult(context.Background(), 7); err != nil {
		t.Fatalf("eligible adult rejected: %v", err)
	}
	repo.record.AgeCheckResult = AgeIneligible
	if err := usecase.RequireAdult(context.Background(), 7); !errors.Is(err, ErrAdultIdentityRequired) {
		t.Fatalf("underage identity accepted: %v", err)
	}
}
