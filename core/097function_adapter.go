package core

import "fmt"

type FunctionAdapter struct {
	name     string
	instance Function
	StatementAdapter
	ValueAdapter
}

func (this *FunctionAdapter) init(name string, obj Function) {
	this.name = name
	this.instance = obj
	this.StatementAdapter.initStatement(obj)
}
func (this *FunctionAdapter) setArgs([]Value) {}

func (this *FunctionAdapter) varList() Variables {
	return nil
}
func (this *FunctionAdapter) parentFrame() Frame {
	return this.instance.getParent()
}

func (this *FunctionAdapter) val() interface{} {
	return this.instance
}
func (this *FunctionAdapter) typeName() string {
	return "Function"
}
func (this *FunctionAdapter) isFunction() bool {
	return true
}
func (this *FunctionAdapter) isObject() bool {
	return true
}

func (this *FunctionAdapter) get(key string) Value {
	if key == "type" {
		return callable(func() Value {
			return newQKValue("Function")
		})
	} else {
		return NULL
	}
}

func (this *FunctionAdapter) Pr() {
	fmt.Println(this.String())
}
func (this *FunctionAdapter) String() string {
	return "function " + this.name + "()"
}
