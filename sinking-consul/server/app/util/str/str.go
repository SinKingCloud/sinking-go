package str

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/axgle/mahonia"
	"golang.org/x/crypto/bcrypt"
	"io"
	rand2 "math/rand"
	"strings"
	"time"
)

// StringTool 字符串工具
type StringTool struct {
	string string
}

// NewStringTool 实例化工具类
func NewStringTool(str string) *StringTool {
	return &StringTool{string: str}
}

// NewStrTool 实例化工具类
func NewStrTool() *StringTool {
	return &StringTool{}
}

// GetBetween 获取文本中间字符
func (stringTool *StringTool) GetBetween(start string, end string) string {
	n := strings.Index(stringTool.string, start)
	if n == -1 {
		n = 0
	}
	stringTool.string = string([]byte(stringTool.string)[n:])
	m := strings.Index(stringTool.string, end)
	if m == -1 {
		m = len(stringTool.string)
	}
	stringTool.string = string([]byte(stringTool.string)[len(start):m])
	return stringTool.string
}

// Md5 获取字符串md5值
func (stringTool *StringTool) Md5() string {
	h := md5.New()
	h.Write([]byte(stringTool.string))
	return hex.EncodeToString(h.Sum(nil))
}

// AesCbcEncrypt aes cbc加密
func (stringTool *StringTool) AesCbcEncrypt(key string) string {
	plainByte := []byte(stringTool.string)
	keyByte := []byte(key)
	if len(plainByte)%aes.BlockSize != 0 {
		return ""
	}
	block, err := aes.NewCipher(keyByte)
	if err != nil {
		return ""
	}
	cipherByte := make([]byte, aes.BlockSize+len(plainByte))
	iv := cipherByte[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return ""
	}
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherByte[aes.BlockSize:], plainByte)
	return fmt.Sprintf("%x\n", cipherByte)
}

// AesCbcDecrypt aes cbc解密
func (stringTool *StringTool) AesCbcDecrypt(key string) string {
	cipherByte, _ := hex.DecodeString(stringTool.string)
	keyByte := []byte(key)
	block, err := aes.NewCipher(keyByte)
	if err != nil {
		return ""
	}
	if len(cipherByte) < aes.BlockSize {
		return ""
	}
	iv := cipherByte[:aes.BlockSize]
	cipherByte = cipherByte[aes.BlockSize:]
	if len(cipherByte)%aes.BlockSize != 0 {
		return ""
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(cipherByte, cipherByte)
	return string(cipherByte[:])
}

func (stringTool *StringTool) AesEncrypt(key string) string {
	origData := []byte(stringTool.string)
	k := []byte(key)
	block, _ := aes.NewCipher(k)
	blockSize := aes.BlockSize
	origData = stringTool.pKCS7Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, k[:blockSize])
	cryted := make([]byte, len(origData))
	blockMode.CryptBlocks(cryted, origData)
	return base64.StdEncoding.EncodeToString(cryted)
}

func (stringTool *StringTool) AesDecrypt(key string) string {
	crytedByte, _ := base64.StdEncoding.DecodeString(stringTool.string)
	k := []byte(key)
	block, _ := aes.NewCipher(k)
	blockSize := aes.BlockSize
	blockMode := cipher.NewCBCDecrypter(block, k[:blockSize])
	orig := make([]byte, len(crytedByte))
	blockMode.CryptBlocks(orig, crytedByte)
	orig = stringTool.pKCS7UnPadding(orig)
	return string(orig)
}

func (StringTool) pKCS7Padding(ciphertext []byte, blocksize int) []byte {
	padding := blocksize - len(ciphertext)%blocksize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padText...)
}

func (StringTool) pKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unPadding := int(origData[length-1])
	return origData[:(length - unPadding)]
}

// GetPassword 获取密码
func (stringTool StringTool) GetPassword() string {
	hash, err := bcrypt.GenerateFromPassword([]byte(stringTool.string), bcrypt.DefaultCost) //加密处理
	if err != nil {
		return ""
	}
	return string(hash)
}

// CheckPassword 比对密码
func (stringTool StringTool) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(stringTool.string), []byte(password)) //验证（对比）
	if err != nil {
		return false
	} else {
		return true
	}
}

// GenValidateCode 生成指定位数验证码
func (stringTool StringTool) GenValidateCode(width int) string {
	numeric := [10]byte{1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numeric)
	rand2.Seed(time.Now().UnixNano())
	var sb strings.Builder
	for i := 0; i < width; i++ {
		_, err := fmt.Fprintf(&sb, "%d", numeric[rand2.Intn(r)])
		if err != nil {
			return ""
		}
	}
	return sb.String()
}

func (stringTool StringTool) RandInt64(min, max int64) int64 {
	if min >= max || min == 0 || max == 0 {
		return max
	}
	return rand2.Int63n(max-min) + min
}

// Show 显示首尾数据字符个数
func (stringTool StringTool) Show(left int, right int, str string) string {
	if stringTool.string == "" {
		return ""
	}
	nameRune := []rune(stringTool.string)
	lens := len(nameRune)
	if lens-right < 0 || left > lens {
		return ""
	}
	leftStr := nameRune[0:left]
	rightStr := nameRune[lens-right : lens]
	if len(leftStr)+len(rightStr) >= lens {
		return string(leftStr) + string(rightStr)
	}
	temp := ""
	for i := 0; i < lens-len(leftStr)-len(rightStr); i++ {
		temp += str
	}
	return string(leftStr) + temp + string(rightStr)
}

// ConvertToString 编码转换
func (stringTool StringTool) ConvertToString(srcCode string, tagCode string) string {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(stringTool.string)
	tagCoder := mahonia.NewDecoder(tagCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	result := string(cdata)
	return result
}

// ToJson 转json
func (stringTool StringTool) ToJson(data interface{}) string {
	str, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	return string(str)
}

// ToMap 转map
func (stringTool StringTool) ToMap(data string) map[string]string {
	var d map[string]string
	e := json.Unmarshal([]byte(data), &d)
	if e == nil {
		return d
	}
	return nil
}
