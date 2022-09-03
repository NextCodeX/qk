package core

import (
	"bytes"
	"fmt"
	"strconv"
)

// 新建ByteArray
func (this *InternalFunctionSet) NewBytes() Value {
	return newByteArrayValue(nil)
}
func (this *InternalFunctionSet) Newbs() Value {
	return newByteArrayValue(nil)
}

type ByteArrayValue struct {
	goValue []byte
	ClassObject
}

func emptyByteArray() Value {
	return newByteArrayValue(nil)
}

func newByteArrayValue(raw []byte) Value {
	bs := &ByteArrayValue{goValue: raw}
	bs.ClassObject.initAsClass("ByteArray", &bs)
	return bs
}

func (bs *ByteArrayValue) val() interface{} {
	return bs.goValue
}
func (bs *ByteArrayValue) isByteArray() bool {
	return true
}
func (bs *ByteArrayValue) sub(start, end int) []byte {
	return bs.goValue[start:end]
}

// 获取指定位置的单个字节字符，以字符串形式返回
func (bs *ByteArrayValue) At(index int) string {
	return string(bs.goValue[index])
}

// 字节数组与ByteArray或String比较
func (bs *ByteArrayValue) Equal(arg interface{}) bool {
	var data []byte
	if subBytes, ok := arg.([]byte); ok {
		data = subBytes
	} else if subStr, ok := arg.(string); ok {
		data = []byte(subStr)
	} else {
		return false
	}
	if len(bs.goValue) != len(data) {
		return false
	}
	if len(bs.goValue) == 0 && len(data) == 0 {
		return true
	}
	for i, b := range bs.goValue {
		if data[i] != b {
			return false
		}
	}
	return true
}
func (bs *ByteArrayValue) Eq(arg interface{}) bool {
	return bs.Equal(arg)
}

// 转字符串
func (bs *ByteArrayValue) Str() string {
	return string(bs.goValue)
}

// 转Integer
func (bs *ByteArrayValue) Int() int64 {
	return bytesToInt(bs.goValue)
}
func (bs *ByteArrayValue) I() int64 {
	res, err := strconv.Atoi(string(bs.goValue))
	if err != nil {
		fmt.Println(err)
		return 0
	}
	return int64(res)
}

// 转Float
func (bs *ByteArrayValue) Float() float64 {
	return bytesToFloat(bs.goValue)
}
func (bs *ByteArrayValue) F() float64 {
	res, err := strconv.ParseFloat(string(bs.goValue), 64)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	return res
}

// 是否包含指定ByteArray或String
func (bs *ByteArrayValue) Contain(arg interface{}) bool {
	if subBytes, ok := arg.([]byte); ok {
		return bytes.Contains(bs.goValue, subBytes)
	} else if subStr, ok := arg.(string); ok {
		return bytes.Contains(bs.goValue, []byte(subStr))
	} else {
		return false
	}
}

// 指定ByteArray或String在当前字节数组的位置
func (bs *ByteArrayValue) Index(arg interface{}) int {
	if subBytes, ok := arg.([]byte); ok {
		return bytes.Index(bs.goValue, subBytes)
	} else if subStr, ok := arg.(string); ok {
		return bytes.Index(bs.goValue, []byte(subStr))
	} else {
		return -1
	}
}

// 根据ByteArray, String, Integer或Float对字节数组进行切割
func (bs *ByteArrayValue) Split(arg interface{}) [][]byte {
	var sep []byte
	if bs, ok := arg.([]byte); ok {
		sep = bs
	} else if str, ok := arg.(string); ok {
		sep = []byte(str)
	} else if ival, ok := arg.(int64); ok {
		sep = intToBytes(ival)
	} else if fval, ok := arg.(float64); ok {
		sep = floatToBytes(fval)
	} else {
		fmt.Println("ByteArray.Split() parameter type must be one of ByteArray/String/Number")
	}
	return bytes.Split(bs.goValue, sep)
}

// 往字节数组里添加ByteArray, String, Integer或Float
func (bs *ByteArrayValue) Add(arg interface{}) {
	if subBytes, ok := arg.([]byte); ok {
		bs.goValue = append(bs.goValue, subBytes...)
	} else if subStr, ok := arg.(string); ok {
		bs.goValue = append(bs.goValue, []byte(subStr)...)
	} else if ival, ok := arg.(int64); ok {
		bs.goValue = append(bs.goValue, intToBytes(ival)...)
	} else if fval, ok := arg.(float64); ok {
		bs.goValue = append(bs.goValue, floatToBytes(fval)...)
	} else {
		fmt.Println("ByteArray.add() parameter type must be one of ByteArray/String/Number")
	}
}

// 字节数组大小
func (bs *ByteArrayValue) Size() int {
	return len(bs.goValue)
}

// 以二进制形式打印字节数组
func (bs *ByteArrayValue) Show() {
	showBytes(bs.goValue)
}

func (bs *ByteArrayValue) String() string {
	return fmt.Sprint(bs.goValue)
}
