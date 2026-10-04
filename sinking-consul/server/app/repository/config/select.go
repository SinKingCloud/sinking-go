package config

import (
	"server/app/model"
	"server/app/util/page"
)

// SelectAll 查询所有
func (r *Repository) SelectAll() (list []*Config, err error) {
	err = r.Database.Db.Model(&model.Config{}).Find(&list).Error
	return list, err
}

// SelectInGroupAndName 根据group name查询
func (r *Repository) SelectInGroupAndName(keys []*model.Config) (list []*model.Config, err error) {
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
func (r *Repository) Select(where *SelectConfig, queryPage *page.Query) (*page.Result[*Config], error) {
	query := r.Database.Db.Model(&model.Config{})
	if where != nil {
		if where.Group != nil {
			query = query.Where("`group` = ?", *where.Group)
		}
		if where.Name != nil {
			query = query.Where("`name` = ?", *where.Name)
		}
		if where.Type != nil {
			query = query.Where("`type` = ?", *where.Type)
		}
		if where.Hash != nil {
			query = query.Where("`hash` = ?", *where.Hash)
		}
		if where.Status != nil {
			query = query.Where("`status` = ?", *where.Status)
		}
		if where.CreateTimeStart != nil {
			query = query.Where("`create_time` >= ?", *where.CreateTimeStart)
		}
		if where.CreateTimeEnd != nil {
			query = query.Where("`create_time` <= ?", *where.CreateTimeEnd)
		}
		if where.UpdateTimeStart != nil {
			query = query.Where("`update_time` >= ?", *where.UpdateTimeStart)
		}
		if where.UpdateTimeEnd != nil {
			query = query.Where("`update_time` <= ?", *where.UpdateTimeEnd)
		}
		if where.Content != nil {
			query = query.Where("`content` LIKE ?", "%"+*where.Content+"%")
		}
		if where.Keyword != nil {
			keyword := "%" + *where.Keyword + "%"
			query = query.Where("(`group` LIKE ? OR `name` LIKE ? OR `hash` LIKE ? OR `content` LIKE ?)", keyword, keyword, keyword, keyword)
		}
	}
	return r.Repository.SelectPage(query, queryPage, "group", "name")
}
