package setting

import "errors"

// GetString 获取配置
func (s *service) GetString(key string) string {
	return s.conf.GetString(key)
}

// GetInt 获取配置
func (s *service) GetInt(key string) int {
	return s.conf.GetInt(key)
}

// GetStringSlice 获取配置
func (s *service) GetStringSlice(key string) []string {
	return s.conf.GetStringSlice(key)
}

// GetByGroup 通过group获取配置
func (s *service) GetByGroup(key string) interface{} {
	return s.conf.AllSettings()[key]
}

// Set 设置数据
func (s *service) Set(key string, value string) error {
	s.conf.Set(key, value)
	return s.conf.WriteConfig()
}

// Sets 批量设置数据
func (s *service) Sets(list []*Config) error {
	if list == nil {
		return errors.New("设置数据不能为空")
	}
	for _, v := range list {
		s.conf.Set(v.Key, v.Value)
	}
	return s.conf.WriteConfig()
}
