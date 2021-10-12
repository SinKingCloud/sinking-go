package cache

import (
	goCache "github.com/patrickmn/go-cache"
	"time"
)

var system = goCache.New(
	3600*time.Second,
	time.Second,
)

func SetClient(c *goCache.Cache) {
	system = c
}

func Get(key string) interface{} {
	obj, found := system.Get(key)
	if found {
		return obj
	}
	return nil
}

func Set(key string, value interface{}) {
	system.SetDefault(key, value)
}

func SetWithExpire(key string, value interface{}, second time.Duration) {
	system.Set(key, value, second)
}

func Delete(key string) {
	system.Delete(key)
}

func Remember(key string, fun func() interface{}, second time.Duration) interface{} {
	value, exists := system.Get(key)
	if !exists {
		value = fun()
		SetWithExpire(key, value, second)
	}
	return value
}
