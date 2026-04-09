package cluster

import (
	"server/app/model"
)

// Create 插入数据
func (r *Repository) Create(data *model.Cluster) error {
	return r.Database.Db.Create(&data).Error
}
