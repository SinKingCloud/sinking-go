package czdb

import (
	"fmt"
	"io"
	"os"
)

// RandomAccessFile 结构体扩展了标准库的 os.File，
// 并为文件指针添加了一个偏移量。
// 每当调用 Seek1 方法时，都会添加这个偏移量。
// 这在您想将文件的一部分视为单独文件时很有用。
type RandomAccessFile struct {
	*os.File
	//要添加到文件指针的偏移量
	offset int64
}

// NewRandomAccessFile 创建一个新的 RandomAccessFile 实例。
// 每当调用 Seek1 方法时，偏移量会被添加到文件指针上。
//
// 参数:
//
//	name: 系统依赖的文件名
//	flag: 打开文件的模式（例如 os.O_RDONLY, os.O_RDWR 等）
//	perm: 文件权限（如果创建新文件）
//	offset: 要添加到文件指针的偏移量
//
// 返回值:
//
//	*RandomAccessFile: 新创建的 RandomAccessFile 实例
//	error: 如果在打开或创建文件时发生错误
func NewRandomAccessFile(name string, offset int64) (*RandomAccessFile, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	return &RandomAccessFile{File: file, offset: offset}, nil
}

// Seek1 设置下一次读取或写入操作的文件指针偏移量，从文件开头开始计算。
// 偏移量可以设置在文件末尾之后。设置超过文件末尾的偏移量不会改变文件长度。
// 文件长度只会通过在设置偏移量超过文件末尾后进行写入操作而改变。
//
// 参数:
//
//	offset: 要设置的偏移量，以字节为单位，从文件开头开始计算
//	whence: 偏移量的参考点（io.SeekStart, io.SeekCurrent, 或 io.SeekEnd）
//
// 返回值:
//
//	int64: 新的文件偏移量
//	error: 如果 offset 小于 0 或发生 I/O 错误
func (f *RandomAccessFile) Seek1(offset int64) (int64, error) {
	return f.File.Seek(offset+f.offset, io.SeekStart)
}

// Length 返回文件的大小，考虑了偏移量。
//
// 返回值:
//
//	int64: 文件大小（字节）
//	error: 如果发生 I/O 错误
func (f *RandomAccessFile) Length() (int64, error) {
	info, err := f.File.Stat()
	if err != nil {
		return 0, err
	}
	return info.Size() - f.offset, nil
}

func (f *RandomAccessFile) ReadFully(p []byte) error {
	_, err := io.ReadFull(f.File, p)
	return err
}

func (f *RandomAccessFile) ReadFullyAt(p []byte, off, length int) error {
	if off < 0 || length < 0 || off > len(p)-length {
		return fmt.Errorf("invalid offset or length")
	}
	n, err := io.ReadAtLeast(f.File, p[off:], length)
	if err == io.EOF || err == nil {
		if n == length {
			return nil // 刚好读取完指定长度
		} else {
			return io.ErrUnexpectedEOF
		}
	}
	return err
}
