package core

type StmtListType int
const (
	StmtListTypeFunc StmtListType = 1 << iota
	StmtListTypeIf
	StmtListTypeFor
	StmtListTypeNormal
)

func executeFunctionStatementList(stmts []*Statement, stack *VariableStack) *Value {
	defer stack.pop()

	executeStatementList(stmts, stack, StmtListTypeFunc)
	res := stack.searchVariable(funcResultName)
	return res
}

func executeStatementList(stmts []*Statement, stack *VariableStack, t StmtListType) *StatementResult {
	if t == StmtListTypeFunc {
		stack.addLocalVariable(funcResultName, NULL)
	}
	var res *StatementResult
	for _, stmt := range stmts {
		res = executeStatement(stmt, stack)

		if res.isContinue() {
			if t == StmtListTypeFor {
				res.t = StatementNormal
			}
			break
		} else if res.isReturn() || res.isBreak() {
			break
		}

	}
	return res
}

func executeStatement(stmt *Statement, stack *VariableStack) *StatementResult {
	var res *StatementResult
	if stmt.isExpressionStatement() {
		exprResult := executeExpression(stmt.expr, stack)
		res = newStatementResult(StatementNormal, exprResult)

	} else if stmt.isContinueStatement() {
		res = newStatementResult(StatementContinue, NULL)

	} else if stmt.isBreakStatement() {
		res = newStatementResult(StatementBreak, NULL)

	} else if stmt.isReturnStatement() {
		executeStatementList(stmt.block, stack, StmtListTypeNormal)
		funcResult := stack.searchVariable(funcResultName)
		res = newStatementResult(StatementReturn, funcResult)

	} else if stmt.isIfStatement() {
		res = executeIfStatement(stmt, stack)

	} else if stmt.isForStatement() {
		res = executeForStatement(stmt, stack)

	} else if stmt.isForeachStatement() || stmt.isForIndexStatement() || stmt.isForItemStatement() {
		res = executeForeachStatement(stmt, stack)

	} else {
		runtimeExcption("unknow statememnt:", tokensString(stmt.raw))
	}

	return res
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
	return val.boolean
}

func executeIfStatement(stmt *Statement, stack *VariableStack) (res *StatementResult) {
	for _, condStmts := range stmt.condStmts {
		flag := evalBoolExpression(condStmts.condExpr, stack)
		if flag {
			res = executeStatementList(condStmts.block, stack, StmtListTypeIf)
			return
		}
	}

	if stmt.defStmt != nil {
		return executeStatementList(stmt.defStmt.block, stack, StmtListTypeIf)
	}

	return newStatementResult(StatementNormal, NULL)
}

func executeForStatement(stmt *Statement, stack *VariableStack) (res *StatementResult) {
	if stmt.preExpr != nil {
		executeExpression(stmt.preExpr, stack)
	}
	var flag bool
	if stmt.condExpr != nil {
		flag = evalBoolExpression(stmt.condExpr, stack)
	}
	//for _, st := range stmt.block {
	//	fmt.Println("stmts: ", tokensString(st.raw))
	//}
	for flag {

		res = executeStatementList(stmt.block, stack, StmtListTypeFor)


		if res.isBreak() {
			res.t = StatementNormal
			return
		}else if res.isReturn() {
			return
		}

		if stmt.postExpr != nil {
			executeExpression(stmt.postExpr, stack)
		}
		if stmt.condExpr != nil {
			flag = evalBoolExpression(stmt.condExpr, stack)
		}
	}
	return res
}

func executeForeachStatement(stmt *Statement, stack *VariableStack) (res *StatementResult) {
	fpi := stmt.fpi
	varVal := stack.searchVariable(fpi.iterator)
	itr := toIterator(varVal)
	if itr == nil {
		runtimeExcption(fpi.iterator, "is not iterator!")
	}

	indexs := itr.indexs()
	for _, index := range indexs {

		if !stmt.isForItemStatement() {
			i := newQkValue(index)
			stack.addLocalVariable(fpi.indexName, i)
		}
		if !stmt.isForIndexStatement() {
			item := itr.getItem(index)
			stack.addLocalVariable(fpi.itemName, item)
		}

		res = executeStatementList(stmt.block, stack, StmtListTypeFor)

		if res.isBreak() {
			res.t = StatementNormal
			return
		}else if res.isReturn() {
			return
		}
	}
	return res
}
