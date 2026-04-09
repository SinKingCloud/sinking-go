package config

import "server/app/model"

// FindByGroupAndName 通过group name查询父信息
func (s *service) FindByGroupAndName(group string, name string) (data *model.Config, err error) {
	return s.repository.FindByGroupAndName(group, name)
}
