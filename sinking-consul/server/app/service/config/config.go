package config

import (
	"server/app/model"
	"server/app/repository/config"
	"server/app/util/cache"
	"sync"
)

// Service 配置服务接口
type Service interface {
	SelectAll() ([]*config.Config, error)
	Select(where *config.SelectConfig, orderByField string, orderByType string, page int, pageSize int) ([]*config.Config, int64, error)
	CountByStatus(status int) (int64, error)
	CountAll() (int64, error)
	Save() error
	Each(group string, fun func(value *model.Config))
	Sets(list []*model.Config)
	SetOperateTime(group string)
	GetOperateTime(group string) int64
	CheckIsChange(list []*model.Config) bool
	GetAllConfigs(group string, showContent bool, filterStatus bool) []*model.Config
	FindByGroupAndName(group string, name string) (*model.Config, error)
	DeleteByGroupAndName(keys []*model.Config) error
	UpdateByGroupAndName(keys []*model.Config, data *config.UpdateConfig) error
	Create(data *model.Config) error
}

// service 配置服务
type service struct {
	repository config.Interface
	cache      cache.Interface

	once                      sync.Once
	configPool                map[string]map[string]*model.Config
	configLock                sync.RWMutex
	configLastOperateTime     map[string]int64
	configLastOperateTimeLock sync.RWMutex
}

// NewService 创建配置服务
func NewService(repository config.Interface, cache cache.Interface) Service {
	s := &service{
		repository:            repository,
		cache:                 cache,
		configPool:            make(map[string]map[string]*model.Config),
		configLastOperateTime: make(map[string]int64),
	}
	s.initialize()
	return s
}
