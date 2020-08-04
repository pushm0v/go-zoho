package cache

type Set func(string, interface{}, int) error
type Get func(string) (interface{}, error)

type Option struct {
	SetFunc Set
	GetFunc Get
}

type externalCache struct {
	option Option
}

func NewExternalCache(option Option) Cache {
	return &externalCache{
		option: option,
	}
}

func (ec *externalCache) Set(key string, value interface{}, expire int) error {
	return ec.option.SetFunc(key, value, expire)
}

func (ec *externalCache) Get(key string) (value interface{}, err error) {
	return ec.option.GetFunc(key)
}
