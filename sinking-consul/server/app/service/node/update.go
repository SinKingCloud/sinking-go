package node

import (
	"server/app/model"
	"server/app/util"
	"server/app/util/str"
	"time"
)

// UpdateAll 更新
func (s *Service) UpdateAll(data map[string]interface{}) (err error) {
	data["update_time"] = str.DateTime(time.Now())
	err = util.Database.Db.Model(&model.Node{}).Where("1 = 1").Updates(data).Error
	return
}
