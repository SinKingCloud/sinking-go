package log

import (
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/util/logs"
	"log"
)

func ConnectionLogServer(host string, port int, timeOut int) {
	log.Println("正在链接Logstash...")
	client := logs.New(host, port, timeOut)
	_, err := client.Connect()
	if err != nil {
		log.Println(err)
	} else {
		logs.SetLogClient(client)
		log.Println("链接Logstash成功")
	}
}
