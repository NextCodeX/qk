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
		expr := extractExpression(ts)
		stmt.addExpression(expr)

	case stmt.isIfStatement():
		parseIfStatement(stmt)

	case stmt.isForStatement():
	case stmt.isSwitchStatement():
	case stmt.isReturnStatement():
	}
}

func parseIfStatement(stmt *Statement) {
	for _, condStmt := range stmt.condStmts {
		condStmt.condExpr = extractExpression(condStmt.condExprTokens)
		Compile(condStmt)
	}

	Compile(stmt.defStmt)
}