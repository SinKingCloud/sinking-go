package entity

import (
	"math"
	"server/app/util/ip/czdb/byteUtil"
)

// HeaderBlock 表示数据库中的一个头块。
// 它包含一个索引起始IP地址和一个索引指针。
// 它还提供了一个方法来获取用于数据库存储的字节。
//
// HeaderBlock实例的内存布局:
//
// +-----------------+-----------------+
// | indexStartIp    | indexPtr        |
// | (字节数组)       | (int)           |
// +-----------------+-----------------+
type HeaderBlock struct {
	// indexStartIp 是索引起始IP地址。
	indexStartIp []byte

	// indexPtr 是索引指针。
	indexPtr int
}

// HEADER_LINE_SIZE 定义了头行的大小。
const HEADER_LINE_SIZE = 20

// NewHeaderBlock 使用指定的索引起始IP地址和索引指针构造一个新的HeaderBlock。
func NewHeaderBlock(indexStartIp []byte, indexPtr int) *HeaderBlock {
	return &HeaderBlock{
		indexStartIp: indexStartIp,
		indexPtr:     indexPtr,
	}
}

// GetIndexStartIp 返回此头块的索引起始IP地址。
func (hb *HeaderBlock) GetIndexStartIp() []byte {
	return hb.indexStartIp
}

// SetIndexStartIp 将此头块的索引起始IP地址设置为指定值。
func (hb *HeaderBlock) SetIndexStartIp(indexStartIp []byte) *HeaderBlock {
	hb.indexStartIp = indexStartIp
	return hb
}

// GetIndexPtr 返回此头块的索引指针。
func (hb *HeaderBlock) GetIndexPtr() int {
	return hb.indexPtr
}

// SetIndexPtr 将此头块的索引指针设置为指定值。
func (hb *HeaderBlock) SetIndexPtr(indexPtr int) *HeaderBlock {
	hb.indexPtr = indexPtr
	return hb
}

// GetBytes 返回用于数据库存储的字节。
// 返回的字节数组长度为20字节，前16字节为索引起始IP地址，后4字节为索引指针。
func (hb *HeaderBlock) GetBytes() []byte {
	b := make([]byte, 20)
	copy(b[:int(math.Min(float64(len(hb.indexStartIp)), 16))], hb.indexStartIp)
	byteUtil.WriteIntLong(b, 16, int64(hb.indexPtr))
	return b
}
