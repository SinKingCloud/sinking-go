package enum

import (
	"server/app/enum/cluster_status"
	"server/app/enum/config_status"
	"server/app/enum/config_type"
	"server/app/enum/log_type"
	"server/app/enum/node_online_status"
	"server/app/enum/node_status"
)

// Data 枚举信息
var Data = map[string]interface{}{
	"log": map[string]interface{}{
		"type": log_type.Map(), //日志类型
	},
	"cluster": map[string]interface{}{
		"status": cluster_status.Map(), //在线状态
	},
	"node": map[string]interface{}{
		"online_status": node_online_status.Map(), //在线状态
		"status":        node_status.Map(),        //集群状态
	},
	"config": map[string]interface{}{
		"type":   config_type.Map(),   //配置类型
		"status": config_status.Map(), //是否启用
	},
}
