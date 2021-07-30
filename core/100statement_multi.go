package core

type MultiStatement struct {
    StatementAdapter
}

func newMultiStatement(ts []Token) Statement {
    stmt := &MultiStatement{}
    stmt.ts = ts
    stmt.initStatement(stmt)
    return stmt
}

func (stmt *MultiStatement) parse() {
    Compile(stmt)
}

func (stmt *MultiStatement) execute() StatementResult {
    return stmt.executeStatementList(stmt.block, StmtListTypeNormal)
}