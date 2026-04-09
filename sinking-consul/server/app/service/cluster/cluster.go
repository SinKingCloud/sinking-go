package cluster

import (
	"crypto/tls"
	"net/http"
	"server/app/model"
	"server/app/repository/cluster"
	"server/app/service/config"
	"server/app/service/node"
	"server/app/service/setting"
	"server/app/util/cache"
	"sync"
	"time"
)

// Service 集群服务接口
type Service interface {
	Select(where *cluster.SelectCluster, orderByField string, orderByType string, page int, pageSize int) ([]*model.Cluster, int64, error)
	CountByStatus(status int) (int64, error)
	CountAll() (int64, error)
	Each(fun func(key string, value *model.Cluster) bool)
	Register(address string)
	GetLocalAddr() string
	RegisterRemoteService(remoteAddress string) error
	SynchronizeRemoteData(remoteAddress string) error
	SyncDataLock() error
	SyncDataUnLock() error
	ChangeAllClusterLockStatus(status int) error
	DeleteAllClusterData(configs []*model.Config, nodes []*model.Node)
	UpdateAllClusterData(configs *ConfigUpdateValidate, nodes *NodeUpdateValidate)
	CreateAllClusterData(config *model.Config, node *model.Node)
	Save() error
}

// service 集群服务
type service struct {
	repository    cluster.Interface
	cache         cache.Interface
	setting       setting.Service
	nodeService   node.Service
	configService config.Service

	once                   sync.Once
	syncDataCoroutineCount int64
	clusterPool            sync.Map
	globalClient           *http.Client
}

// NewService 创建集群服务
func NewService(repository cluster.Interface, cache cache.Interface, setting setting.Service, nodeService node.Service, configService config.Service) Service {
	s := &service{
		repository:    repository,
		cache:         cache,
		setting:       setting,
		nodeService:   nodeService,
		configService: configService,
		globalClient: &http.Client{
			Timeout: 10 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        10,
				IdleConnTimeout:     30 * time.Second,
				DisableKeepAlives:   false,
				TLSHandshakeTimeout: 5 * time.Second,
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		},
	}
	s.initialize()
	return s
}
