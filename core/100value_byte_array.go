package core

import (
	"bytes"
	"encoding/base64"
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

func (this *ByteArrayValue) val() interface{} {
	return this.goValue
}
func (this *ByteArrayValue) isByteArray() bool {
	return true
}
func (this *ByteArrayValue) sub(start, end int) []byte {
	return this.goValue[start:end]
}

// 获取指定位置的单个字节字符，以字符串形式返回
func (this *ByteArrayValue) At(index int) string {
	return string(this.goValue[index])
}

// 字节数组与ByteArray或String比较
func (this *ByteArrayValue) Equal(arg interface{}) bool {
	var data []byte
	if subBytes, ok := arg.([]byte); ok {
		data = subBytes
	} else if subStr, ok := arg.(string); ok {
		data = []byte(subStr)
	} else {
		return false
	}
	if len(this.goValue) != len(data) {
		return false
	}
	if len(this.goValue) == 0 && len(data) == 0 {
		return true
	}
	for i, b := range this.goValue {
		if data[i] != b {
			return false
		}
	}
	return true
}
func (this *ByteArrayValue) Eq(arg interface{}) bool {
	return this.Equal(arg)
}

// 转字符串
func (this *ByteArrayValue) Str() string {
	return string(this.goValue)
}

// 转Integer
func (this *ByteArrayValue) Int() int64 {
	return bytesToInt(this.goValue)
}
func (this *ByteArrayValue) I() int64 {
	res, err := strconv.Atoi(string(this.goValue))
	if err != nil {
		fmt.Println(err)
		return 0
	}
	return int64(res)
}

// 转Float
func (this *ByteArrayValue) Float() float64 {
	return bytesToFloat(this.goValue)
}
func (this *ByteArrayValue) F() float64 {
	res, err := strconv.ParseFloat(string(this.goValue), 64)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	return res
}

// 是否包含指定ByteArray或String
func (this *ByteArrayValue) Contain(arg interface{}) bool {
	if subBytes, ok := arg.([]byte); ok {
		return bytes.Contains(this.goValue, subBytes)
	} else if subStr, ok := arg.(string); ok {
		return bytes.Contains(this.goValue, []byte(subStr))
	} else {
		return false
	}
}

// 指定ByteArray或String在当前字节数组的位置
func (this *ByteArrayValue) Index(arg interface{}) int {
	if subBytes, ok := arg.([]byte); ok {
		return bytes.Index(this.goValue, subBytes)
	} else if subStr, ok := arg.(string); ok {
		return bytes.Index(this.goValue, []byte(subStr))
	} else {
		return -1
	}
}

// 根据ByteArray, String, Integer或Float对字节数组进行切割
func (this *ByteArrayValue) Split(arg interface{}) [][]byte {
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
	return bytes.Split(this.goValue, sep)
}

// 往字节数组里添加ByteArray, String, Integer或Float
func (this *ByteArrayValue) Add(arg interface{}) {
	if subBytes, ok := arg.([]byte); ok {
		this.goValue = append(this.goValue, subBytes...)
	} else if subStr, ok := arg.(string); ok {
		this.goValue = append(this.goValue, []byte(subStr)...)
	} else if ival, ok := arg.(int64); ok {
		this.goValue = append(this.goValue, intToBytes(ival)...)
	} else if fval, ok := arg.(float64); ok {
		this.goValue = append(this.goValue, floatToBytes(fval)...)
	} else {
		fmt.Println("ByteArray.add() parameter type must be one of ByteArray/String/Number")
	}
}

func (this *ByteArrayValue) Save(path string) {
	fileSave(path, this.goValue)
}
func (this *ByteArrayValue) Base64() string {
	return base64.StdEncoding.EncodeToString(this.goValue)
}
func (this *ByteArrayValue) Debase64() []byte {
	data, err := base64.StdEncoding.DecodeString(string(this.goValue))
	if err != nil {
		return nil
	}
	return data
}
func (this *ByteArrayValue) Gzip() []byte {
	return gzipEncode(this.goValue)
}
func (this *ByteArrayValue) DeGzip() []byte {
	return gzipDecode(this.goValue)
}

// 字节数组大小
func (this *ByteArrayValue) Size() int {
	return len(this.goValue)
}

// 以二进制形式打印字节数组
func (this *ByteArrayValue) Show() {
	showBytes(this.goValue)
}

func (this *ByteArrayValue) Pr() {
	fmt.Println(this.String())
}
func (this *ByteArrayValue) String() string {
	return fmt.Sprint(this.goValue)
}
