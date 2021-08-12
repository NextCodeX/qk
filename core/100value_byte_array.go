package core

import "fmt"

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

func (bs *ByteArrayValue) Str() string {
	return string(bs.goValue)
}

func (bs *ByteArrayValue) String() string {
	return fmt.Sprint(bs.goValue)
}
