package config

import (
	"server/app/model"
)

// SelectAll 查询所有
func (r *Repository) SelectAll() (list []*Config, err error) {
	err = r.Database.Db.Model(&model.Config{}).Find(&list).Error
	return list, err
}

// SelectInGroupAndName 根据group name查询
func (r *Repository) SelectInGroupAndName(keys []*model.Config) (list []*Config, err error) {
	var conditions [][]interface{}
	for _, key := range keys {
		if key.Group != "" && key.Name != "" {
			conditions = append(conditions, []interface{}{key.Group, key.Name})
		}
	}
	err = r.Database.Db.Model(&model.Config{}).Where("(`group`, `name`) IN (?)", conditions).Find(&list).Error
	return list, err
}

// Select 查询数据
func (r *Repository) Select(where *SelectConfig, orderByField string, orderByType string, page int, pageSize int) (list []*Config, total int64, err error) {
	offset := pageSize * (page - 1)
	query := r.Database.Db.Model(&model.Config{})
	if where != nil {
		if where.Group != "" {
			query = query.Where("`group` = ?", where.Group)
		}
		if where.Name != "" {
			query = query.Where("`name` = ?", where.Name)
		}
		if where.Type != "" {
			query = query.Where("`type` = ?", where.Type)
		}
		if where.Hash != "" {
			query = query.Where("`hash` = ?", where.Hash)
		}
		if where.Content != "" {
			query = query.Where("`content` LIKE ?", "%"+where.Content+"%")
		}
		if where.Status != "" {
			query = query.Where("`status` = ?", where.Status)
		}
		if where.CreateTimeStart != "" {
			query = query.Where("`create_time` >= ?", where.CreateTimeStart)
		}
		if where.CreateTimeEnd != "" {
			query = query.Where("`create_time` <= ?", where.CreateTimeEnd)
		}
		if where.UpdateTimeStart != "" {
			query = query.Where("`update_time` >= ?", where.UpdateTimeStart)
		}
		if where.UpdateTimeEnd != "" {
			query = query.Where("`update_time` <= ?", where.UpdateTimeEnd)
		}
	}
	err = query.Count(&total).Limit(pageSize).Offset(offset).Order(orderByField + " " + orderByType).Find(&list).Error
	return
}
