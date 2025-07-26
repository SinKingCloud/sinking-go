package str

import "encoding/json"

// JsonTool json工具
type JsonTool struct {
	data interface{}
}

// NewJsonTool 实例化工具类
func NewJsonTool(data interface{}) *JsonTool {
	return &JsonTool{data: data}
}

// ToString 转json字符串
func (j *JsonTool) ToString() string {
	str, err := json.Marshal(j.data)
	if err != nil {
		return ""
	}
	return string(str)
}
