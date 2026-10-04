package context

import (
	"server/app/util/page"
	"server/app/util/validator"
	"strings"
)

// Validator 结构体验证
func (c *Context) Validator(data interface{}) (bool, string) {
	return validator.Check(data)
}

// ValidatorAll 参数绑定
func (c *Context) ValidatorAll(data interface{}) (bool, string) {
	if c.BindAll(data) != nil {
		return false, "参数绑定失败"
	}
	return validator.Check(data)
}

// ValidatorJson json参数验证
func (c *Context) ValidatorJson(data interface{}) (bool, string) {
	if c.BindJson(data) != nil {
		return false, "json参数绑定失败"
	}
	return validator.Check(data)
}

// ValidatorPost post参数验证
func (c *Context) ValidatorPost(data interface{}) (bool, string) {
	if c.BindForm(data) != nil {
		return false, "post参数绑定失败"
	}
	return validator.Check(data)
}

// ValidatorGet get参数验证
func (c *Context) ValidatorGet(data interface{}) (bool, string) {
	if c.BindQuery(data) != nil {
		return false, "get参数绑定失败"
	}
	return validator.Check(data)
}

// ValidatorPath path参数验证
func (c *Context) ValidatorPath(data interface{}) (bool, string) {
	if c.BindParam(data) != nil {
		return false, "path参数绑定失败"
	}
	return validator.Check(data)
}

// ValidatePage 验证分页和排序参数，并返回分页查询对象。
// cursorAllowField 表示游标分页允许的排序字段，pageAllowFields 表示普通分页允许的排序字段。
// 不传普通分页白名单时沿用游标白名单；白名单为空或 "*" 时不限制排序字段。
// page 大于 0 时使用普通分页，否则使用游标分页。
func (c *Context) ValidatePage(defaultField, defaultType, cursorAllowField string, pageAllowFields ...string) *page.Query {
	const defaultPageSize = 20
	pageAllowField := strings.Join(pageAllowFields, ",")
	defaultType = strings.ToLower(defaultType)
	if defaultType != "" && defaultType != "desc" && defaultType != "asc" {
		defaultType = "desc"
	}
	var pageInfo struct {
		Page         int    `json:"page" label:"页码编号"`
		PageSize     int    `json:"page_size" label:"每页数量"`
		Field        string `json:"order_by_field" label:"排序字段"`
		Type         string `json:"order_by_type" label:"排序类型"`
		CursorId     string `json:"cursor_id" label:"分页游标排序值"`
		CursorLastId string `json:"cursor_last_id" label:"分页游标定位值"`
	}
	if c.BindAll(&pageInfo) != nil {
		return &page.Query{
			PageSize:     defaultPageSize,
			OrderByField: defaultField,
			OrderByType:  defaultType,
		}
	}
	query := &page.Query{
		Page:         pageInfo.Page,
		PageSize:     pageInfo.PageSize,
		OrderByField: pageInfo.Field,
		OrderByType:  strings.ToLower(pageInfo.Type),
		CursorId:     pageInfo.CursorId,
		CursorLastId: pageInfo.CursorLastId,
	}
	if query.PageSize <= 0 || query.PageSize > 1000 {
		query.PageSize = defaultPageSize
	}
	if query.OrderByField == "" {
		query.OrderByField = defaultField
	}
	if query.OrderByType == "" {
		query.OrderByType = defaultType
	}
	if query.OrderByType != "" && query.OrderByType != "desc" && query.OrderByType != "asc" {
		query.OrderByType = defaultType
	}
	isAllowed := func(fields string) bool {
		if fields == "" || fields == "*" {
			return true
		}
		for _, field := range strings.Fields(strings.ReplaceAll(fields, ",", " ")) {
			if field == query.OrderByField {
				return true
			}
		}
		return false
	}
	allowField := cursorAllowField
	if !query.IsCursor() && pageAllowField != "" {
		allowField = pageAllowField
	}
	if !isAllowed(allowField) {
		query.OrderByField = defaultField
	}
	return query
}
