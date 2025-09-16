package client

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func Test(t *testing.T) {
	getServer()
}

func getServer() {
	cli := NewClient([]string{"43.248.186.86:5679", "101.133.153.115:3455"}, "sinking-cloud", "service-gateway", "127.0.0.1:1000", "E10ADC3949BA59ABBE56E057F20F883E")
	err, list := cli.getConfigList(cli.getAddr(false), 0)
	if err == nil {
		for _, v := range list {
			fmt.Println(v.Group, v.Name, v.Type, v.Status, v.Hash, v.Content)
		}
	}
	_ = cli.register(cli.getAddr(false))
	err2, list2 := cli.getNodeList(cli.getAddr(false), 0)
	if err2 == nil {
		for _, v := range list2 {
			fmt.Println(v.Group, v.Name, v.Address, v.Status, v.LastHeart, v.OnlineStatus)
		}
	}
}

// 演示JSON格式解析
func demoJSON() {
	fmt.Println("===== 演示JSON格式解析 =====")
	jsonConfig := `{
		"server": {
			"host": "localhost",
			"port": 8080,
			"ssl": true,
			"timeout": "30s"
		},
		"database": {
			"maxConnections": 10,
			"enabled": false,
			"params": {
				"charset": "utf8"
			}
		},
		"users": ["admin", "guest"],
		"scores": [90.5, 85.3, 92.0],
		"startTime": "2024-01-01T08:00:00Z"
	}`

	parser, err := NewConfigParser(jsonConfig, "json")
	if err != nil {
		log.Fatalf("创建JSON解析器失败: %v", err)
	}

	// 获取基本类型
	host := parser.GetString("server.host")
	port := parser.GetInt("server.port")
	ssl := parser.GetBool("server.ssl")
	timeout := parser.GetDuration("server.timeout")
	fmt.Printf("服务器配置: %s:%d (SSL: %v, 超时: %v)\n", host, port, ssl, timeout)

	// 获取数组
	users := parser.GetStringSlice("users")
	fmt.Printf("用户列表: %v\n", users)

	// 获取嵌套Map
	dbParams := parser.GetStringMapString("database.params")
	fmt.Printf("数据库参数: %v\n", dbParams)

	// 获取时间
	startTime := parser.GetTime("startTime")
	fmt.Printf("开始时间: %s\n", startTime.Format("2006-01-02 15:04:05"))

	// 检查路径是否存在
	fmt.Printf("是否存在 'database.maxConnections'? %v\n", parser.IsSet("database.maxConnections"))
	fmt.Printf("是否存在 'unknown.key'? %v\n", parser.IsSet("unknown.key"))
	fmt.Println()
}

// 演示YAML格式解析（含嵌套数组）
func demoYAML() {
	fmt.Println("===== 演示YAML格式解析 =====")
	yamlConfig := `
app:
  name: "demo"
  version: 1.0.2
  features:
    - "auth"
    - "log"
    - "cache"
  nestedArray:
    - [10, 20, 30]
    - [40, 50, 60]
  limits:
    maxSize: 1024
    minSize: 64
`

	parser, err := NewConfigParser(yamlConfig, "yaml")
	if err != nil {
		log.Fatalf("创建YAML解析器失败: %v", err)
	}

	// 获取嵌套数组（修复前的版本不支持，修复后可正常使用）
	// 注意：如果使用未修复嵌套数组的版本，此行会报错
	val := parser.GetInt("app.nestedArray[0][1]")
	if err != nil {
		fmt.Printf("获取嵌套数组元素失败（未修复版本会出现此提示）: %v\n", err)
	} else {
		fmt.Printf("嵌套数组元素 app.nestedArray[0][1]: %d\n", val)
	}

	// 获取版本号（float64转string）
	version := parser.GetFloat64("app.version")
	fmt.Printf("应用版本: %.1f\n", version)

	// 获取特性列表
	features := parser.GetStringSlice("app.features")
	fmt.Printf("支持特性: %v\n", features)
	fmt.Println()
}

// 演示INI格式解析
func demoINI() {
	fmt.Println("===== 演示INI格式解析 =====")
	iniConfig := `
# 服务器配置
[server]
host = 127.0.0.1
port = 9000
timeout = 60s

# 日志配置
[log]
level = info
path = /var/log/app.log
enabled = true
`

	parser, err := NewConfigParser(iniConfig, "ini")
	if err != nil {
		log.Fatalf("创建INI解析器失败: %v", err)
	}

	// 获取INI配置（注意INI的section会作为前缀）
	logLevel := parser.GetString("log.level")
	port := parser.GetInt("server.port")
	timeout := parser.GetDuration("server.timeout")
	fmt.Printf("日志级别: %s, 端口: %d, 超时: %v\n", logLevel, port, timeout)
	fmt.Println()
}

// 演示动态更新配置
func demoUpdateConfig() {
	fmt.Println("===== 演示动态更新配置 =====")
	// 初始配置
	oldConfig := `{"name": "oldApp", "version": 1}`
	parser, _ := NewConfigParser(oldConfig, "json")
	name := parser.GetString("name")
	fmt.Printf("更新前名称: %s\n", name)

	// 更新配置
	newConfig := `{"name": "newApp", "version": 2, "author": "dev"}`
	if err := parser.UpdateConfig(newConfig, "json"); err != nil {
		log.Fatalf("更新配置失败: %v", err)
	}

	// 验证更新结果
	newName := parser.GetString("name")
	author := parser.GetString("author")
	fmt.Printf("更新后名称: %s, 作者: %s\n", newName, author)
	fmt.Println()
}

// 演示缓存功能
func demoCache() {
	fmt.Println("===== 演示缓存功能 =====")
	config := `{"a": 100, "b": "test", "c": [1,2,3]}`
	parser, _ := NewConfigParser(config, "json")

	// 首次获取（无缓存）
	start := time.Now()
	_ = parser.GetInt("a")
	fmt.Printf("首次获取耗时: %v\n", time.Since(start))

	// 第二次获取（有缓存）
	start = time.Now()
	_ = parser.GetInt("a")
	fmt.Printf("缓存获取耗时: %v\n", time.Since(start))

	// 禁用缓存
	parser.EnableCache(false)
	start = time.Now()
	_ = parser.GetInt("a")
	fmt.Printf("禁用缓存后耗时: %v\n", time.Since(start))

	// 清空指定路径缓存
	parser.EnableCache(true)
	parser.InvalidateCache("a") // 仅清空a的缓存
	fmt.Println()
}
