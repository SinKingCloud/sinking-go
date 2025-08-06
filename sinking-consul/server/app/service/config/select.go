package config

import (
	"server/app/model"
	"server/app/util"
)

// SelectAll 查询所有
func (s *Service) SelectAll() (list []*Config, err error) {
	err = util.Database.Db.Model(&model.Config{}).Find(&list).Error
	return list, err
}
