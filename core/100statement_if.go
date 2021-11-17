package core

type IfStatement struct {
	condExpr  Expression
	condStmts []Statement
	defStmt   Statement
	StatementAdapter
}

func newMultiIfStatement(condStmts []Statement, defStmt Statement) Statement {
	stmt := &IfStatement{}
	stmt.condStmts = condStmts
	stmt.defStmt = defStmt
	stmt.initStatement(stmt)
	return stmt
}

func newSingleIfStatement(condTokens, bodyTokens []Token) Statement {
	stmt := &IfStatement{}
	stmt.condExpr = extractExpression(condTokens)
	stmt.ts = bodyTokens
	stmt.initStatement(stmt)
	return stmt
}

func (stmt *IfStatement) parse() {
	for _, condStmt := range stmt.condStmts {
		ifStmt := condStmt.(*IfStatement)
		ifStmt.condExpr.setParent(stmt.getStack())

		condStmt.setParent(stmt.getStack())
		Compile(condStmt)
	}

	if stmt.defStmt == nil {
		return
	}
	stmt.defStmt.setParent(stmt.getStack())
	Compile(stmt.defStmt)
}

func (stmt *IfStatement) execute() StatementResult {
	for _, condStmt := range stmt.condStmts {
		ifstmt := condStmt.(*IfStatement)
		flag := toBoolean(ifstmt.condExpr.execute())
		if flag {
			return stmt.executeStatementList(condStmt.stmts(), StmtListTypeIf)
		}
	}

	if stmt.defStmt != nil {
		return stmt.executeStatementList(stmt.defStmt.stmts(), StmtListTypeIf)
	}

	return newStatementResult(StatementNormal, NULL)
}
