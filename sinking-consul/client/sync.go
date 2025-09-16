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
		c.mu.Lock()
		// 清空现有节点缓存，重新构建（接口返回全量数据）
		c.nodes = make(map[string][]*Node)
		if len(nodes) > 0 {
			// 按服务名分组节点
			for _, node := range nodes {
				if c.nodes[node.Name] == nil {
					c.nodes[node.Name] = make([]*Node, 0)
				}
				c.nodes[node.Name] = append(c.nodes[node.Name], node)
			}
		}
		c.mu.Unlock()
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
		c.mu.Lock()
		defer c.mu.Unlock()
		// 清空现有配置缓存，重新构建（接口返回全量数据）
		c.configs = make(map[string]*Config)
		c.parsers = make(map[string]*ConfigParser)
		if len(configs) > 0 {
			for _, config := range configs {
				// 更新配置缓存
				c.configs[config.Name] = config
				// 创建ConfigParser
				parser, err := NewConfigParser(config.Content, config.Type)
				if err != nil {
					// 解析失败，删除该配置
					delete(c.configs, config.Name)
					continue
				}
				c.parsers[config.Name] = parser
			}
		}
		atomic.StoreInt64(&c.configLastSyncTime, time.Now().Unix())
	}
}
