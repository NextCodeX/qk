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

// 初始化内置变量/函数
func (this *CustomFunction) setInternalVars(vars map[string]Value) {
	for name, val := range vars {
		this.localVars.add(name, val)
	}
	// 防止内置变量函数被覆盖无法使用，多提供一种方法以供调用
	this.localVars.add("_qk", jsonObject(vars))
}

// 调用自定义函数前， 初始化函数本地变量池
func (this *CustomFunction) setArgs(args []Value) {
	// 每次执行自定义函数前，初始化本地变量池
	this.localVars = newVariables()
	this.args = args

	argsLen := len(this.args)
	for i, paramName := range this.paramNames {
		var arg Value
		if i >= argsLen {
			arg = NULL
		} else {
			arg = this.args[i]
		}
		this.localVars.add(paramName, arg)
	}
}

func (this *CustomFunction) execute() StatementResult {
	if this.name != "main" && len(this.paramNames) < 1 {
		// 每次匿名函数调用时，重新初始化局部变量表
		this.localVars = newVariables()
	}

	return this.executeStatementList(this.block, StmtListTypeFunc)
}

func (this *CustomFunction) varList() Variables {
	return this.localVars
}

func (this *CustomFunction) ptr() string {
	return fmt.Sprintf("%p", this)
}
