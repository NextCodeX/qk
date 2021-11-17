package core

type ReturnStatement struct {
	expr Expression
	StatementAdapter
}

func newReturnStatement() Statement {
	stmt := &ReturnStatement{}
	stmt.initStatement(stmt)
	return stmt
}

func (stmt *ReturnStatement) parse() {
	if len(stmt.ts) > 0 {
		stmt.expr = extractExpression(stmt.ts)
	}
	if stmt.expr != nil {
		stmt.expr.setParent(stmt.getStack())
	}
}

func (stmt *ReturnStatement) execute() StatementResult {
	if stmt.expr == nil {
		return newStatementResult(StatementReturn, NULL)
	}
	res := stmt.expr.execute()
	return newStatementResult(StatementReturn, res)
}
