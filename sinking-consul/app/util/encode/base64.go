package encode

import (
	"encoding/base64"
	"unsafe"
)

var base64Table = "IJjkKLMNO567PQX12RVW3YZaDEFGbcdefghiABCHlSTUmnopqrxyz04stuvw89+/"

func Base64Encode(data string) string {
	content := *(*[]byte)(unsafe.Pointer(&data))
	coder := base64.NewEncoding(base64Table)
	return coder.EncodeToString(content)
}

func Base64Decode(data string) string {
	coder := base64.NewEncoding(base64Table)
	result, _ := coder.DecodeString(data)
	return *(*string)(unsafe.Pointer(&result))
}
