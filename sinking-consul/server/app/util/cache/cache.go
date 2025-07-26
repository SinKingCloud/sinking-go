package cache

import (
	goCache "github.com/patrickmn/go-cache"
	"sync"
	"time"
)

func NewCache(defaultExpiration, cleanupInterval time.Duration) *Cache {
	c := &Cache{
		memory: goCache.New(
			defaultExpiration,
			cleanupInterval,
		),
		locks:           make(map[string]*sync.Mutex),
		lockExpirations: make(map[string]time.Time),
		mu:              sync.Mutex{},
	}
	c.once.Do(func() {
		go func() {
			for {
				c.CleanExpiredLock()
				time.Sleep(cleanupInterval)
			}
		}()
	})
	return c
}

type Cache struct {
	memory          *goCache.Cache
	locks           map[string]*sync.Mutex
	lockExpirations map[string]time.Time
	mu              sync.Mutex
	once            sync.Once
}

func (c *Cache) Get(key string) interface{} {
	obj, found := c.memory.Get(key)
	if found {
		return obj
	}
	return nil
}

func (c *Cache) Set(key string, value interface{}) {
	c.memory.SetDefault(key, value)
}

func (c *Cache) SetWithExpire(key string, value interface{}, second time.Duration) {
	c.memory.Set(key, value, second)
}

func (c *Cache) Delete(key string) {
	c.memory.Delete(key)
}

func (c *Cache) Remember(key string, fun func() interface{}, second time.Duration) interface{} {
	value, exists := c.memory.Get(key)
	if !exists {
		value = fun()
		c.SetWithExpire(key, value, second)
	}
	return value
}

// Lock 尝试获取锁，如果锁已被占用或者已过期则返回false，否则返回true并获取锁
func (c *Cache) Lock(key string, expiration time.Duration) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	expirationTime, ok := c.lockExpirations[key]
	if ok && expirationTime.After(time.Now()) {
		return false
	}
	if _, ok = c.locks[key]; !ok {
		c.locks[key] = &sync.Mutex{}
	}
	c.lockExpirations[key] = time.Now().Add(expiration)
	c.locks[key].Lock()
	return true
}

// UnLock 释放锁
func (c *Cache) UnLock(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if mutex, ok := c.locks[key]; ok {
		mutex.Unlock()
		delete(c.locks, key)
		delete(c.lockExpirations, key)
	}
}

// CleanExpiredLock 清理过期的锁
func (c *Cache) CleanExpiredLock() {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := time.Now()
	for key, expiration := range c.lockExpirations {
		if expiration.Before(now) {
			if mutex, ok := c.locks[key]; ok {
				mutex.Unlock()
			}
			delete(c.locks, key)
			delete(c.lockExpirations, key)
		}
	}
}
