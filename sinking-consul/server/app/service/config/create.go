package config

import (
	"errors"
	"server/app/constant"
	"server/app/enum/config_status"
	"server/app/enum/config_type"
	"server/app/model"
)

// Create 插入数据
func (s *service) Create(data *model.Config) (err error) {
	if data == nil {
		return errors.New("创建数据不能为空")
	}
	if _, ok := config_type.Map()[data.Type]; !ok {
		return errors.New("配置类型不合法")
	}
	if _, ok := config_status.Map()[data.Status]; !ok {
		return errors.New("配置状态不合法")
	}
	key := constant.LockConfigCreate
	if !s.cache.Lock(key, constant.LockTimeConfigCreate) {
		return errors.New("系统繁忙,请稍后重试")
	}
	defer s.cache.UnLock(key)
	info, err := s.FindByGroupAndName(data.Group, data.Name)
	if err != nil {
		return err
	}
	if info != nil {
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
