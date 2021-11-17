package core

type FunctionPrimaryExpression struct {
	name       string
	paramNames []string
	bodyTokens []Token
	PrimaryExpressionImpl
}

func newFunctionPrimaryExpression(name string, paramNames []string, bodyTokens []Token) PrimaryExpression {
	expr := &FunctionPrimaryExpression{}
	expr.t = FunctionPrimaryExpressionType
	expr.name = name
	expr.paramNames = paramNames
	expr.bodyTokens = bodyTokens

	expr.doExec = expr.doExecute
	return expr
}

func (priExpr *FunctionPrimaryExpression) doExecute() Value {
	fn := newCustomFunction(priExpr.name, priExpr.bodyTokens, priExpr.paramNames)
	fn.setParent(priExpr.parent)
	Compile(fn)

	funcName := priExpr.name
	if funcName != "" {
		priExpr.setVar(funcName, fn)
	}

	return fn
}
