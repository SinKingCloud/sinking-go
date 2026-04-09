package node

import (
	"server/app/model"
	repositoryNode "server/app/repository/node"
	"server/app/util/cache"
	"sync"
)

// Service 节点服务接口
type Service interface {
	SelectInAddress(addresses []string) ([]*model.Node, error)
	Select(where *repositoryNode.SelectNode, orderByField string, orderByType string, page int, pageSize int) ([]*model.Node, int64, error)
	CountByOnlineStatus(onlineStatus int) (int64, error)
	CountAll() (int64, error)
	Save() error
	Each(group string, fun func(value *Node))
	Sets(list []*Node)
	SetOperateTime(group string)
	GetOperateTime(group string) int64
	Register(group string, name string, address string)
	GetLocalNodes() []*model.Node
	GetAllOnlineNodes(group string) []*model.Node
	DeleteByAddress(addresses []string) error
	UpdateByAddresses(addresses []string, data *repositoryNode.UpdateNode) error
	Create(data *model.Node) error
}

// service 节点服务
type service struct {
	repository repositoryNode.Interface
	cache      cache.Interface

	once                    sync.Once
	nodePool                map[string]map[string]*Node
	nodeLock                sync.RWMutex
	nodeLastOperateTime     map[string]int64
	nodeLastOperateTimeLock sync.RWMutex
}

// NewService 创建节点服务
func NewService(repository repositoryNode.Interface, cache cache.Interface) *service {
	s := &service{
		repository:          repository,
		cache:               cache,
		nodePool:            make(map[string]map[string]*Node),
		nodeLastOperateTime: make(map[string]int64),
	}
	s.initialize()
	return s
}
