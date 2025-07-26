package str

import (
	"time"
)

// TimeTool 时间工具
type TimeTool struct {
}

// NewTimeTool 实例化工具类
func NewTimeTool() *TimeTool {
	return &TimeTool{}
}

// GetFirstDateOfMonth 获取某月第一天
func (t *TimeTool) GetFirstDateOfMonth(d time.Time) time.Time {
	d = d.AddDate(0, 0, -d.Day()+1)
	return t.GetZeroTime(d)
}

// GetLastDateOfMonth 获取传入的时间所在月份的最后一天，即某月最后一天的0点。如传入time.Now(), 返回当前月份的最后一天0点时间。
func (t *TimeTool) GetLastDateOfMonth(d time.Time) time.Time {
	return t.GetFirstDateOfMonth(d).AddDate(0, 1, -1)
}

// GetZeroTime 获取某一天的0点时间
func (*TimeTool) GetZeroTime(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
}
