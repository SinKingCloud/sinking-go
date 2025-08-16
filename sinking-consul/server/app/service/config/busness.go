package config

import (
	"server/app/model"
	"time"
)

func (s *Service) Init() {
	configOnce.Do(func() {
		all, e := s.SelectAll()
		if e == nil && all != nil {
			for _, v := range all {
				s.Set(v.Group, v.Name, &Config{
					Config: &model.Config{
						Group:      v.Group,
						Name:       v.Name,
						Type:       v.Type,
						Hash:       v.Hash,
						Content:    v.Content,
						CreateTime: v.CreateTime,
						UpdateTime: v.UpdateTime,
					},
				})
			}
		}
	})
}

// GetGroup 获取集群组节点信息
func (s *Service) GetGroup(group string) map[string]*Config {
	configLock.RLock()
	defer configLock.RUnlock()
	return configPool[group]
}

// Get 获取集群节点信息
func (s *Service) Get(group string, key string) *Config {
	configLock.RLock()
	defer configLock.RUnlock()
	return configPool[group][key]
}

// Each 遍历集群信息
func (s *Service) Each(group string, fun func(value *Config)) {
	configLock.RLock()
	defer configLock.RUnlock()
	if group == "*" || group == "" {
		for _, g := range configPool {
			for _, value := range g {
				fun(value)
			}
		}
	} else {
		if value, ok := configPool[group]; ok {
			for _, v := range value {
				fun(v)
			}
		}
	}
}

// Set 设置集群配置信息
func (s *Service) Set(group string, key string, value *Config) {
	configLock.Lock()
	defer configLock.Unlock()
	if _, ok := configPool[group]; !ok {
		configPool[group] = make(map[string]*Config)
	}
	configPool[group][key] = value
}

// Sets 批量设置集群配置信息
func (s *Service) Sets(list []*Config) {
	configLock.Lock()
	defer configLock.Unlock()
	for _, v := range list {
		if _, ok := configPool[v.Group]; !ok {
			configPool[v.Group] = make(map[string]*Config)
		}
		if value, exists := configPool[v.Group][v.Name]; exists {
			if time.Time(value.UpdateTime).Unix() < time.Time(v.UpdateTime).Unix() {
				configPool[v.Group][v.Name] = v
			}
		} else {
			configPool[v.Group][v.Name] = v
		}
	}
}

// SetOperateTime 设置操作时间
func (s *Service) SetOperateTime(group string) {
	configLock.Lock()
	defer configLock.Unlock()
	configLastOperateTime[group] = time.Now().Unix()
}

// GetOperateTime 获取上次操作时间
func (s *Service) GetOperateTime(group string) int64 {
	configLock.Lock()
	defer configLock.Unlock()
	return configLastOperateTime[group]
}

// CheckIsChange 检查配置是否有变更
func (s *Service) CheckIsChange(list []*Config) bool {
	configLock.Lock()
	defer configLock.Unlock()
	change := false
	for _, v := range list {
		if value, ok := configPool[v.Group]; ok {
			if config, exists := value[v.Name]; exists {
				if config.Hash != v.Hash || time.Time(config.UpdateTime).Unix() < time.Time(v.UpdateTime).Unix() {
					change = true
					break
				}
			} else {
				change = true
				break
			}
		}
	}
	return change
}

// Delete 删除集群节点信息
func (s *Service) Delete(group string, key string) {
	configLock.Lock()
	defer configLock.Unlock()
	if key == "*" || key == "" {
		if _, ok := configPool[group]; ok {
			delete(configPool, group)
		}
	} else {
		if value, ok := configPool[group]; ok {
			if _, ok = value[key]; ok {
				delete(configPool[group], key)
			}
		}
	}
}

// GetAllConfigs 获取本地配置信息
func (s *Service) GetAllConfigs(group string, showContent bool) []*Config {
	configLock.RLock()
	defer configLock.RUnlock()
	count := 0
	for _, g := range configPool {
		count += len(g)
	}
	list := make([]*Config, 0, count)
	for _, g := range configPool {
		for _, value := range g {
			if group != "" && group != "*" {
				if value.Group != group {
					continue
				}
			}
			if showContent {
				list = append(list, value)
			} else {
				list = append(list, &Config{
					Config: &model.Config{
						Group:      value.Group,
						Name:       value.Name,
						Type:       value.Type,
						Hash:       value.Hash,
						CreateTime: value.CreateTime,
						UpdateTime: value.UpdateTime,
					},
				})
			}
		}
	}
	return list
}
