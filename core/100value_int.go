package core

import "fmt"

type IntValue struct {
	goValue int64
	ClassObject
}

func newIntValue(raw int64) Value {
	i := &IntValue{goValue: raw}
	i.ClassObject.initAsClass("Int", &i)
	return i
}

func (ival *IntValue) val() interface{} {
	return ival.goValue
}

func (ival *IntValue) isInt() bool {
	return true
}

func (ival *IntValue) Bytes() []byte {
	return intToBytes(ival.goValue)
}

func (ival *IntValue) String() string {
	return fmt.Sprint(ival.goValue)
}