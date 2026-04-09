package cluster

import (
	"server/app/model"
	"server/app/util"
)

// CountByStatus 统计status数量
func (s *service) CountByStatus(status interface{}) (total int64, err error) {
	err = util.Database.Db.Model(&model.Cluster{}).Where("`status` = ? ", status).Count(&total).Error
	return
}

// CountAll 统计status数量
func (s *service) CountAll() (total int64, err error) {
	err = util.Database.Db.Model(&model.Cluster{}).Count(&total).Error
	return
}
