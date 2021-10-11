package setting

import "github.com/spf13/viper"

var (
	system        *viper.Viper
	systemSetting *SystemConfig
)

// SetSetting 设置数据
func SetSetting(setting *viper.Viper) {
	system = setting
}

// GetConfig 获取设置数据
func GetConfig() *viper.Viper {
	return system
}

// SetSystemSetting 设置数据
func SetSystemSetting(setting *SystemConfig) {
	systemSetting = setting
}

// GetSystemConfig 获取设置数据
func GetSystemConfig() *SystemConfig {
	return systemSetting
}

// SystemConfig 设置实体
type SystemConfig struct {
	Name               string `json:"name"`
	Description        string `json:"description"`
	MinimumCoreVersion string `json:"minimumCoreVersion"`
	Title              string `json:"title"`
	Version            string `json:"version"`
	Author             string `json:"author"`
	App                struct {
		Debug bool   `json:"debug"`
		Ip    string `json:"ip"`
		Port  int    `json:"port"`
	} `json:"app"`
	Database struct {
		Sql    string `json:"sql"`
		Sqlite struct {
			Database string `json:"database"`
			Prefix   string `json:"prefix"`
		} `json:"sqlite"`
		Mysql struct {
			Host     string `json:"host"`
			Port     string `json:"port"`
			User     string `json:"user"`
			Pwd      string `json:"pwd"`
			Database string `json:"database"`
			Prefix   string `json:"prefix"`
		} `json:"mysql"`
		Logstash struct {
			Host    string `json:"host"`
			Port    int    `json:"port"`
			Timeout int    `json:"timeout"`
		} `json:"logstash"`
	} `json:"database"`
	Servers struct {
		Cluster   string `json:"cluster"`
		TokenName string `json:"token-name"`
		Token     string `json:"token"`
	} `json:"servers"`
}
