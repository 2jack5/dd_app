package infra

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"dd_web/back/internal/biz"
	"dd_web/back/internal/config"
)

type TokenManager struct {
	issuer     string
	secret     []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
	now        func() time.Time
}

func NewTokenManager(cfg config.Auth) *TokenManager {
	if cfg.Issuer == "" {
		cfg.Issuer = "dongdong-business-api"
	}
	if cfg.AccessTokenTTL <= 0 {
		cfg.AccessTokenTTL = 2 * time.Hour
	}
	if cfg.RefreshTokenTTL <= 0 {
		cfg.RefreshTokenTTL = 30 * 24 * time.Hour
	}
	return &TokenManager{issuer: cfg.Issuer, secret: []byte(cfg.TokenSecret), accessTTL: cfg.AccessTokenTTL, refreshTTL: cfg.RefreshTokenTTL, now: time.Now}
}

type jwtClaims struct {
	Issuer      string `json:"iss"`
	Subject     string `json:"sub"`
	Type        string `json:"typ"`
	AuthVersion int32  `json:"av"`
	IssuedAt    int64  `json:"iat"`
	ExpiresAt   int64  `json:"exp"`
	JWTID       string `json:"jti"`
	SessionID   string `json:"sid"`
}

func (m *TokenManager) Issue(userID int64, authVersion int32) (biz.TokenPair, error) {
	return m.IssueForSession(userID, authVersion, randomID())
}

func (m *TokenManager) IssueForSession(userID int64, authVersion int32, sessionID string) (biz.TokenPair, error) {
	if len(m.secret) < 32 {
		return biz.TokenPair{}, fmt.Errorf("%w: auth.token_secret must contain at least 32 bytes", biz.ErrAuthenticationUnavailable)
	}
	now := m.now().UTC()
	if strings.TrimSpace(sessionID) == "" {
		return biz.TokenPair{}, biz.ErrInvalidToken
	}
	accessExpires, refreshExpires := now.Add(m.accessTTL), now.Add(m.refreshTTL)
	access, err := m.sign(jwtClaims{m.issuer, strconv.FormatInt(userID, 10), "access", authVersion, now.Unix(), accessExpires.Unix(), randomID(), sessionID})
	if err != nil {
		return biz.TokenPair{}, err
	}
	refresh, err := m.sign(jwtClaims{m.issuer, strconv.FormatInt(userID, 10), "refresh", authVersion, now.Unix(), refreshExpires.Unix(), randomID(), sessionID})
	if err != nil {
		return biz.TokenPair{}, err
	}
	hash := sha256.Sum256([]byte(refresh))
	return biz.TokenPair{AccessToken: access, RefreshToken: refresh, SessionID: sessionID, AccessExpiresAt: accessExpires, RefreshExpiresAt: refreshExpires, RefreshTokenHash: hex.EncodeToString(hash[:])}, nil
}

func (m *TokenManager) VerifyAccess(token string) (biz.TokenClaims, error) {
	return m.verifyType(token, "access")
}

func (m *TokenManager) VerifyRefresh(token string) (biz.TokenClaims, error) {
	return m.verifyType(token, "refresh")
}

func (m *TokenManager) verifyType(token, tokenType string) (biz.TokenClaims, error) {
	claims, err := m.verify(token)
	if err != nil || claims.Type != tokenType || strings.TrimSpace(claims.SessionID) == "" {
		return biz.TokenClaims{}, errOrInvalidToken(err)
	}
	userID, err := strconv.ParseInt(claims.Subject, 10, 64)
	if err != nil || userID <= 0 {
		return biz.TokenClaims{}, biz.ErrInvalidToken
	}
	return biz.TokenClaims{UserID: userID, AuthVersion: claims.AuthVersion, SessionID: claims.SessionID}, nil
}

func (m *TokenManager) sign(claims jwtClaims) (string, error) {
	header, _ := json.Marshal(map[string]string{"alg": "HS256", "typ": "JWT"})
	payload, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}
	unsigned := encodeJWT(header) + "." + encodeJWT(payload)
	mac := hmac.New(sha256.New, m.secret)
	_, _ = mac.Write([]byte(unsigned))
	return unsigned + "." + encodeJWT(mac.Sum(nil)), nil
}

func (m *TokenManager) verify(token string) (jwtClaims, error) {
	var claims jwtClaims
	if len(m.secret) < 32 || len(token) > 4096 {
		if len(m.secret) < 32 {
			return claims, biz.ErrAuthenticationUnavailable
		}
		return claims, biz.ErrInvalidToken
	}
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return claims, biz.ErrInvalidToken
	}
	mac := hmac.New(sha256.New, m.secret)
	_, _ = mac.Write([]byte(parts[0] + "." + parts[1]))
	signature, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil || subtle.ConstantTimeCompare(signature, mac.Sum(nil)) != 1 {
		return claims, biz.ErrInvalidToken
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil || json.Unmarshal(payload, &claims) != nil || claims.Issuer != m.issuer || claims.ExpiresAt <= m.now().Unix() {
		return jwtClaims{}, biz.ErrInvalidToken
	}
	return claims, nil
}

func encodeJWT(value []byte) string { return base64.RawURLEncoding.EncodeToString(value) }

func errOrInvalidToken(err error) error {
	if err != nil {
		return err
	}
	return biz.ErrInvalidToken
}

func randomID() string {
	var raw [16]byte
	if _, err := rand.Read(raw[:]); err != nil {
		return strconv.FormatInt(time.Now().UnixNano(), 36)
	}
	return hex.EncodeToString(raw[:])
}
