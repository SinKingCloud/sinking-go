package logs

import (
	"encoding/json"
	"log"
)

var logs *Logstash

// SetLogClient 设置log对象
func SetLogClient(client *Logstash) *Logstash {
	logs = client
	return logs
}

// GetLogClient 获取log对象
func GetLogClient() *Logstash {
	return logs
}

func Println(data ...interface{}) {
	message, err := json.Marshal(data)
	if err == nil && logs != nil {
		err = logs.Writeln(string(message))
		if err != nil {
			log.Println("Save log err:", err, data)
		} else {
			log.Println(data)
		}
	} else {
		log.Println(data)
	}
}
