package config

import (
	"server/app/model"
	"server/app/util"
)

// CountByStatus 统计status数量
func (s *service) CountByStatus(status interface{}) (total int64, err error) {
	err = util.Database.Db.Model(&model.Config{}).Where("`status` = ? ", status).Count(&total).Error
	return
}

// CountAll 统计status数量
func (s *service) CountAll() (total int64, err error) {
	err = util.Database.Db.Model(&model.Config{}).Count(&total).Error
	return
}
