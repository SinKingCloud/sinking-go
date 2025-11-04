package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"sync"
)

// Rpc RPC服务器
type Rpc struct {
	client   *Client
	mu       sync.RWMutex
	handlers map[string]RpcHandlerFunc
}

// NewRpc 创建RPC服务器
func NewRpc(client *Client) *Rpc {
	return &Rpc{
		client:   client,
		handlers: make(map[string]RpcHandlerFunc),
	}
}

// Register 注册RPC服务
func (s *Rpc) Register(action string, handler RpcHandlerFunc) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.handlers[action] = handler
}

// ServeHTTP 实现http.Handler接口，处理RPC请求
func (s *Rpc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json;charset=utf-8;")
	w.WriteHeader(http.StatusOK)

	// 定义响应结构
	var resp struct {
		Code    int             `json:"code"`
		Message string          `json:"message"`
		Data    json.RawMessage `json:"data"`
	}

	// 验证token
	if r.Header.Get(TokenName) != s.client.token {
		resp.Code = ResponseFail
		resp.Message = "认证失败"
		_ = json.NewEncoder(w).Encode(resp)
		return
	}
	// 只接受POST请求
	if r.Method != http.MethodPost {
		resp.Code = ResponseFail
		resp.Message = "只支持POST请求"
		_ = json.NewEncoder(w).Encode(resp)
		return
	}
	// 读取请求体
	body, err := io.ReadAll(r.Body)
	if err != nil {
		resp.Code = ResponseFail
		resp.Message = "读取请求失败: " + err.Error()
		_ = json.NewEncoder(w).Encode(resp)
		return
	}
	defer r.Body.Close()
	// 解析请求
	var req struct {
		Action string          `json:"action"`
		Params json.RawMessage `json:"params"`
	}
	if err = json.Unmarshal(body, &req); err != nil {
		resp.Code = ResponseFail
		resp.Message = "解析请求失败: " + err.Error()
		_ = json.NewEncoder(w).Encode(resp)
		return
	}
	// 验证方法名不能为空
	if req.Action == "" {
		resp.Code = ResponseFail
		resp.Message = "方法名不能为空"
		_ = json.NewEncoder(w).Encode(resp)
		return
	}
	// 查找处理器
	s.mu.RLock()
	handler, exists := s.handlers[req.Action]
	s.mu.RUnlock()
	if !exists {
		resp.Code = ResponseFail
		resp.Message = "方法不存在: " + req.Action
		_ = json.NewEncoder(w).Encode(resp)
		return
	}
	// 调用处理器
	result, err := handler(req.Params)
	if err != nil {
		resp.Code = ResponseFail
		resp.Message = err.Error()
		_ = json.NewEncoder(w).Encode(resp)
		return
	}
	// 发送成功响应
	dataBytes, _ := json.Marshal(result)
	resp.Code = int(ResponseSuccess)
	resp.Message = "success"
	resp.Data = dataBytes
	_ = json.NewEncoder(w).Encode(resp)
}

// Call 调用远程RPC服务
// serviceName: 服务名称（自动通过服务发现路由）
// params: 请求参数
// result: 返回结果（必须是指针）
// options: 可选参数 [loadBalance, hashKey]
func (s *Rpc) Call(serviceName string, action string, params interface{}, result interface{}, options ...interface{}) error {
	// 解析可选参数
	loadBalance := Poll
	hashKey := ""
	n := len(options)
	if n > 0 {
		if lb, ok := options[0].(Type); ok {
			loadBalance = lb
		}
	}
	if n > 1 {
		if hk, ok := options[1].(string); ok {
			hashKey = hk
		}
	}
	// 获取服务节点
	var node *Node
	var err error
	if loadBalance == Hash {
		if hashKey == "" {
			return errors.New("Hash模式需要提供hashKey")
		}
		node, err = s.client.GetService(serviceName, loadBalance, hashKey)
	} else {
		node, err = s.client.GetService(serviceName, loadBalance)
	}
	if err != nil {
		return err
	}
	// 构建URL（注册中心的地址已包含完整路径）
	address := node.Address
	if !strings.HasPrefix(address, "http://") && !strings.HasPrefix(address, "https://") {
		address = "http://" + address
	}
	// 序列化请求
	reqData, err := json.Marshal(map[string]interface{}{"action": action, "params": params})
	if err != nil {
		return errors.New("序列化请求失败: " + err.Error())
	}
	// 创建HTTP请求
	httpReq, err := http.NewRequest(http.MethodPost, address, bytes.NewReader(reqData))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set(TokenName, s.client.token)
	// 发送请求
	httpResp, err := globalClient.Do(httpReq)
	if err != nil {
		return errors.New("RPC调用失败: " + err.Error())
	}
	defer httpResp.Body.Close()
	// 读取并解析响应
	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return errors.New("读取响应失败: " + err.Error())
	}
	var resp struct {
		Code    int             `json:"code"`
		Message string          `json:"message"`
		Data    json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return errors.New("解析响应失败: " + err.Error())
	}
	// 检查响应码
	if resp.Code != int(ResponseSuccess) {
		return errors.New(resp.Message)
	}
	// 绑定返回结果
	if result != nil && len(resp.Data) > 0 {
		if err := json.Unmarshal(resp.Data, result); err != nil {
			return errors.New("绑定返回结果失败: " + err.Error())
		}
	}
	return nil
}
