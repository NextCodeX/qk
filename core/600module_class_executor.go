package core

import (
	"fmt"
	"reflect"
)


// 类对象: 用于内置对象执行
type ClassExecutor struct {
	raw interface{}
	name    string // class name
	fields  map[string]*FieldInfo
	methods map[string]*FunctionExecutor
}

var ClassType = reflect.TypeOf(&ClassExecutor{})

func newClassExecutor(name string, objPtr interface{}, objDoublePtr interface{}) Value {
	fs := collectFieldInfo(objPtr)
	mts := collectFunctionInfo(objDoublePtr)
	m := make(map[string]Value)
	for name, f := range fs {
		m[name] = newQKValue(f.get())
	}
	for name, mt := range mts {
		m[name] = newModuleFunc(name, mt)
	}
	return newClass(name, m)
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

//func (clazz *ClassExecutor) callMethod(name string, vals []interface{}) Value {
//	mt, exist := clazz.methods[name]
//	assert(!exist, fmt.Sprintf("method %v.%v is undefined!", clazz.name, name))
//	return callFunctionExecutor(mt, vals)
//}

//func evalClassMethod(any interface{}, name string, vals []interface{}) Value {
//	clazz := any.(*ClassExecutor)
//	return clazz.callMethod(name, vals)
//}

func evalClassField(any interface{}, attrname string) Value {
	clazz := any.(*ClassExecutor)
	return newQKValue(clazz.fieldValue(attrname))
}

func (clazz *ClassExecutor) String() string {
	si, ok := clazz.raw.(fmt.Stringer)
	if ok {
		return si.String()
	}
	return fmt.Sprint("class ", clazz.name)
}
