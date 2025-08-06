package config

type Type string

const (
	JSON Type = "json"
	YML       = "yml"
	INI       = "ini"
)

// Types 类型数据
func (s *Service) Types() map[Type]string {
	return map[Type]string{
		JSON: "JSON",
		YML:  "YML",
		INI:  "INI",
	}
}
