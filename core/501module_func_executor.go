package core

import (
	"reflect"
)

// 用于内置函数的执行
type FunctionExecutor struct {
	name string // 函数名称
	ins []reflect.Type //入参类型
	outs []reflect.Type // 出参类型
	obj reflect.Value // 函数对象
}

func (f *FunctionExecutor) Run(args []reflect.Value) interface{} {
	resList := f.obj.Call(args)
	if len(resList) < 1 {
		return nil
	}
	res := resList[0]
	return res.Interface()
}

func (f *FunctionExecutor) InNum() int {
	return len(f.ins)
}

func (f *FunctionExecutor) InLastIndex() int {
	return len(f.ins) - 1
}

func (f *FunctionExecutor) OutNum() int {
	return len(f.outs)
}

