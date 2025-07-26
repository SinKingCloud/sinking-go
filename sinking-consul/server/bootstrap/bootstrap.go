package bootstrap

func Load() {
	LoadConf()
	LoadLog()
	LoadCache()
	LoadDatabase()
}
