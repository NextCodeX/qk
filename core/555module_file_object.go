package core

import (
	"fmt"
	"os"
)

// 打开文件，返回一个文件对象
func (fns *InternalFunctionSet) Fopen(dir string) Value {
	file, err := os.OpenFile(dir, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	obj := &QKFile{file}
	return newClass("File", &obj)
}

// 文件读写操作是基于文件指针所在位置进行
// 每次对文件执行读写操作后，文件指针位置会偏移
// 使用Seek()方法，可以使文件指针偏移到我们想要的位置，然后进行读写
type QKFile struct {
	raw *os.File
}

// 从特定位置，向文件写入字节数组/字符串/Int/Float(默认从文件开始位置写入)
func (f *QKFile) Write(args []interface{}) {
	if len(args) < 1 {
		runtimeExcption("write(data[, off]) parameter data is required")
		return
	}
	raw := args[0]

	var err error
	var data []byte
	if bs, ok := raw.([]byte); ok {
		data = bs

	} else if str, ok := raw.(string); ok {
		data = []byte(str)

	} else if ival, ok := raw.(int64); ok {
		data = intToBytes(ival)

	} else if fval, ok := raw.(float64); ok {
		data = floatToBytes(fval)

	} else {
		runtimeExcption("write(data, off) data type must be one of ByteArray/String/Int/Float")
	}
	var off int64
	if len(args) > 1 {
		off = args[1].(int64)
	}

	_, err = f.raw.WriteAt(data, off)
	if err != nil {
		fmt.Println(err)
	}
	f.raw.Sync()
}

// 默认从文件开始位置，读取整个文件的数据。
// 但也可以从从指定位置读取指定大小的数据
func (f *QKFile) Read(args []interface{}) []byte {
	var off, length int64
	if len(args) > 0 {
		off = args[0].(int64)
	}
	if len(args) > 1 {
		length = args[1].(int64)
	} else {
		length = f.Size() - off
	}

	f.raw.Seek(off, 0)
	data := make([]byte, length)
	_, err := f.raw.Read(data)
	if err != nil {
		fmt.Println(err)
	}
	return data
}

// 获取文件大小
func (f *QKFile) Size() int64 {
	stat, err := f.raw.Stat()
	if err != nil {
		fmt.Println(err)
		return 0
	}
	return stat.Size()
}

// flushing the file system's in-memory copy of recently written data to disk.
func (f *QKFile) Flush() {
	err := f.raw.Sync()
	if err != nil {
		fmt.Println(err)
	}
}
