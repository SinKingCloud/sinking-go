package entity

import (
	"server/app/util/ip/czdb/byteUtil"
	"server/app/util/ip/czdb/constant"
)

// IndexBlock 表示数据库中的索引块。
// 索引块包含起始IP、结束IP、数据指针和数据长度。
// 起始IP和结束IP用于确定数据块覆盖的IP地址范围。
// 数据指针用于定位数据库中的数据块。
// 数据长度用于从数据库中读取数据块。
type IndexBlock struct {
	// startIp 是数据块覆盖范围的起始IP地址。
	// 它是一个长度为4(IPv4)或16(IPv6)的字节数组。
	startIp []byte

	// endIp 是数据块覆盖范围的结束IP地址。
	// 它是一个长度为4(IPv4)或16(IPv6)的字节数组。
	endIp []byte

	// dataPtr 是数据库中数据块的指针。
	// 它是一个整数，表示数据块相对于数据库开始的偏移量。
	dataPtr int

	// dataLen 是数据块的字节长度。
	// 它是一个整数，表示从数据指针开始需要读取的字节数。
	dataLen int

	dbType string
}

// NewIndexBlock 是IndexBlock结构体的构造函数。
// 它使用提供的值初始化起始IP、结束IP、数据指针和数据长度。
//
// 参数:
//
//	startIp: 数据块覆盖范围的起始IP地址。
//	endIp: 数据块覆盖范围的结束IP地址。
//	dataPtr: 数据库中数据块的指针。
//	dataLen: 数据块的字节长度。
//	dbType: 数据库类型（IPv4或IPv6）。
func NewIndexBlock(startIp, endIp []byte, dataPtr, dataLen int, dbType string) *IndexBlock {
	return &IndexBlock{
		startIp: startIp,
		endIp:   endIp,
		dataPtr: dataPtr,
		dataLen: dataLen,
		dbType:  dbType,
	}
}

// Getter和Setter方法

func (ib *IndexBlock) GetStartIp() []byte {
	return ib.startIp
}

func (ib *IndexBlock) SetStartIp(startIp []byte) *IndexBlock {
	ib.startIp = startIp
	return ib
}

func (ib *IndexBlock) GetEndIp() []byte {
	return ib.endIp
}

func (ib *IndexBlock) SetEndIp(endIp []byte) *IndexBlock {
	ib.endIp = endIp
	return ib
}

func (ib *IndexBlock) GetDataPtr() int {
	return ib.dataPtr
}

func (ib *IndexBlock) SetDataPtr(dataPtr int) *IndexBlock {
	ib.dataPtr = dataPtr
	return ib
}

func (ib *IndexBlock) GetDataLen() int {
	return ib.dataLen
}

func (ib *IndexBlock) SetDataLen(dataLen int) *IndexBlock {
	ib.dataLen = dataLen
	return ib
}

// GetIndexBlockLength 返回索引块的长度
func GetIndexBlockLength(dbType string) int {
	// 如果是IPv6，则起始IP和结束IP各占16字节
	// 如果是IPv4，则起始IP和结束IP各占4字节
	// 再加上4字节的数据指针和1字节的数据长度
	if dbType == constant.IPV4 {
		return 13
	}
	return 37
}

// GetBytes 返回表示索引块的字节数组。
// 字节数组的结构如下：
// +------------+-----------+-----------+
// | 4/16 bytes | 4/16bytes | 4bytes    | 1bytes
// +------------+-----------+-----------+
//
//	起始IP       结束IP      数据指针    长度
//
// 返回值: 表示索引块的字节数组。
func (ib *IndexBlock) GetBytes() []byte {
	ipBytesLength := 4
	if ib.dbType == constant.IPV6 {
		ipBytesLength = 16
	}

	b := make([]byte, GetIndexBlockLength(ib.dbType))
	copy(b[:ipBytesLength], ib.startIp)
	copy(b[ipBytesLength:ipBytesLength*2], ib.endIp)

	// 写入数据指针和长度
	byteUtil.WriteIntLong(b, ipBytesLength*2, int64(ib.dataPtr))
	byteUtil.Write(b, ipBytesLength*2+4, int64(ib.dataLen), 1)
	return b
}
