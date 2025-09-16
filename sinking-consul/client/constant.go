package client

import (
	"crypto/tls"
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
)

// ResponseCode 响应码
type ResponseCode int

const (
	ResponseSuccess ResponseCode = 200 //请求成功
	ResponseFail                 = 500 //请求失败
	ResponseError                = 400 //请求错误
)

var (
	// globalClient 全局请求client
	globalClient = &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        10,
			IdleConnTimeout:     30 * time.Second,
			DisableKeepAlives:   false,
			TLSHandshakeTimeout: 5 * time.Second,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
)
