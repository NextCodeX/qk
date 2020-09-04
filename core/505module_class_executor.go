package core

import (
	"fmt"
	"reflect"
)


// 类对象: 用于内置对象执行
type ClassExecutor struct {
	name    string // class name
	fields  map[string]*FieldInfo
	methods map[string]*FunctionExecutor
}

func classType() reflect.Type {
	t := reflect.TypeOf(&ClassExecutor{})
	return t
}

func newClassExecutor(name string, objPtr interface{}, objDoublePtr interface{}) *ClassExecutor {
	fs := collectFieldInfo(objPtr)
	mts := collectFunctionInfo(objDoublePtr)
	res := &ClassExecutor{
		name: name,
		fields: fs,
		methods: mts,
	}
	return res
}

func (clazz *ClassExecutor) fieldValue(name string) interface{} {
	field, exist := clazz.fields[name]
	assert(!exist, fmt.Sprintf("field %v.%v is undefined!", clazz.name, name))
	return field.get()
}

func (clazz *ClassExecutor) setField(name string, rawVal interface{}) {
	field, exist := clazz.fields[name]
	assert(!exist, fmt.Sprintf("field %v.%v is undefined!", clazz.name, name))
	field.set(rawVal)
}

func (clazz *ClassExecutor) callMethod(name string, vals []interface{}) *Value {
	mt, exist := clazz.methods[name]
	assert(!exist, fmt.Sprintf("method %v.%v is undefined!", clazz.name, name))
	return callFunctionExecutor(mt, vals)
}

func evalClassMethod(any interface{}, name string, vals []interface{}) *Value {
	clazz := any.(*ClassExecutor)
	return clazz.callMethod(name, vals)
}

func evalClassField(any interface{}, attrname string) *Value {
	clazz := any.(*ClassExecutor)
	return toQKValue(clazz.fieldValue(attrname))
}
