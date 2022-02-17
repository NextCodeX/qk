package core

import (
	"fmt"
	"reflect"
)

type InternalFunction struct {
	rawArgs    []interface{}     // 实参
	moduleFunc *FunctionExecutor // 函数对象
	FunctionAdapter
}

func newInternalFunc(name string, f *FunctionExecutor) Function {
	fn := &InternalFunction{moduleFunc: f}
	fn.init(name, fn)
	return fn
}

func extractModuleFuncArgs(f *FunctionExecutor, args []interface{}) []reflect.Value {
	var res []reflect.Value
	var bs []byte
	var ss []string
	if len(f.ins) == 1 && f.ins[0].Kind() == reflect.Slice {
		if f.ins[0] == reflect.TypeOf(bs) {
			// []byte, []string
			res = append(res, reflect.ValueOf(args[0]))
		} else if f.ins[0] == reflect.TypeOf(ss) {
			strArr := make([]string, 0, len(args))
			for _, v := range args {
				strArr = append(strArr, fmt.Sprint(v))
			}
			res = append(res, reflect.ValueOf(strArr))
		} else {
			res = append(res, reflect.ValueOf(args))
		}

		return res
	}

	if len(args) < len(f.ins) {
		errorf("execute %v(): arguments is too less, require %v have %v", f.name, len(f.ins), len(args))
		return nil
	}

	for i, t := range f.ins {
		arg := args[i]
		if t.Kind() == reflect.Int && reflect.TypeOf(arg).Kind() == reflect.Int64 {
			i := arg.(int64)
			res = append(res, reflect.ValueOf(int(i)))
			continue
		}
		if t.Kind() != reflect.Interface && t != reflect.TypeOf(arg) {
			runtimeExcption("execute", f.name, "(), arguments type is not match!", t, reflect.TypeOf(arg))
			return nil
		}
		res = append(res, reflect.ValueOf(arg))
	}
	return res
}

// 设置函数参数
func (fn *InternalFunction) setGoArgs(rawArgs []interface{}) {
	fn.rawArgs = rawArgs
}

func (fn *InternalFunction) execute() StatementResult {
	params := extractModuleFuncArgs(fn.moduleFunc, fn.rawArgs)
	res := fn.moduleFunc.Run(params)
	return newStatementResult(StatementNormal, newQKValue(res))
}
