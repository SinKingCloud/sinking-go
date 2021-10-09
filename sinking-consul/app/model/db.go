package model

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"strings"
	"time"
)

var (
	Db       *gorm.DB
	DbError  error
	DbPrefix string
)

// DateTime 返回日期格式
type DateTime time.Time

func (t DateTime) MarshalJSON() ([]byte, error) {
	var stamp = fmt.Sprintf("\"%s\"", time.Time(t).Format("2006-01-02 15:04:05"))
	return []byte(stamp), nil
}
func (t *DateTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	var err error
	str := string(data)
	timeStr := strings.Trim(str, "\"")
	t1, err := time.Parse("2006-01-02 15:04:05", timeStr)
	*t = DateTime(t1)
	return err
}

func (t DateTime) Value() (driver.Value, error) {
	tTime := time.Time(t)
	return tTime.Format("2006-01-02 15:04:05"), nil
}

func (t *DateTime) Scan(v interface{}) error {
	switch vt := v.(type) {
	case time.Time:
		*t = DateTime(vt)
	default:
		return errors.New("类型处理错误")
	}
	return nil
}

func (t *DateTime) String() string {
	return fmt.Sprintf("hhh:%s", time.Time(*t).String())
}
