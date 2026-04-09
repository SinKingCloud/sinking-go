package config

import (
	"server/app/model"
)

// DeleteByGroupAndName 通过集群和名称删除
func (r *Repository) DeleteByGroupAndName(keys []*model.Config) error {
	var conditions [][]interface{}
	for _, key := range keys {
		if key.Group != "" && key.Name != "" {
			conditions = append(conditions, []interface{}{key.Group, key.Name})
		}
	}
	return r.Database.Db.Where("(`group`, `name`) IN (?)", conditions).Delete(&model.Config{}).Error
}
