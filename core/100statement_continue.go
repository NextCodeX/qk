package core

type ContinueStatement struct {
    StatementAdapter
}

func newContinueStatement() Statement {
    return &ContinueStatement{}
}

func (stmt *ContinueStatement) parse() {

}

func (stmt *ContinueStatement) execute() StatementResult {
    return newStatementResult(StatementContinue, NULL)
}