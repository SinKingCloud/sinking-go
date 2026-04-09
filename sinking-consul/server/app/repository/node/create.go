package node

import "server/app/model"

// Create 插入数据
func (r *Repository) Create(data *model.Node) error {
	return r.Database.Db.Create(&data).Error
}
