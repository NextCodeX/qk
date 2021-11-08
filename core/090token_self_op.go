package core


// 自增
type SelfIncrToken struct {
	TokenAdapter
	raw Token
}
func newSelfIncrToken(raw Token) Token {
	t := &SelfIncrToken{raw:raw}
	t.typName = "SelfIncr"
	return t
}
func (t *SelfIncrToken) toExpr() PrimaryExpression {
	return newSelfIncrPrimaryExpression(t.raw.toExpr())
}
func (t *SelfIncrToken) String() string {
	return t.raw.String() + "++"
}


// 自减
type SelfDecrToken struct {
	TokenAdapter
	raw Token
}
func newSelfDecrToken(raw Token) Token {
	t := &SelfDecrToken{raw:raw}
	t.typName = "SelfDecr"
	return t
}
func (t *SelfDecrToken) toExpr() PrimaryExpression {
	return newSelfDecrPrimaryExpression(t.raw.toExpr())
}
func (t *SelfDecrToken) String() string {
	return t.raw.String() + "--"
}