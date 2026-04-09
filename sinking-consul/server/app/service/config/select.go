package config

import (
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
	return s.repository.Select(where, orderByField, orderByType, page, pageSize)
}
