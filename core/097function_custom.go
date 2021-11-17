package core

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
	f.StatementAdapter.ts = ts
	f.paramNames = paramNames
	return f
}

// 初始化内置变量
func (fn *CustomFunction) setInternalVars(vars map[string]Value) {
	// 初始化main函数本地变量池
	fn.localVars = newVariables()
	for name, val := range vars {
		fn.localVars.add(name, val)
	}
}

// 调用自定义函数前， 初始化函数本地变量池
func (fn *CustomFunction) setQkArgs(args []Value) {
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
	if fn.localVars == nil || (fn.name != "main" && len(fn.paramNames) < 1) {
		fn.localVars = newVariables()
	}

	return fn.executeStatementList(fn.block, StmtListTypeFunc)
}

func (fn *CustomFunction) varList() Variables {
	return fn.localVars
}
