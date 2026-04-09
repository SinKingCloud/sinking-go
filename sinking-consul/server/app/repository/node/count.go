package node

import (
	"server/app/model"
)

// CountByOnlineStatus 统计online_status数量
func (r *Repository) CountByOnlineStatus(onlineStatus int) (total int64, err error) {
	err = r.Database.Db.Model(&model.Node{}).Where("`online_status` = ? ", onlineStatus).Count(&total).Error
	return
}

// CountAll 统计status数量
func (r *Repository) CountAll() (total int64, err error) {
	err = r.Database.Db.Model(&model.Node{}).Count(&total).Error
	return
}
