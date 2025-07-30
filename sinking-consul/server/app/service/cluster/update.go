package cluster

import (
	"server/app/model"
	"server/app/util"
	"server/app/util/str"
	"time"
)

// updateAll 更新
func (s *Service) updateAll(data map[string]interface{}) (err error) {
	data["update_time"] = str.DateTime(time.Now())
	err = util.Database.Db.Model(&model.Cluster{}).Updates(data).Error
	return
}
