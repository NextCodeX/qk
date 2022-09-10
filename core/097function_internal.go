package core

import (
	"fmt"
	"reflect"
)

var EMPTY_BYTE_ARRAY []byte
var EMPTY_STR_ARRAY []string
var DEFAULT_BOOL_VALUE bool
var DEFAULT_INT_VALUE int
var DEFAULT_I32_VALUE int32
var DEFAULT_I64_VALUE int64
var DEFAULT_F32_VALUE float32
var DEFAULT_F64_VALUE float64
var DEFAULT_STR_VALUE string
var JSON_ARRAY_TYPE = reflect.TypeOf(emptyArray())
var JSON_OBJECT_TYPE = reflect.TypeOf(emptyJsonObject())
var BYTE_ARRAY_TYPE = reflect.TypeOf(EMPTY_BYTE_ARRAY)
var STR_ARRAY_TYPE = reflect.TypeOf(DEFAULT_STR_VALUE)
var BOOL_TYPE = reflect.TypeOf(DEFAULT_BOOL_VALUE)
var INT_TYPE = reflect.TypeOf(DEFAULT_INT_VALUE)
var I32_TYPE = reflect.TypeOf(DEFAULT_I32_VALUE)
var I64_TYPE = reflect.TypeOf(DEFAULT_I64_VALUE)
var F32_TYPE = reflect.TypeOf(DEFAULT_F32_VALUE)
var F64_TYPE = reflect.TypeOf(DEFAULT_F64_VALUE)
var STR_TYPE = reflect.TypeOf(DEFAULT_STR_VALUE)

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

func extractModuleFuncArgs(f *FunctionExecutor, args []any) []reflect.Value {
	// 数组类型处理: []byte, []string, []any
	if len(f.ins) == 1 && f.ins[0].Kind() == reflect.Slice {
		return handleArrayArgs(f, args)
	}

	var res []reflect.Value
	funcName := f.name
	for i, t := range f.ins {
		arg := findFuncArg(funcName, i, t, args)
		res = append(res, arg)
	}
	return res
}

func handleArrayArgs(f *FunctionExecutor, args []any) []reflect.Value {
	var res []reflect.Value
	if f.ins[0] == BYTE_ARRAY_TYPE {
		// []byte
		res = append(res, reflect.ValueOf(args[0]))
	} else if f.ins[0] == STR_ARRAY_TYPE {
		// []string
		strArr := make([]string, 0, len(args))
		for _, v := range args {
			strArr = append(strArr, fmt.Sprint(v))
		}
		res = append(res, reflect.ValueOf(strArr))
	} else {
		// [] any
		res = append(res, reflect.ValueOf(args))
	}
	return res
}

func findFuncArg(funcName string, index int, t reflect.Type, args []any) reflect.Value {
	// 允许实参数量小于形参数量
	var resTmp any
	if index < len(args) {
		resTmp = args[index]
	} else {
		resTmp = nil
	}

	if resTmp == nil {
		switch {
		case JSON_ARRAY_TYPE.AssignableTo(t):
			resTmp = emptyArray()
		case JSON_OBJECT_TYPE.AssignableTo(t):
			resTmp = emptyJsonObject()
		case BYTE_ARRAY_TYPE.AssignableTo(t):
			resTmp = EMPTY_BYTE_ARRAY
		case BOOL_TYPE.AssignableTo(t):
			resTmp = DEFAULT_BOOL_VALUE
		case INT_TYPE.AssignableTo(t):
			resTmp = DEFAULT_INT_VALUE
		case I32_TYPE.AssignableTo(t):
			resTmp = DEFAULT_I32_VALUE
		case I64_TYPE.AssignableTo(t):
			resTmp = DEFAULT_I64_VALUE
		case F32_TYPE.AssignableTo(t):
			resTmp = DEFAULT_F32_VALUE
		case F64_TYPE.AssignableTo(t):
			resTmp = DEFAULT_F64_VALUE
		case STR_TYPE.AssignableTo(t):
			resTmp = DEFAULT_STR_VALUE
		}
	} else {
		if reflect.TypeOf(resTmp).Kind() == reflect.Int64 {
			if t.Kind() == reflect.Int {
				return reflect.ValueOf(int(resTmp.(int64)))
			}
			if t.Kind() == reflect.Int32 {
				return reflect.ValueOf(int32(resTmp.(int64)))
			}
		}
		if t.Kind() != reflect.Interface && t != reflect.TypeOf(resTmp) {
			runtimeExcption("execute", funcName, "(), arguments type is not match!", t, reflect.TypeOf(resTmp))
			return reflect.ValueOf(nil)
		}
	}
	return reflect.ValueOf(resTmp)
}

// 设置函数参数
func (this *InternalFunction) setGoArgs(rawArgs []interface{}) {
	this.rawArgs = rawArgs
}

func (this *InternalFunction) execute() StatementResult {
	params := extractModuleFuncArgs(this.moduleFunc, this.rawArgs)
	res := this.moduleFunc.Run(params)
	return newStatementResult(StatementNormal, newQKValue(res))
}
