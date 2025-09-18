package client

import (
	"sync/atomic"
	"time"
)

// registerTask 注册任务（每5秒执行一次）
func (c *Client) registerTask() {
	defer c.wg.Done()
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	// 立即执行一次注册
	_ = c.register()
	for {
		select {
		case <-c.ctx.Done():
			return
		case <-ticker.C:
			_ = c.register()
		}
	}
}

// syncNodeTask 节点同步任务（每5秒执行一次）
func (c *Client) syncNodeTask() {
	defer c.wg.Done()
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	// 立即执行一次同步
	c.doSyncNodes()
	for {
		select {
		case <-c.ctx.Done():
			return
		case <-ticker.C:
			c.doSyncNodes()
		}
	}
}

// doSyncNodes 执行节点同步
func (c *Client) doSyncNodes() {
	lastSyncTime := atomic.LoadInt64(&c.nodeLastSyncTime)
	nodes, err := c.getNodeList(lastSyncTime)
	// 更新缓存
	if err == nil && nodes != nil {
		// 先在锁外构建新的节点缓存
		newNodes := make(map[string][]*Node)
		if len(nodes) > 0 {
			// 按服务名分组节点
			for _, node := range nodes {
				if newNodes[node.Name] == nil {
					newNodes[node.Name] = make([]*Node, 0)
				}
				newNodes[node.Name] = append(newNodes[node.Name], node)
			}
		}

		// 加锁快速更新缓存
		c.nodesMu.Lock()
		c.nodes = newNodes
		c.nodesMu.Unlock()

		atomic.StoreInt64(&c.nodeLastSyncTime, time.Now().Unix())
	}
}

// syncConfigTask 配置同步任务（每10秒执行一次）
func (c *Client) syncConfigTask() {
	defer c.wg.Done()
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	// 立即执行一次同步
	c.doSyncConfigs()
	for {
		select {
		case <-c.ctx.Done():
			return
		case <-ticker.C:
			c.doSyncConfigs()
		}
	}
}

// doSyncConfigs 执行配置同步
func (c *Client) doSyncConfigs() {
	lastSyncTime := atomic.LoadInt64(&c.configLastSyncTime)
	configs, err := c.getConfigList(lastSyncTime)
	// 更新缓存
	if err == nil && configs != nil {
		// 先在锁外构建新的配置缓存和解析器
		newConfigs := make(map[string]*Config)
		newParsers := make(map[string]*ConfigParser)

		if len(configs) > 0 {
			for _, config := range configs {
				// 创建ConfigParser
				parser, err := NewConfigParser(config.Content, config.Type)
				if err != nil {
					// 解析失败，跳过该配置
					continue
				}
				// 只有解析成功才添加到缓存
				newConfigs[config.Name] = config
				newParsers[config.Name] = parser
			}
		}

		// 加锁快速更新缓存
		c.configsMu.Lock()
		c.configs = newConfigs
		c.parsers = newParsers
		c.configsMu.Unlock()

		atomic.StoreInt64(&c.configLastSyncTime, time.Now().Unix())
	}
}
