package node

import (
	"server/app/model"

	"gorm.io/gorm"
)

// DeleteByAddress 通过address删除
func (r *Repository) DeleteByAddress(addresses []string) error {
	return r.Database.Transaction(func(tx *gorm.DB) error {
		return r.Database.BatchExecute(addresses, 1000, func(batch interface{}) error {
			return tx.Where("`address` IN ?", batch.([]string)).Delete(&model.Node{}).Error
		})
	})
}
