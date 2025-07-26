package str

import (
	"crypto/md5"
	"encoding/hex"
)

// ByteTool 字节工具
type ByteTool struct {
	byte []byte
}

// NetByteTool 实例化工具类
func NetByteTool(str []byte) *ByteTool {
	return &ByteTool{byte: str}
}

// Md5 获取md5
func (s *ByteTool) Md5() string {
	hash := md5.New()
	hash.Write(s.byte)
	md5Str := hex.EncodeToString(hash.Sum(nil))
	return md5Str
}
