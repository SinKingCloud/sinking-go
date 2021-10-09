package log

import (
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/util/logs"
	"log"
	"os"
	"strconv"
)

func ConnectionLogServer(host string, port int, timeOut int) {
	log.SetPrefix("thread-" + strconv.Itoa(os.Getppid()) + " ")
	log.SetFlags(log.Lmicroseconds | log.Ldate | log.Lmsgprefix)
	log.Println("正在链接Logstash...")
	client := logs.New(host, port, timeOut)
	_, err := client.Connect()
	if err != nil {
		log.Println(err)
	}
	logs.SetLogClient(client)
	log.Println("链接Logstash成功")
}
