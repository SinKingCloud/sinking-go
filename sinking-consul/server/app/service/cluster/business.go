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

// GetAllClusters 获取所有集群数据
func (s *Service) GetAllClusters() []*Cluster {
	list := make([]*Cluster, 100)
	clusterPool.Range(func(key, value any) bool {
		v := value.(*Cluster)
		if v.IsDelete == int(False) {
			list = append(list, value.(*Cluster))
		}
		return true
	})
	return list
}

// Register 注册集群信息
func (s *Service) Register(address string) {
	data := s.Get(address)
	if data == nil {
		s.Set(address, &Cluster{
			Cluster: &model.Cluster{
				Address:      address,
				OnlineStatus: int(Online),
				Status:       int(Normal),
				LastHeart:    time.Now().Unix(),
			},
		})
	} else {
		data.OnlineStatus = int(Online)
		data.LastHeart = time.Now().Unix()
	}
}

// Init 初始化服务
func (s *Service) Init() {
	clusterOnce.Do(func() {
		list := util.Conf.GetStringSlice(constant.ClusterNodes)
		for _, v := range list {
			d, e := s.FindByAddress(v)
			if e != nil || d == nil {
				_ = s.create(&model.Cluster{
					Address:      v,
					OnlineStatus: int(Offline),
					Status:       int(Normal),
					LastHeart:    0,
				})
			}
		}
		_ = s.UpdateAll(map[string]interface{}{"online_status": Offline})
		all, e := s.SelectAll()
		if e == nil && all != nil {
			for _, v := range all {
				s.Set(v.Address, &Cluster{
					Cluster: &model.Cluster{
						Address:      v.Address,
						OnlineStatus: v.OnlineStatus,
						Status:       v.Status,
						LastHeart:    v.LastHeart,
						CreateTime:   v.CreateTime,
						UpdateTime:   v.UpdateTime,
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

// getLocalAddr 获取本机地址
func (s *Service) getLocalAddr() string {
	str := util.Cache.Remember(constant.CacheNameWithLocalIp, func() interface{} {
		local := util.Conf.GetString(constant.ClusterLocal)
		if local == "" {
			_, port := util.ServerAddr()
			local = ip.GetLocalIP() + ":" + strconv.Itoa(port)
		}
		return local
	}, constant.CacheTimeWithLocalIp*time.Second)
	return str.(string)
}

// RegisterRemoteService 向集群节点发送注册请求
func (s *Service) RegisterRemoteService(remoteAddress string) error {
	body := map[string]string{
		"address": s.getLocalAddr(),
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

// SynchronizeRemoteData 同步集群信息
func (s *Service) SynchronizeRemoteData(remoteAddress string) error {
	code, _, data, err := s.request(remoteAddress, "api/cluster/node", nil)
	if err == nil && code == 200 {
		var list []*node.Node
		if err = json.Unmarshal(data, &list); err == nil && list != nil {
			node.GetIns().Sets(list)
		}
	}
	code, _, data, err = s.request(remoteAddress, "api/cluster/list", nil)
	if err == nil && code == 200 {
		var list []*Cluster
		if err = json.Unmarshal(data, &list); err == nil && list != nil {
			s.Sets(list)
		}
	}
	body := map[string]interface{}{
		"show_content": false,
	}
	code, _, data, err = s.request(remoteAddress, "api/cluster/config", body)
	if err == nil && code == 200 {
		var list []*config.Config
		if err = json.Unmarshal(data, &list); err == nil && list != nil {
			isChange := config.GetIns().CheckIsChange(list)
			if isChange {
				body = map[string]interface{}{
					"show_content": true,
				}
				code, _, data, err = s.request(remoteAddress, "api/cluster/config", body)
				if err == nil && code == 200 {
					if err = json.Unmarshal(data, &list); err == nil && list != nil {
						config.GetIns().Sets(list)
					}
				}
			}
		}
	}
	return nil
}
