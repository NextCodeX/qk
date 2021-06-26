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
	return callFunctionExecutor(f, args)
}

func callFunctionExecutor(f *FunctionExecutor, args []interface{}) *Value {
	params := extractModuleFuncArgs(f, args)
	res := f.Run(params)
	return toQKValue(res)
}

func extractModuleFuncArgs(f *FunctionExecutor, args []interface{}) []reflect.Value {
	var res []reflect.Value
	if len(f.ins) == 1 && f.ins[0].Kind() == reflect.Slice {
		res = append(res, reflect.ValueOf(args))
		return res
	}

	if len(args) < len(f.ins) {
		runtimeExcption("execute", f.name, ", arguments is too less")
		return nil
	}

	for i, t := range f.ins {
		arg := args[i]
		if t.Kind() != reflect.Interface && t != reflect.TypeOf(arg) {
			runtimeExcption("execute", f.name, ", arguments type is not match!", t, reflect.TypeOf(arg))
			return nil
		}
		res = append(res, reflect.ValueOf(arg))
	}
	return res
}