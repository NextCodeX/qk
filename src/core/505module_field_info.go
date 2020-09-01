package core

import (
	"reflect"
)


// 成员变量相关信息
type FieldInfo struct {
	name    string // field name
	t reflect.Type // field go type
	v reflect.Value // field go internal value
}

func (f *FieldInfo) set(val interface{}) {
	f.v = reflect.ValueOf(val)
}

func (f *FieldInfo) get() interface{} {
	return f.v.Interface()
}



