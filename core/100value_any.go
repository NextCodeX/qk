package core

import "reflect"

type AnyValue struct {
	goValue interface{}
	ValueAdapter
}

func newAnyValue(raw interface{}) *AnyValue {
	return &AnyValue{goValue: raw}
}

func (any *AnyValue) val() interface{} {
	//si, ok := any.goValue.(fmt.Stringer)
	//if ok {
	//	return si.String()
	//}
	return any.goValue
}
func (any *AnyValue) isAny() bool {
	return true
}
func (any *AnyValue) isClass() bool {
	return reflect.TypeOf(any.goValue).AssignableTo(ClassType)
}

