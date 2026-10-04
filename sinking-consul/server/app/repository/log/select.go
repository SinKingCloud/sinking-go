package log

import (
	"server/app/model"
	"server/app/util/page"
)

// Select 查询数据
func (r *Repository) Select(where *SelectLog, queryPage *page.Query) (*page.Result[*model.Log], error) {
	query := r.Database.Db.Model(&model.Log{})
	if where != nil {
		if where.Type != nil {
			query = query.Where("`type` = ?", *where.Type)
		}
		if where.CreateTimeStart != nil {
			query = query.Where("`create_time` >= ?", *where.CreateTimeStart)
		}
		if where.CreateTimeEnd != nil {
			query = query.Where("`create_time` <= ?", *where.CreateTimeEnd)
		}
		if where.UpdateTimeStart != nil {
			query = query.Where("`update_time` >= ?", *where.UpdateTimeStart)
		}
		if where.UpdateTimeEnd != nil {
			query = query.Where("`update_time` <= ?", *where.UpdateTimeEnd)
		}
		if where.Ip != nil {
			query = query.Where("`ip` like ?", "%"+*where.Ip+"%")
		}
		if where.Title != nil {
			query = query.Where("`title` like ?", "%"+*where.Title+"%")
		}
		if where.Content != nil {
			query = query.Where("`content` like ?", "%"+*where.Content+"%")
		}
		if where.Keyword != nil {
			keyword := "%" + *where.Keyword + "%"
			query = query.Where("(`ip` LIKE ? OR `title` LIKE ? OR `content` LIKE ?)", keyword, keyword, keyword)
		}
	}
	return r.Repository.SelectPage(query, queryPage)
}
