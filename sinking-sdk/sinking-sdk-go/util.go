package sinking_sdk_go

import (
	"encoding/json"
	"net/http"
	"time"
)

// Param 参数构建
type Param map[string]interface{}

// toJson 转json
func toJson(data interface{}) string {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	return string(jsonBytes)
}

// client 请求对象
var client = &http.Client{
	Timeout: 3 * time.Second, //超时时间
}

// setHttpHeader 批量设置header
func setHttpHeader(headers map[string]string, curl *http.Request) *http.Request {
	for k, v := range headers {
		curl.Header.Set(k, v)
	}
	return curl
}
