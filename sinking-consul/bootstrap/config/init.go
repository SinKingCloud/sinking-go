package config

func Init() {
	LoadConfig("./config/", "system.json", "json")
}
