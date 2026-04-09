package cluster

import (
	"server/app/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Create 插入数据
func (r *Repository) Create(data *model.Cluster) error {
	return r.Database.Db.Create(&data).Error
}

// Save 保存数据
func (r *Repository) Save(clusters []*model.Cluster) error {
	return r.Database.Transaction(func(tx *gorm.DB) error {
		return r.Database.BatchExecute(clusters, 1000, func(batch interface{}) error {
			batchClusters := batch.([]*model.Cluster)
			return tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "address"}},
				DoUpdates: clause.AssignmentColumns([]string{"status", "last_heart"}),
			}).CreateInBatches(batchClusters, len(batchClusters)).Error
		})
	})
}
