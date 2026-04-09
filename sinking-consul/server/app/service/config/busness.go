package config

import (
	"server/app/enum/config_status"
	"server/app/model"
	"time"
)

// initialize 初始化服务
func (s *service) initialize() {
	s.once.Do(func() {
		all, e := s.repository.SelectAll()
		if e == nil && all != nil {
			for _, v := range all {
				info, err := s.FindByGroupAndName(v.Group, v.Name)
				if err == nil && info != nil {
					s.Set(info.Group, info.Name, info)
				}
			}
		}
	})
}

// GetGroup 获取集群组节点信息
func (s *service) GetGroup(group string) map[string]*model.Config {
	s.configLock.RLock()
	defer s.configLock.RUnlock()
	return s.configPool[group]
}

// Get 获取集群节点信息
func (s *service) Get(group string, key string) *model.Config {
	s.configLock.RLock()
	defer s.configLock.RUnlock()
	return s.configPool[group][key]
}

// Each 遍历集群信息
func (s *service) Each(group string, fun func(value *model.Config)) {
	s.configLock.RLock()
	defer s.configLock.RUnlock()
	if group == "*" || group == "" {
		for _, g := range s.configPool {
			for _, value := range g {
				fun(value)
			}
		}
	} else {
		if value, ok := s.configPool[group]; ok {
			for _, v := range value {
				fun(v)
			}
		}
	}
}

// Set 设置集群配置信息
func (s *service) Set(group string, key string, value *model.Config) {
	s.configLock.Lock()
	defer s.configLock.Unlock()
	if _, ok := s.configPool[group]; !ok {
		s.configPool[group] = make(map[string]*model.Config)
	}
	temp := s.configPool[group][key]
	if temp == nil || temp.Status != value.Status || temp.Hash != value.Hash {
		s.SetOperateTime(group)
	}
	s.configPool[group][key] = value
}

// Sets 批量设置集群配置信息
func (s *service) Sets(list []*model.Config) {
	s.configLock.Lock()
	defer s.configLock.Unlock()
	groups := make(map[string]int64)
	for _, v := range list {
		if _, ok := s.configPool[v.Group]; !ok {
			s.configPool[v.Group] = make(map[string]*model.Config)
		}
		if value, exists := s.configPool[v.Group][v.Name]; exists {
			if time.Time(value.UpdateTime).Unix() < time.Time(v.UpdateTime).Unix() {
				if v.Status != value.Status || v.Hash != value.Hash {
					groups[v.Group] = 1
				}
				s.configPool[v.Group][v.Name] = v
			}
		} else {
			groups[v.Group] = 1
			s.configPool[v.Group][v.Name] = v
		}
	}
	for group := range groups {
		s.SetOperateTime(group)
	}
}

// SetOperateTime 设置操作时间
func (s *service) SetOperateTime(group string) {
	s.configLastOperateTimeLock.Lock()
	defer s.configLastOperateTimeLock.Unlock()
	s.configLastOperateTime[group] = time.Now().UnixMicro()
}

// GetOperateTime 获取上次操作时间
func (s *service) GetOperateTime(group string) int64 {
	s.configLastOperateTimeLock.RLock()
	defer s.configLastOperateTimeLock.RUnlock()
	return s.configLastOperateTime[group]
}

// CheckIsChange 检查配置是否有变更
func (s *service) CheckIsChange(list []*model.Config) bool {
	s.configLock.Lock()
	defer s.configLock.Unlock()
	change := false
	for _, v := range list {
		if value, ok := s.configPool[v.Group]; ok {
			if config, exists := value[v.Name]; exists {
				if config.Hash != v.Hash || time.Time(config.UpdateTime).Unix() < time.Time(v.UpdateTime).Unix() {
					change = true
					break
				}
			} else {
				change = true
				break
			}
		} else {
			change = true
			break
		}
	}
	return change
}

// Delete 删除配置信息
func (s *service) Delete(group string, key string) {
	s.configLock.Lock()
	defer s.configLock.Unlock()
	s.SetOperateTime(group)
	if key == "*" || key == "" {
		if _, ok := s.configPool[group]; ok {
			delete(s.configPool, group)
		}
	} else {
		if value, ok := s.configPool[group]; ok {
			if _, ok = value[key]; ok {
				delete(s.configPool[group], key)
			}
		}
	}
}

// GetAllConfigs 获取本地配置信息
func (s *service) GetAllConfigs(group string, showContent bool, filterStatus bool) []*model.Config {
	s.configLock.RLock()
	defer s.configLock.RUnlock()
	count := 0
	for _, g := range s.configPool {
		count += len(g)
	}
	list := make([]*model.Config, 0, count)
	for _, g := range s.configPool {
		for _, value := range g {
			if filterStatus && value.Status == config_status.Stop {
				continue
			}
			if group != "" && group != "*" && value.Group != group {
				continue
			}
			if showContent {
				list = append(list, value)
			} else {
				list = append(list, &model.Config{
					Group:      value.Group,
					Name:       value.Name,
					Type:       value.Type,
					Hash:       value.Hash,
					CreateTime: value.CreateTime,
					UpdateTime: value.UpdateTime,
				})
			}
		}
	}
	return list
}
