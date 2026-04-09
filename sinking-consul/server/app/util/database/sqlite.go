package database

import (
	"errors"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

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
