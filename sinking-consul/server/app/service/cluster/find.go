package cluster

import (
	"server/app/model"
	"server/app/util"
)

func (s *Service) FindByAddress(address string) (user *Cluster, err error) {
	err = util.Database.Db.Model(&model.Cluster{}).Where("`address` = ?", address).First(&user).Error
	return
}
