package database

import (
	"fmt"
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/model"
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/util/logs"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"sync"
	"time"
)

type MySqlPool struct{}

var instance *MySqlPool
var once sync.Once

// GetMysqlInstance 单例模式
func GetMysqlInstance() *MySqlPool {
	once.Do(func() {
		instance = &MySqlPool{}
	})
	return instance
}

// ConnectionMysql 数据库初始函数
func (pool *MySqlPool) ConnectionMysql(host string, port string, user string, pwd string, database string, prefix string) bool {
	logs.Println("正在链接数据库...")
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), //io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		logger.Config{
			SlowThreshold:             time.Second,   // 慢 SQL 阈值
			LogLevel:                  logger.Silent, // 日志级别
			IgnoreRecordNotFoundError: true,          // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  false,         // 禁用彩色打印
		},
	)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, pwd, host, port, database)
	model.Db, model.DbError = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if model.DbError != nil {
		logs.Println(model.DbError.Error())
		return false
	}
	model.DbPrefix = prefix
	logs.Println("数据库链接成功")
	return true
}
