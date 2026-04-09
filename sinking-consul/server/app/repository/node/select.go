package node

import (
	"server/app/model"
)

// SelectAll 查询所有
func (r *Repository) SelectAll() (list []*model.Node, err error) {
	err = r.Database.Db.Model(&model.Node{}).Find(&list).Error
	return
}

// SelectInAddress 根据节点地址查询
func (r *Repository) SelectInAddress(addresses []string) (list []*model.Node, err error) {
	if len(addresses) == 0 {
		return nil, nil
	}
	err = r.Database.BatchExecute(addresses, 1000, func(batch interface{}) error {
		var batchList []*model.Node
		err = r.Database.Db.Model(&model.Node{}).
			Where("`address` IN ?", batch.([]string)).
			Find(&batchList).Error
		if err != nil {
			return err
		}
		list = append(list, batchList...)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return list, nil
}

// Select 查询数据
func (r *Repository) Select(where *SelectNode, orderByField string, orderByType string, page int, pageSize int) (list []*model.Node, total int64, err error) {
	offset := pageSize * (page - 1)
	query := r.Database.Db.Model(&model.Node{})
	if where != nil {
		if where.Group != "" {
			query = query.Where("`group` = ?", where.Group)
		}
		if where.Name != "" {
			query = query.Where("`name` = ?", where.Name)
		}
		if where.OnlineStatus != "" {
			query = query.Where("`online_status` = ?", where.OnlineStatus)
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
