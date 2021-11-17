package core

type CallableFunction struct {
	action func() Value
	FunctionAdapter
}

func callable(fn func() Value) Function {
	f := &CallableFunction{action: fn}
	f.init("callable", f)
	return f
}

func (fn *CallableFunction) execute() StatementResult {
	return newStatementResult(StatementNormal, fn.action())
}

type RunnableFunction struct {
	action func()
	FunctionAdapter
}

func runnable(fn func()) Function {
	f := &RunnableFunction{action: fn}
	f.init("runnable", f)
	return f
}

func (fn *RunnableFunction) execute() StatementResult {
	fn.action()
	return newStatementResult(StatementNormal, NULL)
}
