package core

import (
	"reflect"
	"sync"
)

// 用于实现多种形式的方法调用。
type Object interface {
	get(key string) Value
	Value
}

type ClassObject struct {
	name string
	raw interface{}
	show *reflect.Value
	methods map[string]Function
	mux sync.Mutex
	ValueAdapter
}

func newClass(name string, raw interface{}) Value {
	return &ClassObject{name: name, raw: raw}
}

func (clazz *ClassObject) String() string {
	clazz.initMethods()
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
	clazz.initMethods()
	mt, ok := clazz.methods[key]
	if !ok {
		if key == "type" {
			return newAnonymousFunc(func() Value {
				return newQKValue(clazz.name)
			})
		} else {
			return NULL
		}
	} else {
		return mt
	}
}

func (clazz *ClassObject) initMethods() {
	if clazz.methods != nil {
		return
	}
	clazz.mux.Lock()
	if clazz.methods != nil {
		return
	}

	mts := collectFunctionInfo(clazz.raw)
	clazz.methods = make(map[string]Function)
	for name, mt := range mts {
		if name == "string" {
			clazz.show = &mt.obj
		}
		clazz.methods[name] = newModuleFunc(name, mt)
	}

	defer clazz.mux.Unlock()
}



