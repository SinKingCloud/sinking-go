package log

import (
	"server/app/model"
	"server/app/util"
)

// create 插入数据
func (s *Service) create(data *model.Log) (err error) {
	err = util.Database.Db.Create(&data).Error
	return
}

// Create 插入数据
func (s *Service) Create(ip string, types Type, title string, content string) {
	go func(types Type, ip string, title string, content string) {
		_ = s.create(&model.Log{
			Type:    int(types),
			Ip:      ip,
			Title:   title,
			Content: content,
		})
	}(types, ip, title, content)
}
