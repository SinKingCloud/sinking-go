package config

import (
	"server/app/model"
	"sync"
)

// Service 单例对象
type Service struct {
}

var (
	//实例对象
	obj = &Service{}
	//原子锁
	configOnce = &sync.Once{}
	// 节点池 组[配置名]
	configPool            = make(map[string]map[string]*Config)
	configLock            = sync.RWMutex{}
	configLastOperateTime = make(map[string]int64)
)

// GetIns 获取单例
func GetIns() *Service {
	return obj
}

// Config 配置信息
type Config struct {
	*model.Config
}
