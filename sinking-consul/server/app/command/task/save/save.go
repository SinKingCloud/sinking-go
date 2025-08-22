package save

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"server/app/service"
	"server/app/service/cluster"
	"server/app/service/config"
	"server/app/service/node"
	"server/app/util"
	"time"
)

var (
	saveInterval = time.Minute // 保存间隔时间
	batchSize    = 1000        // 批量保存大小
)

// Init 定时保存数据
func Init() {
	go func() {
		ticker := time.NewTicker(saveInterval)
		defer ticker.Stop()
		for range ticker.C {
			saveClusters()
			saveNodes()
			saveConfig()
		}
	}()
}

// saveClusters 保存集群数据
func saveClusters() {
	var clusters []*cluster.Cluster
	service.Cluster.Each(func(key string, value *cluster.Cluster) bool {
		clusters = append(clusters, value)
		return true
	})
	if len(clusters) == 0 {
		return
	}
	totalClusters := len(clusters)
	for i := 0; i < totalClusters; i += batchSize {
		end := i + batchSize
		if end > totalClusters {
			end = totalClusters
		}
		if i >= totalClusters || end <= i {
			continue
		}
		batch := clusters[i:end]
		err := util.Database.Db.Transaction(func(tx *gorm.DB) error {
			for _, clusterData := range batch {
				tx.Clauses(clause.OnConflict{
					Columns:   []clause.Column{{Name: "address"}},
					DoUpdates: clause.AssignmentColumns([]string{"status", "last_heart"}),
				}).Create(clusterData.Cluster)
			}
			return nil
		})
		if err != nil {
			continue
		}
	}
}

// saveNodes 保存节点数据
func saveNodes() {
	var nodes []*node.Node
	service.Node.Each("*", func(value *node.Node) {
		nodes = append(nodes, value)
	})
	if len(nodes) == 0 {
		return
	}
	totalNodes := len(nodes)
	for i := 0; i < totalNodes; i += batchSize {
		end := i + batchSize
		if end > totalNodes {
			end = totalNodes
		}
		if i >= totalNodes || end <= i {
			continue
		}
		batch := nodes[i:end]
		err := util.Database.Db.Debug().Transaction(func(tx *gorm.DB) error {
			for _, nodeData := range batch {
				tx.Clauses(clause.OnConflict{
					Columns:   []clause.Column{{Name: "group"}, {Name: "name"}, {Name: "address"}},
					DoUpdates: clause.AssignmentColumns([]string{"name", "online_status", "status", "last_heart"}),
				}).Create(nodeData.Node)
			}
			return nil
		})
		if err != nil {
			continue
		}
	}
}

func saveConfig() {
	dbConfigs, err := service.Config.SelectAll()
	if err != nil {
		return
	}
	dbConfigMap := make(map[string]*config.Config)
	for _, dbConfig := range dbConfigs {
		key := dbConfig.Group + ":" + dbConfig.Name
		dbConfigMap[key] = dbConfig
	}
	var configsToUpdate []*config.Config
	service.Config.Each("*", func(localConfig *config.Config) {
		key := localConfig.Group + ":" + localConfig.Name
		dbConfig, exists := dbConfigMap[key]
		if !exists || dbConfig.Hash != localConfig.Hash || time.Time(dbConfig.UpdateTime).Unix() < time.Time(localConfig.UpdateTime).Unix() {
			configsToUpdate = append(configsToUpdate, localConfig)
		}
	})
	if len(configsToUpdate) == 0 {
		return
	}
	totalConfigs := len(configsToUpdate)
	for i := 0; i < totalConfigs; i += batchSize {
		end := i + batchSize
		if end > totalConfigs {
			end = totalConfigs
		}
		if i >= totalConfigs || end <= i {
			continue
		}
		batch := configsToUpdate[i:end]
		err := util.Database.Db.Transaction(func(tx *gorm.DB) error {
			for _, configData := range batch {
				tx.Clauses(clause.OnConflict{
					Columns:   []clause.Column{{Name: "group"}, {Name: "name"}},
					DoUpdates: clause.AssignmentColumns([]string{"type", "hash", "content", "is_delete"}),
				}).Create(configData.Config)
			}
			return nil
		})
		if err != nil {
			continue
		}
	}
}
