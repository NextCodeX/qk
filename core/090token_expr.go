package core

type ExprToken struct {
	TokenAdapter
	ts []Token
}

func newExprToken(ts []Token) Token {
	t := &ExprToken{ts: ts}
	t.typName = "Expr"
	return t
}

func (t *ExprToken) toExpr() PrimaryExpression {
	rawExpr := extractExpression(t.ts)
	return newNestedPrimaryExpression(rawExpr)
}
func (t *ExprToken) String() string {
	res := "(" + tokensShow10(t.ts) + ")"
	return res
}
