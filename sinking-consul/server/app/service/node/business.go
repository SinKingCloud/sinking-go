package node

import (
	"server/app/enum/node_online_status"
	"server/app/enum/node_status"
	"server/app/model"
	repositoryNode "server/app/repository/node"
	"time"
)

// initialize 初始化服务
func (s *service) initialize() {
	s.once.Do(func() {
		_ = s.UpdateAll(&repositoryNode.UpdateNode{OnlineStatus: node_online_status.Offline})
		all, e := s.repository.SelectAll()
		if e == nil && all != nil {
			for _, v := range all {
				s.Set(v.Group, v.Address, &Node{
					Node: &model.Node{
						Group:        v.Group,
						Name:         v.Name,
						Address:      v.Address,
						OnlineStatus: node_online_status.Offline,
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
func (s *service) GetGroup(group string) map[string]*Node {
	s.nodeLock.RLock()
	defer s.nodeLock.RUnlock()
	return s.nodePool[group]
}

// Get 获取集群节点信息
func (s *service) Get(group string, key string) *Node {
	s.nodeLock.RLock()
	defer s.nodeLock.RUnlock()
	return s.nodePool[group][key]
}

// Each 遍历集群信息
func (s *service) Each(group string, fun func(value *Node)) {
	s.nodeLock.RLock()
	defer s.nodeLock.RUnlock()
	if group == "*" || group == "" {
		for _, g := range s.nodePool {
			for _, value := range g {
				fun(value)
			}
		}
	} else {
		if value, ok := s.nodePool[group]; ok {
			for _, v := range value {
				fun(v)
			}
		}
	}
}

// Set 设置集群节点信息
func (s *service) Set(group string, key string, value *Node) {
	s.nodeLock.Lock()
	defer s.nodeLock.Unlock()
	if _, ok := s.nodePool[group]; !ok {
		s.nodePool[group] = make(map[string]*Node)
	}
	temp := s.nodePool[group][key]
	if temp == nil || temp.Status != value.Status || temp.OnlineStatus != value.OnlineStatus {
		s.SetOperateTime(group)
	}
	s.nodePool[group][key] = value
}

// Sets 批量设置集群节点信息
func (s *service) Sets(list []*Node) {
	s.nodeLock.Lock()
	defer s.nodeLock.Unlock()
	groups := make(map[string]int64)
	for _, v := range list {
		if _, ok := s.nodePool[v.Group]; !ok {
			s.nodePool[v.Group] = make(map[string]*Node)
		}
		if value, ok := s.nodePool[v.Group][v.Address]; ok {
			v.IsLocal = value.IsLocal
			if v.Status != value.Status || v.OnlineStatus != value.OnlineStatus {
				groups[v.Group] = 1
			}
		} else {
			groups[v.Group] = 1
		}
		s.nodePool[v.Group][v.Address] = v
	}
	for group := range groups {
		s.SetOperateTime(group)
	}
}

// SetOperateTime 设置上次操作时间
func (s *service) SetOperateTime(group string) {
	s.nodeLastOperateTimeLock.Lock()
	defer s.nodeLastOperateTimeLock.Unlock()
	s.nodeLastOperateTime[group] = time.Now().UnixMicro()
}

// GetOperateTime 获取上次操作时间
func (s *service) GetOperateTime(group string) int64 {
	s.nodeLastOperateTimeLock.RLock()
	defer s.nodeLastOperateTimeLock.RUnlock()
	return s.nodeLastOperateTime[group]
}

// Delete 删除集群节点信息
func (s *service) Delete(group string, key string) {
	s.nodeLock.Lock()
	defer s.nodeLock.Unlock()
	s.SetOperateTime(group)
	if key == "*" || key == "" {
		if _, ok := s.nodePool[group]; ok {
			delete(s.nodePool, group)
		}
	} else {
		if value, ok := s.nodePool[group]; ok {
			if _, ok = value[key]; ok {
				delete(s.nodePool[group], key)
			}
		}
	}
}

// Register 注册节点信息
func (s *service) Register(group string, name string, address string) {
	data := s.Get(group, address)
	if data == nil {
		s.Set(group, address, &Node{
			Node: &model.Node{
				Group:        group,
				Name:         name,
				Address:      address,
				OnlineStatus: node_online_status.Online,
				Status:       node_status.Normal,
				LastHeart:    time.Now().Unix(),
			},
			IsLocal: true,
		})
	} else {
		if data.OnlineStatus != node_online_status.Online {
			s.SetOperateTime(group)
		}
		data.IsLocal = true
		data.OnlineStatus = node_online_status.Online
		data.LastHeart = time.Now().Unix()
	}
}

// GetLocalNodes 获取本地服务信息
func (s *service) GetLocalNodes() []*model.Node {
	s.nodeLock.RLock()
	defer s.nodeLock.RUnlock()
	count := 0
	for _, g := range s.nodePool {
		count += len(g)
	}
	list := make([]*model.Node, 0, count)
	for _, g := range s.nodePool {
		for _, value := range g {
			if value.IsLocal {
				list = append(list, value.Node)
			}
		}
	}
	return list
}

// GetAllOnlineNodes 获取正常服务数据信息
func (s *service) GetAllOnlineNodes(group string) []*model.Node {
	s.nodeLock.RLock()
	defer s.nodeLock.RUnlock()
	count := 0
	for _, g := range s.nodePool {
		count += len(g)
	}
	list := make([]*model.Node, 0, count)
	for _, g := range s.nodePool {
		for _, value := range g {
			if value.OnlineStatus == node_online_status.Online && value.Status == node_status.Normal && value.Group == group {
				list = append(list, value.Node)
			}
		}
	}
	return list
}
