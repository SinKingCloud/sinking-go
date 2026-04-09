package config

import (
	"errors"
	"server/app/constant"
	"server/app/model"
)

// Create 插入数据
func (s *service) Create(data *model.Config) (err error) {
	key := constant.LockConfigCreate
	if !s.cache.Lock(key, constant.LockTimeConfigCreate) {
		return errors.New("系统繁忙,请稍后重试")
	}
	defer s.cache.UnLock(key)
	info, err := s.FindByGroupAndName(data.Group, data.Name)
	if err == nil && info != nil {
		return errors.New("配置已存在")
	}
	err = s.repository.Create(data)
	if err == nil {
		info, err = s.FindByGroupAndName(data.Group, data.Name)
		if err == nil && info != nil {
			s.Set(info.Group, info.Name, info)
		}
	}
	return err
}
