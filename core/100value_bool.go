package core

import "fmt"

type BooleanValue struct {
	goValue bool
	ClassObject
}

func newBooleanValue(raw bool) Value {
	bl := &BooleanValue{goValue: raw}
	bl.ClassObject.raw = &bl
	bl.ClassObject.name = "Boolean"
	return bl
}

func (boolVal *BooleanValue) val() interface{} {
	return boolVal.goValue
}

func (boolVal *BooleanValue) isBoolean() bool {
	return true
}

func (boolVal *BooleanValue) String() string {
	return fmt.Sprint(boolVal.goValue)
}