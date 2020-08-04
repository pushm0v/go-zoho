package storage

import (
	cache2 "github.com/pushm0v/go-zoho/cache"
)

type TokenStorage interface {
	SaveToken(accessToken, refreshToken string, expire int)
	AccessToken() string
	RefreshToken() string
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

func (t *tokenStorage) loadToken(tokenName string) string {
	token, err := t.cache.Get(tokenName)
	if err == nil && token != nil {
		return token.(string)
	}

	return ""
}

func (t *tokenStorage) SaveToken(accessToken, refreshToken string, expire int) {
	//t.cache.Set("access_token", accessToken, time.Duration(expire)*time.Second)
	t.cache.Set("access_token", accessToken, expire)
	t.cache.Set("refresh_token", refreshToken, expire)
}
