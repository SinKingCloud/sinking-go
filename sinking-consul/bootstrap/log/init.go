package log

import "github.com/SinKingCloud/sinking-go/sinking-consul/app/util/setting"

func Init() {
	if setting.GetConfig() != nil &&
		setting.GetConfig().GetString("database.logstash.host") != "" &&
		setting.GetConfig().GetInt("database.logstash.port") != 0 {
		ConnectionLogServer(
			setting.GetConfig().GetString("database.logstash.host"),
			setting.GetConfig().GetInt("database.logstash.port"),
			setting.GetConfig().GetInt("database.logstash.timeout"),
		)
	}
}
