package context

import (
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

// ValidatePage 分页参数验证
func (c *Context) ValidatePage() (page int, pageSize int) {
	var pageInfo struct {
		Page     int `json:"page" default:"1" validate:"required,gte=1" label:"页码编号"`
		PageSize int `json:"page_size" default:"20" validate:"required,lte=1000" label:"每页数量"`
	}
	if ok, _ := c.ValidatorAll(&pageInfo); !ok {
		pageInfo.Page = 1
		pageInfo.PageSize = 20
	}
	return pageInfo.Page, pageInfo.PageSize
}

// ValidateOrderBy 排序参数验证
func (c *Context) ValidateOrderBy(defaultField string, defaultType string, allowField string) (field string, sort string) {
	var orderBy struct {
		Field string `json:"order_by_field" label:"排序字段"`
		Type  string `json:"order_by_type" label:"排序类型"`
	}
	if c.BindAll(&orderBy) != nil {
		return defaultField, defaultType
	}
	if orderBy.Type == "" {
		orderBy.Type = defaultType
	}
	if orderBy.Field == "" {
		orderBy.Field = defaultField
	}
	orderBy.Type = strings.ToLower(orderBy.Type)
	if orderBy.Type != "desc" && orderBy.Type != "asc" {
		orderBy.Type = "desc"
	}
	fields := strings.Fields(strings.ReplaceAll(allowField, ",", " "))
	for _, v := range fields {
		if v == orderBy.Field {
			return orderBy.Field, orderBy.Type
		}
	}
	return defaultField, orderBy.Type
}
