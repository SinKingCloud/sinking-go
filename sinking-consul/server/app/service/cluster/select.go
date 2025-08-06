package cluster

import (
	"server/app/model"
	"server/app/util"
)

// SelectAll 查询所有
func (s *Service) SelectAll() (list []*Cluster, err error) {
	err = util.Database.Db.Model(&model.Cluster{}).Find(&list).Error
	return list, err
}
