package database

import (
	"errors"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

// Database 数据库连接
type Database struct {
	Db      *gorm.DB
	DbError error
}

func newLogger() logger.Interface {
	return logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), //io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		logger.Config{
			SlowThreshold:             time.Second,   // 慢 SQL 阈值
			LogLevel:                  logger.Silent, // 日志级别
			IgnoreRecordNotFoundError: true,          // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  false,         // 禁用彩色打印
		},
	)
}

// NewSqlite 实例化一个sqlite连接
func NewSqlite(file string) *Database {
	Db, DbError := gorm.Open(sqlite.Open(file), &gorm.Config{
		Logger: newLogger(),
	})
	if DbError != nil {
		return &Database{Db: Db, DbError: errors.New("sql connect error")}
	}
	return &Database{Db: Db, DbError: DbError}
}
