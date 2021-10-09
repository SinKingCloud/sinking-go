package database

import (
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/util/logs"
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/util/setting"
	"os"
)

func Init() {
	if setting.GetConfig() != nil {
		res := GetMysqlInstance().ConnectionMysql(
			setting.GetConfig().GetString("database.mysql.host"),
			setting.GetConfig().GetString("database.mysql.port"),
			setting.GetConfig().GetString("database.mysql.user"),
			setting.GetConfig().GetString("database.mysql.pwd"),
			setting.GetConfig().GetString("database.mysql.database"),
			setting.GetConfig().GetString("database.mysql.prefix"),
		)
		if res {
			os.Exit(0)
		}
	} else {
		logs.Println("获取配置失败，请检查配置文件")
		os.Exit(0)
	}
}
