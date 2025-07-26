package entity

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/vmihailenco/msgpack/v5"
)

// DataBlock 表示数据库中的一个数据块。
// 它包含一个区域和一个指向数据库文件中数据的指针。
// DataBlock 对象的内存结构如下：
// +----------------+-----------+
// | []byte         | int       |
// +----------------+-----------+
// | region         | dataPtr   |
// +----------------+-----------+
type DataBlock struct {
	// region 是一个字节切片，表示该数据块覆盖的地理区域。
	region []byte

	// dataPtr 是一个整数，表示数据在数据库文件开始处的偏移量。
	dataPtr int
}

// NewDataBlock 创建一个新的 DataBlock，使用给定的 region 和 dataPtr。
func NewDataBlock(region []byte, dataPtr int) *DataBlock {
	return &DataBlock{
		region:  region,
		dataPtr: dataPtr,
	}
}

// GetRegion 返回此数据块的区域。
func (db *DataBlock) GetRegion(geoMapData []byte, columnSelection int64) string {
	result, err := db.unpack(geoMapData, columnSelection)
	if err != nil {
		return "null"
	}
	return result
}

// SetRegion 将此数据块的区域设置为指定值。
func (db *DataBlock) SetRegion(region []byte) *DataBlock {
	db.region = region
	return db
}

// GetDataPtr 返回此数据块的数据指针。
func (db *DataBlock) GetDataPtr() int {
	return db.dataPtr
}

// SetDataPtr 将此数据块的数据指针设置为指定值。
func (db *DataBlock) SetDataPtr(dataPtr int) *DataBlock {
	db.dataPtr = dataPtr
	return db
}

func (db *DataBlock) unpack(geoMapData []byte, columnSelection int64) (string, error) {
	var geoPosMixSize int64
	var otherData string

	decoder := msgpack.NewDecoder(bytes.NewReader(db.region))
	if err := decoder.Decode(&geoPosMixSize); err != nil {
		return "", err
	}
	if err := decoder.Decode(&otherData); err != nil {
		return "", err
	}

	if geoPosMixSize == 0 {
		return otherData, nil
	}

	dataLen := int((geoPosMixSize >> 24) & 0xFF)
	dataPtr := int(geoPosMixSize & 0x00FFFFFF)

	regionData := geoMapData[dataPtr : dataPtr+dataLen]
	var sb strings.Builder

	decoder = msgpack.NewDecoder(bytes.NewReader(regionData))
	columnNumber, err := decoder.DecodeArrayLen()
	if err != nil {
		return "", err
	}

	for i := 0; i < columnNumber; i++ {
		columnSelected := (columnSelection>>(i+1))&1 == 1
		var value string
		if err = decoder.Decode(&value); err != nil {
			return "", err
		}
		if value == "" {
			value = "null"
		}

		if columnSelected {
			sb.WriteString(value)
			sb.WriteString("\t")
		}
	}

	return fmt.Sprintf("%s%s", sb.String(), otherData), nil
}
