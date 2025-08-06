package save

import "time"

// Init 定时保存数据
func Init() {
	go func() {
		for {
			time.Sleep(time.Minute)
			// 执行保存操作

		}
	}()
}
