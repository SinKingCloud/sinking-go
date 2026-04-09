package config

import (
	"server/app/model"
)

// CountByStatus 统计status数量
func (r *Repository) CountByStatus(status int) (total int64, err error) {
	err = r.Database.Db.Model(&model.Config{}).Where("`status` = ? ", status).Count(&total).Error
	return
}

// CountAll 统计status数量
func (r *Repository) CountAll() (total int64, err error) {
	err = r.Database.Db.Model(&model.Config{}).Count(&total).Error
	return
}
