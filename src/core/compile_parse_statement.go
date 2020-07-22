package core

import "fmt"

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
	fmt.Println("stmt.condStmts:", len(stmt.condStmts))
	for _, condStmt := range stmt.condStmts {
		fmt.Println("condStmt.condExprTokens:", tokensString(condStmt.condExprTokens))
		condStmt.condExpr = extractExpression(condStmt.condExprTokens)
		fmt.Println("condStmt", tokensString(condStmt.raw))
		Compile(condStmt)
	}

	fmt.Println("stmt.defStmt:", stmt.defStmt, )
	if stmt.defStmt==nil {
		return
	}
	fmt.Println("compile if def:", tokensString(stmt.defStmt.raw), len(stmt.defStmt.raw))
	Compile(stmt.defStmt)
}