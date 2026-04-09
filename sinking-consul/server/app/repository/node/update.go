package node

import (
	"server/app/model"
	"server/app/util/str"
	"time"

	"gorm.io/gorm"
)

// UpdateAll 更新
func (r *Repository) UpdateAll(node *UpdateNode) error {
	if node == nil {
		return nil
	}
	updates := make(map[string]interface{})
	if node.Group != nil {
		updates["group"] = node.Group
	}
	if node.Name != nil {
		updates["name"] = node.Name
	}
	if node.Address != nil {
		updates["address"] = node.Address
	}
	if node.OnlineStatus != nil {
		updates["online_status"] = node.OnlineStatus
	}
	if node.Status != nil {
		updates["status"] = node.Status
	}
	if node.LastHeart != nil {
		updates["last_heart"] = node.LastHeart
	}
	if len(updates) == 0 {
		return nil
	}
	updates["update_time"] = str.DateTime(time.Now())
	return r.Database.Db.Model(&model.Node{}).Where("1 = 1").Updates(updates).Error
}

// UpdateByAddresses 通过节点地址更新
func (r *Repository) UpdateByAddresses(addresses []string, node *UpdateNode) error {
	if node == nil {
		return nil
	}
	updates := make(map[string]interface{})
	if node.Group != nil {
		updates["group"] = node.Group
	}
	if node.Name != nil {
		updates["name"] = node.Name
	}
	if node.Address != nil {
		updates["address"] = node.Address
	}
	if node.OnlineStatus != nil {
		updates["online_status"] = node.OnlineStatus
	}
	if node.Status != nil {
		updates["status"] = node.Status
	}
	if node.LastHeart != nil {
		updates["last_heart"] = node.LastHeart
	}
	if len(updates) == 0 {
		return nil
	}
	updates["update_time"] = str.DateTime(time.Now())
	return r.Database.Transaction(func(tx *gorm.DB) error {
		return r.Database.BatchExecute(addresses, 1000, func(batch interface{}) error {
			return tx.Model(&model.Node{}).Where("`address` in ? ", batch.([]string)).Updates(updates).Error
		})
	})
}
