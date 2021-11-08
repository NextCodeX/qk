package core

import "bytes"

type ChainCallToken struct {
	TokenAdapter
	header Token
	subs []Token
}

func newChainCallToken(header Token, subs []Token) Token {
	t := &ChainCallToken{}
	t.header = header
	t.subs = subs
	t.typName = "Chain Call"
	return t
}
func (t *ChainCallToken) toExpr() PrimaryExpression {
	var priExprs []PrimaryExpression
	for _, tk := range t.subs {
		priExpr := tk.toExpr()
		priExprs = append(priExprs, priExpr)
	}

	headExpr := t.header.toExpr()

	return newChainCallPrimaryExpression(headExpr, priExprs)
}
func (t *ChainCallToken) String() string {
	var buf bytes.Buffer
	buf.WriteString(t.header.String())
	for _, token := range t.subs {
		buf.WriteString(".")
		buf.WriteString(token.String())
	}
	return buf.String()
}
