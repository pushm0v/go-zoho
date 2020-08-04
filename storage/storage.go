package storage

import "github.com/pushm0v/go-zoho/cache"

type Storage struct {
	Token TokenStorage
}

func NewStorage(cache cache.Cache) *Storage {
	return &Storage{
		Token: NewTokenStorage(cache),
	}
}
