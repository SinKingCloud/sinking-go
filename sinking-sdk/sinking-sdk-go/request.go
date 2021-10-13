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
func (r *RequestServer) sendRequest(req *http.Request) []byte {
	r.getHttpHeader(req)
	resp, err := client.Do(req)
	if err != nil {
		return nil
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
		return nil
	}
	return body
}

// registerResult 结果
type registerResult struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

// registerServer 服务注册
func (r *RequestServer) registerServer(name string, appName string, envName string, groupName string, addr string) *registerResult {
	url := fmt.Sprintf("http://%s/api/service/register", r.Server)
	post := toJson(param{
		"name":       name,
		"app_name":   appName,
		"env_name":   envName,
		"group_name": groupName,
		"addr":       addr,
	})
	req, err := http.NewRequest("POST", url, strings.NewReader(post))
	if err != nil {
		return nil
	}
	body := r.sendRequest(req)
	var res *registerResult
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil
	}
	return res
}

// getServerListResult 服务列表结果
type getServerListResult struct {
	Code    int        `json:"code"`
	Data    []*Service `json:"data"`
	Message string     `json:"message"`
}

// getServerList 拉取服务列表
func (r *RequestServer) getServerList() *getServerListResult {
	url := fmt.Sprintf("http://%s/api/service/list", r.Server)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil
	}
	body := r.sendRequest(req)
	var res *getServerListResult
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil
	}
	return res
}

// changeServerStatus 更改服务状态
func (r *RequestServer) changeServerStatus(serviceHash string, status int) *registerResult {
	url := fmt.Sprintf("http://%s/api/service/status", r.Server)
	post := toJson(param{
		"service_hash": serviceHash,
		"status":       status,
	})
	req, err := http.NewRequest("POST", url, strings.NewReader(post))
	if err != nil {
		return nil
	}
	body := r.sendRequest(req)
	var res *registerResult
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil
	}
	return res
}
