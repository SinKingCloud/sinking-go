package entity

import "server/app/util/ip/czdb/byteUtil"

// HyperHeaderBlock 表示HyperHeader结构的头块。
// HyperHeaderBlock类封装了以下信息：
// 1. version: 一个4字节长整型，表示HyperHeaderBlock的版本。
// 2. clientId: 一个4字节长整型，表示客户端ID。
// 3. encryptedBlockSize: 一个4字节整型，表示加密数据块的大小。
// 4. encryptedData: 一个大小为encryptedBlockSize的字节数组，表示加密数据。
// 5. 随机字节: 一些随机字节，其大小保存在encryptedData中。
//
// HyperHeaderBlock的内存结构如下：
// |------------------|------------------|----------------------|---------------------------|-------------------|
// | 版本 (4字节)      | 客户端ID (4字节)   | 加密块大小 (4字节)     | 加密数据 (可变长度)         | 随机字节           |
// |------------------|------------------|----------------------|---------------------------|-------------------|

const HeaderSize = 12

type HyperHeaderBlock struct {
	Version            int64
	ClientId           int64
	EncryptedBlockSize int
	EncryptedData      []byte
	DecryptedBlock     *DecryptedBlock
}

// GetVersion 获取HyperHeaderBlock的版本。
func (h *HyperHeaderBlock) GetVersion() int64 {
	return h.Version
}

// SetVersion 设置HyperHeaderBlock的版本。
func (h *HyperHeaderBlock) SetVersion(version int64) {
	h.Version = version
}

// GetClientId 获取HyperHeaderBlock的客户端ID。
func (h *HyperHeaderBlock) GetClientId() int64 {
	return h.ClientId
}

// SetClientId 设置HyperHeaderBlock的客户端ID。
func (h *HyperHeaderBlock) SetClientId(clientId int64) {
	h.ClientId = clientId
}

// GetEncryptedBlockSize 获取加密数据块的大小。
func (h *HyperHeaderBlock) GetEncryptedBlockSize() int {
	return h.EncryptedBlockSize
}

// SetEncryptedBlockSize 设置加密数据块的大小。
func (h *HyperHeaderBlock) SetEncryptedBlockSize(encryptedBlockSize int) {
	h.EncryptedBlockSize = encryptedBlockSize
}

// GetEncryptedData 获取加密数据。
func (h *HyperHeaderBlock) GetEncryptedData() []byte {
	return h.EncryptedData
}

// SetEncryptedData 设置加密数据。
func (h *HyperHeaderBlock) SetEncryptedData(encryptedData []byte) {
	h.EncryptedData = encryptedData
}

// GetDecryptedBlock 获取解密后的块。
func (h *HyperHeaderBlock) GetDecryptedBlock() *DecryptedBlock {
	return h.DecryptedBlock
}

// SetDecryptedBlock 设置解密后的块。
func (h *HyperHeaderBlock) SetDecryptedBlock(decryptedBlock *DecryptedBlock) {
	h.DecryptedBlock = decryptedBlock
}

// ToBytes 将HyperHeaderBlock实例转换为字节数组。
// 此方法将HyperHeaderBlock实例序列化为字节数组，可用于存储或传输。
// 字节数组的结构如下：
// - 前4个字节表示HyperHeaderBlock的版本。
// - 接下来的4个字节表示客户端ID。
// - 再接下来的4个字节表示加密数据的长度。
func (h *HyperHeaderBlock) ToBytes() []byte {
	bytes := make([]byte, 12)
	byteUtil.WriteIntLong(bytes, 0, h.Version)
	byteUtil.WriteIntLong(bytes, 4, h.ClientId)
	byteUtil.WriteIntLong(bytes, 8, int64(h.EncryptedBlockSize))
	return bytes
}

// FromBytes 从12字节长的字节数组反序列化HyperHeaderBlock实例。
// 此方法接受一个字节数组并从中构造一个HyperHeaderBlock实例。
// 字节数组的结构预期如下：
// - 前4个字节表示HyperHeaderBlock的版本。
// - 接下来的4个字节表示客户端ID。
// - 再接下来的4个字节表示加密数据的长度。
func FromBytes(bytes []byte) *HyperHeaderBlock {
	version := byteUtil.GetIntLong(bytes, 0)
	clientId := byteUtil.GetIntLong(bytes, 4)
	encryptedBlockSize := byteUtil.GetIntLong(bytes, 8)

	headerBlock := &HyperHeaderBlock{
		Version:            version,
		ClientId:           clientId,
		EncryptedBlockSize: int(encryptedBlockSize),
	}

	return headerBlock
}

// GetHeaderSize 返回HyperHeaderBlock的总大小。
// 大小计算为以下各项之和：
// - 头部的大小（12字节）
// - 加密数据块的大小
// - 随机字节的大小
func (h *HyperHeaderBlock) GetHeaderSize() int {
	return HeaderSize + int(h.EncryptedBlockSize) + h.DecryptedBlock.GetRandomSize()
}
