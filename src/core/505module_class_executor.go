package core

import (
	"fmt"
	"reflect"
)


// 类对象
type ClassExecutor struct {
	name    string // class name
	fields  map[string]*FieldInfo
	methods map[string]*FunctionExecutor
}

func (clazz *ClassExecutor) fieldValue(name string) interface{} {
	val, exist := clazz.fields[name]
	assert(!exist, fmt.Sprintf("%v.%v is undefined!", clazz.name, name))
	return val.v.Interface()
}

func (clazz *ClassExecutor) setField(name string, rawVal interface{}) {
	goVal := reflect.ValueOf(rawVal)
	clazz.fields[name].v = goVal
}
