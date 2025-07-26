package log

type Type int //日志类型

const (
	User Type = iota //用户操作
)

// Types 类型数据
func (s *Service) Types() map[Type]string {
	return map[Type]string{
		User: "用户操作",
	}
}
