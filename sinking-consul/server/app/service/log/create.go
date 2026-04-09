package log

import (
	"server/app/model"
)

// Create 插入数据
func (s *service) Create(ip string, types int, title string, content string) {
	go func(types int, ip string, title string, content string) {
		_ = s.repositoryLog.Create(&model.Log{
			Type:    types,
			Ip:      ip,
			Title:   title,
			Content: content,
		})
	}(types, ip, title, content)
}
