package log

import (
	"server/app/model"
	"server/app/util"
)

// Select 获取数据
func (s *Service) Select(where map[string]string, orderByField string, orderByType string, page int, pageSize int) (list []*model.Log, total int64, err error) {
	offset := pageSize * (page - 1)
	query := util.Database.Db.Model(&model.Log{})
	if where["type"] != "" {
		query.Where("`type` = ?", where["type"])
	}
	if where["ip"] != "" {
		query.Where("`ip` like ?", "%"+where["ip"]+"%")
	}
	if where["title"] != "" {
		query.Where("`user` like ?", "%"+where["title"]+"%")
	}
	if where["content"] != "" {
		query.Where("`content` like ?", "%"+where["content"]+"%")
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
