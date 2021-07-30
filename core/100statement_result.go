package core


type StatementResultType int

const (
	StatementReturn StatementResultType = 1 << iota
	StatementContinue
	StatementBreak
	StatementNormal
)

type StatementResult interface {
	isReturn() bool
	isContinue() bool
	isBreak() bool
	isNormal() bool
	setType(t StatementResultType)
	value() Value
}


type StatementResultImpl struct {
	t StatementResultType
	val Value
}

func newStatementResult(t StatementResultType, val Value) StatementResult {
	return &StatementResultImpl{t, val}
}

func (sr *StatementResultImpl) setType(t StatementResultType) {
	sr.t = t
}

func (sr *StatementResultImpl) value() Value {
	return sr.val
}

func (sr *StatementResultImpl) isReturn() bool {
	return (sr.t & StatementReturn) == StatementReturn
}

func (sr *StatementResultImpl) isContinue() bool {
	return (sr.t & StatementContinue) == StatementContinue
}

func (sr *StatementResultImpl) isBreak() bool {
	return (sr.t & StatementBreak) == StatementBreak
}

func (sr *StatementResultImpl) isNormal() bool {
	return (sr.t & StatementNormal) == StatementNormal
}
