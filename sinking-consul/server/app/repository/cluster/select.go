package cluster

import (
	"server/app/model"
)

// SelectAll 查询所有
func (r *Repository) SelectAll() (list []*model.Cluster, err error) {
	err = r.Database.Db.Model(&model.Cluster{}).Find(&list).Error
	return list, err
}

// Select 查询数据
func (r *Repository) Select(where *SelectCluster, orderByField string, orderByType string, page int, pageSize int) (list []*model.Cluster, total int64, err error) {
	offset := pageSize * (page - 1)
	query := r.Database.Db.Model(&model.Cluster{})
	if where != nil {
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
		if where.Address != nil {
			query = query.Where("`address` LIKE ?", "%"+*where.Address+"%")
		}
		if where.Keyword != nil {
			query = query.Where("`address` LIKE ?", "%"+*where.Keyword+"%")
		}
	}
	err = query.Count(&total).Limit(pageSize).Offset(offset).Order(orderByField + " " + orderByType).Find(&list).Error
	return
}
