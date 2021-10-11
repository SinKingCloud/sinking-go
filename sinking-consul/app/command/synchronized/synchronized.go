package synchronized

import (
	"fmt"
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/service"
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/util/logs"
	"time"
)

func Run() {
	logs.Println("开始执行数据同步任务")
	for {
		for _, v := range service.Clusters {
			test, _ := v.LastHeartTime.Value()
			fmt.Println(test, v.Status)
		}
		time.Sleep(time.Second)
	}
}
