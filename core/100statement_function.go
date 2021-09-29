package core

import (
	"bytes"
	"fmt"
	"reflect"
)

type Function interface {
	setRaw(ts []Token)
	setParamNames(paramNames []string)

	getName() string
	getParent() Function
	getCurrentStack() Function
	getLocalVars() Variables

	setInternalVars(vars map[string]Value)
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
	local            Variables // 当前作用域的变量列表
	fns              map[string]Function
	name             string   // 函数名
	paramNames       []string // 形参名列表
	args             []Value
	rawArgs          []interface{}
	moduleFunc       *FunctionExecutor
	anonymousFunc    func() Value
	internalFunc     func([]interface{}) interface{}
	internalFuncFlag bool
	StatementAdapter
	ValueAdapter
}

func newFuncWithoutTokens(name string) Function {
	return newFunc(name, nil, nil)
}

func newFunc(name string, ts []Token, paramNames []string) Function {
	f := &FunctionImpl{name: name}
	f.StatementAdapter.ts = ts
	f.paramNames = paramNames
	f.fns = make(map[string]Function)
	f.initStatement(f)
	return f
}

func newAnonymousFunc(anonymousFunc func() Value) Function {
	f := &FunctionImpl{}
	f.anonymousFunc = anonymousFunc
	f.internalFuncFlag = true
	f.initStatement(f)
	return f
}

func newInternalFunc(name string, internalFunc func([]interface{}) interface{}) Function {
	f := &FunctionImpl{name: name}
	f.internalFunc = internalFunc
	f.internalFuncFlag = true
	f.initStatement(f)
	return f
}

func newModuleFunc(name string, moduleFunc *FunctionExecutor) Function {
	f := &FunctionImpl{name: name}
	f.moduleFunc = moduleFunc
	f.internalFuncFlag = true
	f.initStatement(f)
	return f
}

func extractModuleFuncArgs(f *FunctionExecutor, args []interface{}) []reflect.Value {
	var res []reflect.Value
	var bs []byte
	if len(f.ins) == 1 && f.ins[0].Kind() == reflect.Slice {
		if f.ins[0] == reflect.TypeOf(bs) {
			res = append(res, reflect.ValueOf(args[0]))
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

func (f *FunctionImpl) setInternalVars(vars map[string]Value) {
	// 初始化main函数本地变量池
	f.local = newVariables()
	for name, val := range vars {
		f.local.add(name, val)
	}
}

func (f *FunctionImpl) setArgs(args []Value) {
	// 每次执行自定义函数前，初始化本地变量池
	f.local = newVariables()
	f.args = args

	for i, paramName := range f.paramNames {
		var arg Value
		if i >= len(f.args) {
			arg = NULL
		} else {
			arg = f.args[i]
		}
		f.local.add(paramName, arg)
	}
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
		} else {
		}
		return newStatementResult(StatementNormal, newQKValue(res))
	}

	if f.local == nil || (f.name != "main" && len(f.paramNames) < 1) {
		f.local = newVariables()
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

func (f *FunctionImpl) showArgs() string {
	var res bytes.Buffer
	for _, name := range f.paramNames {
		res.WriteString(fmt.Sprint(f.getVar(name)))
		res.WriteString(", ")
	}
	return res.String()
}
