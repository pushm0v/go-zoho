package storage

import (
	"fmt"
	"time"

	"github.com/pushm0v/go-zoho/common"

	cache2 "github.com/pushm0v/go-zoho/cache"
)

type TokenStorage interface {
	SaveToken(accessToken, refreshToken string, expire int)
	SaveAccessToken(accessToken string, expire int)
	AccessToken() string
	RefreshToken() string
	ExpireTime() float64
	IsTokenExpired() bool
}

type tokenStorage struct {
	cache cache2.Cache
}

func NewTokenStorage(cache cache2.Cache) TokenStorage {
	return &tokenStorage{
		cache: cache,
	}
}

func (t *tokenStorage) RefreshToken() string {
	return t.loadToken("refresh_token")
}

func (t *tokenStorage) AccessToken() string {
	return t.loadToken("access_token")
}

func (t *tokenStorage) ExpireTime() float64 {
	expTime, err := t.getExpireTime()
	if err != nil {
		return -1
	}
	timenow := time.Now().UTC()
	return expTime.Sub(timenow).Seconds()
}

func (t *tokenStorage) loadToken(tokenName string) string {
	token, err := t.cache.Get(tokenName)
	if err == nil && token != nil {
		return token.(string)
	}

	return ""
}

func (t *tokenStorage) SaveToken(accessToken, refreshToken string, expire int) {
	t.cache.Set("access_token", accessToken, expire)
	t.cache.Set("refresh_token", refreshToken, expire)
	timenow := time.Now().Add(time.Second * time.Duration(int64(expire))).UTC()
	t.cache.Set("expire_time", timenow.Format(common.TIME_FORMAT), -1)
}

func (t *tokenStorage) SaveAccessToken(accessToken string, expire int) {
	t.cache.Set("access_token", accessToken, expire)
	timenow := time.Now().Add(time.Second * time.Duration(int64(expire))).UTC()
	t.cache.Set("expire_time", timenow.Format(common.TIME_FORMAT), -1)
}

func (t *tokenStorage) getExpireTime() (expTime time.Time, err error) {
	timenow := time.Now().UTC()
	expireTimeString, err := t.cache.Get("expire_time")
	if err != nil {
		return timenow, err
	}

	if expireTimeString.(string) == "" {
		return timenow, fmt.Errorf("Expire time not found")
	}

	expTime, err = time.Parse(common.TIME_FORMAT, expireTimeString.(string))
	return
}

func (t *tokenStorage) IsTokenExpired() bool {
	timenow := time.Now().UTC()
	expTime, err := t.getExpireTime()
	if err != nil {
		return true
	}
	return expTime.Before(timenow)
}
