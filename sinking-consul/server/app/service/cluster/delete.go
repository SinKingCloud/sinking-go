package cluster

import (
	"server/app/model"
	"server/app/util"
)

// DeleteAll 删除全部数据
func (s *Service) DeleteAll() (err error) {
	err = util.Database.Db.Where("1 = 1").Delete(&model.Cluster{}).Error
	return
}
