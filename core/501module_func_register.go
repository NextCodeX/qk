package core

import (
	"fmt"
	"reflect"
	"strings"
)

// internal function register
func init() {
	fns := &InternalFunctionSet{}
	fmap := collectFunctionInfo(&fns)
	functionRegister(fmap)
}

// 通过这个类型将所有内部函数串在一起
type InternalFunctionSet struct{}

// 通过反射收集函数信息
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

// 将内部函数注册到主函数(main)的内部栈中
func functionRegister(fmap map[string]*FunctionExecutor) {
	for fname, f := range fmap {
		addModuleFunc(fname, f)
	}
}

// 将函数名第一个字母转小写
func formatName(methodName string) string {
	return fmt.Sprintf("%v%v", strings.ToLower(methodName[:1]), methodName[1:])
}



