package node

import (
	"server/app/model"
	"server/app/util"
)

// SelectAll 查询所有
func (s *Service) SelectAll() (list []*Node, err error) {
	err = util.Database.Db.Model(&model.Node{}).Find(&list).Error
	return list, err
}

// SelectInAddress 根据节点地址查询
func (s *Service) SelectInAddress(addresses []string) (list []*Node, err error) {
	err = util.Database.Db.Model(&model.Node{}).Where("`address` in ? ", addresses).Find(&list).Error
	return list, err
}

// Select 获取数据
func (s *Service) Select(where map[string]string, orderByField string, orderByType string, page int, pageSize int) (list []*model.Node, total int64, err error) {
	offset := pageSize * (page - 1)
	query := util.Database.Db.Model(&model.Node{})
	if where["group"] != "" {
		query.Where("`group` = ?", where["group"])
	}
	if where["name"] != "" {
		query.Where("`name` = ?", where["name"])
	}
	if where["online_status"] != "" {
		query.Where("`online_status` = ?", where["online_status"])
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
