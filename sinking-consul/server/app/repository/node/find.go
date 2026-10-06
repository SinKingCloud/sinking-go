package node

import (
	"errors"
	"server/app/model"

	"gorm.io/gorm"
)

// FindByGroupAndAddress 通过分组和地址查询节点信息
func (r *Repository) FindByGroupAndAddress(group string, address string) (*model.Node, error) {
	var data *model.Node
	err := r.Database.Db.
		Model(&model.Node{}).
		Where("`group` = ? AND `address` = ?", group, address).
		First(&data).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if data == nil || data.Group == "" || data.Name == "" || data.Address == "" {
		return nil, nil
	}
	return data, nil
}
