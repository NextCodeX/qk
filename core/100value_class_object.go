package core

import (
	"reflect"
)

// 用于实现多种形式的方法调用。
type Object interface {
	get(key string) Value
	Value
}

type ClassObject struct {
	name    string
	raw     interface{}
	show    *reflect.Value
	methods map[string]Function
	ValueAdapter
}

func newClass(name string, raw interface{}) Value {
	clazz := &ClassObject{name: name, raw: raw}
	clazz.initMethods()
	return clazz
}

func (clazz *ClassObject) initAsClass(name string, raw interface{}) {
	clazz.name = name
	clazz.raw = raw
	clazz.initMethods()
}

func (clazz *ClassObject) String() string {
	if clazz.show != nil {
		var args []reflect.Value
		resList := clazz.show.Call(args)
		res := resList[0].Interface()
		return res.(string)
	} else {
		return "class " + clazz.name
	}
}

func (clazz *ClassObject) val() interface{} {
	return clazz
}
func (clazz *ClassObject) typeName() string {
	return clazz.name
}
func (clazz *ClassObject) isObject() bool {
	return true
}

func (clazz *ClassObject) get(key string) Value {
	mt, ok := clazz.methods[key]
	if ok {
		return mt
	}

	if key != "type" {
		return NULL
	}

	return callable(func() Value {
		return newQKValue(clazz.name)
	})
}

func (clazz *ClassObject) initMethods() {
	mts := collectFunctionInfo(clazz.raw)
	clazz.methods = make(map[string]Function)
	for name, mt := range mts {
		if name == "string" {
			clazz.show = &mt.obj
		}
		clazz.methods[name] = newInternalFunc(name, mt)
	}
}
