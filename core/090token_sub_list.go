package core

type SubListToken struct {
	TokenAdapter
	start []Token
	end []Token
}
func newSubListToken(start, end []Token) Token {
	t := &SubListToken{}
	t.start = start
	t.end = end
	t.typName = "SubList"
	return t
}
func (t *SubListToken) toExpr() PrimaryExpression {
	start := extractExpression(t.start)
	end := extractExpression(t.end)
	return newSubListPrimaryExpression(start, end)
}

func (t *SubListToken) String() string {
	return "["+tokensString(t.start)+":"+tokensString(t.end)+"]"
}
