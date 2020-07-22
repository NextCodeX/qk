package core

import "fmt"

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
	fmt.Println("executeIfstmt")
	for _, condStmts := range stmt.condStmts {
		flag := executeExpression(condStmts.condExpr, stack)
		fmt.Println("flag:", flag.bool_value, flag.isBooleanValue())
		if flag.bool_value {
			executeStatementList(condStmts.block, stack)
			return
		}
	}

	fmt.Println("execute if stmt def:", stmt.defStmt != nil, stmt.defStmt.block)
	if stmt.defStmt != nil {
		executeStatementList(stmt.defStmt.block, stack)
	}

	return
}
