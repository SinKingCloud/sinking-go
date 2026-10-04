package node

import (
	"server/app/model"
	"server/app/util/page"
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
func (r *Repository) Select(where *SelectNode, queryPage *page.Query) (*page.Result[*model.Node], error) {
	query := r.Database.Db.Model(&model.Node{})
	if where != nil {
		if where.Group != nil {
			query = query.Where("`group` = ?", *where.Group)
		}
		if where.Name != nil {
			query = query.Where("`name` = ?", *where.Name)
		}
		if where.OnlineStatus != nil {
			query = query.Where("`online_status` = ?", *where.OnlineStatus)
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
		if where.Address != nil {
			query = query.Where("`address` LIKE ?", "%"+*where.Address+"%")
		}
		if where.Keyword != nil {
			keyword := "%" + *where.Keyword + "%"
			query = query.Where("(`group` LIKE ? OR `name` LIKE ? OR `address` LIKE ?)", keyword, keyword, keyword)
		}
	}
	return r.Repository.SelectPage(query, queryPage, "group", "name", "address")
}
