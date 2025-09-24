package node

import (
	"server/app/model"
	"server/app/util"
	"server/app/util/str"
	"time"
)

// UpdateAll 更新
func (s *Service) UpdateAll(data map[string]interface{}) (err error) {
	data["update_time"] = str.DateTime(time.Now())
	err = util.Database.Db.Model(&model.Node{}).Where("1 = 1").Updates(data).Error
	return
}

// UpdateByAddresses 通过节点地址更新
func (s *Service) UpdateByAddresses(addresses []string, data map[string]interface{}) (err error) {
	data["update_time"] = str.DateTime(time.Now())
	err = util.Database.Db.Model(&model.Node{}).Where("`address` in ? ", addresses).Updates(data).Error
	list, err := s.SelectInAddress(addresses)
	if err == nil {
		s.Sets(list)
		g := make(map[string]int64)
		for _, v := range list {
			g[v.Group] = 1
		}
		for group := range g {
			s.SetOperateTime(group)
		}
	}
	return err
}
