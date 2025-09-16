package client

import (
	"context"
	"errors"
	"math/rand"
	"sync"
	"sync/atomic"
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
	mu      sync.RWMutex
	running int32 // 原子操作标识是否运行中
	// 数据缓存
	nodes   map[string][]*Node       // 节点缓存 key: serviceName, value: []*Node
	configs map[string]*Config       // 配置缓存 key: name, value: *Config
	parsers map[string]*ConfigParser // 配置解析器缓存
	// 同步时间戳
	nodeLastSyncTime   int64
	configLastSyncTime int64
	// 轮询计数器
	pollCounter uint64
}

// NewClient 实例化client
// group 服务分组
// token 请求密钥
func NewClient(address []string, group string, name string, addr string, token string) *Client {
	ctx, cancel := context.WithCancel(context.Background())
	return &Client{
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
	c.mu.Lock()
	c.nodes = make(map[string][]*Node)
	c.configs = make(map[string]*Config)
	c.parsers = make(map[string]*ConfigParser)
	c.mu.Unlock()
	return nil
}

// GetService 获取服务节点
// name 服务名称
// types 获取方式
func (c *Client) GetService(name string, types Type) (*Node, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
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
		index := rand.Intn(num)
		return serviceNodes[index], nil
	default:
		return serviceNodes[0], nil
	}
}

// GetAllService 获取所有服务节点
func (c *Client) GetAllService() ([]*Node, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	var allNodes []*Node
	for _, serviceNodes := range c.nodes {
		allNodes = append(allNodes, serviceNodes...)
	}
	if len(allNodes) == 0 {
		return nil, errors.New("没有可用的服务节点")
	}
	return allNodes, nil
}

// GetServiceNodes 获取指定服务的所有节点
func (c *Client) GetServiceNodes(name string) ([]*Node, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	serviceNodes, exists := c.nodes[name]
	if !exists || len(serviceNodes) == 0 {
		return nil, errors.New("服务不存在")
	}
	return serviceNodes, nil
}

// GetConfig 获取配置信息
// name 配置名
func (c *Client) GetConfig(name string) (*ConfigParser, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	parser, exists := c.parsers[name]
	if !exists {
		return nil, errors.New("配置不存在或已禁用")
	}
	return parser, nil
}

// GetAllConfigs 获取所有配置信息
func (c *Client) GetAllConfigs() ([]*ConfigParser, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	var parsers []*ConfigParser
	for _, parser := range c.parsers {
		parsers = append(parsers, parser)
	}
	return parsers, nil
}
