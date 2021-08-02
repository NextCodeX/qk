package core

type ExpressionStatement struct {
	expr Expression
	StatementAdapter
}

func newExpressionStatement(ts []Token) Statement {
	stmt := &ExpressionStatement{}
	stmt.StatementAdapter.ts = ts
	stmt.initStatement(stmt)
	return stmt
}

func (exprStmt *ExpressionStatement) isExpressionStatement() bool {
	return true
}

func (exprStmt *ExpressionStatement) parse() {
	expr := extractExpression(exprStmt.ts)
	exprStmt.expr = expr

	// 将stack从statement下传至expression
	// 复合expression可以通过重载 setStack()方法将stack 下传至 各子expression
	expr.setStack(exprStmt.getStack())
}

func (exprStmt *ExpressionStatement) execute() StatementResult {
	res := exprStmt.expr.execute()
	return newStatementResult(StatementNormal, res)
}

