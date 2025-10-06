package constant

const (
	BasePath = "." //基础目录

	DBPath = BasePath + "/config" //数据库文件目录
	DBFile = "server.db"          //数据库文件

	ConfPath = BasePath + "/config" //配置文件目录
	ConfFile = "application"        //配置文件

	ServerMode = "server.mode" //运行模式
	ServerHost = "server.host" //服务器地址
	ServerPort = "server.port" //服务器端口

	AuthAccount    = "auth.account"  //登录账号
	AuthPassword   = "auth.password" //登录密码
	AuthApiToken   = "auth.token"    //API访问token
	AuthLoginToken = "auth.session"  //登录token
	AuthExpire     = "auth.expire"   //登录过期时间

	ClusterLocal = "cluster.local" //本机地址
	ClusterNodes = "cluster.nodes" //集群服务器列表

	WebTitle    = "web.title"    //网站标题
	WebName     = "web.name"     //网站名称
	WebKeyWords = "web.keywords" //网站关键字
	WebDescribe = "web.describe" //网站描述

	UiLayout    = "ui.layout"    //界面布局
	UiWaterMark = "ui.watermark" //水印内容
	UiTheme     = "ui.theme"     //主题
	UiCompact   = "ui.compact"   //紧凑模式
	UiColor     = "ui.color"     //主题色
	UiRadius    = "ui.radius"    //主题圆角
)
