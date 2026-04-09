package config

import (
	"server/app/model"
	repositoryConfig "server/app/repository/config"
)

// UpdateByGroupAndName 通过group name更新
func (s *service) UpdateByGroupAndName(keys []*model.Config, data *repositoryConfig.UpdateConfig) (err error) {
	err = s.repository.UpdateByGroupAndName(keys, data)
	if err == nil {
		list, err2 := s.repository.SelectInGroupAndName(keys)
		if err2 == nil {
			s.Sets(list)
		}
	}
	return err
}
