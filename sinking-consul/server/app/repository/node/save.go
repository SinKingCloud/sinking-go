package node

import (
	"server/app/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Save 保存数据
func (r *Repository) Save(nodes []*model.Node) error {
	return r.Database.Transaction(func(tx *gorm.DB) error {
		return r.Database.BatchExecute(nodes, 1000, func(batch interface{}) error {
			batchNodes := batch.([]*model.Node)
			return tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "group"}, {Name: "name"}, {Name: "address"}},
				DoUpdates: clause.AssignmentColumns([]string{"name", "online_status", "status", "last_heart"}),
			}).CreateInBatches(batchNodes, len(batchNodes)).Error
		})
	})
}
