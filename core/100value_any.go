package core

import (
	"fmt"
)

type AnyValue struct {
	goValue interface{}
	ValueAdapter
}

func newAnyValue(raw interface{}) Value {
	return &AnyValue{goValue: raw}
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

