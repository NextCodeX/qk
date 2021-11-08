package core

import "bytes"

type ElemFCallToken struct {
	TokenAdapter
	header Token
	subs []Token
}

func newElemFCallToken(header Token, subs []Token) Token {
	t := &ElemFCallToken{}
	t.header = header
	t.subs = subs
	t.typName = "Elem & Func call"
	return t
}
func (t *ElemFCallToken) toExpr() PrimaryExpression {
	if nameToken, ok := t.header.(*NameToken); ok && len(t.subs) ==1 {
		if funcToken, ok := t.subs[0].(*FuncLiteralToken); ok {
			funcToken.name = nameToken.name()
			return funcToken.toExpr()
		}
	}

	var priExprs []PrimaryExpression
	for _, tk := range t.subs {
		priExpr := tk.toExpr()
		priExprs = append(priExprs, priExpr)
	}

	headExpr := t.header.toExpr()

	return newElemFunctionCallPrimaryExpression(headExpr, priExprs)
}

func (t *ElemFCallToken) String() string {
	var buf bytes.Buffer
	buf.WriteString(t.header.String())
	for _, token := range t.subs {
		buf.WriteString(token.String())
	}
	return buf.String()
}