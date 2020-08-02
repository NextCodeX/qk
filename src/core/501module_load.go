package core

import (
	"fmt"
	"reflect"
	"strings"
)

var funcs =  make(map[string]*FunctionExecutor)


func collectFunctionInfo(obj interface{}, moduleName string)  {
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
			funcExe.Ins = append(funcExe.Ins, argType)
		}

		// out params
		outcount := methodType.NumOut()
		for ii := 0; ii < outcount; ii++ {
			argType := methodType.Out(ii)
			funcExe.Outs = append(funcExe.Outs, argType)
		}

		funcExe.Obj = methodObject
		funcExe.Name = standardName(moduleName, methodName)
		funcs[funcExe.Name] = funcExe
	}
}

func standardName(moduleName, methodName string) string {
	return fmt.Sprintf("%v_%v%v", moduleName, strings.ToLower(methodName[:1]), methodName[1:])
}

func Load() map[string]*FunctionExecutor {
	fmt.Println("register all module...")
	fileModuleInit()
	return funcs
}
