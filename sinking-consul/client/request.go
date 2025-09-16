package client

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
)

// getAddr 获取请求地址
// change 是否更换节点
func (c *Client) getAddr(change bool) string {
	if !change && c.server != "" {
		return c.server
	}
	addr := ""
	if len(c.address) == 1 {
		addr = c.address[0]
	} else {
		num := uint64(len(c.address))
		hash := sha256.Sum256([]byte(c.group + c.name + c.addr))
		hashValue := binary.BigEndian.Uint64(hash[:8])
		index := hashValue % num
		//判断节点是否可用
		for i := uint64(0); i < num; i++ {
			addr = c.address[index]
			err := c.testing(addr)
			if err == nil {
				break
			}
			index++
			if index == num {
				index = 0
			}
		}
	}
	c.server = addr
	return addr
}

// testing 测试请求
func (c *Client) testing(addr string) error {
	code, message, _, err := c.request(addr, "api/cluster/testing", nil)
	if err != nil {
		return err
	}
	if code != ResponseSuccess {
		return errors.New(message)
	}
	return nil
}

// register 注册请求
func (c *Client) register(addr string) error {
	body := map[string]string{
		"group":   c.group,
		"name":    c.name,
		"address": c.addr,
	}
	code, message, _, err := c.request(addr, "api/node/register", body)
	if code == ResponseError {
		c.getAddr(true)
	}
	if err != nil {
		return err
	}
	if code != ResponseSuccess {
		return errors.New(message)
	}
	return nil
}

// getNodeList 获取节点信息
func (c *Client) getNodeList(addr string, lastSyncTime int64) (error, []*Node) {
	body := map[string]interface{}{
		"group":          c.group,
		"last_sync_time": lastSyncTime,
	}
	code, message, data, err := c.request(addr, "api/node/sync", body)
	if code == ResponseError {
		c.getAddr(true)
	}
	if err != nil {
		return err, nil
	}
	if code != ResponseSuccess {
		return errors.New(message), nil
	}
	var list []*Node
	err = json.Unmarshal(data, &list)
	if err != nil {
		return err, nil
	}
	return nil, list
}

// getConfigList 获取配置信息
func (c *Client) getConfigList(addr string, lastSyncTime int64) (error, []*Config) {
	body := map[string]interface{}{
		"group":          c.group,
		"last_sync_time": lastSyncTime,
	}
	code, message, data, err := c.request(addr, "api/config/sync", body)
	if code == ResponseError {
		c.getAddr(true)
	}
	if err != nil {
		return err, nil
	}
	if code != ResponseSuccess {
		return errors.New(message), nil
	}
	var list []*Config
	err = json.Unmarshal(data, &list)
	if err != nil {
		return err, nil
	}
	return nil, list
}

// request 向集群节点发送请求
func (c *Client) request(address string, action string, body interface{}) (ResponseCode, string, []byte, error) {
	if !strings.HasPrefix(address, "http://") {
		address = "http://" + address
	}
	if !strings.HasSuffix(address, "/") {
		address += "/"
	}
	if strings.HasPrefix(action, "/") {
		action = strings.TrimPrefix(action, "/")
	}
	address += action
	var reader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return ResponseFail, "json转换失败", nil, err
		}
		reader = bytes.NewReader(b)
	}
	req, err := http.NewRequest(http.MethodPost, address, reader)
	req.Header.Set("content-type", "application/json")
	req.Header.Set(TokenName, c.token)
	resp, err := globalClient.Do(req)
	if err != nil {
		return ResponseError, "请求集群失败" + err.Error(), nil, errors.New("请求集群失败: " + err.Error())
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	all, err := io.ReadAll(resp.Body)
	if err != nil {
		return ResponseFail, "读取响应失败" + err.Error(), nil, errors.New("读取响应失败: " + err.Error())
	}
	var response Response
	if err = json.Unmarshal(all, &response); err != nil {
		return ResponseFail, "解析响应失败" + err.Error(), nil, errors.New("解析响应失败: " + err.Error())
	}
	code := ResponseCode(response.Code)
	if code == ResponseSuccess {
		return ResponseSuccess, "", response.Data, nil
	}
	return ResponseCode(response.Code), response.Message, response.Data, nil
}
