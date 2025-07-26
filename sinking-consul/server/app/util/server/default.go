package server

import "strconv"

// GetStringWithDefault 获取string内容(不存在则返回默认内容)
func (c *Context) GetStringWithDefault(value string, defaultValue string) string {
	if value != "" {
		return value
	}
	return defaultValue
}

// GetIntWithDefault 获取int内容(不存在则返回默认内容)
func (c *Context) GetIntWithDefault(value string, defaultValue int) int {
	if value != "" {
		num, err := strconv.Atoi(value)
		if err == nil && num >= 0 {
			return num
		}
	}
	return defaultValue
}

// GetInt64WithDefault 获取int64内容(不存在则返回默认内容)
func (c *Context) GetInt64WithDefault(value string, defaultValue int64) int64 {
	if value != "" {
		num, err := strconv.ParseInt(value, 64, 10)
		if err == nil && num >= 0 {
			return num
		}
	}
	return defaultValue
}

// GetFloatWithDefault 获取float内容(不存在则返回默认内容)
func (c *Context) GetFloatWithDefault(value string, defaultValue float64) float64 {
	if value != "" {
		num, err := strconv.ParseFloat(value, 64)
		if err == nil && num >= 0 {
			return num
		}
	}
	return defaultValue
}

// GetBoolWithDefault 获取bool内容(不存在则返回默认内容)
func (c *Context) GetBoolWithDefault(value string, defaultValue bool) bool {
	if value != "" {
		if value == "true" || value == "1" {
			return true
		} else {
			return false
		}
	}
	return defaultValue
}
