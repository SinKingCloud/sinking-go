package node

import (
	"errors"
	"server/app/constant"
	"server/app/model"
	"server/app/util"
)

// Create 插入数据
func (s *Service) Create(data *model.Node) (err error) {
	key := constant.LockNodeCreate
	if !util.Cache.Lock(key, constant.LockTimeNodeCreate) {
		return errors.New("系统繁忙,请稍后重试")
	}
	defer util.Cache.UnLock(key)
	info, err := s.FindByGroupAndAddress(data.Group, data.Address)
	if err == nil && info != nil {
		return errors.New("配置已存在")
	}
	err = util.Database.Db.Create(&data).Error
	if err == nil {
		info, err = s.FindByGroupAndAddress(data.Group, data.Address)
		if err == nil && info != nil {
			s.Set(info.Group, info.Name, &Node{Node: info})
		}
	}
	return err
}
