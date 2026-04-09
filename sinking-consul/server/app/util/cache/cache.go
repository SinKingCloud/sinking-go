package cache

import "time"

type Interface interface {
	Get(key string) interface{}
	Set(key string, value interface{})
	SetWithExpire(key string, value interface{}, expiration time.Duration)
	Delete(key string)
	Remember(key string, fn func() interface{}, expiration time.Duration) interface{}
	Lock(key string, expiration time.Duration) bool
	UnLock(key string)
	IsLock(key string) bool
}
