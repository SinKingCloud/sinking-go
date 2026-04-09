package config

import (
	"server/app/model"
	repositoryConfig "server/app/repository/config"
	"time"
)

// Save 保存数据
func (s *service) Save() error {
	var configs []*model.Config
	s.Each("*", func(value *model.Config) {
		if value != nil {
			configs = append(configs, value)
		}
	})
	if len(configs) == 0 {
		return nil
	}
	dbConfigs, err := s.repository.SelectAll()
	if err != nil {
		return err
	}
	dbConfigMap := make(map[string]*repositoryConfig.Config, len(dbConfigs))
	for _, dbConfig := range dbConfigs {
		if dbConfig != nil {
			key := dbConfig.Group + ":" + dbConfig.Name
			dbConfigMap[key] = dbConfig
		}
	}
	configsToSave := make([]*model.Config, 0, len(configs))
	for _, localConfig := range configs {
		if localConfig == nil {
			continue
		}
		key := localConfig.Group + ":" + localConfig.Name
		dbConfig, exists := dbConfigMap[key]
		if !exists || dbConfig.Hash != localConfig.Hash || time.Time(dbConfig.UpdateTime).Unix() < time.Time(localConfig.UpdateTime).Unix() {
			configsToSave = append(configsToSave, localConfig)
		}
	}
	if len(configsToSave) == 0 {
		return nil
	}
	return s.repository.Save(configsToSave)
}
