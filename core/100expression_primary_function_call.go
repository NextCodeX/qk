package core

type FunctionCallPrimaryExpression struct {
	args []Expression // 函数调用参数
	PrimaryExpressionImpl
}

func newFunctionCallPrimaryExpression(args []Expression) PrimaryExpression {
	expr := &FunctionCallPrimaryExpression{}
	expr.t = FunctionCallPrimaryExpressionType
	expr.args = args
	expr.doExec = expr.doExecute
	return expr
}

func (priExpr *FunctionCallPrimaryExpression) setParent(p Function) {
	priExpr.ExpressionAdapter.setParent(p)

	for _, subExpr := range priExpr.args {
		subExpr.setParent(p)
	}
}

func (priExpr *FunctionCallPrimaryExpression) doExecute() Value {
	runtimeExcption("running FunctionCallPrimaryExpression.doExecute is error")
	return nil
}

func (priExpr *FunctionCallPrimaryExpression) callFunc(fn Function) Value {
	if f1, ok := fn.(*InternalFunction); ok {
		f1.setGoArgs(evalGoValues(priExpr.args))
	} else if f2, ok := fn.(*CustomFunction); ok {
		f2.setArgs(evalQKValues(priExpr.args))
	}

	return fn.execute().value()
}
