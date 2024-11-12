package config

func Init() {
	LoadConfig("./config/", "application.yml", "yml")
}
