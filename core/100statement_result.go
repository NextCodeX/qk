package core


type StatementResultType int

const (
	StatementReturn StatementResultType = 1 << iota
	StatementContinue
	StatementBreak
	StatementNormal
)

type StatementResult struct {
	t StatementResultType
	val Value
}

func newStatementResult(t StatementResultType, val Value) *StatementResult {
	return &StatementResult{t, val}
}

func (sr *StatementResult) isReturn() bool {
	return (sr.t & StatementReturn) == StatementReturn
}

func (sr *StatementResult) isContinue() bool {
	return (sr.t & StatementContinue) == StatementContinue
}

func (sr *StatementResult) isBreak() bool {
	return (sr.t & StatementBreak) == StatementBreak
}

func (sr *StatementResult) isNormal() bool {
	return (sr.t & StatementNormal) == StatementNormal
}
