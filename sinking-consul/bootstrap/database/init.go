package database

import (
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/util/logs"
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/util/setting"
	"os"
)

func Init() {
	if setting.GetConfig() != nil {
		res := false
		if setting.GetConfig().GetString("database.sql") == "mysql" {
			res = mysqlConn()
		} else {
			res = sqliteConn()
		}
		if !res {
			os.Exit(0)
		}
	} else {
		logs.Println("获取配置失败，请检查配置文件")
		os.Exit(0)
	}
}

func mysqlConn() bool {
	return GetMysqlInstance().ConnectionMysql(
		setting.GetConfig().GetString("database.mysql.host"),
		setting.GetConfig().GetString("database.mysql.port"),
		setting.GetConfig().GetString("database.mysql.user"),
		setting.GetConfig().GetString("database.mysql.pwd"),
		setting.GetConfig().GetString("database.mysql.database"),
		setting.GetConfig().GetString("database.mysql.prefix"),
	)
}

func sqliteConn() bool {
	return ConnectionSqlite(
		setting.GetConfig().GetString("database.sqlite.database"),
		setting.GetConfig().GetString("database.sqlite.prefix"),
	)
}
