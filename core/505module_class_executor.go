package core

import (
	"fmt"
	"reflect"
)

// 类对象: 用于内置对象执行
type ClassExecutor struct {
	raw     interface{}
	name    string // class name
	fields  map[string]*FieldInfo
	methods map[string]*FunctionExecutor
}

var ClassType = reflect.TypeOf(&ClassExecutor{})

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

func evalClassField(obj any, attrname string) Value {
	clazz := obj.(*ClassExecutor)
	return newQKValue(clazz.fieldValue(attrname))
}

func (clazz *ClassExecutor) String() string {
	si, ok := clazz.raw.(fmt.Stringer)
	if ok {
		return si.String()
	}
	return fmt.Sprint("class ", clazz.name)
}
