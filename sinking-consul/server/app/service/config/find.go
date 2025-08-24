package config

import (
	"server/app/model"
	"server/app/util"
)

// FindByGroupAndName 通过group name查询父信息
func (s *Service) FindByGroupAndName(group string, name string) (data *model.Config, err error) {
	err = util.Database.Db.Where("`group` = ? AND `name` = ?", group, name).First(&data).Error
	return
}
