package cache

import (
	"fmt"
	"time"

	"github.com/patrickmn/go-cache"
)

type localCache struct {
	gocache *cache.Cache
}

func NewLocalCache() Cache {
	return &localCache{
		gocache: cache.New(-1, -1), //Never Expire
	}
}

func (lc *localCache) Set(key string, value interface{}, expire int) error {
	lc.gocache.Set(key, value, time.Duration(expire)*time.Second)

	return nil
}

func (lc *localCache) Get(key string) (value interface{}, err error) {
	val, found := lc.gocache.Get(key)
	if found {
		return val, nil
	}

	return nil, fmt.Errorf("Not found")
}
