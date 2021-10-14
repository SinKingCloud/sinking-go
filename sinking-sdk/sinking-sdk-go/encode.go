package sinking_sdk_go

import (
	"crypto/md5"
	"encoding/hex"
)

// Md5Encode md5加密
func Md5Encode(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}
