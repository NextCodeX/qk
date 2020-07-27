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
	val *Value
}

func newStatementResult(t StatementResultType, val *Value) *StatementResult {
	return &StatementResult{t, val}
}

func (this *StatementResult) isStatementReturn() bool {
	return (this.t & StatementReturn) == StatementReturn
}

func (this *StatementResult) isStatementContinue() bool {
	return (this.t & StatementContinue) == StatementContinue
}

func (this *StatementResult) isStatementBreak() bool {
	return (this.t & StatementBreak) == StatementBreak
}

func (this *StatementResult) isStatementNormal() bool {
	return (this.t & StatementNormal) == StatementNormal
}
