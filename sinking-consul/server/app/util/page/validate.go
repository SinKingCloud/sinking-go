package page

import "server/app/util/server"

type ValidatePage struct {
	Page     int `json:"page" default:"1" validate:"required,gte=1" label:"页码编号"`
	PageSize int `json:"page_size" default:"20" validate:"required,lte=1000" label:"每页数量"`
}

// ValidatePageDefault 默认分页验证
func ValidatePageDefault(c *server.Context) *ValidatePage {
	pageInfo := &ValidatePage{}
	if ok, _ := c.ValidatorAll(pageInfo); !ok {
		pageInfo.Page = 1
		pageInfo.PageSize = 20
	}
	return pageInfo
}
