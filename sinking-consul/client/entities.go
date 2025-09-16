package client

import "encoding/json"

// Config 配置表
type Config struct {
	Group   string `json:"group"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	Hash    string `json:"hash"`
	Content string `json:"content"`
	Status  int    `json:"status"`
}

// Node 服务列表
type Node struct {
	Group        string `json:"group"`
	Name         string `json:"name"`
	Address      string `json:"address"`
	OnlineStatus int    `json:"online_status"`
	Status       int    `json:"status"`
	LastHeart    int64  `json:"last_heart"`
}

// Response 接口响应
type Response struct {
	Code      int             `json:"code"`
	Data      json.RawMessage `json:"data"`
	Message   string          `json:"message"`
	RequestId string          `json:"request_id"`
}
