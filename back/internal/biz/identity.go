package biz

import (
	"context"
	"errors"
	"strings"
	"time"
	"unicode/utf8"
)

var ErrIdentityUnavailable = errors.New("identity verification unavailable")
var ErrAdultIdentityRequired = errors.New("verified adult identity required")

const (
	IdentityStatusPassed = "PASSED"
	IdentityStatusFailed = "FAILED"
	AgeEligible          = "ELIGIBLE"
	AgeIneligible        = "INELIGIBLE"
	AgeUnknown           = "UNKNOWN"
)

type IdentityVerifyInput struct {
	Provider       string
	ConsentVersion string
	RealName       string
	IdentityNumber string
}

type IdentityProviderResult struct {
	RequestID string
	Matched   bool
	ResultRef string
	ExpiresAt *time.Time
}

type IdentityRecord struct {
	ID                int64
	UserID            int64
	Provider          string
	ProviderRequestID string
	ConsentVersion    string
	Status            string
	NameMasked        string
	AgeCheckResult    string
	ResultRef         string
	VerifiedAt        *time.Time
	ExpiresAt         *time.Time
	CreatedAt         time.Time
}

type IdentityVerifier interface {
	Verify(context.Context, IdentityVerifyInput) (IdentityProviderResult, error)
}

type IdentityRepository interface {
	CreateIdentityCheck(context.Context, IdentityRecord) (IdentityRecord, error)
	CurrentIdentityCheck(context.Context, int64) (IdentityRecord, error)
	ValidAdultIdentityCheck(context.Context, int64, time.Time) (IdentityRecord, error)
}

type IdentityUsecase struct {
	verifier IdentityVerifier
	repo     IdentityRepository
	now      func() time.Time
}

func NewIdentityUsecase(verifier IdentityVerifier, repo IdentityRepository) *IdentityUsecase {
	return &IdentityUsecase{verifier: verifier, repo: repo, now: time.Now}
}

func (u *IdentityUsecase) Verify(ctx context.Context, userID int64, input IdentityVerifyInput) (IdentityRecord, error) {
	input.Provider = strings.TrimSpace(strings.ToUpper(input.Provider))
	input.ConsentVersion = strings.TrimSpace(input.ConsentVersion)
	input.RealName = strings.TrimSpace(input.RealName)
	input.IdentityNumber = strings.ToUpper(strings.TrimSpace(input.IdentityNumber))
	if userID <= 0 || input.Provider == "" || len(input.Provider) > 32 || input.ConsentVersion == "" || len(input.ConsentVersion) > 64 || utf8.RuneCountInString(input.RealName) < 2 || utf8.RuneCountInString(input.RealName) > 64 {
		return IdentityRecord{}, ErrInvalidArgument
	}
	birthDate, ok := parsePRCIdentityNumber(input.IdentityNumber)
	if !ok {
		return IdentityRecord{}, ErrInvalidArgument
	}
	if u == nil || u.verifier == nil || u.repo == nil {
		return IdentityRecord{}, ErrIdentityUnavailable
	}
	providerResult, err := u.verifier.Verify(ctx, input)
	if err != nil {
		return IdentityRecord{}, err
	}
	providerResult.RequestID = strings.TrimSpace(providerResult.RequestID)
	if providerResult.RequestID == "" || len(providerResult.RequestID) > 128 {
		return IdentityRecord{}, ErrIdentityUnavailable
	}
	now := u.now().UTC()
	record := IdentityRecord{
		UserID: userID, Provider: input.Provider, ProviderRequestID: providerResult.RequestID,
		ConsentVersion: input.ConsentVersion, NameMasked: maskRealName(input.RealName),
		Status: IdentityStatusFailed, AgeCheckResult: AgeUnknown,
		ResultRef: strings.TrimSpace(providerResult.ResultRef), ExpiresAt: providerResult.ExpiresAt,
	}
	if providerResult.Matched {
		record.Status = IdentityStatusPassed
		record.AgeCheckResult = ageCheckResult(birthDate, now)
		record.VerifiedAt = &now
	}
	return u.repo.CreateIdentityCheck(ctx, record)
}

func (u *IdentityUsecase) Current(ctx context.Context, userID int64) (IdentityRecord, error) {
	if userID <= 0 {
		return IdentityRecord{}, ErrInvalidArgument
	}
	if u == nil || u.repo == nil {
		return IdentityRecord{}, ErrIdentityUnavailable
	}
	return u.repo.CurrentIdentityCheck(ctx, userID)
}

func (u *IdentityUsecase) RequireAdult(ctx context.Context, userID int64) error {
	if userID <= 0 {
		return ErrAdultIdentityRequired
	}
	if u == nil || u.repo == nil {
		return ErrIdentityUnavailable
	}
	now := u.now().UTC()
	record, err := u.repo.ValidAdultIdentityCheck(ctx, userID, now)
	if err != nil {
		return err
	}
	if record.Status != IdentityStatusPassed || record.AgeCheckResult != AgeEligible || record.VerifiedAt == nil || record.ExpiresAt != nil && !record.ExpiresAt.After(now) {
		return ErrAdultIdentityRequired
	}
	return nil
}

func parsePRCIdentityNumber(value string) (time.Time, bool) {
	if len(value) != 18 {
		return time.Time{}, false
	}
	weights := [...]int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}
	checks := "10X98765432"
	sum := 0
	for i := 0; i < 17; i++ {
		if value[i] < '0' || value[i] > '9' {
			return time.Time{}, false
		}
		sum += int(value[i]-'0') * weights[i]
	}
	if value[17] != checks[sum%11] {
		return time.Time{}, false
	}
	birthDate, err := time.Parse("20060102", value[6:14])
	if err != nil || birthDate.After(time.Now().UTC()) {
		return time.Time{}, false
	}
	return birthDate, true
}

func ageCheckResult(birthDate, now time.Time) string {
	age := now.Year() - birthDate.Year()
	if now.Month() < birthDate.Month() || now.Month() == birthDate.Month() && now.Day() < birthDate.Day() {
		age--
	}
	if age >= 18 {
		return AgeEligible
	}
	return AgeIneligible
}

func maskRealName(value string) string {
	runes := []rune(value)
	if len(runes) == 0 {
		return ""
	}
	if len(runes) == 1 {
		return "*"
	}
	if len(runes) == 2 {
		return string(runes[0]) + "*"
	}
	return string(runes[0]) + strings.Repeat("*", len(runes)-2) + string(runes[len(runes)-1])
}
