package database

import (
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/util/logs"
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/util/setting"
	"os"
)

func Init() {
	if setting.GetSystemConfig() != nil {
		res := false
		if setting.GetSystemConfig().Database.Sql == "mysql" {
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
		setting.GetSystemConfig().Database.Mysql.Host,
		setting.GetSystemConfig().Database.Mysql.Port,
		setting.GetSystemConfig().Database.Mysql.User,
		setting.GetSystemConfig().Database.Mysql.Pwd,
		setting.GetSystemConfig().Database.Mysql.Database,
		setting.GetSystemConfig().Database.Mysql.Prefix,
	)
}

func sqliteConn() bool {
	return ConnectionSqlite(
		setting.GetSystemConfig().Database.Sqlite.Database,
		setting.GetSystemConfig().Database.Sqlite.Prefix,
	)
}
