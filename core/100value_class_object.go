package core

import (
	"fmt"
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

func (this *ClassObject) initAsClass(name string, raw interface{}) {
	this.name = name
	this.raw = raw
	this.initMethods()
}

func (this *ClassObject) Pr() {
	fmt.Println(this.String())
}
func (this *ClassObject) String() string {
	if this.show != nil {
		var args []reflect.Value
		resList := this.show.Call(args)
		res := resList[0].Interface()
		return res.(string)
	} else {
		return "class " + this.name
	}
}

func (this *ClassObject) val() interface{} {
	return this
}
func (this *ClassObject) typeName() string {
	return this.name
}
func (this *ClassObject) isObject() bool {
	return true
}

func (this *ClassObject) get(key string) Value {
	mt, ok := this.methods[key]
	if ok {
		return mt
	}

	if key == "type" {
		return callable(func() Value {
			return newQKValue(this.name)
		})
	} else if key == "pr" {
		return runnable(func() {
			fmt.Println(this.String())
		})
	} else {
		return NULL
	}
}

func (this *ClassObject) initMethods() {
	mts := collectFunctionInfo(this.raw)
	this.methods = make(map[string]Function)
	for name, mt := range mts {
		if name == "string" {
			this.show = &mt.obj
		}
		this.methods[name] = newInternalFunc(name, mt)
	}
}
