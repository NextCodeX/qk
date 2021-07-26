package core

import (
	"fmt"
	"reflect"
)

type AnyValue struct {
	goValue interface{}
}

func newAnyValue(raw interface{}) *AnyValue {
	return &AnyValue{raw}
}

func (any *AnyValue) val() interface{} {
	si, ok := any.goValue.(fmt.Stringer)
	if ok {
		return si.String()
	}
	return any.goValue
}
func (any *AnyValue) isNULL() bool {
	return false
}
func (any *AnyValue) isInt() bool {
	return false
}
func (any *AnyValue) isFloat() bool {
	return false
}
func (any *AnyValue) isBoolean() bool {
	return false
}
func (any *AnyValue) isString() bool {
	return false
}
func (any *AnyValue) isAny() bool {
	return true
}
func (any *AnyValue) isClass() bool {
	return reflect.TypeOf(any.goValue).AssignableTo(ClassType)
}
func (any *AnyValue) isJsonArray() bool {
	return false
}
func (any *AnyValue) isJsonObject() bool {
	return false
}