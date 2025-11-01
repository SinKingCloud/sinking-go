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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// 验证token
	if r.Header.Get(TokenName) != s.client.token {
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"code": ResponseFail, "message": "认证失败"})
		return
	}
	// 只接受POST请求
	if r.Method != http.MethodPost {
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"code": ResponseFail, "message": "只支持POST请求"})
		return
	}
	// 读取请求体
	body, err := io.ReadAll(r.Body)
	if err != nil {
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"code": ResponseFail, "message": "读取请求失败: " + err.Error()})
		return
	}
	defer r.Body.Close()
	// 解析请求
	var req struct {
		Action string          `json:"action"`
		Params json.RawMessage `json:"params"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"code": ResponseFail, "message": "解析请求失败: " + err.Error()})
		return
	}
	// 验证服务名不能为空
	if req.Action == "" {
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"code": ResponseFail, "message": "方法名不能为空"})
		return
	}
	// 查找处理器
	s.mu.RLock()
	handler, exists := s.handlers[req.Action]
	s.mu.RUnlock()
	if !exists {
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"code": ResponseFail, "message": "方法不存在: " + req.Action})
		return
	}
	// 调用处理器
	result, err := handler(req.Params)
	if err != nil {
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"code": ResponseFail, "message": err.Error()})
		return
	}
	// 发送成功响应
	dataBytes, _ := json.Marshal(result)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"code":    ResponseSuccess,
		"message": "success",
		"data":    json.RawMessage(dataBytes),
	})
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
	if len(options) > 0 {
		if lb, ok := options[0].(Type); ok {
			loadBalance = lb
		}
	}
	if len(options) > 1 {
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
