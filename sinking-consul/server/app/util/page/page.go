package page

// Result 分页查询结果。
type Result[T any] struct {
	Total            *int64  `json:"total,omitempty"`
	Page             *int    `json:"page,omitempty"`
	PageSize         int     `json:"page_size"`
	List             []T     `json:"list"`
	HasMore          *bool   `json:"has_more,omitempty"`
	NextCursorId     *string `json:"next_cursor_id,omitempty"`
	NextCursorLastId *string `json:"next_cursor_last_id,omitempty"`
}

// New 创建分页查询结果。
func New[T any](query *Query, total int64, list []T) *Result[T] {
	if list == nil {
		list = make([]T, 0)
	}
	if query == nil {
		pageNumber := 1
		return &Result[T]{Total: &total, Page: &pageNumber, PageSize: 20, List: list}
	}
	result := &Result[T]{PageSize: query.PageSize, List: list}
	if query.IsCursor() {
		hasMore := query.HasMore
		result.HasMore = &hasMore
		if query.HasMore {
			nextCursorId := query.NextCursorId
			result.NextCursorId = &nextCursorId
			if query.NextCursorLastId != "" {
				nextCursorLastId := query.NextCursorLastId
				result.NextCursorLastId = &nextCursorLastId
			}
		}
		return result
	}
	pageNumber := query.Page
	result.Total = &total
	result.Page = &pageNumber
	return result
}
