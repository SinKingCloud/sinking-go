package cluster

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"server/app/constant"
	"server/app/model"
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
		_ = s.updateAll(map[string]interface{}{"online_status": Offline})
		all, e := s.SelectAll()
		if e == nil && all != nil {
			for _, v := range all {
				s.Set(v.Address, &Cluster{
					Cluster: &model.Cluster{
						Address:      v.Address,
						OnlineStatus: v.OnlineStatus,
						Status:       v.Status,
						LastHeart:    v.LastHeart,
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
	b, err := json.Marshal(body)
	if err != nil {
		return 500, "json转换失败", nil, err
	}
	req, err := http.NewRequest(http.MethodPost, address, bytes.NewReader(b))
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

// RegisterService 向集群节点发送注册请求
func (s *Service) RegisterService(remoteAddress string) error {
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

// SynchronizeData 同步集群信息并储存更新数据库
func (s *Service) SynchronizeData(remoteAddress string) error {
	//1. 获取集群列表并发送注册集群请求(/Cluster/heart接口用于心跳和注册)
	//2. 定时发送心跳并同步集群内的服务和配置(/Cluster/sync接口用于获取集群所有数据,请求所有集群节点获取数据，并保存存活的节点信息，配置信息和节点重复的话保存最新的心跳和修改时间,所有数据保存在内存中，定时落盘)
	//3. 定时保存数据到数据库进行持久化(如集群信息、服务信息、配置等)
	//获取所有集群节点的数据保存到本机数据库和同步到内存缓存，

	return nil
}
