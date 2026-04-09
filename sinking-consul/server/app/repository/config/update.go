package config

import (
	"server/app/model"
	"server/app/util/str"
	"time"
)

// UpdateByGroupAndName 根据 code 更新权限
func (r *Repository) UpdateByGroupAndName(keys []*model.Config, config *UpdateConfig) error {
	if config == nil {
		return nil
	}
	updates := make(map[string]interface{})
	if config.Group != nil {
		updates["group"] = config.Group
	}
	if config.Name != nil {
		updates["name"] = config.Name
	}
	if config.Type != nil {
		updates["type"] = config.Type
	}
	if config.Hash != nil {
		updates["hash"] = config.Hash
	}
	if config.Content != nil {
		updates["content"] = config.Content
	}
	if config.Status != nil {
		updates["status"] = config.Status
	}
	if len(updates) == 0 {
		return nil
	}
	var conditions [][]interface{}
	for _, key := range keys {
		if key.Group != "" && key.Name != "" {
			conditions = append(conditions, []interface{}{key.Group, key.Name})
		}
	}
	updates["update_time"] = str.DateTime(time.Now())
	return r.Database.Db.Model(&model.Config{}).Where("(`group`, `name`) IN (?)", conditions).Updates(updates).Error
}
