package config

import "server/app/model"

// DeleteByGroupAndName 通过集群和名称删除
func (s *service) DeleteByGroupAndName(keys []*model.Config) (err error) {
	err = s.repository.DeleteByGroupAndName(keys)
	if err == nil {
		for _, key := range keys {
			if key.Group != "" && key.Name != "" {
				s.Delete(key.Group, key.Name)
			}
		}
	}
	return
}
