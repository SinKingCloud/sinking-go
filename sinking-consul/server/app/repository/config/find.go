package config

import (
	"errors"
	"server/app/model"

	"gorm.io/gorm"
)

// FindByGroupAndName 通过group name查询父信息
func (r *Repository) FindByGroupAndName(group string, name string) (*model.Config, error) {
	var data *model.Config
	err := r.Database.Db.Model(&model.Config{}).Where("`group` = ? AND `name` = ?", group, name).First(&data).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if data.Group == "" {
		return nil, nil
	}
	return data, nil
}
