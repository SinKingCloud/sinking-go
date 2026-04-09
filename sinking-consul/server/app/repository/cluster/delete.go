package cluster

import (
	"server/app/model"
)

// DeleteAll 删除全部数据
func (r *Repository) DeleteAll() (err error) {
	err = r.Database.Db.Where("1 = 1").Delete(&model.Cluster{}).Error
	return
}
