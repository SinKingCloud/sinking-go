package config

type Status int

const (
	Normal Status = iota //正常
	Stop                 //暂停
)

// Status 状态数据
func (s *Service) Status() map[Status]string {
	return map[Status]string{
		Normal: "正常",
		Stop:   "暂停",
	}
}

type Type string

const (
	JSON Type = "json"
	YAML      = "yaml"
	INI       = "ini"
)

// Types 类型数据
func (s *Service) Types() map[Type]string {
	return map[Type]string{
		JSON: "JSON",
		YAML: "YAML",
		INI:  "INI",
	}
}
