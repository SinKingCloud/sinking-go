package config

import (
	"server/app/model"
	"server/app/util/str"
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

// SelectConfig 搜索结果
type SelectConfig struct {
	Group      string       `gorm:"column:group" json:"group"`
	Name       string       `gorm:"column:name" json:"name"`
	Type       string       `gorm:"column:type" json:"type"`
	Hash       string       `gorm:"column:hash" json:"hash"`
	Status     int          `gorm:"column:status" json:"status"`
	CreateTime str.DateTime `gorm:"column:create_time" json:"create_time"`
	UpdateTime str.DateTime `gorm:"column:update_time" json:"update_time"`
}
