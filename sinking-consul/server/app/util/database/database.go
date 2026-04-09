package database

import (
	"log"
	"os"
	"reflect"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Database 数据库连接
type Database struct {
	Db      *gorm.DB
	DbError error
}

// Transaction 事务执行
func (d *Database) Transaction(fc func(tx *gorm.DB) error) error {
	if d == nil || d.Db == nil {
		return gorm.ErrInvalidDB
	}
	return d.Db.Transaction(fc)
}

// BatchExecute	批量执行
func (d *Database) BatchExecute(slice interface{}, batchSize int, fn func(batch interface{}) error) error {
	val := reflect.ValueOf(slice)
	if val.Kind() != reflect.Slice {
		return nil
	}
	total := val.Len()
	if total == 0 {
		return nil
	}
	if batchSize <= 0 {
		batchSize = 1000
	}
	for i := 0; i < total; i += batchSize {
		end := i + batchSize
		if end > total {
			end = total
		}
		batch := val.Slice(i, end).Interface()
		if err := fn(batch); err != nil {
			return err
		}
	}
	return nil
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
