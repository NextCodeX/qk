package core

type ExpressionStatement struct {
	expr Expression
	StatementAdapter
}

func newExpressionStatement(ts []Token) Statement {
	stmt := &ExpressionStatement{}
	stmt.setTokenList(ts)
	stmt.initStatement(stmt)
	return stmt
}

func (exprStmt *ExpressionStatement) isExpressionStatement() bool {
	return true
}

func (exprStmt *ExpressionStatement) parse() {
	expr := extractExpression(exprStmt.tokenList())
	exprStmt.expr = expr

	// 将stack从statement下传至expression
	// 复合expression可以通过重载 setFrame()方法将stack 下传至 各子expression
	expr.setParent(exprStmt.getStack())

	// 函数定义
	if fnExpr, ok := expr.(*FunctionPrimaryExpression); ok {
		fnExpr.declareFunction()
	}
}

func (exprStmt *ExpressionStatement) execute() StatementResult {
	res := exprStmt.expr.execute()
	return newStatementResult(StatementNormal, res)
}
