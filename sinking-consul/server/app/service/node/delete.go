package node

import (
	"server/app/model"
	"server/app/util"
)

// DeleteByAddress 通过address删除
func (s *Service) DeleteByAddress(addresses []string) (err error) {
	list2, err2 := s.SelectInAddress(addresses)
	err = util.Database.Db.Where("`address` IN ?", addresses).Delete(&model.Node{}).Error
	if list2 != nil && err2 == nil {
		g := make(map[string]int64)
		for _, key := range list2 {
			g[key.Group] = 1
			if key.Group != "" && key.Name != "" {
				s.Delete(key.Group, key.Address)
			}
		}
		for group := range g {
			s.SetOperateTime(group)
		}
	}
	return
}
