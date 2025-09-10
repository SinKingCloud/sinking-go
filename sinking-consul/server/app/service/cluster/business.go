package cluster

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"server/app/constant"
	"server/app/model"
	"server/app/service/config"
	"server/app/service/node"
	"server/app/util"
	"server/app/util/ip"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// Get 获取集群节点信息
func (s *Service) Get(key string) *Cluster {
	if v, ok := clusterPool.Load(key); ok && v != nil {
		return v.(*Cluster)
	}
	return nil
}

// Set 设置集群节点信息
func (s *Service) Set(key string, value *Cluster) {
	clusterPool.Store(key, value)
}

// Sets 批量设置集群信息
func (s *Service) Sets(list []*Cluster) {
	for _, v := range list {
		clusterPool.Store(v.Address, v)
	}
}

// Delete 删除集群节点
func (s *Service) Delete(key string) {
	clusterPool.Delete(key)
}

// Each 遍历集群信息
func (s *Service) Each(fun func(key string, value *Cluster) bool) {
	clusterPool.Range(func(key, value any) bool {
		return fun(key.(string), value.(*Cluster))
	})
}

// Register 注册集群信息
func (s *Service) Register(address string) {
	data := s.Get(address)
	if data == nil {
		s.Set(address, &Cluster{
			Cluster: &model.Cluster{
				Address:   address,
				Status:    int(Online),
				LastHeart: time.Now().Unix(),
			},
		})
	} else {
		data.Status = int(Online)
		data.LastHeart = time.Now().Unix()
	}
}

// Init 初始化服务
func (s *Service) Init() {
	clusterOnce.Do(func() {
		_ = s.DeleteAll()
		list := util.Conf.GetStringSlice(constant.ClusterNodes)
		for _, v := range list {
			d, e := s.FindByAddress(v)
			if e != nil || d == nil {
				_ = s.create(&model.Cluster{
					Address:   v,
					Status:    int(Offline),
					LastHeart: 0,
				})
			}
		}
		all, e := s.SelectAll()
		if e == nil && all != nil {
			for _, v := range all {
				s.Set(v.Address, &Cluster{
					Cluster: &model.Cluster{
						Address:    v.Address,
						Status:     v.Status,
						LastHeart:  v.LastHeart,
						CreateTime: v.CreateTime,
						UpdateTime: v.UpdateTime,
					},
				})
			}
		}
	})
}

// request 向集群节点发送请求
func (s *Service) request(address string, action string, body interface{}) (int, string, []byte, error) {
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
			return 500, "json转换失败", nil, err
		}
		reader = bytes.NewReader(b)
	}
	req, err := http.NewRequest(http.MethodPost, address, reader)
	req.Header.Set("content-type", "application/json")
	req.Header.Set(constant.JwtTokenName, util.Conf.GetString(constant.AuthApiToken))
	resp, err := globalClient.Do(req)
	if err != nil {
		return 500, "请求集群失败" + err.Error(), nil, errors.New("请求集群失败: " + err.Error())
	}
	defer resp.Body.Close()
	all, err := io.ReadAll(resp.Body)
	if err != nil {
		return 500, "读取响应失败" + err.Error(), nil, errors.New("读取响应失败: " + err.Error())
	}
	type Response struct {
		Code      int             `json:"code"`
		Data      json.RawMessage `json:"data"`
		Message   string          `json:"message"`
		RequestId string          `json:"request_id"`
	}
	var response Response
	if err = json.Unmarshal(all, &response); err != nil {
		return 500, "解析响应失败" + err.Error(), nil, errors.New("解析响应失败: " + err.Error())
	}
	return response.Code, response.Message, response.Data, nil
}

// GetLocalAddr 获取本机地址
func (s *Service) GetLocalAddr() string {
	str := util.Cache.Remember(constant.CacheNameWithLocalIp, func() interface{} {
		local := util.Conf.GetString(constant.ClusterLocal)
		if local == "" {
			_, port := util.ServerAddr()
			local = ip.GetLocalIP() + ":" + strconv.Itoa(port)
		}
		return local
	}, constant.CacheTimeWithLocalIp)
	return str.(string)
}

// RegisterRemoteService 向集群节点发送注册请求
func (s *Service) RegisterRemoteService(remoteAddress string) error {
	body := map[string]string{
		"address": s.GetLocalAddr(),
	}
	code, message, _, err := s.request(remoteAddress, "api/cluster/register", body)
	if err != nil {
		return err
	}
	if code != 200 {
		return errors.New(message)
	}
	return nil
}

// RemoteLock 远程分布式锁
// status 0加锁 1解锁
func (s *Service) RemoteLock(remoteAddress string, status int) error {
	body := map[string]int{
		"status": status,
	}
	code, message, _, err := s.request(remoteAddress, "api/cluster/lock", body)
	if err != nil {
		return err
	}
	if code != 200 {
		return errors.New(message)
	}
	return nil
}

// RemoteDeleteData 远程删除数据
func (s *Service) RemoteDeleteData(remoteAddress string, configs []*model.Config, nodes []*model.Node) error {
	body := map[string]interface{}{}
	if configs != nil {
		body["configs"] = configs
	}
	if nodes != nil {
		body["nodes"] = nodes
	}
	code, message, _, err := s.request(remoteAddress, "api/cluster/delete", body)
	if err != nil {
		return err
	}
	if code != 200 {
		return errors.New(message)
	}
	return nil
}

// RemoteCreateData 远程创建数据
func (s *Service) RemoteCreateData(remoteAddress string, config *model.Config, node *model.Node) error {
	body := map[string]interface{}{}
	if config != nil {
		body["config"] = config
	}
	if node != nil {
		body["node"] = node
	}
	code, message, _, err := s.request(remoteAddress, "api/cluster/create", body)
	if err != nil {
		return err
	}
	if code != 200 {
		return errors.New(message)
	}
	return nil
}

// RemoteUpdateData 远程更新数据
func (s *Service) RemoteUpdateData(remoteAddress string, configs *ConfigUpdateValidate, nodes *NodeUpdateValidate) error {
	body := map[string]interface{}{}
	if configs != nil {
		body["configs"] = configs
	}
	if nodes != nil {
		body["nodes"] = nodes
	}
	code, message, _, err := s.request(remoteAddress, "api/cluster/update", body)
	if err != nil {
		return err
	}
	if code != 200 {
		return errors.New(message)
	}
	return nil
}

// SynchronizeRemoteData 同步集群信息
func (s *Service) SynchronizeRemoteData(remoteAddress string) error {
	for {
		if util.Cache.IsLock(constant.LockSyncData) {
			time.Sleep(time.Second)
		} else {
			break
		}
	}
	atomic.AddInt64(&syncDataCoroutineCount, 1)
	defer atomic.AddInt64(&syncDataCoroutineCount, -1)
	var nodeList []*node.Node
	needSetNode := false
	code, _, data, err := s.request(remoteAddress, "api/cluster/node", nil)
	if err == nil && code == 200 {
		if err = json.Unmarshal(data, &nodeList); err == nil && nodeList != nil {
			needSetNode = true
		}
	}
	var configList []*config.Config
	needSetConfig := false
	body := map[string]interface{}{
		"show_content": false,
	}
	code, _, data, err = s.request(remoteAddress, "api/cluster/config", body)
	if err == nil && code == 200 {
		if err = json.Unmarshal(data, &configList); err == nil && configList != nil {
			isChange := config.GetIns().CheckIsChange(configList)
			if isChange {
				body = map[string]interface{}{
					"show_content": true,
				}
				code, _, data, err = s.request(remoteAddress, "api/cluster/config", body)
				if err == nil && code == 200 {
					if err = json.Unmarshal(data, &configList); err == nil && configList != nil {
						needSetConfig = true
					}
				}
			}
		}
	}
	if needSetNode && nodeList != nil {
		node.GetIns().Sets(nodeList)
	}
	if needSetConfig && configList != nil {
		config.GetIns().Sets(configList)
	}
	return nil
}

// SyncDataLock 加锁同步数据
func (s *Service) SyncDataLock() error {
	key := constant.LockSyncData
	if !util.Cache.Lock(key, constant.LockTimeSyncData) {
		return errors.New("系统繁忙,请稍后重试(-1)")
	}
	success := false
	for i := 0; i < 15; i++ {
		if atomic.LoadInt64(&syncDataCoroutineCount) >= 1 {
			time.Sleep(time.Second)
		} else {
			success = true
			break
		}
	}
	if success {
		return nil
	}
	return errors.New("系统繁忙,请稍后重试(-2)")
}

// SyncDataUnLock 解锁同步数据
func (s *Service) SyncDataUnLock() error {
	key := constant.LockSyncData
	util.Cache.UnLock(key)
	return nil
}

// ChangeAllClusterLockStatus 远程分布式锁(0加锁 1解锁)
func (s *Service) ChangeAllClusterLockStatus(status int) error {
	type lockResult struct {
		cluster *Cluster
		err     error
	}
	list := make([]*Cluster, 0)
	s.Each(func(key string, value *Cluster) bool {
		if value.Status == int(Online) {
			list = append(list, value)
		}
		return true
	})
	// 使用工作池模式，限制最大并发数
	maxWorkers := 10
	clusterChan := make(chan *Cluster, len(list))
	resultChan := make(chan *lockResult, len(list))
	// 启动工作goroutine
	var wg sync.WaitGroup
	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for cluster := range clusterChan {
				err := s.RemoteLock(cluster.Address, status)
				resultChan <- &lockResult{
					cluster: cluster,
					err:     err,
				}
			}
		}()
	}
	// 发送任务到channel
	for _, cluster := range list {
		clusterChan <- cluster
	}
	close(clusterChan)
	// 等待所有工作完成
	go func() {
		wg.Wait()
		close(resultChan)
	}()
	// 收集结果
	successList := make([]*Cluster, 0)
	var lastError error
	for result := range resultChan {
		if result.err == nil {
			successList = append(successList, result.cluster)
		} else {
			lastError = errors.New("集群节点" + result.cluster.Address + "操作失败:" + result.err.Error())
			// 如果是加锁操作，遇到错误立即停止其他操作
			if status == 0 {
				break
			}
		}
	}
	if status == 0 && lastError != nil {
		// 加锁失败,解锁已加锁成功的节点
		var rollbackWg sync.WaitGroup
		rollbackChan := make(chan *Cluster, len(successList))
		// 启动回滚工作goroutine
		for i := 0; i < maxWorkers; i++ {
			rollbackWg.Add(1)
			go func() {
				defer rollbackWg.Done()
				for cluster := range rollbackChan {
					_ = s.RemoteLock(cluster.Address, 1) // 1表示解锁
				}
			}()
		}
		// 发送回滚任务
		for _, cluster := range successList {
			rollbackChan <- cluster
		}
		close(rollbackChan)
		// 等待所有回滚操作完成
		rollbackWg.Wait()
		return lastError
	}
	return lastError
}

// DeleteAllClusterData 删除所有节点数据
func (s *Service) DeleteAllClusterData(configs []*model.Config, nodes []*model.Node) {
	list := make([]*Cluster, 0)
	s.Each(func(key string, value *Cluster) bool {
		if value.Status == int(Online) {
			list = append(list, value)
		}
		return true
	})
	tasks := make([]func() error, 0)
	for _, cluster := range list {
		tasks = append(tasks, func() error {
			return s.RemoteDeleteData(cluster.Address, configs, nodes)
		})
	}
	s.executeWithWorkers(10, tasks)
}

// UpdateAllClusterData 更新所有节点数据
func (s *Service) UpdateAllClusterData(configs *ConfigUpdateValidate, nodes *NodeUpdateValidate) {
	list := make([]*Cluster, 0)
	s.Each(func(key string, value *Cluster) bool {
		if value.Status == int(Online) {
			list = append(list, value)
		}
		return true
	})
	tasks := make([]func() error, 0)
	for _, cluster := range list {
		tasks = append(tasks, func() error {
			return s.RemoteUpdateData(cluster.Address, configs, nodes)
		})
	}
	s.executeWithWorkers(10, tasks)
}

// CreateAllClusterData 创建所有节点数据
func (s *Service) CreateAllClusterData(config *model.Config, node *model.Node) {
	list := make([]*Cluster, 0)
	s.Each(func(key string, value *Cluster) bool {
		if value.Status == int(Online) {
			list = append(list, value)
		}
		return true
	})
	tasks := make([]func() error, 0)
	for _, cluster := range list {
		tasks = append(tasks, func() error {
			return s.RemoteCreateData(cluster.Address, config, node)
		})
	}
	s.executeWithWorkers(10, tasks)
}

// ExecuteWithWorkers 使用固定数量的工作goroutine并发执行任务
// workers: 工作goroutine数量
// tasks: 要执行的任务函数列表
// 返回: 每个任务的错误结果，顺序与输入任务顺序一致
func (s *Service) executeWithWorkers(workers int, tasks []func() error) []error {
	if workers <= 0 {
		workers = 1
	}
	if len(tasks) == 0 {
		return nil
	}
	// 创建任务通道和结果通道
	taskChan := make(chan int, len(tasks))
	resultChan := make(chan error, len(tasks))
	var wg sync.WaitGroup
	// 启动工作goroutine
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for taskIndex := range taskChan {
				resultChan <- tasks[taskIndex]()
			}
		}()
	}
	// 发送任务索引到通道
	for i := range tasks {
		taskChan <- i
	}
	close(taskChan)
	// 等待所有工作完成
	wg.Wait()
	close(resultChan)
	// 收集结果
	results := make([]error, len(tasks))
	for i := 0; i < len(tasks); i++ {
		results[i] = <-resultChan
	}
	return results
}
