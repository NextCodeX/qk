package core

import "fmt"

type ByteArrayValue struct {
	goValue []byte
	ValueAdapter
}

func newByteArrayValue(raw []byte) Value {
	return &ByteArrayValue{goValue: raw}
}

func (bs *ByteArrayValue) val() interface{} {
	return bs.goValue
}
func (bs *ByteArrayValue) isByteArray() bool {
	return true
}

func (bs *ByteArrayValue) String() string {
	return fmt.Sprint(bs.goValue)
}
