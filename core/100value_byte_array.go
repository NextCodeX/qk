package core

import (
	"bytes"
	"fmt"
)

func (fns *InternalFunctionSet) NewBytes() Value {
	return newByteArrayValue(nil)
}

type ByteArrayValue struct {
	goValue []byte
	ClassObject
}

func newByteArrayValue(raw []byte) Value {
	bs := &ByteArrayValue{goValue: raw}
	bs.ClassObject.raw = &bs
	bs.ClassObject.name = "ByteArray"
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

func (bs *ByteArrayValue) Str() string {
	return string(bs.goValue)
}

func (bs *ByteArrayValue) Int() int64 {
	return bytesToInt(bs.goValue)
}

func (bs *ByteArrayValue) Float() float64 {
	return bytesToFloat(bs.goValue)
}

func (bs *ByteArrayValue) Contain(arg interface{}) bool {
	if subBytes, ok := arg.([]byte); ok {
		return bytes.Contains(bs.goValue, subBytes)
	} else if subStr, ok := arg.(string); ok {
		return bytes.Contains(bs.goValue, []byte(subStr))
	} else {
		return false
	}
}

func (bs *ByteArrayValue) Index(arg interface{}) int {
	if subBytes, ok := arg.([]byte); ok {
		return bytes.Index(bs.goValue, subBytes)
	} else if subStr, ok := arg.(string); ok {
		return bytes.Index(bs.goValue, []byte(subStr))
	} else {
		return -1
	}
}

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

func (bs *ByteArrayValue) Size() int {
	return len(bs.goValue)
}

func (bs *ByteArrayValue) Show() {
	showBytes(bs.goValue)
}

func (bs *ByteArrayValue) String() string {
	return fmt.Sprint(bs.goValue)
}
