package core

import (
	"fmt"
)

type AnyValue struct {
	goValue interface{}
	ClassObject
}

func newAnyValue(raw interface{}) Value {
	any := &AnyValue{goValue: raw}
	any.initAsClass("Anything", &any)
	return any
}


func (any *AnyValue) val() interface{} {
	return any.goValue
}
func (any *AnyValue) isAny() bool {
	return true
}

func (any *AnyValue) String() string {
	return fmt.Sprint(any.goValue)
}

