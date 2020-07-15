package storage

import (
	"time"

	"github.com/patrickmn/go-cache"
	cache2 "github.com/pushm0v/go-zoho/cache"
)

type TokenStorage interface {
	SaveToken(accessToken, refreshToken string, expire int)
	AccessToken() string
	RefreshToken() string
}

type tokenStorage struct {
	cache *cache.Cache
}

func NewTokenStorage() TokenStorage {
	return &tokenStorage{
		cache: cache2.NewLocalCache(),
	}
}

func (t *tokenStorage) RefreshToken() string {
	return t.loadToken("refresh_token")
}

func (t *tokenStorage) AccessToken() string {
	return t.loadToken("access_token")
}

func (t *tokenStorage) loadToken(tokenName string) string {
	token, found := t.cache.Get(tokenName)
	if found {
		return token.(string)
	}

	return ""
}

func (t *tokenStorage) SaveToken(accessToken, refreshToken string, expire int) {
	t.cache.Set("access_token", accessToken, time.Duration(expire)*time.Second)
	t.cache.Set("refresh_token", refreshToken, time.Duration(expire)*time.Second)
}
