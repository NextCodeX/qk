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

func toQKValue(v interface{}) *Value {
	typ := reflect.TypeOf(v)
	kind := typ.Kind()
	switch kind {
	case reflect.Map:
		mapRes := make(map[string]*Value)
		m := v.(map[string]interface{})
		for key, value := range m {
			var qkVal *Value
			if isDecomposable(value) {
				qkVal = toQKValue(value)
			} else {
				qkVal = newVal(value)
			}
			mapRes[key] = qkVal
		}
		tmp := toJSONObject(mapRes)
		return newVal(tmp)

	case reflect.Slice:
		var arrRes []*Value
		list := v.([]interface{})
		for _, item := range list {
			var qkVal *Value
			if isDecomposable(item) {
				qkVal = toQKValue(item)
			} else {
				qkVal = newVal(item)
			}
			arrRes = append(arrRes, qkVal)
		}
		tmp := toJSONArray(arrRes)
		return newVal(tmp)

	default:
		return newVal(v)
	}
}

func isDecomposable(v interface{}) bool {
	kind := reflect.TypeOf(v).Kind()
	return kind == reflect.Map || kind == reflect.Slice
}

func extractModuleFuncArgs(f *FunctionExecutor, args []interface{}) []reflect.Value {
	if len(args) < len(f.Ins) {
		runtimeExcption("execute", f.Name, ", arguments is too less")
		return nil
	}
	var res []reflect.Value
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