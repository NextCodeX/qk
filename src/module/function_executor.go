package module

import (
	"reflect"
)

type FunctionExecutor struct {
	Name string
	Ins []reflect.Type
	Outs []reflect.Type
	Obj reflect.Value
}

func (f *FunctionExecutor) Run(args []reflect.Value) interface{} {
	resList := f.Obj.Call(args)
	if len(resList) < 1 {
		return nil
	}
	res := resList[0]
	return res.Interface()
}

