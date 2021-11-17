package core

type ForStatement struct {
	preExprTokens  []Token    //
	condExprTokens []Token    //
	postExprTokens []Token    //
	preExpr        Expression //
	condExpr       Expression //
	postExpr       Expression //
	StatementAdapter
}

func newForStatement(pre, cond, post []Token) Statement {
	stmt := &ForStatement{preExprTokens: pre, condExprTokens: cond, postExprTokens: post}
	stmt.initStatement(stmt)
	return stmt
}

func (stmt *ForStatement) parse() {
	if stmt.preExprTokens != nil {
		stmt.preExpr = extractExpression(stmt.preExprTokens)
		stmt.preExpr.setParent(stmt.getStack())
	}
	if stmt.condExprTokens != nil {
		stmt.condExpr = extractExpression(stmt.condExprTokens)
		stmt.condExpr.setParent(stmt.getStack())
	}
	if stmt.postExprTokens != nil {
		stmt.postExpr = extractExpression(stmt.postExprTokens)
		stmt.postExpr.setParent(stmt.getStack())
	}
	Compile(stmt)
}

func (stmt *ForStatement) execute() StatementResult {
	if stmt.preExpr != nil {
		stmt.preExpr.execute()
	}
	var flag bool
	if stmt.condExpr != nil {
		flag = toBoolean(stmt.condExpr.execute())
	} else {
		flag = true
	}
	var res StatementResult
	for flag {
		res = stmt.executeStatementList(stmt.block, StmtListTypeFor)

		if res.isBreak() {
			res.setType(StatementNormal)
			return res
		} else if res.isReturn() {
			return res
		}

		if stmt.postExpr != nil {
			stmt.postExpr.execute()
		}
		if stmt.condExpr != nil {
			flag = toBoolean(stmt.condExpr.execute())
		}
	}
	// fix: empty loop exception
	if res == nil {
		res = newStatementResult(StatementNormal, NULL)
	}
	return res
}
