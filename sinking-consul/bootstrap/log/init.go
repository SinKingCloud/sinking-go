package log

import "github.com/SinKingCloud/sinking-go/sinking-consul/app/util/setting"

func Init() {
	if setting.GetSystemConfig() != nil &&
		setting.GetSystemConfig().Database.Logstash.Host != "" &&
		setting.GetSystemConfig().Database.Logstash.Port != 0 {
		ConnectionLogServer(
			setting.GetSystemConfig().Database.Logstash.Host,
			setting.GetSystemConfig().Database.Logstash.Port,
			setting.GetSystemConfig().Database.Logstash.Timeout,
		)
	}
}
