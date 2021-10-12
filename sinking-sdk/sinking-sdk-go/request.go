package sinking_sdk_go

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type RequestServer struct {
	Server    string //注册中心
	TokenName string //通信密匙名称
	Token     string //通信密匙
}

var client = &http.Client{}

type param map[string]interface{}

// setHttpHeader 批量设置header
func setHttpHeader(headers map[string]string, curl *http.Request) *http.Request {
	for k, v := range headers {
		curl.Header.Set(k, v)
	}
	return curl
}

// getHttpHeader 获取通用请求头
func (r *RequestServer) getHttpHeader(curl *http.Request) *http.Request {
	return setHttpHeader(map[string]string{
		"content-type": "application/json",
		r.TokenName:    r.Token,
	}, curl)
}

// toJson 转json
func toJson(data interface{}) string {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	return string(jsonBytes)
}

// sendRequest 发送请求
func (r *RequestServer) sendRequest(req *http.Request) string {
	r.getHttpHeader(req)
	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	if resp != nil && resp.Body != nil {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				return
			}
		}(resp.Body)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	return string(body)
}

// registerServer 服务注册
func (r *RequestServer) registerServer(name string, appName string, envName string, groupName string, addr string) {
	url := fmt.Sprintf("http://%s/api/service/register", r.Server)
	post := toJson(param{
		"Name":      name,
		"AppName":   appName,
		"EnvName":   envName,
		"GroupName": groupName,
		"Addr":      addr,
	})
	req, err := http.NewRequest("POST", url, strings.NewReader(post))
	if err != nil {
		return
	}
	body := r.sendRequest(req)
	fmt.Println(body)
}
