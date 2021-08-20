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

type QKFile struct {
	raw *os.File
}

// 从特定位置，向文件写入字节数组
func (f *QKFile) WriteAt(raw interface{}, off int64) {
	var err error
	if data, ok := raw.([]byte); ok {
		_, err = f.raw.WriteAt(data, off)
	} else if data, ok := raw.(string); ok {
		_, err = f.raw.WriteAt([]byte(data), off)
	} else {
		runtimeExcption("data type must be ByteArray or String")
	}

	if err != nil {
		fmt.Println(err)
	}
	f.raw.Sync()
}

// 向文件写入字节数组/字符串
func (f *QKFile) Write(raw interface{}) {
	var err error
	if data, ok := raw.([]byte); ok {
		_, err = f.raw.Write(data)
	} else if data, ok := raw.(string); ok {
		_, err = f.raw.WriteString(data)
	} else {
		runtimeExcption("data type must be ByteArray or String")
	}

	if err != nil {
		fmt.Println(err)
	}
	f.raw.Sync()
}

// 读取文件
func (f *QKFile) Read() []byte {
	data := make([]byte, f.Size())
	_, err := f.raw.Read(data)
	if err != nil {
		fmt.Println(err)
	}
	return data
}

// 从指定位置读数据
func (f *QKFile) ReadAt(off int64, length int) []byte {
	data := make([]byte, length)
	_, err := f.raw.ReadAt(data, off)
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