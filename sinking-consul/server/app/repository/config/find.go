package config

import (
	"errors"
	"server/app/model"

	"gorm.io/gorm"
)

// FindByGroupAndName 通过分组和名称查询配置信息
func (r *Repository) FindByGroupAndName(group string, name string) (*model.Config, error) {
	var data *model.Config
	err := r.Database.Db.
		Model(&model.Config{}).
		Where("`group` = ? AND `name` = ?", group, name).
		First(&data).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if data == nil || data.Group == "" || data.Name == "" {
		return nil, nil
	}
	return data, nil
}
