package core

func executeFunctionStatementList(stmts []*Statement, stack *VariableStack) {
	stack.push()
	executeStatementList(stmts, stack)
	stack.pop()
}

func executeStatementList(stmts []*Statement, stack *VariableStack) {
	for _, stmt := range stmts {
		executeStatement(stmt, stack)
	}
}

func executeStatement(stmt *Statement, stack *VariableStack) *StatementResultType {
	if stmt.isExpressionStatement() {
		for _, expr := range stmt.exprs {
			executeExpression(expr, stack)
		}
	}
	if stmt.isIfStatement() {
		executeIfStatement(stmt, stack)
	}

	return nil
}

func executeExpression(expr *Expression, stack *VariableStack) (res *Value) {
	tmpVars := newVariables()
	exprExecutor := newExpressionExecutor(expr, stack, &tmpVars)
	return exprExecutor.run()
}

func executeIfStatement(stmt *Statement, stack *VariableStack) (res *StatementResultType) {
	for _, condStmts := range stmt.condStmts {
		flag := executeExpression(condStmts.condExpr, stack)
		if flag.bool_value {
			executeStatementList(condStmts.block, stack)
			return
		}
	}

	if stmt.defStmt != nil {
		executeStatementList(stmt.defStmt.block, stack)
	}

	return
}
