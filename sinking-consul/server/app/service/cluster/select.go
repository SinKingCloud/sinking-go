package cluster

import (
	"server/app/model"
	"server/app/util"
)

// Select 获取数据
func (s *Service) Select(where map[string]string, orderByField string, orderByType string, page int, pageSize int) (list []*Cluster, total int64, err error) {
	offset := pageSize * (page - 1)
	query := util.Database.Db.Model(&model.Cluster{})
	if where["status"] != "" {
		query.Where("`status` = ?", where["status"])
	}
	if where["create_time_start"] != "" {
		query.Where("`create_time` >= ?", where["create_time_start"])
	}
	if where["create_time_end"] != "" {
		query.Where("`create_time` <= ?", where["create_time_end"])
	}
	if where["update_time_start"] != "" {
		query.Where("`update_time` >= ?", where["update_time_start"])
	}
	if where["update_time_end"] != "" {
		query.Where("`update_time` <= ?", where["update_time_end"])
	}
	err = query.Count(&total).Limit(pageSize).Offset(offset).Order(orderByField + " " + orderByType).Find(&list).Error
	return
}
