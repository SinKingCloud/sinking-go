package czdb

import (
	"fmt"
	"net"
	"os"
	"server/app/util/ip/czdb/byteUtil"
	"server/app/util/ip/czdb/constant"
	"server/app/util/ip/czdb/entity"
	"server/app/util/ip/czdb/hyperHeaderDecoder"
)

// DbSearcher 提供在数据库中搜索数据的方法。
// 它支持三种类型的搜索算法:内存、二进制和B树。
// 数据库类型(IPv4或IPv6)和查询类型(MEMORY、BINARY、BTREE)在运行时确定。
// 该类还提供根据查询类型初始化搜索参数的方法,以及通过IP地址获取地区的方法。
// DbSearcher类使用RandomAccessFile从数据库文件读取和写入。
// 对于B树搜索,它使用2D字节数组和整数数组来表示每个索引块的起始IP和数据指针。
// 对于内存和二进制搜索,它使用字节数组来表示数据库的原始二进制字符串。
// 该类还提供关闭数据库的方法。
type DbSearcher struct {
	dbType        string
	ipBytesLength int
	queryType     string

	totalHeaderBlockSize int64

	raf *RandomAccessFile

	// 仅用于B树搜索
	headerSip    [][]byte
	headerPtr    []int
	headerLength int

	// 用于内存和二进制搜索
	firstIndexPtr    int64
	totalIndexBlocks int

	// 仅用于内存搜索
	dbBinStr []byte

	columnSelection int64
	geoMapData      []byte
}

// NewDbSearcher 创建一个新的DbSearcher实例
func NewDbSearcher(dbFile string, queryType string, key string) (*DbSearcher, error) {
	ds := &DbSearcher{
		queryType: queryType,
	}
	file, err := os.Open(dbFile)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()
	headerBlock, err := hyperHeaderDecoder.Decrypt(file, key)
	if err != nil {
		return nil, err
	}
	ds.raf, err = NewRandomAccessFile(dbFile, int64(headerBlock.GetHeaderSize()))
	if err != nil {
		return nil, err
	}

	// 设置数据库类型
	_, err = ds.raf.Seek1(0)
	if err != nil {
		return nil, err
	}
	superBytes := make([]byte, constant.SUPER_PART_LENGTH)
	err = ds.raf.ReadFully(superBytes)
	if err != nil {
		return nil, err
	}

	if (superBytes[0] & 1) == 0 {
		ds.dbType = constant.IPV4
		ds.ipBytesLength = 4
	} else {
		ds.dbType = constant.IPV6
		ds.ipBytesLength = 16
	}

	// 加载地理设置
	ds.loadGeoSetting(ds.raf, key)

	if queryType == "MEMORY" {
		err = ds.initializeForMemorySearch()
	} else if queryType == "BTREE" {
		err = ds.initBtreeModeParam(ds.raf)
	}

	if err != nil {
		return nil, err
	}

	return ds, nil
}

// 加载地理位置映射表。
func (ds *DbSearcher) loadGeoSetting(raf *RandomAccessFile, key string) {
	//设置位置到结束索引ptr + ip字节长度 * 2 + 4
	_, err := raf.Seek1(constant.END_INDEX_PTR)
	if err != nil {
		return
	}
	data := make([]byte, 4)
	err = raf.ReadFully(data)
	if err != nil {
		return
	}

	endIndexPtr := byteUtil.GetIntLong(data, 0)
	columnSelectionPtr := endIndexPtr + int64(entity.GetIndexBlockLength(ds.dbType))
	_, err = raf.Seek1(columnSelectionPtr)
	if err != nil {
		return
	}
	err = raf.ReadFully(data)
	if err != nil {
		return
	}
	ds.columnSelection = byteUtil.GetIntLong(data, 0)
	if ds.columnSelection == 0 {
		return
	}
	geoMapPtr := columnSelectionPtr + 4
	_, err = raf.Seek1(geoMapPtr)
	if err != nil {
		return
	}
	err = raf.ReadFully(data)
	if err != nil {
		return
	}
	geoMapSize := int(byteUtil.GetIntLong(data, 0))

	_, err = raf.Seek1(geoMapPtr + 4)
	if err != nil {
		return
	}
	ds.geoMapData = make([]byte, geoMapSize)
	err = raf.ReadFully(ds.geoMapData)
	if err != nil {
		return
	}

	decrypt := NewDecrypted(key)
	ds.geoMapData = decrypt.decrypt(ds.geoMapData)
	return
}

/**
*初始化DbSearcher实例进行内存搜索。
*将整个数据库文件读取到内存中，然后初始化内存或二进制搜索的参数。
 */
func (ds *DbSearcher) initializeForMemorySearch() error {
	length, err := ds.raf.Length()
	if err != nil {
		return err
	}
	ds.dbBinStr = make([]byte, length)
	_, err = ds.raf.Seek1(0)
	if err != nil {
		return err
	}
	err = ds.raf.ReadFully(ds.dbBinStr)
	if err != nil {
		return err
	}
	err = ds.raf.Close()
	if err != nil {
		return err
	}
	return ds.initMemoryOrBinaryModeParam(ds.dbBinStr, length)
}

func (ds *DbSearcher) initMemoryOrBinaryModeParam(bytes []byte, fileSize int64) error {
	ds.totalHeaderBlockSize = byteUtil.GetIntLong(bytes, constant.HEADER_BLOCK_PTR)
	fileSizeInFile := byteUtil.GetIntLong(bytes, constant.FILE_SIZE_PTR)
	if fileSizeInFile != fileSize {
		return fmt.Errorf("db file size error, excepted [%d], real [%d]", fileSizeInFile, fileSize)
	}
	ds.firstIndexPtr = byteUtil.GetIntLong(bytes, constant.FIRST_INDEX_PTR)
	lastIndexPtr := byteUtil.GetIntLong(bytes, constant.END_INDEX_PTR)
	ds.totalIndexBlocks = (int(lastIndexPtr-ds.firstIndexPtr) / entity.GetIndexBlockLength(ds.dbType)) + 1

	b := make([]byte, ds.totalHeaderBlockSize)
	copy(b, bytes[constant.SUPER_PART_LENGTH:])
	return ds.initHeaderBlock(b)
}

func (ds *DbSearcher) initBtreeModeParam(raf *RandomAccessFile) error {
	_, err := raf.Seek1(0)
	if err != nil {
		return err
	}
	superBytes := make([]byte, constant.SUPER_PART_LENGTH)
	err = raf.ReadFully(superBytes)
	if err != nil {
		return err
	}
	ds.totalHeaderBlockSize = byteUtil.GetIntLong(superBytes, constant.HEADER_BLOCK_PTR)

	fileSizeInFile := byteUtil.GetIntLong(superBytes, constant.FILE_SIZE_PTR)
	realFileSize, err := ds.raf.Length()
	if err != nil {
		return err
	}
	if fileSizeInFile != realFileSize {
		return fmt.Errorf("db file size error, excepted [%d], real [%d]", fileSizeInFile, realFileSize)
	}

	b := make([]byte, ds.totalHeaderBlockSize)
	err = ds.raf.ReadFully(b)
	if err != nil {
		return err
	}

	return ds.initHeaderBlock(b)
}

func (ds *DbSearcher) initHeaderBlock(headerBytes []byte) error {
	indexLength := 20
	length := len(headerBytes) / indexLength
	idx := 0
	ds.headerSip = make([][]byte, length)
	for i := range ds.headerSip {
		ds.headerSip[i] = make([]byte, 16)
	}
	ds.headerPtr = make([]int, length)
	var dataPtr int64
	for i := 0; i < len(headerBytes); i += indexLength {
		dataPtr = byteUtil.GetIntLong(headerBytes, i+16)
		if dataPtr == 0 {
			break
		}
		copy(ds.headerSip[idx][:16], headerBytes[i:i+16])
		ds.headerPtr[idx] = int(dataPtr)
		idx++
	}
	ds.headerLength = idx
	return nil
}

/*Search
 * 此方法用于根据提供的IP地址在数据库中搜索区域。
 * 它支持三种类型的搜索算法：内存搜索、二分查找和B树搜索。
 * 搜索算法的类型由DbSearcher实例的queryType属性决定。
 * 该方法首先将IP地址转换为字节数组，然后根据查询类型执行搜索。
 * 如果搜索成功，它将返回找到的数据块的区域信息。
 * 如果搜索不成功，则返回null。
 *
 * @param ip 要搜索的IP地址。它是一个标准IP地址格式的字符串。
 * @return 如果搜索成功，返回找到的数据块的区域信息；否则返回null。
 * @throws IpFormatException 如果提供的IP地址格式不正确，则抛出此异常。
 * @throws IOException 如果在搜索过程中发生I/O错误，则抛出此异常。
 */
func (ds *DbSearcher) Search(ip string) (string, error) {
	ipBytes, err := ds.getIpBytes(ip)
	if err != nil {
		return "", err
	}

	var dataBlock *entity.DataBlock

	switch ds.queryType {
	case "MEMORY":
		dataBlock = ds.memorySearch(ipBytes)
	case "BTREE":
		dataBlock, err = ds.bTreeSearch(ipBytes)
		if err != nil {
			return "", err
		}
	default:
		return "", fmt.Errorf("unsupported query type")
	}

	if dataBlock == nil {
		return "", nil
	}
	return dataBlock.GetRegion(ds.geoMapData, ds.columnSelection), nil
}

/**
 * 此方法执行内存搜索，根据提供的IP地址在数据库中查找数据块。
 * 它使用二分查找算法来搜索索引块并找到数据。
 * 如果搜索成功，它将返回包含区域信息和数据指针的数据块。
 * 如果搜索不成功，则返回null。
 *
 * @param ip 要搜索的IP地址。它是一个表示IP地址的字节数组。
 * @return 如果搜索成功，返回包含区域信息和数据指针的数据块；否则返回null。
 */
func (ds *DbSearcher) memorySearch(ip []byte) *entity.DataBlock {
	blockLen := entity.GetIndexBlockLength(ds.dbType)

	spurNeptune := ds.searchInHeader(ip)
	spur, entry := spurNeptune[0], spurNeptune[1]

	if spur == 0 {
		return nil
	}

	l, h := 0, (entry-spur)/blockLen

	sip := make([]byte, ds.ipBytesLength)
	eip := make([]byte, ds.ipBytesLength)

	var dataPtr int
	var dataLen int
	for l <= h {
		m := (l + h) >> 1
		p := spur + m*blockLen

		copy(sip, ds.dbBinStr[p:p+ds.ipBytesLength])
		copy(eip, ds.dbBinStr[p+ds.ipBytesLength:p+2*ds.ipBytesLength])

		cmpStart := compareBytes(ip, sip, ds.ipBytesLength)
		cmpEnd := compareBytes(ip, eip, ds.ipBytesLength)

		if cmpStart >= 0 && cmpEnd <= 0 {
			dataPtr = int(byteUtil.GetIntLong(ds.dbBinStr, p+ds.ipBytesLength*2))
			dataLen = byteUtil.GetInt1(ds.dbBinStr, p+ds.ipBytesLength*2+4)
			break
		} else if cmpStart < 0 {
			h = m - 1
		} else {
			l = m + 1
		}
	}

	if dataPtr == 0 {
		return nil
	}

	region := make([]byte, dataLen)
	copy(region, ds.dbBinStr[dataPtr:dataPtr+dataLen])

	return entity.NewDataBlock(region, dataPtr)
}

// searchInHeader 在头部搜索
func (ds *DbSearcher) searchInHeader(ip []byte) [2]int {
	l, h := 0, ds.headerLength-1
	var spur, entry int

	for l <= h {
		m := (l + h) >> 1
		cmp := compareBytes(ip, ds.headerSip[m], ds.ipBytesLength)

		if cmp < 0 {
			h = m - 1
		} else if cmp > 0 {
			l = m + 1
		} else {
			if m > 0 {
				spur = ds.headerPtr[m-1]
			} else {
				spur = ds.headerPtr[m]
			}
			entry = ds.headerPtr[m]
			break
		}
	}

	// less than header range
	if l == 0 && h <= 0 {
		return [2]int{0, 0}
	}

	if l > h {
		if l < ds.headerLength {
			spur = ds.headerPtr[l-1]
			entry = ds.headerPtr[l]
		} else if h >= 0 && h+1 < ds.headerLength {
			spur = ds.headerPtr[h]
			entry = ds.headerPtr[h+1]
		} else {
			spur = ds.headerPtr[ds.headerLength-1]
			entry = spur + entity.GetIndexBlockLength(ds.dbType)
		}
	}

	return [2]int{spur, entry}
}

// bTreeSearch 执行B树搜索
func (ds *DbSearcher) bTreeSearch(ip []byte) (*entity.DataBlock, error) {
	spurNeptune := ds.searchInHeader(ip)
	spar, entry := spurNeptune[0], spurNeptune[1]

	if spar == 0 {
		return nil, nil
	}

	blockLen := entry - spar
	blen := entity.GetIndexBlockLength(ds.dbType)

	iBuffer := make([]byte, blockLen+blen)
	_, err := ds.raf.Seek1(int64(spar))
	if err != nil {
		return nil, err
	}
	err = ds.raf.ReadFully(iBuffer)
	if err != nil {
		return nil, err
	}

	l, h := 0, blockLen/blen

	sip := make([]byte, ds.ipBytesLength)
	eip := make([]byte, ds.ipBytesLength)

	dataPtr := 0
	dataLen := 0

	for l <= h {
		m := (l + h) >> 1
		p := m * blen

		copy(sip, iBuffer[p:p+ds.ipBytesLength])
		copy(eip, iBuffer[p+ds.ipBytesLength:p+2*ds.ipBytesLength])

		cmpStart := compareBytes(ip, sip, ds.ipBytesLength)
		cmpEnd := compareBytes(ip, eip, ds.ipBytesLength)

		if cmpStart >= 0 && cmpEnd <= 0 {
			dataPtr = int(byteUtil.GetIntLong(iBuffer, p+ds.ipBytesLength*2))
			dataLen = byteUtil.GetInt1(iBuffer, p+ds.ipBytesLength*2+4)
			break
		} else if cmpStart < 0 {
			h = m - 1
		} else {
			l = m + 1
		}
	}
	if dataPtr == 0 {
		return nil, nil
	}

	_, err = ds.raf.Seek1(int64(dataPtr))
	if err != nil {
		return nil, err
	}
	region := make([]byte, dataLen)
	err = ds.raf.ReadFully(region)
	if err != nil {
		return nil, err
	}

	return entity.NewDataBlock(region, dataPtr), nil
}

func (ds *DbSearcher) getByIndexPtr(ptr int64) (*entity.DataBlock, error) {
	_, err := ds.raf.Seek1(ptr)
	if err != nil {
		return nil, err
	}
	buffer := make([]byte, 36)
	err = ds.raf.ReadFullyAt(buffer, 0, 36)
	if err != nil {
		return nil, err
	}
	extra := byteUtil.GetIntLong(buffer, 32)

	dataLen := int((extra >> 24) & 0xFF)
	dataPtr := int(extra & 0x00FFFFFF)

	_, err = ds.raf.Seek1(int64(dataPtr))
	if err != nil {
		return nil, err
	}
	region := make([]byte, dataLen)
	err = ds.raf.ReadFully(region)
	if err != nil {
		return nil, err
	}
	return entity.NewDataBlock(region, dataPtr), nil
}

func (ds *DbSearcher) getDbType() string {
	return ds.dbType
}
func (ds *DbSearcher) getQueryType() string {
	return ds.queryType
}
func (ds *DbSearcher) Close() {
	ds.headerSip = nil
	ds.headerPtr = nil
	ds.dbBinStr = nil
	if ds.raf != nil {
		_ = ds.raf.Close()
	}
}

// getIpBytes 将IP地址字符串转换为字节切片
func (ds *DbSearcher) getIpBytes(ip string) ([]byte, error) {
	ipBytes := net.ParseIP(ip)
	if ipBytes == nil {
		return nil, fmt.Errorf("invalid IP format")
	}
	if ds.dbType == constant.IPV4 {
		return ipBytes.To4(), nil
	}
	return ipBytes.To16(), nil
}

/**
 * 此方法比较两个字节数组，最多比较指定的长度。
 * 它用于比较以字节数组格式表示的IP地址。
 * 比较是逐字节进行的，一旦发现差异，方法立即返回结果。
 * 如果两个数组当前位置的字节都是正数或都是负数，方法会比较它们的值。
 * 如果两个数组当前位置的字节符号不同，方法会将负字节视为较大。
 * 如果当前位置的一个字节是零而另一个不是，方法会将零字节视为较小。
 * 如果方法比较到指定长度都没有发现差异，它会比较字节数组的长度。
 * 如果长度相等，则认为字节数组相等。
 * 如果一个字节数组比另一个长，则认为它更大。
 *
 * @param bytes1 要比较的第一个字节数组。它代表一个IP地址。
 * @param bytes2 要比较的第二个字节数组。它代表一个IP地址。
 * @param length 要在每个字节数组中比较的字节数。
 * @return 如果第一个字节数组小于第二个，返回负整数；如果它们相等，返回零；如果第一个字节数组大于第二个，返回正整数。
 */
func compareBytes(bytes1, bytes2 []byte, length int) int {
	for i := 0; i < len(bytes1) && i < len(bytes2) && i < length; i++ {
		b1 := int(bytes1[i])
		b2 := int(bytes2[i])
		if b1*b2 > 0 {
			if b1 < b2 {
				return -1
			} else if b1 > b2 {
				return 1
			}
		} else if b1*b2 < 0 {
			// When the signs are different, the negative byte is considered larger
			if b1 > 0 {
				return -1
			} else {
				return 1
			}
		} else if b1*b2 == 0 && b1+b2 != 0 {
			// When one byte is zero and the other is not, the zero byte is considered smaller
			if b1 == 0 {
				return -1
			} else {
				return 1
			}
		}
	}
	if len(bytes1) >= length && len(bytes2) >= length {
		return 0
	} else {
		return len(bytes1) - len(bytes2)
	}
}
