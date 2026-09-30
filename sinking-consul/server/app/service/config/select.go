package config

import (
	"errors"

	"server/app/enum/config_status"
	"server/app/enum/config_type"
	"server/app/repository/config"
)

// SelectAll 查询所有
func (s *service) SelectAll() (list []*config.Config, err error) {
	all, err := s.repository.SelectAll()
	if err != nil {
		return nil, err
	}
	return all, nil
}

// Select 获取数据
func (s *service) Select(where *config.SelectConfig, orderByField string, orderByType string, page int, pageSize int) (list []*config.Config, total int64, err error) {
	if where != nil {
		if where.Type != nil {
			if _, ok := config_type.Map()[*where.Type]; !ok {
				return nil, 0, errors.New("配置类型参数不合法")
			}
		}
		if where.Status != nil {
			if _, ok := config_status.Map()[*where.Status]; !ok {
				return nil, 0, errors.New("配置状态参数不合法")
			}
		}
	}
	return s.repository.Select(where, orderByField, orderByType, page, pageSize)
}
