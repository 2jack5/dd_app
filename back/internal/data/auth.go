package data

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"dd_web/back/internal/biz"
	"dd_web/back/internal/model"
	"gorm.io/gorm"
)

type AuthRepo struct{ data *Data }

func NewAuthRepo(data *Data) *AuthRepo { return &AuthRepo{data: data} }

func (r *AuthRepo) FindOrCreateWechatUser(ctx context.Context, identity biz.WechatIdentity) (biz.LoginUser, bool, error) {
	for attempt := 0; attempt < 2; attempt++ {
		user, created, err := r.findOrCreateWechatUser(ctx, identity)
		if !errors.Is(err, gorm.ErrDuplicatedKey) {
			return user, created, err
		}
	}
	return biz.LoginUser{}, false, fmt.Errorf("create wechat identity: %w", gorm.ErrDuplicatedKey)
}

func (r *AuthRepo) findOrCreateWechatUser(ctx context.Context, value biz.WechatIdentity) (result biz.LoginUser, created bool, err error) {
	err = r.data.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var identity model.AuthIdentity
		err := tx.Where("provider = ? AND app_id = ? AND subject = ?", value.Provider, value.AppID, value.OpenID).First(&identity).Error
		if err == nil {
			var user model.User
			if err := tx.First(&user, identity.UserID).Error; err != nil {
				return err
			}
			now := time.Now().UTC()
			if err := tx.Model(&user).Update("last_login_at", now).Error; err != nil {
				return err
			}
			user.LastLoginAt = &now
			result = loginUser(user)
			return nil
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		userNo, err := randomUserNo()
		if err != nil {
			return err
		}
		now := time.Now().UTC()
		user := model.User{UserNo: userNo, Nickname: "微信用户" + userNo[len(userNo)-6:], Status: model.UserStatusActive, AuthVersion: 1, ContentVersion: 1, LastLoginAt: &now}
		if err := tx.Create(&user).Error; err != nil {
			return err
		}
		var unionID *string
		if value.UnionID != "" {
			unionID = &value.UnionID
		}
		identity = model.AuthIdentity{UserID: user.ID, Provider: model.AuthIdentityProvider(value.Provider), AppID: value.AppID, Subject: value.OpenID, UnionID: unionID, BoundAt: now}
		if err := tx.Create(&identity).Error; err != nil {
			return err
		}
		result, created = loginUser(user), true
		return nil
	})
	return
}

func (r *AuthRepo) UserForAuthentication(ctx context.Context, userID int64) (biz.LoginUser, error) {
	var user model.User
	if err := r.data.db.WithContext(ctx).First(&user, userID).Error; err != nil {
		return biz.LoginUser{}, err
	}
	return loginUser(user), nil
}

func loginUser(user model.User) biz.LoginUser {
	return biz.LoginUser{ID: user.ID, UserNo: user.UserNo, Nickname: user.Nickname, Status: string(user.Status), AuthVersion: user.AuthVersion}
}

func randomUserNo() (string, error) {
	var raw [8]byte
	if _, err := rand.Read(raw[:]); err != nil {
		return "", fmt.Errorf("generate user number: %w", err)
	}
	return "DD" + hex.EncodeToString(raw[:]), nil
}
