package entity

import (
	"encoding/base64"
	"server/app/util/ip/czdb/aesECB"
	"server/app/util/ip/czdb/byteUtil"
)

/*DecryptedBlock
 * 表示加密前的解密数据块,包含客户端ID、过期日期和IP数据库的起始指针。
 * 该结构体提供了将这些数据序列化为字节数组的功能,并使用AES加密。
 *
 * +----------------+----------------+----------------+----------------+
 * | clientId       | expirationDate |                |                |
 * | (12 位)        | (20 位)        |                |                |
 * +----------------+----------------+----------------+----------------+
 * | randomSize (32 位)                                                |
 * +--------------------------------------------------------------------+
 * | 保留 (64 位)                                                       |
 * +--------------------------------------------------------------------+
 */
type DecryptedBlock struct {
	// clientId 占据4字节段的前12位
	clientId int
	// expirationDate 占据4字节段的后20位,格式为yyMMdd
	expirationDate int
	// 随机字节的大小
	randomSize int
}

// GetClientId 获取客户端ID
func (d *DecryptedBlock) GetClientId() int {
	return d.clientId
}

// SetClientId 设置客户端ID
func (d *DecryptedBlock) SetClientId(clientId int) {
	d.clientId = clientId
}

// GetExpirationDate 获取过期日期
func (d *DecryptedBlock) GetExpirationDate() int {
	return d.expirationDate
}

// SetExpirationDate 设置过期日期
func (d *DecryptedBlock) SetExpirationDate(expirationDate int) {
	d.expirationDate = expirationDate
}

// GetRandomSize 获取随机字节的大小
// 此方法返回加密过程中使用的随机字节的大小
// 随机字节的大小对于确保加密的唯一性和安全性至关重要
func (d *DecryptedBlock) GetRandomSize() int {
	return d.randomSize
}

// SetRandomSize 设置随机字节的大小
// 此方法允许设置加密过程中要使用的随机字节的大小
// 调整随机字节的大小可能会影响加密的安全性和性能
func (d *DecryptedBlock) SetRandomSize(randomSize int) {
	d.randomSize = randomSize
}

// ToBytes 将DecryptedBlock实例序列化为字节数组
// 数组结构如下：前4字节包含客户端ID和过期日期,
// 接下来的4字节包含随机大小,最后8字节保留并初始化为0
func (d *DecryptedBlock) ToBytes() []byte {
	b := make([]byte, 16)
	byteUtil.WriteIntLong(b, 0, int64((d.clientId<<20)|d.expirationDate))
	byteUtil.WriteIntLong(b, 4, int64(d.randomSize))
	// 保留的8字节默认已初始化为0
	return b
}

// Encrypt 使用指定的密钥对提供的字节数组进行AES加密
// 密钥应为表示AES密钥的base64编码字符串
func (d *DecryptedBlock) Encrypt(data []byte, key string) ([]byte, error) {
	// 解码base64编码的密钥
	keyBytes, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return nil, err
	}
	res, err := aesECB.AesEncrypt(data, keyBytes)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// ToEncryptedBytes 将DecryptedBlock加密并返回加密后的字节
func (d *DecryptedBlock) ToEncryptedBytes(key string) ([]byte, error) {
	return d.Encrypt(d.ToBytes(), key)
}

// Decrypt 解密加密的字节并返回DecryptedBlock
func Decrypt(key string, encryptedBytes []byte) (*DecryptedBlock, error) {
	// 解码base64编码的密钥
	keyBytes, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return nil, err
	}

	decryptedBytes, err := aesECB.AesDecrypt(encryptedBytes, keyBytes)
	if err != nil {
		return nil, err
	}

	// 解析解密后的字节
	decryptedBlock := &DecryptedBlock{}
	decryptedBlock.SetClientId(int(byteUtil.GetIntLong(decryptedBytes, 0) >> 20))
	decryptedBlock.SetExpirationDate(int(byteUtil.GetIntLong(decryptedBytes, 0) & 0xFFFFF))
	decryptedBlock.SetRandomSize(int(byteUtil.GetIntLong(decryptedBytes, 4)))

	return decryptedBlock, nil
}
