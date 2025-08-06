package node

import (
	"server/app/model"
	"server/app/util"
)

// SelectAll 查询所有
func (s *Service) SelectAll() (list []*Node, err error) {
	err = util.Database.Db.Model(&model.Node{}).Find(&list).Error
	return list, err
}
