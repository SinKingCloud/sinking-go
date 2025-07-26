package page

// NewPage 实例化新page
func NewPage(total int64, page int, pageSize int, list interface{}, other ...map[string]interface{}) map[string]interface{} {
	data := map[string]interface{}{
		"total":     total,
		"page":      page,
		"page_size": pageSize,
		"list":      list,
	}
	if len(other) > 0 {
		for _, v := range other {
			for k2, v2 := range v {
				data[k2] = v2
			}
		}
	}
	return data
}
