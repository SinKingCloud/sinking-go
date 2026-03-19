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
	localLastOperateTime := atomic.LoadInt64(&c.nodeLastOperateTime)
	lastOperateTime, nodes, err := c.getNodeList(localLastOperateTime)
	if err == nil && localLastOperateTime != lastOperateTime {
		newNodes := make(map[string][]*Node)
		for _, node := range nodes {
			if newNodes[node.Name] == nil {
				newNodes[node.Name] = make([]*Node, 0)
			}
			newNodes[node.Name] = append(newNodes[node.Name], node)
		}
		c.nodesMu.Lock()
		c.nodes = newNodes
		c.nodesMu.Unlock()
		atomic.StoreInt64(&c.nodeLastOperateTime, lastOperateTime)
	}
}

// syncConfigTask 配置同步任务（每5秒执行一次）
func (c *Client) syncConfigTask() {
	defer c.wg.Done()
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-c.ctx.Done():
			return
		case <-ticker.C:
			c.doSyncConfigs()
		}
	}
}

// doSyncConfigs 执行配置同步（支持增删改）
func (c *Client) doSyncConfigs() {
	localLastOperateTime := atomic.LoadInt64(&c.configLastOperateTime)
	lastOperateTime, configs, err := c.getConfigList(localLastOperateTime)
	if err == nil && localLastOperateTime != lastOperateTime {
		if configs == nil {
			configs = []*Config{}
		}
		// 第一步：构建服务端配置名称集合，用于增量删除
		serverConfigNames := make(map[string]bool, len(configs))
		for _, config := range configs {
			serverConfigNames[config.Name] = true
		}
		// 第二步：找出需要删除的配置（本地有但服务端没有）
		c.configsMu.Lock()
		var deleteNames []string
		for name := range c.configs {
			if !serverConfigNames[name] {
				deleteNames = append(deleteNames, name)
			}
		}
		// 执行删除
		for _, name := range deleteNames {
			delete(c.configs, name)
			delete(c.parsers, name)
		}
		c.configsMu.Unlock()
		// 第三步：处理新增和更新
		type updateAction struct {
			config         *Config
			existingParser *ConfigParser
			newParser      *ConfigParser
			actionType     string
		}
		var actions []updateAction
		c.configsMu.RLock()
		for _, config := range configs {
			if existingParser, exists := c.parsers[config.Name]; exists {
				actions = append(actions, updateAction{
					config:         config,
					existingParser: existingParser,
					actionType:     "update",
				})
			} else {
				actions = append(actions, updateAction{
					config:     config,
					actionType: "create",
				})
			}
		}
		c.configsMu.RUnlock()
		for i := range actions {
			action := &actions[i]
			switch action.actionType {
			case "update":
				if err = action.existingParser.UpdateConfig(action.config.Content, action.config.Type); err != nil {
					if newParser, parseErr := NewConfigParser(action.config.Content, action.config.Type); parseErr == nil {
						action.newParser = newParser
						action.actionType = "replace"
					} else {
						action.actionType = "skip"
					}
				}
			case "create":
				if newParser, parseErr := NewConfigParser(action.config.Content, action.config.Type); parseErr == nil {
					action.newParser = newParser
				} else {
					action.actionType = "skip"
				}
			}
		}
		c.configsMu.Lock()
		for _, action := range actions {
			switch action.actionType {
			case "update":
				c.configs[action.config.Name] = action.config
			case "replace":
				c.configs[action.config.Name] = action.config
				c.parsers[action.config.Name] = action.newParser
			case "create":
				c.configs[action.config.Name] = action.config
				c.parsers[action.config.Name] = action.newParser
			case "skip":
				continue
			}
		}
		c.configsMu.Unlock()
		atomic.StoreInt64(&c.configLastOperateTime, lastOperateTime)
	}
}
