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

type IsDelete int

const (
	False IsDelete = iota
	True
)

// IsDelete 是否删除
func (s *Service) IsDelete() map[IsDelete]string {
	return map[IsDelete]string{
		False: "未删除",
		True:  "已删除",
	}
}
