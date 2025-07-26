package config

import (
	"server/app/model"
	"server/app/util"
)

// selectByGroup 查看配置数据
func (*Service) selectByGroup(group string) map[string]string {
	var configs []*model.Config
	temp := make(map[string]string)
	query := util.Database.Db.Model(&model.Config{})
	if group != "" {
		query.Where("`key` like ?", group+"%")
	}
	query.Find(&configs)
	for _, v := range configs {
		temp[v.Key] = v.Value
	}
	return temp
}
