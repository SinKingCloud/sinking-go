package sinking_consul

import (
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/util/setting"
	"github.com/SinKingCloud/sinking-go/sinking-consul/bootstrap/config"
	"github.com/SinKingCloud/sinking-go/sinking-consul/bootstrap/database"
	"github.com/SinKingCloud/sinking-go/sinking-consul/bootstrap/log"
)

func main() {
	//加载系统配置
	config.LoadConfig("./sinking-spider/config/", "system.json", "json")
	//链接elk日志
	log.ConnectionLogServer(
		setting.GetConfig().GetString("database.logstash.host"),
		setting.GetConfig().GetInt("database.logstash.port"),
		setting.GetConfig().GetInt("database.logstash.timeout"),
	)
	//链接数据库
	database.GetMysqlInstance().ConnectionMysql(
		setting.GetConfig().GetString("database.mysql.host"),
		setting.GetConfig().GetString("database.mysql.port"),
		setting.GetConfig().GetString("database.mysql.user"),
		setting.GetConfig().GetString("database.mysql.pwd"),
		setting.GetConfig().GetString("database.mysql.database"),
		setting.GetConfig().GetString("database.mysql.prefix"),
	)
	//启动应用
}
