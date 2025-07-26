package config

import (
	"server/app/model"
	"server/app/util"
	"server/app/util/str"
	"time"
)

// updateByKey 通过KEY更新
func (s *Service) updateByKey(key string, value string) (err error) {
	err = util.Database.Db.Model(&model.Config{}).Where("`key` = ? ", key).Updates(map[string]any{
		"update_time": str.DateTime(time.Now()),
		"value":       value,
	}).Error
	return
}
