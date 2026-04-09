package node

import (
	"server/app/model"
)

// FindByGroupAndAddress 通过address查询信息
func (r *Repository) FindByGroupAndAddress(group string, address string) (data *model.Node, err error) {
	err = r.Database.Db.Where("`group` = ? AND `address` = ?", group, address).First(&data).Error
	return
}
