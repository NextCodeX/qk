package core

import "bytes"

// 函数调用
type FunctionCallToken struct {
	TokenAdapter
	pts []Token  // 调用参数Token列表
}
func newFunctionCallToken(ts []Token) Token {
	t :=  &FunctionCallToken{pts: ts}
	t.typName = "Func Call"
	return t
}
func (t *FunctionCallToken) toExpr() PrimaryExpression {
	exprs := getArgExprsFromToken(t.pts)
	return newFunctionCallPrimaryExpression(exprs)
}
func (t *FunctionCallToken) String() string {
	return "("+tokensString(t.pts)+")"
}

// 函数字面值
type FuncLiteralToken struct {
	TokenAdapter
	name string  // 函数名称
	pts []Token  // 形参Token列表
	bodyTokens []Token
}
func newFuncLiteralToken(name string, pts, bodyTokens []Token) Token {
	t :=  &FuncLiteralToken{
		name:name,
		pts:pts,
		bodyTokens: bodyTokens}
	t.typName = "Func Literal"
	return t
}
func (t *FuncLiteralToken) toExpr() PrimaryExpression {
	paramNames := getFuncParamNames(t.pts)
	return newFunctionPrimaryExpression(t.name, paramNames, t.bodyTokens)
}
func (t *FuncLiteralToken) String() string {
	var buf bytes.Buffer
	buf.WriteString(t.name)
	buf.WriteString("(")
	buf.WriteString(tokensString(t.pts))
	buf.WriteString(")")
	buf.WriteString("{")
	for _, token := range t.bodyTokens {
		buf.WriteString(token.String())
	}
	buf.WriteString("}")
	return buf.String()
}