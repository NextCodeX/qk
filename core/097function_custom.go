package core

import "fmt"

type CustomFunction struct {
	paramNames []string // 形参名列表
	args       []Value  // 实参列表

	localVars Variables // 当前作用域的变量列表
	FunctionAdapter
}

func newMainFunction() *CustomFunction {
	return newCustomFunction("main", nil, nil)
}

func newCustomFunction(name string, ts []Token, paramNames []string) *CustomFunction {
	f := &CustomFunction{}
	f.init(name, f)
	f.setTokenList(ts)
	f.paramNames = paramNames
	f.localVars = newVariables()
	return f
}

// 初始化内置变量
func (fn *CustomFunction) setInternalVars(vars map[string]Value) {
	for name, val := range vars {
		fn.localVars.add(name, val)
	}
}

// 调用自定义函数前， 初始化函数本地变量池
func (fn *CustomFunction) setArgs(args []Value) {
	// 每次执行自定义函数前，初始化本地变量池
	fn.localVars = newVariables()
	fn.args = args

	argsLen := len(fn.args)
	for i, paramName := range fn.paramNames {
		var arg Value
		if i >= argsLen {
			arg = NULL
		} else {
			arg = fn.args[i]
		}
		fn.localVars.add(paramName, arg)
	}
}

func (fn *CustomFunction) execute() StatementResult {
	if fn.name != "main" && len(fn.paramNames) < 1 {
		// 每次匿名函数调用时，重新初始化局部变量表
		fn.localVars = newVariables()
	}

	return fn.executeStatementList(fn.block, StmtListTypeFunc)
}

func (fn *CustomFunction) varList() Variables {
	return fn.localVars
}

func (fn *CustomFunction) ptr() string {
	return fmt.Sprintf("%p", fn)
}
