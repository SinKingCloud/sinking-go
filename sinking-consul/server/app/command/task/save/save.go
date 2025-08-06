package save

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"server/app/service"
	"server/app/service/cluster"
	"server/app/service/node"
	"server/app/util"
	"time"
)

var (
	saveInterval = time.Second // 保存间隔时间
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
		}
	}()
}

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
				err := tx.Clauses(clause.OnConflict{
					Columns:   []clause.Column{{Name: "address"}},
					DoUpdates: clause.AssignmentColumns([]string{"online_status", "status", "last_heart", "update_time"}),
				}).Create(clusterData.Cluster).Error
				if err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			continue
		}
	}
}

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
		err := util.Database.Db.Transaction(func(tx *gorm.DB) error {
			for _, nodeData := range batch {
				err := tx.Clauses(clause.OnConflict{
					Columns:   []clause.Column{{Name: "group"}, {Name: "address"}},
					DoUpdates: clause.AssignmentColumns([]string{"name", "online_status", "status", "last_heart", "update_time"}),
				}).Create(nodeData.Node).Error
				if err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			continue
		}
	}
}
