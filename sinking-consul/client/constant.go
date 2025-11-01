package client

import (
	"crypto/tls"
	"encoding/json"
	"net/http"
	"time"
)

const (
	TokenName = "token" //请求token名称
)

// Type 获取service类型
type Type int

const (
	Poll Type = iota //轮询
	Rand             //随机
	Hash             //哈希
)

// ResponseCode 响应码
type ResponseCode int

const (
	ResponseSuccess ResponseCode = 200 //请求成功
	ResponseFail                 = 500 //请求失败
	ResponseError                = 400 //请求错误
)

// RpcHandlerFunc RPC处理函数
type RpcHandlerFunc func(params json.RawMessage) (interface{}, error)

var (
	// globalClient 全局请求client，针对高并发场景优化连接池配置
	globalClient = &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:          100,              // 总空闲连接数，高并发下避免频繁创建连接
			MaxIdleConnsPerHost:   50,               // 每个host的空闲连接数
			MaxConnsPerHost:       100,              // 每个host的最大连接数，防止连接泄漏
			IdleConnTimeout:       90 * time.Second, // 空闲连接超时时间
			DisableKeepAlives:     false,
			TLSHandshakeTimeout:   5 * time.Second,
			DisableCompression:    false, // 启用压缩
			ForceAttemptHTTP2:     true,  // 尝试使用 HTTP/2
			ExpectContinueTimeout: 1 * time.Second,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
)
