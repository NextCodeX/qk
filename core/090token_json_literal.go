package core

// json array
type ArrLiteralToken struct {
	TokenAdapter
	ts []Token
}
func newArrLiteralToken(ts []Token) Token {
	t := &ArrLiteralToken{ts:ts}
	t.typName = "Array Literal"
	return t
}
func (t *ArrLiteralToken) toExpr() PrimaryExpression {
	return newArrayPrimaryExpression(t.ts)
}
func (t *ArrLiteralToken) String() string {
	return "["+tokensString(t.ts)+"]"
}

// json object
type ObjLiteralToken struct {
	TokenAdapter
	ts []Token
}
func newObjLiteralToken(ts []Token) Token {
	t := &ObjLiteralToken{ts:ts}
	t.typName = "Object Literal"
	return t
}
func (t *ObjLiteralToken) toExpr() PrimaryExpression {
	return newObjectPrimaryExpression(t.ts)
}
func (t *ObjLiteralToken) String() string {
	return "{"+tokensString(t.ts)+"}"
}
