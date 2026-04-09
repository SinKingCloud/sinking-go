package node

import (
	"server/app/model"
	"sync"
)

// service 单例对象
type service struct {
}

var (
	//实例对象
	obj = &service{}
	//原子锁
	nodeOnce = &sync.Once{}
	// 节点池 组[节点地址]
	nodePool                = make(map[string]map[string]*Node)
	nodeLock                = &sync.RWMutex{}
	nodeLastOperateTimeLock = &sync.RWMutex{}
	nodeLastOperateTime     = make(map[string]int64)
)

// GetIns 获取单例
func GetIns() *service {
	return obj
}

// Node 集群列表
type Node struct {
	*model.Node
	IsLocal bool `gorm:"-" json:"is_local"`
}
