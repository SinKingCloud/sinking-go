package node

import (
	"server/app/model"
	"server/app/util"
)

// FindByGroupAndAddress 通过address查询信息
func (s *Service) FindByGroupAndAddress(group string, address string) (data *model.Node, err error) {
	err = util.Database.Db.Where("`group` = ? AND `address` = ?", group, address).First(&data).Error
	return
}
