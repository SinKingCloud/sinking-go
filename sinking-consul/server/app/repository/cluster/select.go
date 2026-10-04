package cluster

import (
	"server/app/model"
	"server/app/util/page"
)

// SelectAll 查询所有
func (r *Repository) SelectAll() (list []*model.Cluster, err error) {
	err = r.Database.Db.Model(&model.Cluster{}).Find(&list).Error
	return list, err
}

// Select 查询数据
func (r *Repository) Select(where *SelectCluster, queryPage *page.Query) (*page.Result[*model.Cluster], error) {
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
	return r.Repository.SelectPage(query, queryPage, "address")
}
