package client

import (
	"context"
	"errors"
	"hash/fnv"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

type Client struct {
	address []string //注册和配置中心地址
	group   string   //所属群组
	name    string   //服务名称
	addr    string   //服务地址
	token   string   //请求密钥
	server  string   //使用的注册中心
	// 内部状态管理
	ctx     context.Context
	cancel  context.CancelFunc
	wg      sync.WaitGroup
	running int32 // 原子操作标识是否运行中
	// 数据缓存和对应的独立锁
	nodesMu   sync.RWMutex
	nodes     map[string][]*Node // 节点缓存 key: serviceName, value: []*Node
	configsMu sync.RWMutex
	configs   map[string]*Config       // 配置缓存 key: name, value: *Config
	parsers   map[string]*ConfigParser // 配置解析器缓存
	// 同步时间戳
	nodeLastSyncTime   int64
	configLastSyncTime int64
	// 轮询计数器
	pollCounter uint64
	// RPC服务器
	rpc *Rpc
}

// NewClient 实例化client
// group 服务分组
// token 请求密钥
func NewClient(address []string, group string, name string, addr string, token string) *Client {
	ctx, cancel := context.WithCancel(context.Background())
	c := &Client{
		address: address,
		group:   group,
		name:    name,
		addr:    addr,
		token:   token,
		ctx:     ctx,
		cancel:  cancel,
		nodes:   make(map[string][]*Node),
		configs: make(map[string]*Config),
		parsers: make(map[string]*ConfigParser),
	}
	c.rpc = NewRpc(c)
	return c
}

// Connect 自动注册
func (c *Client) Connect() error {
	if !atomic.CompareAndSwapInt32(&c.running, 0, 1) {
		return errors.New("客户端已经在运行中")
	}
	addr := c.getAddr(false)
	if err := c.testing(addr); err != nil {
		atomic.StoreInt32(&c.running, 0)
		return errors.New("连接服务器失败: " + err.Error())
	}

	// 先执行一次注册和同步，确保客户端连接后立即有可用数据
	_ = c.register()
	c.doSyncNodes()
	c.doSyncConfigs()

	// 启动异步同步任务
	c.wg.Add(1)
	go c.registerTask()
	c.wg.Add(1)
	go c.syncNodeTask()
	c.wg.Add(1)
	go c.syncConfigTask()
	return nil
}

// Close 销毁实例
func (c *Client) Close() error {
	if !atomic.CompareAndSwapInt32(&c.running, 1, 0) {
		return errors.New("客户端未在运行中")
	}
	c.cancel()
	c.wg.Wait()
	// 分别清理节点和配置缓存
	c.nodesMu.Lock()
	c.nodes = make(map[string][]*Node)
	c.nodesMu.Unlock()
	c.configsMu.Lock()
	c.configs = make(map[string]*Config)
	c.parsers = make(map[string]*ConfigParser)
	c.configsMu.Unlock()
	return nil
}

// GetService 获取服务节点
// name 服务名称
// types 获取方式
// options 可选参数，当types为Hash时，options[0]应为hashKey
func (c *Client) GetService(name string, types Type, options ...string) (*Node, error) {
	c.nodesMu.RLock()
	defer c.nodesMu.RUnlock()
	serviceNodes, exists := c.nodes[name]
	num := len(serviceNodes)
	if !exists || num == 0 {
		return nil, errors.New("服务不存在")
	}
	switch types {
	case Poll: // 轮询
		counter := atomic.AddUint64(&c.pollCounter, 1)
		index := counter % uint64(num)
		return serviceNodes[index], nil
	case Rand: // 随机
		index := int(time.Now().UnixNano()) % num
		return serviceNodes[index], nil
	case Hash: // 哈希
		if len(options) == 0 || options[0] == "" {
			return nil, errors.New("Hash模式需要提供hashKey")
		}
		h := fnv.New32a()
		_, _ = h.Write([]byte(options[0]))
		index := int(h.Sum32()) % num
		return serviceNodes[index], nil
	default:
		return serviceNodes[0], nil
	}
}

// GetAllService 获取所有服务节点
func (c *Client) GetAllService() ([]*Node, error) {
	c.nodesMu.RLock()
	defer c.nodesMu.RUnlock()
	// 预先计算总节点数，避免多次扩容
	totalCount := 0
	for _, serviceNodes := range c.nodes {
		totalCount += len(serviceNodes)
	}
	if totalCount == 0 {
		return nil, errors.New("没有可用的服务节点")
	}
	// 预分配足够容量的切片
	allNodes := make([]*Node, 0, totalCount)
	for _, serviceNodes := range c.nodes {
		allNodes = append(allNodes, serviceNodes...)
	}
	return allNodes, nil
}

// GetServiceNodes 获取指定服务的所有节点
func (c *Client) GetServiceNodes(name string) ([]*Node, error) {
	c.nodesMu.RLock()
	defer c.nodesMu.RUnlock()
	serviceNodes, exists := c.nodes[name]
	if !exists || len(serviceNodes) == 0 {
		return nil, errors.New("服务不存在")
	}
	// 返回切片副本，避免指针引用问题
	result := make([]*Node, len(serviceNodes))
	copy(result, serviceNodes)
	return result, nil
}

// GetConfig 获取配置信息
// name 配置名
func (c *Client) GetConfig(name string) (*ConfigParser, error) {
	c.configsMu.RLock()
	defer c.configsMu.RUnlock()
	parser, exists := c.parsers[name]
	if !exists {
		return nil, errors.New("配置不存在或已禁用")
	}
	return parser, nil
}

// GetAllConfigs 获取所有配置信息
func (c *Client) GetAllConfigs() ([]*ConfigParser, error) {
	c.configsMu.RLock()
	defer c.configsMu.RUnlock()
	var parsers []*ConfigParser
	for _, parser := range c.parsers {
		parsers = append(parsers, parser)
	}
	return parsers, nil
}

// Register 注册RPC服务
func (c *Client) Register(service string, handler RpcHandlerFunc) {
	c.rpc.Register(service, handler)
}

// RpcHandle 获取RPC HTTP处理器
func (c *Client) RpcHandle() http.Handler {
	return c.rpc
}

// RpcCall 调用远程RPC服务
func (c *Client) RpcCall(serviceName string, params interface{}, result interface{}, options ...interface{}) error {
	return c.rpc.Call(serviceName, params, result, options...)
}
