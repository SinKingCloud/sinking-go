package log

import (
	"server/app/enum/log_type"
	"server/app/model"
)

// Create 插入数据
func (s *service) Create(ip string, types int, title string, content string) {
	if _, ok := log_type.Map()[types]; !ok {
		return
	}
	go func(types int, ip string, title string, content string) {
		_ = s.repositoryLog.Create(&model.Log{
			Type:    types,
			Ip:      ip,
			Title:   title,
			Content: content,
		})
	}(types, ip, title, content)
}
