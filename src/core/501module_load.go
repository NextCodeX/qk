package core

import (
	"fmt"
	"reflect"
	"strings"
)

var funcs =  make(map[string]*FunctionExecutor)

type ModuleRegister struct{}

// register all module...
func Load() map[string]*FunctionExecutor {
	mr := &ModuleRegister{}
	v1 := reflect.ValueOf(&mr).Elem()
	for i := 0; i < v1.NumMethod(); i++ {
		v1.Method(i).Call(nil)
	}
	return funcs
}

func collectFieldInfo(obj interface{}) (res map[string]reflect.Type) {
	res = make(map[string]reflect.Type)
	v := reflect.ValueOf(obj).Elem()
	k := v.Type()
	for i := 0; i < v.NumField(); i++ {
		key := k.Field(i)
		val := v.Field(i)
		if !val.CanInterface() { //CanInterface(): 判断该成员变量是否能被获取值
			continue
		}
		res[key.Name] = val.Type()
	}
	return res
}

func collectFunctionInfo(obj interface{}) map[string]*FunctionExecutor {
	res := make(map[string]*FunctionExecutor)
	v1 := reflect.ValueOf(obj).Elem()
	k1 := v1.Type()
	for i := 0; i < v1.NumMethod(); i++ {
		funcExe := &FunctionExecutor{}

		methodName := k1.Method(i).Name
		methodObject := v1.Method(i)

		methodType := methodObject.Type()
		// in params
		incount := methodType.NumIn()
		for ii := 0; ii < incount; ii++ {
			argType := methodType.In(ii)
			funcExe.ins = append(funcExe.ins, argType)
		}

		// out params
		outcount := methodType.NumOut()
		for ii := 0; ii < outcount; ii++ {
			argType := methodType.Out(ii)
			funcExe.outs = append(funcExe.outs, argType)
		}

		funcExe.obj = methodObject
		res[methodName] = funcExe
	}
	return res
}

func functionRegister(module string, fmap map[string]*FunctionExecutor) {
	for name, f := range fmap {
		var funcKey string
		if module == "" {
			funcKey = formatName(module)
		} else {
			funcKey = standardName(module, name)
		}
		funcs[funcKey] = f
	}
}

func standardName(moduleName, methodName string) string {
	return fmt.Sprintf("%v_%v%v", moduleName, strings.ToLower(methodName[:1]), methodName[1:])
}

func formatName(methodName string) string {
	return fmt.Sprintf("%v%v", strings.ToLower(methodName[:1]), methodName[1:])
}


