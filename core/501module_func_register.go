package core

import (
	"fmt"
	"reflect"
	"strings"
)


func collectFunctionInfo(objDoublePtr interface{}) (res map[string]*FunctionExecutor) {
	res = make(map[string]*FunctionExecutor)
	v1 := reflect.ValueOf(objDoublePtr).Elem()
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
		funcExe.name = formatName(methodName)
		res[funcExe.name] = funcExe
	}
	return res
}

func functionRegister(module string, fmap map[string]*FunctionExecutor) {
	for name, f := range fmap {
		if module == "" {
			addModuleFunc(name, f)
			continue
		}

		funcKey := standardName(module, name)
		addModuleFunc(funcKey, f)
	}
}

func standardName(moduleName, methodName string) string {
	return fmt.Sprintf("%v_%v", moduleName, methodName)
}

func formatName(methodName string) string {
	return fmt.Sprintf("%v%v", strings.ToLower(methodName[:1]), methodName[1:])
}



