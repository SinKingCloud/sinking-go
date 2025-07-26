package hyperHeaderDecoder

import (
	"fmt"
	"io"
	"server/app/util/ip/czdb/byteUtil"
	"server/app/util/ip/czdb/entity"
	"strconv"
	"time"
)

// package HyperHeaderDecoder 提供用于解码HyperHeaderBlock对象的实用方法。

// Decrypt 从io.Reader读取数据并将其反序列化为HyperHeaderBlock。
// 该方法首先读取头部字节并提取版本、客户端ID和加密块大小。
// 然后读取加密字节并将其解密为DecryptedBlock。
// 它检查DecryptedBlock中的clientId和expirationDate是否与HyperHeaderBlock中的clientId和version匹配。
// 如果不匹配，则抛出异常。
// 如果DecryptedBlock中的expirationDate小于当前日期，则抛出异常。
// 最后，它用读取和解密的数据创建一个新的HyperHeaderBlock并返回。
//
// 参数：
//
//	reader: 用于读取数据的io.Reader。
//	key: 用于解密的密钥。
//
// 返回：
//
//	解密后的HyperHeaderBlock和可能的错误。
func Decrypt(reader io.Reader, key string) (*entity.HyperHeaderBlock, error) {
	headerBytes := make([]byte, entity.HeaderSize)
	_, err := io.ReadFull(reader, headerBytes)
	if err != nil {
		return nil, fmt.Errorf("读取头部失败: %w", err)
	}

	version := byteUtil.GetIntLong(headerBytes, 0)
	clientID := int(byteUtil.GetIntLong(headerBytes, 4))
	encryptedBlockSize := int(byteUtil.GetIntLong(headerBytes, 8))

	encryptedBytes := make([]byte, encryptedBlockSize)
	_, err = io.ReadFull(reader, encryptedBytes)
	if err != nil {
		return nil, fmt.Errorf("读取加密块失败: %w", err)
	}

	decryptedBlock, err := entity.Decrypt(key, encryptedBytes)
	if err != nil {
		return nil, fmt.Errorf("解密块失败: %w", err)
	}

	// 检查DecryptedBlock中的clientId是否与HyperHeaderBlock中的clientId匹配
	if decryptedBlock.GetClientId() != clientID {
		return nil, fmt.Errorf("客户端ID错误")
	}

	// 检查DecryptedBlock中的expirationDate是否小于当前日期
	currentDate, _ := strconv.Atoi(time.Now().Format("060102"))
	if decryptedBlock.GetExpirationDate() < currentDate {
		return nil, fmt.Errorf("数据库已过期")
	}

	hyperHeaderBlock := &entity.HyperHeaderBlock{
		Version:            version,
		ClientId:           int64(clientID),
		EncryptedBlockSize: encryptedBlockSize,
		DecryptedBlock:     decryptedBlock,
	}

	return hyperHeaderBlock, nil
}
