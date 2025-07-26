package server

import "server/app/util/validator"

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
