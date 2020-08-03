package core

import (
	"reflect"
)

var moduleFuncs map[string]*FunctionExecutor

func init()  {
	moduleFuncs = Load()
}

func isModuleFunc(funcName string) bool {
	_, ok := moduleFuncs[funcName]
	return ok
}

func executeModuleFunc(funcName string, args []interface{}) *Value {
	f := moduleFuncs[funcName]
	params := extractModuleFuncArgs(f, args)
	res := f.Run(params)
	return toQKValue(res)
}

func extractModuleFuncArgs(f *FunctionExecutor, args []interface{}) []reflect.Value {
	var res []reflect.Value
	if len(f.Ins) == 1 && f.Ins[0].Kind() == reflect.Slice {
		res = append(res, reflect.ValueOf(args))
		return res
	}

	if len(args) < len(f.Ins) {
		runtimeExcption("execute", f.Name, ", arguments is too less")
		return nil
	}

	for i, t := range f.Ins {
		arg := args[i]
		if t != reflect.TypeOf(arg) {
			runtimeExcption("execute", f.Name, ", arguments type is not match!")
			return nil
		}
		res = append(res, reflect.ValueOf(arg))
	}
	return res
}