package database

import (
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/model"
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/util/logs"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

// ConnectionSqlite 数据库初始函数
func ConnectionSqlite(file string, prefix string) bool {
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
	model.Db, model.DbError = gorm.Open(sqlite.Open(file), &gorm.Config{
		Logger: newLogger,
	})
	if model.DbError != nil {
		log.Println(model.DbError.Error())
		return false
	}
	model.DbPrefix = prefix
	logs.Println("数据库链接成功")
	return true
}
