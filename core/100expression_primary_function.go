package core

type FunctionPrimaryExpression struct {
	name       string
	paramNames []string
	bodyTokens []Token

	cache []Statement
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
	return priExpr.funcDefinition()
}

func (priExpr *FunctionPrimaryExpression) funcDefinition() Function {
	fn := newCustomFunction(priExpr.name, priExpr.bodyTokens, priExpr.paramNames)
	fn.setParent(priExpr.parent)
	Compile(fn)
	//if priExpr.cache == nil {
	//	Compile(fn)
	//	priExpr.cache = fn.stmts()
	//} else {
	//	fn.setStatements(priExpr.cache)
	//	for _, stmt := range fn.stmts() {
	//		stmt.setParent(fn)
	//		stmt.parse()
	//	}
	//}
	return fn
}

func (priExpr *FunctionPrimaryExpression) declareFunction() {
	funcName := priExpr.name
	if funcName != "" {
		fn := priExpr.funcDefinition()
		priExpr.setVar(funcName, fn)
	}
}
