package cluster

import (
	"errors"
	"server/app/model"

	"gorm.io/gorm"
)

// FindByAddress 通过地址查询
func (r *Repository) FindByAddress(address string) (*model.Cluster, error) {
	var data *model.Cluster
	err := r.Database.Db.Model(&model.Cluster{}).Where("`address` = ?", address).First(&data).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if data.Address == "" {
		return nil, nil
	}
	return data, nil
}
