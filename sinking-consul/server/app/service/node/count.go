package node

import (
	"server/app/model"
	"server/app/util"
)

// CountByOnlineStatus 统计online_status数量
func (s *Service) CountByOnlineStatus(onlineStatus interface{}) (total int64, err error) {
	err = util.Database.Db.Model(&model.Node{}).Where("`online_status` = ? ", onlineStatus).Count(&total).Error
	return
}

// CountAll 统计status数量
func (s *Service) CountAll() (total int64, err error) {
	err = util.Database.Db.Model(&model.Node{}).Count(&total).Error
	return
}
