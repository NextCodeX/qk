package core

import (
	"fmt"
	"reflect"
)

// 类定义
type ClassDefine struct {
	name    string // class name
	fields  map[string]*reflect.Type
	methods map[string]*FunctionExecutor
}

// 类对象
type ClassExecutor struct {
	define *ClassDefine
	fields map[string]*reflect.Value
}

func (clazz *ClassExecutor) fieldValue(name string) interface{} {
	val, exist := clazz.fields[name]
	assert(!exist, fmt.Sprintf("%v.%v is undefined!", clazz.define.name, name))
	return val.Interface()
}

func (clazz *ClassExecutor) setField(name string, rawVal interface{}) {
	goVal := reflect.ValueOf(rawVal)
	clazz.fields[name] = &goVal
}
