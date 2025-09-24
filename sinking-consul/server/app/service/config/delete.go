package config

import (
	"server/app/model"
	"server/app/util"
)

// DeleteByGroupAndName 通过集群和名称删除
func (s *Service) DeleteByGroupAndName(keys []*model.Config) (err error) {
	var conditions [][]interface{}
	for _, key := range keys {
		if key.Group != "" && key.Name != "" {
			conditions = append(conditions, []interface{}{key.Group, key.Name})
		}
	}
	err = util.Database.Db.Where("(`group`, `name`) IN (?)", conditions).Delete(&model.Config{}).Error
	if err == nil {
		g := make(map[string]int64)
		for _, key := range keys {
			g[key.Group] = 1
			if key.Group != "" && key.Name != "" {
				s.Delete(key.Group, key.Name)
			}
		}
		for group := range g {
			s.SetOperateTime(group)
		}
	}
	return
}
