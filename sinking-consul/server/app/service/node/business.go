package node

import (
	"server/app/model"
	"time"
)

// Init 初始化服务
func (s *Service) Init() {
	nodeOnce.Do(func() {
		_ = s.UpdateAll(map[string]interface{}{"online_status": Offline})
		all, e := s.SelectAll()
		if e == nil && all != nil {
			for _, v := range all {
				s.Set(v.Group, v.Address, &Node{
					Node: &model.Node{
						Group:        v.Group,
						Name:         v.Name,
						Address:      v.Address,
						OnlineStatus: int(Offline),
						Status:       v.Status,
						LastHeart:    v.LastHeart,
						CreateTime:   v.CreateTime,
						UpdateTime:   v.UpdateTime,
					},
					IsLocal: false,
				})
			}
		}
	})
}

// GetGroup 获取集群组节点信息
func (s *Service) GetGroup(group string) map[string]*Node {
	nodeLock.RLock()
	defer nodeLock.RUnlock()
	return nodePool[group]
}

// Get 获取集群节点信息
func (s *Service) Get(group string, key string) *Node {
	nodeLock.RLock()
	defer nodeLock.RUnlock()
	return nodePool[group][key]
}

// Each 遍历集群信息
func (s *Service) Each(group string, fun func(value *Node)) {
	nodeLock.RLock()
	defer nodeLock.RUnlock()
	if group == "*" || group == "" {
		for _, g := range nodePool {
			for _, value := range g {
				fun(value)
			}
		}
	} else {
		if value, ok := nodePool[group]; ok {
			for _, v := range value {
				fun(v)
			}
		}
	}
}

// Set 设置集群节点信息
func (s *Service) Set(group string, key string, value *Node) {
	nodeLock.Lock()
	defer nodeLock.Unlock()
	if _, ok := nodePool[group]; !ok {
		nodePool[group] = make(map[string]*Node)
	}
	nodePool[group][key] = value
}

// Sets 批量设置集群节点信息
func (s *Service) Sets(list []*Node) {
	nodeLock.Lock()
	defer nodeLock.Unlock()
	for _, v := range list {
		if _, ok := nodePool[v.Group]; !ok {
			nodePool[v.Group] = make(map[string]*Node)
		}
		if value, ok := nodePool[v.Group][v.Address]; ok {
			v.IsLocal = value.IsLocal
		}
		nodePool[v.Group][v.Address] = v
	}
}

// Delete 删除集群节点信息
func (s *Service) Delete(group string, key string) {
	nodeLock.Lock()
	defer nodeLock.Unlock()
	if key == "*" || key == "" {
		if _, ok := nodePool[group]; ok {
			delete(nodePool, group)
		}
	} else {
		if value, ok := nodePool[group]; ok {
			if _, ok = value[key]; ok {
				delete(nodePool[group], key)
			}
		}
	}
}

// Register 注册节点信息
func (s *Service) Register(group string, name string, address string) {
	data := s.Get(group, address)
	if data == nil {
		s.Set(group, address, &Node{
			Node: &model.Node{
				Group:        group,
				Name:         name,
				Address:      address,
				OnlineStatus: int(Online),
				Status:       int(Normal),
				LastHeart:    time.Now().Unix(),
			},
			IsLocal: true,
		})
	} else {
		data.IsLocal = true
		data.OnlineStatus = int(Online)
		data.LastHeart = time.Now().Unix()
	}
}

// GetLocalNodes 获取本地服务信息
func (s *Service) GetLocalNodes() []*model.Node {
	nodeLock.RLock()
	defer nodeLock.RUnlock()
	count := 0
	for _, g := range nodePool {
		count += len(g)
	}
	list := make([]*model.Node, 0, count)
	for _, g := range nodePool {
		for _, value := range g {
			if value.IsLocal {
				list = append(list, value.Node)
			}
		}
	}
	return list
}
