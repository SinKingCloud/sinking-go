package cluster

import (
	"server/app/model"
	"server/app/util"
)

// SelectAll 查询所有
func (s *Service) SelectAll() (cluster []*Cluster, err error) {
	err = util.Database.Db.Model(&model.Cluster{}).Find(&cluster).Error
	return cluster, err
}
