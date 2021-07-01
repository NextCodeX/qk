package core

func parseStatementList(stmts []*Statement) {
	for _, stmt := range stmts {
		parseStatement(stmt)
	}
}

func parseStatement(stmt *Statement) {
	ts := stmt.raw
	switch {
	case stmt.isExpressionStatement():
		stmt.expr = extractExpression(ts)

	case stmt.isIfStatement():
		parseIfStatement(stmt)

	case stmt.isForStatement():
		parseForStatement(stmt)

	case stmt.isForeachStatement() || stmt.isForIndexStatement() || stmt.isForItemStatement():
		parseForPlusStatement(stmt)

	case stmt.isSwitchStatement():
	case stmt.isReturnStatement():
		parseReturnStatement(stmt)
	}
}

func parseIfStatement(stmt *Statement) {
	for _, condStmt := range stmt.condStmts {
		condStmt.condExpr = extractExpression(condStmt.condExprTokens)
		Compile(condStmt)
	}

	if stmt.defStmt==nil {
		return
	}
	Compile(stmt.defStmt)
}

func parseForStatement(stmt *Statement) {
	if stmt.preExprTokens != nil {
		stmt.preExpr = extractExpression(stmt.preExprTokens)
	}
	if stmt.condExprTokens != nil {
		stmt.condExpr = extractExpression(stmt.condExprTokens)
	}
	if stmt.postExprTokens != nil {
		stmt.postExpr = extractExpression(stmt.postExprTokens)
	}
	Compile(stmt)
}

func parseForPlusStatement(stmt *Statement)  {
	Compile(stmt)
}

func parseReturnStatement(stmt *Statement)  {
	if len(stmt.raw) < 1 {
		stmt.raw = insert(newToken("0", Int), stmt.raw)
	}
	stmt.raw = insert(symbolToken("="), stmt.raw)
	stmt.raw = insert(varToken(funcResultName), stmt.raw)

	Compile(stmt)
}