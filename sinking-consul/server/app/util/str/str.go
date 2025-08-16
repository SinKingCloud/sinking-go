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
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	rand2 "math/rand"
	"strings"
)

// StringTool 字符串工具类（无状态，所有方法通过参数接收输入）
type StringTool struct{}

// NewStringTool 实例化工具类（单例模式，工具类无状态）
func NewStringTool() *StringTool {
	return &StringTool{}
}

// GetBetween 获取两个子串之间的内容（不含起止子串）
// 例如："abc123def" 中获取 "123"，start="abc", end="def"
// 返回：找到的内容，未找到时返回空字符串（非错误）
func (t *StringTool) GetBetween(str, start, end string) string {
	startIdx := strings.Index(str, start)
	if startIdx == -1 {
		return ""
	}
	// 跳过 start 子串
	startPos := startIdx + len(start)
	endIdx := strings.Index(str[startPos:], end)
	if endIdx == -1 {
		// 未找到 end，返回从 start 后到结尾的内容
		return str[startPos:]
	}
	// 截取 [startPos, startPos+endIdx)
	return str[startPos : startPos+endIdx]
}

// Md5 计算字符串的 MD5 哈希（32位小写）
func (t *StringTool) Md5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// AesCBCEncrypt AES-CBC 模式加密（自动 PKCS7 填充，返回 Base64 编码）
// key 长度必须为 16/24/32 字节（对应 AES-128/192/256）
// 返回：加密后的 Base64 字符串，错误信息
func (t *StringTool) AesCBCEncrypt(plaintext, key string) (string, error) {
	keyBytes := []byte(key)
	// 校验 key 长度
	switch len(keyBytes) {
	case 16, 24, 32:
	default:
		return "", fmt.Errorf("aes key length must be 16, 24, or 32 bytes")
	}

	// PKCS7 填充
	plainBytes := t.pkcs7Padding([]byte(plaintext), aes.BlockSize)

	// 创建加密块
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", fmt.Errorf("create cipher failed: %w", err)
	}

	// 生成随机 IV（长度=块大小）
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", fmt.Errorf("generate iv failed: %w", err)
	}

	// 加密
	ciphertext := make([]byte, len(plainBytes))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, plainBytes)

	// 拼接 IV 和密文，再 Base64 编码（IV 用于解密）
	result := append(iv, ciphertext...)
	return base64.StdEncoding.EncodeToString(result), nil
}

// AesCBCDecrypt AES-CBC 模式解密（Base64 解码输入，自动 PKCS7 去填充）
// key 长度必须为 16/24/32 字节（对应 AES-128/192/256）
// 返回：解密后的明文，错误信息
func (t *StringTool) AesCBCDecrypt(ciphertextBase64, key string) (string, error) {
	keyBytes := []byte(key)
	// 校验 key 长度
	switch len(keyBytes) {
	case 16, 24, 32:
	default:
		return "", fmt.Errorf("aes key length must be 16, 24, or 32 bytes")
	}

	// Base64 解码
	ciphertext, err := base64.StdEncoding.DecodeString(ciphertextBase64)
	if err != nil {
		return "", fmt.Errorf("base64 decode failed: %w", err)
	}

	// 分离 IV 和密文（IV 长度=块大小）
	if len(ciphertext) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	// 创建解密块
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", fmt.Errorf("create cipher failed: %w", err)
	}

	// 解密
	plaintext := make([]byte, len(ciphertext))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(plaintext, ciphertext)

	// 去填充
	plaintext, err = t.pkcs7Unpadding(plaintext)
	if err != nil {
		return "", fmt.Errorf("unpadding failed: %w", err)
	}

	return string(plaintext), nil
}

// PKCS7 填充（满足块大小要求）
func (t *StringTool) pkcs7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

// PKCS7 去填充
func (t *StringTool) pkcs7Unpadding(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("data is empty")
	}
	padding := int(data[length-1])
	if padding > length {
		return nil, errors.New("invalid padding")
	}
	return data[:length-padding], nil
}

// BcryptHash 使用 bcrypt 对密码进行哈希（自动生成盐值）
// 返回：哈希后的字符串，错误信息
func (t *StringTool) BcryptHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("bcrypt hash failed: %w", err)
	}
	return string(hash), nil
}

// BcryptVerify 验证密码与 bcrypt 哈希是否匹配
// 返回：true（匹配）/false（不匹配），错误信息（哈希无效等）
func (t *StringTool) BcryptVerify(password, hash string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return false, nil // 密码不匹配（非错误）
	}
	if err != nil {
		return false, fmt.Errorf("bcrypt verify failed: %w", err)
	}
	return true, nil
}

// GenValidateCode 生成指定位数的数字验证码（1-9，不含0）
// width: 验证码长度（必须>0）
// 返回：验证码字符串，错误信息
func (t *StringTool) GenValidateCode(width int) (string, error) {
	if width <= 0 {
		return "", errors.New("width must be greater than 0")
	}
	numeric := []byte{'1', '2', '3', '4', '5', '6', '7', '8', '9'}
	var buf bytes.Buffer
	for i := 0; i < width; i++ {
		buf.WriteByte(numeric[rand2.Intn(len(numeric))])
	}
	return buf.String(), nil
}

// RandInt64 生成 [min, max) 范围内的随机整数（min < max）
func (t *StringTool) RandInt64(min, max int64) (int64, error) {
	if min >= max {
		return 0, errors.New("min must be less than max")
	}
	return rand2.Int63n(max-min) + min, nil
}

// Mask 对字符串进行掩码处理（如手机号中间用*替换）
// left: 保留左侧字符数
// right: 保留右侧字符数
// mask: 掩码字符（如"*"）
// 返回：处理后的字符串，错误信息（参数无效时）
func (t *StringTool) Mask(str string, left, right int, mask string) (string, error) {
	runes := []rune(str)
	lens := len(runes)
	if left < 0 || right < 0 || left+right > lens {
		return "", errors.New("invalid left/right parameters")
	}
	if left == 0 && right == 0 {
		return strings.Repeat(mask, lens), nil
	}
	// 左侧保留部分
	leftPart := runes[:left]
	// 右侧保留部分
	rightPart := runes[lens-right:]
	// 中间掩码部分
	maskLen := lens - left - right
	maskPart := strings.Repeat(mask, maskLen)
	return string(leftPart) + maskPart + string(rightPart), nil
}

// ConvertEncoding 编码转换（替换过时的 macholib）
// srcEncoding: 源编码（如"gbk"、"utf-8"）
// dstEncoding: 目标编码（如"utf-8"、"gbk"）
// 返回：转换后的字符串，错误信息
func (t *StringTool) ConvertEncoding(str, srcEncoding, dstEncoding string) (string, error) {
	// 这里以常见的 GBK→UTF-8 为例，可扩展其他编码（需引入对应编码库）
	// 推荐使用 golang.org/x/text/encoding 下的编码实现
	var srcDecoder encoding.Encoding
	var dstDecoder encoding.Encoding

	switch srcEncoding {
	case "utf-8", "utf8":
		srcDecoder = encoding.Nop // 无操作（已为UTF-8）
	case "gbk":
		srcDecoder = simplifiedchinese.GBK // 需要引入：golang.org/x/text/encoding/simplification
	default:
		return "", fmt.Errorf("unsupported source encoding: %s", srcEncoding)
	}

	switch dstEncoding {
	case "utf-8", "utf8":
		dstDecoder = encoding.Nop
	case "gbk":
		dstDecoder = simplifiedchinese.GBK
	default:
		return "", fmt.Errorf("unsupported destination encoding: %s", dstEncoding)
	}

	// 先将源编码转为 UTF-8，再转为目标编码
	// 步骤1：源编码 → UTF-8
	utf8Bytes, _, err := transform.Bytes(srcDecoder.NewDecoder(), []byte(str))
	if err != nil {
		return "", fmt.Errorf("convert to utf8 failed: %w", err)
	}

	// 步骤2：UTF-8 → 目标编码
	dstBytes, _, err := transform.Bytes(dstDecoder.NewEncoder(), utf8Bytes)
	if err != nil {
		return "", fmt.Errorf("convert to destination encoding failed: %w", err)
	}

	return string(dstBytes), nil
}

// ToJSON 将任意数据序列化为 JSON 字符串
func (t *StringTool) ToJSON(data interface{}) (string, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("json marshal failed: %w", err)
	}
	return string(b), nil
}

// FromJSON 将 JSON 字符串反序列化为指定类型（推荐使用指针接收）
func (t *StringTool) FromJSON(jsonStr string, v interface{}) error {
	if err := json.Unmarshal([]byte(jsonStr), v); err != nil {
		return fmt.Errorf("json unmarshal failed: %w", err)
	}
	return nil
}

// ToNumber 将字符串转换为唯一数字（紧凑标识）
func (t *StringTool) ToNumber(s string, max uint64) uint64 {
	var hash uint64 = 14695981039346656037
	for i, c := range s {
		if i >= 32 {
			break
		}
		hash ^= uint64(c)
		hash *= 1099511628211
	}
	if max > 0 {
		hash ^= hash >> 20
		return hash % (max + 1)
	}
	return hash
}
