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
	if stmt.isReturnStatement() {
		executeStatementList(stmt.block, stack)
	}
	if stmt.isIfStatement() {
		executeIfStatement(stmt, stack)
	}
	if stmt.isForStatement() {
		executeForStatement(stmt, stack)
	}
	if stmt.isForeachStatement() || stmt.isForIndexStatement() || stmt.isForItemStatement() {
		executeForeachStatement(stmt, stack)
	}

	return nil
}

func executeExpression(expr *Expression, stack *VariableStack) (res *Value) {
	tmpVars := newVariables()
	exprExecutor := newExpressionExecutor(expr, stack, &tmpVars)
	return exprExecutor.run()
}

func evalBoolExpression(expr *Expression, stack *VariableStack) bool {
	val := executeExpression(expr, stack)
	if !val.isBooleanValue() {
		runtimeExcption(tokensString(expr.raw), " is not bool expression!")
	}
	return val.bool_value
}

func executeIfStatement(stmt *Statement, stack *VariableStack) (res *StatementResultType) {
	for _, condStmts := range stmt.condStmts {
		flag := evalBoolExpression(condStmts.condExpr, stack)
		if flag {
			executeStatementList(condStmts.block, stack)
			return
		}
	}

	if stmt.defStmt != nil {
		executeStatementList(stmt.defStmt.block, stack)
	}

	return
}

func executeForStatement(stmt *Statement, stack *VariableStack) (res *StatementResultType) {
	if stmt.preExpr != nil {
		executeExpression(stmt.preExpr, stack)
	}
	var flag bool
	if stmt.condExpr != nil {
		flag = evalBoolExpression(stmt.condExpr, stack)
	}
	for flag {

		executeStatementList(stmt.block, stack)

		if stmt.postExpr != nil {
			executeExpression(stmt.postExpr, stack)
		}
		if stmt.condExpr != nil {
			flag = evalBoolExpression(stmt.condExpr, stack)
		}
	}
	return
}

func executeForeachStatement(stmt *Statement, stack *VariableStack) (res *StatementResultType) {
	fpi := stmt.fpi
	varVal := stack.searchVariable(fpi.iterator)
	itr := toIterator(varVal)
	if itr == nil {
		runtimeExcption(fpi.iterator, "is not iterator!")
	}

	indexs := itr.indexs()
	for _, index := range indexs {

		if !stmt.isForItemStatement() {
			i := newVal(index)
			stack.addLocalVariable(fpi.indexName, i)
		}
		if !stmt.isForIndexStatement() {
			item := itr.getItem(index)
			stack.addLocalVariable(fpi.itemName, item)
		}

		executeStatementList(stmt.block, stack)
	}
	return
}
