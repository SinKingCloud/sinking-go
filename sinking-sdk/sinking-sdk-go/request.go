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

// getHttpHeader 获取通用请求头
func (r *RequestServer) getHttpHeader(curl *http.Request) *http.Request {
	return setHttpHeader(map[string]string{
		"content-type": "application/json",
		r.TokenName:    r.Token,
	}, curl)
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
	post := toJson(Param{
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

// getAllServerListResult 所有服务列表
type getAllServerListResult struct {
	Code    int                              `json:"code"`
	Data    map[string]map[string][]*Service `json:"data"`
	Message string                           `json:"message"`
}

// getServerList 拉取服务列表
func (r *RequestServer) getServerList(appName string, envName string, groupName string, name string) *getServerListResult {
	url := fmt.Sprintf("http://%s/api/service/list", r.Server)
	post := toJson(Param{
		"app_name":   appName,
		"env_name":   envName,
		"group_name": groupName,
		"name":       name,
	})
	req, err := http.NewRequest("POST", url, strings.NewReader(post))
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

// getAllServerList 拉取所有服务列表
func (r *RequestServer) getAllServerList(appName string, envName string) *getAllServerListResult {
	url := fmt.Sprintf("http://%s/api/service/all_list", r.Server)
	post := toJson(Param{
		"app_name": appName,
		"env_name": envName,
	})
	req, err := http.NewRequest("POST", url, strings.NewReader(post))
	if err != nil {
		return nil
	}
	body := r.sendRequest(req)
	var res *getAllServerListResult
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil
	}
	return res
}

// changeServerStatus 更改服务状态
func (r *RequestServer) changeServerStatus(serviceHash string, status int) *registerResult {
	url := fmt.Sprintf("http://%s/api/service/status", r.Server)
	post := toJson(Param{
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

// configResult 获取
type configResult struct {
	Code    int       `json:"code"`
	Data    []*Config `json:"data"`
	Message string    `json:"message"`
}

// getConfigs 获取系统配置
func (r *RequestServer) getConfigs(appName string, envName string) *configResult {
	url := fmt.Sprintf("http://%s/api/config/list", r.Server)
	post := toJson(Param{
		"app_name": appName,
		"env_name": envName,
	})
	req, err := http.NewRequest("POST", url, strings.NewReader(post))
	if err != nil {
		return nil
	}
	body := r.sendRequest(req)
	var res *configResult
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil
	}
	return res
}
