package client

import (
	"errors"
)

type Client struct {
	address []string //注册和配置中心地址
	group   string   //所属群组
	name    string   //服务名称
	addr    string   //服务地址
	token   string   //请求密钥

	server string //使用的注册中心
}

// NewClient 实例化client
// group 服务分组
// token 请求密钥
func NewClient(address []string, group string, name string, addr string, token string) *Client {
	return &Client{
		address: address,
		group:   group,
		name:    name,
		addr:    addr,
		token:   token,
	}
}

// Connect 自动注册
func (c *Client) Connect() error {
	return nil
}

// Close 销毁实例
func (c *Client) Close() error {
	return nil
}

// GetService 获取服务地址
// name 服务名称
// types 获取方式
func (c *Client) GetService(name string, types Type) (string, error) {
	return "", errors.New("获取失败")
}

// GetAllService 获取服务所有地址
func (c *Client) GetAllService() ([]string, error) {
	return nil, errors.New("获取失败")
}

// GetConfig 获取配置信息
// name 配置名
func (c *Client) GetConfig(name string) (*ConfigParser, error) {
	return nil, errors.New("获取失败")
}

// GetAllConfigs 获取所有配置信息
func (c *Client) GetAllConfigs() ([]*ConfigParser, error) {
	return nil, errors.New("获取失败")
}
