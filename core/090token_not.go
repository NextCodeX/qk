package core


type NotToken struct {
	TokenAdapter
	toBool bool
	ts []Token
}

func newNotToken(toBoolFlag bool, ts ...Token) Token {
	t := &NotToken{ts:ts}
	t.toBool = toBoolFlag
	t.typName = "Not"
	return t
}

func (t *NotToken) toExpr() PrimaryExpression {
	rawExpr := extractExpression(t.ts)
	return newNotPrimaryExpression(t.toBool, rawExpr)
}
func (t *NotToken) String() string {
	res := "!"+ tokensShow10(t.ts)
	if t.toBool {
		res = "!" + res
	}
	return res
}
