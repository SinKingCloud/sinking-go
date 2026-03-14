package config

import (
	"server/app/model"
	"server/app/util"
	"server/app/util/str"
	"time"
)

// UpdateByGroupAndName 通过group name更新
func (s *Service) UpdateByGroupAndName(keys []*model.Config, data map[string]interface{}) (err error) {
	data["update_time"] = str.DateTime(time.Now())
	var conditions [][]interface{}
	for _, key := range keys {
		if key.Group != "" && key.Name != "" {
			conditions = append(conditions, []interface{}{key.Group, key.Name})
		}
	}
	err = util.Database.Db.Model(&model.Config{}).Where("(`group`, `name`) IN (?)", conditions).Updates(data).Error
	if err == nil {
		list, err2 := s.SelectInGroupAndName(keys)
		if err2 == nil {
			s.Sets(list)
		}
	}
	return err
}
