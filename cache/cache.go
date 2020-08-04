package cache

type Cache interface {
	Set(key string, value interface{}, expire int) error
	Get(key string) (value interface{}, err error)
}

type CacheOption func(*CacheStruct)

type CacheStruct struct {
	WithCache Cache
}

func NewCache(opts ...CacheOption) Cache {
	var c = &CacheStruct{}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c *CacheStruct) Set(key string, value interface{}, expire int) (err error) {
	return c.WithCache.Set(key, value, expire)
}

func (c *CacheStruct) Get(key string) (value interface{}, err error) {
	return c.WithCache.Get(key)
}

func WithLocalCache() CacheOption {
	return func(a *CacheStruct) {
		a.WithCache = NewLocalCache()
	}
}

func WithExternalCache(externalOption Option) CacheOption {
	return func(a *CacheStruct) {
		a.WithCache = NewExternalCache(externalOption)
	}
}
