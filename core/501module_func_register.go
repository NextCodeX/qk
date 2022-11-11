package core

import (
	"fmt"
	"reflect"
	"strings"
)

// 通过这个类型将所有内部函数串在一起
type InternalFunctionSet struct {
	owner            *Interpreter
	internalFuntions map[string]Value // 包含所有内部函数的一个集合
}

func newInternalFunctionSet(parent *Interpreter) *InternalFunctionSet {
	// 通过反射收集内部函数信息
	fns := &InternalFunctionSet{
		owner:            parent,
		internalFuntions: make(map[string]Value),
	}
	fmap := collectFunctionInfo(&fns)
	for fname, moduleFunc := range fmap {
		fns.internalFuntions[fname] = newInternalFunc(fname, moduleFunc)
	}
	return fns
}

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

// 将函数名第一个字母转小写
func formatName(methodName string) string {
	return fmt.Sprintf("%v%v", strings.ToLower(methodName[:1]), methodName[1:])
}
