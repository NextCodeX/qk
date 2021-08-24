package core

import (
	"fmt"
	"reflect"
)

type Function interface {
	setRaw(ts []Token)
	setParamNames(paramNames []string)

	setPreVar(key string, value Value)
	getName() string
	getParent() Function
	getLocalVars() Variables

	setArgs(args []Value)
	setRawArgs(rawArgs []interface{})

	isInternalFunc() bool

	getLocalFunctions() map[string]Function
	addFunc(functionName string, internalFunc Function)
	String() string
	Statement
	Value
}

type FunctionImpl struct {
	local      Variables        // 当前作用域的变量列表
	fns 	   map[string]Function
	name       string           // 函数名
	paramNames []string         // 形参名列表
	args []Value
	rawArgs []interface{}
	moduleFunc *FunctionExecutor
	anonymousFunc func()Value
	internalFunc func([]interface{})interface{}
	internalFuncFlag bool
	preVars JSONObject // 预设变量列表
	StatementAdapter
	ValueAdapter
}

func newFuncWithoutTokens(name string) Function {
	return newFunc(name, nil, nil)
}

func newFunc(name string, ts []Token, paramNames []string) Function {
	f := &FunctionImpl{name:name}
	f.StatementAdapter.ts = ts
	f.paramNames = paramNames
	f.fns = make(map[string]Function)
	f.initStatement(f)
	return f
}

func newAnonymousFunc(anonymousFunc func()Value) Function {
	f := &FunctionImpl{}
	f.anonymousFunc = anonymousFunc
	f.internalFuncFlag = true
	f.initStatement(f)
	return f
}

func newInternalFunc(name string, internalFunc func([]interface{})interface{}) Function {
	f := &FunctionImpl{name:name}
	f.internalFunc = internalFunc
	f.internalFuncFlag = true
	f.initStatement(f)
	return f
}

func newModuleFunc(name string, moduleFunc *FunctionExecutor) Function {
	f := &FunctionImpl{name:name}
	f.moduleFunc = moduleFunc
	f.internalFuncFlag = true
	f.initStatement(f)
	return f
}

func extractModuleFuncArgs(f *FunctionExecutor, args []interface{}) []reflect.Value {
	var res []reflect.Value
	if len(f.ins) == 1 && f.ins[0].Kind() == reflect.Slice {
		res = append(res, reflect.ValueOf(args))
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

// 预设变量
func (f *FunctionImpl) setPreVar(key string, value Value) {
	if f.preVars == nil {
		m := make(map[string]Value)
		f.preVars = jsonObject(m)
	}
	f.preVars.put(key, value)
}

func (f *FunctionImpl) isInternalFunc() bool {
	return f.internalFuncFlag
}

func (f *FunctionImpl) getLocalFunctions() map[string]Function {
	return f.fns
}

func (f *FunctionImpl) addFunc(functionName string, localFunc Function) {
	f.fns[functionName] = localFunc
}

func (f *FunctionImpl) setParamNames(paramNames []string) {
	f.paramNames = paramNames
}

func (f *FunctionImpl) getName() string {
	return f.name
}
func (f *FunctionImpl) getLocalVars() Variables {
	return f.local
}

func (f *FunctionImpl) parse() {
	fmt.Println("parse function？？？")
}

func (f *FunctionImpl) setRawArgs(rawArgs []interface{}) {
	f.rawArgs = rawArgs
}

func (f *FunctionImpl) setArgs(args []Value) {
	f.args = args
}

func (f *FunctionImpl) execute() StatementResult {
	if f.internalFuncFlag {
		var res interface{}
		if f.internalFunc != nil {
			res = f.internalFunc(f.rawArgs)
		} else if f.moduleFunc != nil {
			params := extractModuleFuncArgs(f.moduleFunc, f.rawArgs)
			res = f.moduleFunc.Run(params)
		} else if f.anonymousFunc != nil {
			res = f.anonymousFunc()
		} else {}
		return newStatementResult(StatementNormal, newQKValue(res))
	}

	// 每次执行自定义函数前，初始化本地变量池
	f.local = newVariables()
	defer func() {f.local = nil}()

	// 初始化预设变量(仅main函数使用)
	if f.preVars != nil {
		obj := f.preVars
		for _, key := range obj.keys() {
			f.local.add(key, obj.get(key))
		}
	}

	for i, paramName := range f.paramNames {
		var arg Value
		if i >= len(f.args) {
			arg = NULL
		} else {
			arg = f.args[i]
		}
		f.local.add(paramName, arg)
	}
	return f.executeStatementList(f.block, StmtListTypeFunc)
}


func (f *FunctionImpl) val() interface{} {
	return f
}
func (f *FunctionImpl) typeName() string {
	return "Function"
}
func (f *FunctionImpl) isFunction() bool {
	return true
}
func (f *FunctionImpl) isObject() bool {
	return true
}

func (f *FunctionImpl) get(key string) Value {
	if key == "type" {
		return newAnonymousFunc(func() Value {
			return newQKValue("Function")
		})
	} else {
		return NULL
	}
}

func (f *FunctionImpl) String() string {
	return "func:" + f.name + "()"
}