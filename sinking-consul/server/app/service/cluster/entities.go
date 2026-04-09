package cluster

import "server/app/model"

// ConfigUpdateValidate 配置更新验证器
type ConfigUpdateValidate struct {
	Keys    []*model.Config `json:"keys" default:"" validate:"required,min=1,max=1000" label:"配置列表"`
	Type    string          `json:"type" default:"" validate:"omitempty" label:"配置类型"`
	Content string          `json:"content" default:"" validate:"omitempty" label:"配置内容"`
	Status  string          `json:"status" default:"" validate:"omitempty,numeric" label:"状态"`
}

// NodeUpdateValidate 节点更新验证器
type NodeUpdateValidate struct {
	Addresses []string `json:"addresses" default:"" validate:"required,min=1,max=1000,unique" label:"节点列表"`
	Status    string   `json:"status" default:"" validate:"omitempty,numeric" label:"状态"`
}
