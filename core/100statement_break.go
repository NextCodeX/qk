package core

type BreakStatement struct {
    StatementAdapter
}

func newBreakStatement() Statement {
    return &BreakStatement{}
}

func (stmt *BreakStatement) parse() {

}

func (stmt *BreakStatement) execute() StatementResult {
    return newStatementResult(StatementBreak, NULL)
}