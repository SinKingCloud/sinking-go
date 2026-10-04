package page

// Query 分页和排序查询参数。
type Query struct {
	Page         int
	PageSize     int
	OrderByField string
	OrderByType  string
	CursorId     string
	CursorLastId string

	HasMore          bool
	NextCursorId     string
	NextCursorLastId string
}

// IsCursor 是否使用游标分页。
func (q *Query) IsCursor() bool {
	return q != nil && q.Page <= 0
}
