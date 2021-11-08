package core

type ElemToken struct {
	TokenAdapter
	ts []Token
}

func newElemToken(ts []Token) Token {
	t := &ElemToken{}
	t.ts = ts
	t.typName = "Element"
	return t
}

func (t *ElemToken) toExpr() PrimaryExpression {
	expr := extractExpression(t.ts)
	return newElementPrimaryExpression(expr)
}

func (t *ElemToken) String() string {
	return "["+tokensString(t.ts)+"]"
}
