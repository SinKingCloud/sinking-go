package config

import (
	"errors"
	"server/app/constant"
	"server/app/util"
	"sync"
)

// Service 单例对象
type Service struct {
}

// obj 单例对象
var (
	obj  *Service
	once sync.Once
)

// GetIns 获取单例
func GetIns() *Service {
	once.Do(func() {
		obj = &Service{}
	})
	return obj
}

// Group 获取group所有数据
func (s *Service) Group(group string) map[string]string {
	key := constant.CacheNameWithSysConfig + group
	value := util.Cache.Get(key)
	if value != nil {
		return value.(map[string]string)
	}
	temp := s.selectByGroup(group)
	util.Cache.SetWithExpire(key, temp, constant.CacheTimeWithSysConfig)
	return temp
}

// Set 设置数据
func (s *Service) Set(group string, key string, value string) error {
	lock := constant.LockConfigSet + group + key
	if !util.Cache.Lock(lock, constant.LockTimeConfigSet) {
		return errors.New("获取并发锁失败")
	}
	defer util.Cache.UnLock(lock)
	defer util.Cache.Delete(group)
	if n, e := s.countByKey(key); e == nil && n > 0 {
		return s.updateByKey(key, value)
	}
	return s.create(key, value)
}

// Get 获取数据
func (s *Service) Get(group string, key string) string {
	temp := s.Group(group)
	return temp[key]
}
