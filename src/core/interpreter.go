package core


func Interpret() {
    stack := newVariableStack()
    executeStatementList(stack)
}

func executeStatementList(stack *VariableStack) {
	stack.push()
    for _, stmt := range mainFunc.block {
        executeStatement(stmt, stack)
    }
    stack.pop()
}

func executeStatement(stmt *Statement, stack *VariableStack) *StatementResultType {
    if stmt.isExpressionStatement() {
        for _, expr := range stmt.exprs {
            executeExpression(expr, stack)
        }
    }

    return nil
}

func executeExpression(expr *Expression, stack *VariableStack) (res *Value) {
	tmpVars := newVariables()
	exprExecutor := newExpressionExecutor(expr, stack, &tmpVars)
	return exprExecutor.run()
}

