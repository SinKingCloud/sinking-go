package log

import "server/app/model"

// Select 查询数据
func (r *Repository) Select(where *SelectLog, orderByField string, orderByType string, page int, pageSize int) (list []*model.Log, total int64, err error) {
	offset := pageSize * (page - 1)
	query := r.Database.Db.Model(&model.Log{})
	if where != nil {
		if where.Type != "" {
			query = query.Where("`type` = ?", where.Type)
		}
		if where.Ip != "" {
			query = query.Where("`ip` like ?", "%"+where.Ip+"%")
		}
		if where.Title != "" {
			query = query.Where("`user` like ?", "%"+where.Title+"%")
		}
		if where.Content != "" {
			query = query.Where("`content` like ?", "%"+where.Content+"%")
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
