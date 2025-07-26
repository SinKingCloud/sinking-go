package byteUtil

// Package Utils 提供了用于操作字节数组的工具函数

// Write 从给定的偏移量开始向字节数组写入指定字节
//
// 参数:
//
//	b: 要写入的字节数组
//	offset: 数组中开始写入的位置
//	v: 要写入的值
//	bytes: 要写入的字节数
func Write(b []byte, offset int, v int64, bytes int) {
	for i := 0; i < bytes; i++ {
		b[offset+i] = byte((v >> (8 * i)) & 0xFF)
	}
}

// WriteIntLong 向字节数组写入一个整数
//
// 参数:
//
//	b: 要写入的字节数组
//	offset: 数组中开始写入的位置
//	v: 要写入的值
func WriteIntLong(b []byte, offset int, v int64) {
	b[offset] = byte(v & 0xFF)
	b[offset+1] = byte((v >> 8) & 0xFF)
	b[offset+2] = byte((v >> 16) & 0xFF)
	b[offset+3] = byte((v >> 24) & 0xFF)
}

// GetIntLong 从指定偏移量开始从字节数组中获取一个整数
//
// 参数:
//
//	b: 要读取的字节数组
//	offset: 数组中开始读取的位置
//
// 返回:
//
//	从字节数组中读取的整数值
func GetIntLong(b []byte, offset int) int64 {
	return int64(b[offset]) |
		int64(b[offset+1])<<8 |
		int64(b[offset+2])<<16 |
		int64(b[offset+3])<<24
}

// GetInt3 从指定偏移量开始从字节数组中获取一个3字节整数
//
// 参数:
//
//	b: 要读取的字节数组
//	offset: 数组中开始读取的位置
//
// 返回:
//
//	从字节数组中读取的整数值
func GetInt3(b []byte, offset int) int {
	return int(b[offset]) |
		int(b[offset+1])<<8 |
		int(b[offset+2])<<16
}

// GetInt2 从指定偏移量开始从字节数组中获取一个2字节整数
//
// 参数:
//
//	b: 要读取的字节数组
//	offset: 数组中开始读取的位置
//
// 返回:
//
//	从字节数组中读取的整数值
func GetInt2(b []byte, offset int) int {
	return int(b[offset]) |
		int(b[offset+1])<<8
}

// GetInt1 从指定偏移量开始从字节数组中获取一个1字节整数
//
// 参数:
//
//	b: 要读取的字节数组
//	offset: 数组中开始读取的位置
//
// 返回:
//
//	从字节数组中读取的整数值
func GetInt1(b []byte, offset int) int {
	return int(b[offset])
}
