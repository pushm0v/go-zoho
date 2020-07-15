package cache

import (
	"github.com/patrickmn/go-cache"
)

func NewLocalCache() *cache.Cache {
	return cache.New(-1, -1) //Never Expire
}
