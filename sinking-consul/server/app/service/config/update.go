package config

import (
	"errors"

	"server/app/enum/config_status"
	"server/app/enum/config_type"
	"server/app/model"
	repositoryConfig "server/app/repository/config"
)

// UpdateByGroupAndName 通过group name更新
func (s *service) UpdateByGroupAndName(keys []*model.Config, data *repositoryConfig.UpdateConfig) (err error) {
	if data == nil {
		return errors.New("更新数据不能为空")
	}
	if data.Type != nil {
		if _, ok := config_type.Map()[*data.Type]; !ok {
			return errors.New("配置类型不合法")
		}
	}
	if data.Status != nil {
		if _, ok := config_status.Map()[*data.Status]; !ok {
			return errors.New("配置状态不合法")
		}
	}
	err = s.repository.UpdateByGroupAndName(keys, data)
	if err == nil {
		list, err2 := s.repository.SelectInGroupAndName(keys)
		if err2 == nil {
			s.Sets(list)
		}
	}
	return err
}
