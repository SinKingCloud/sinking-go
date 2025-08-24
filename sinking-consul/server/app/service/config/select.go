package config

import (
	"server/app/model"
	"server/app/util"
)

// SelectAll 查询所有
func (s *Service) SelectAll() (list []*Config, err error) {
	err = util.Database.Db.Model(&model.Config{}).Find(&list).Error
	return list, err
}

// SelectInGroupAndName 根据group name查询
func (s *Service) SelectInGroupAndName(keys []*model.Config) (list []*Config, err error) {
	var conditions [][]interface{}
	for _, key := range keys {
		if key.Group != "" && key.Name != "" {
			conditions = append(conditions, []interface{}{key.Group, key.Name})
		}
	}
	err = util.Database.Db.Model(&model.Node{}).Where("(`group`, `name`) IN (?)", conditions).Find(&list).Error
	return list, err
}

// Select 获取数据
func (s *Service) Select(where map[string]string, orderByField string, orderByType string, page int, pageSize int) (list []*SelectConfig, total int64, err error) {
	offset := pageSize * (page - 1)
	query := util.Database.Db.Model(&model.Config{})
	if where["group"] != "" {
		query.Where("`group` = ?", where["group"])
	}
	if where["name"] != "" {
		query.Where("`name` = ?", where["name"])
	}
	if where["type"] != "" {
		query.Where("`type` = ?", where["type"])
	}
	if where["hash"] != "" {
		query.Where("`hash` = ?", where["hash"])
	}
	if where["content"] != "" {
		query.Where("`content` LIKE ?", "%"+where["content"]+"%")
	}
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
