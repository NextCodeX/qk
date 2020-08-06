package core

import (
	"reflect"
)

type FunctionExecutor struct {
	name string
	ins []reflect.Type
	outs []reflect.Type
	obj reflect.Value
}

func (f *FunctionExecutor) Run(args []reflect.Value) interface{} {
	resList := f.obj.Call(args)
	if len(resList) < 1 {
		return nil
	}
	res := resList[0]
	return res.Interface()
}

