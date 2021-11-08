package core


type NestedPrimaryExpression struct {
	expr Expression
	PrimaryExpressionImpl
}

func newNestedPrimaryExpression(raw Expression) PrimaryExpression {
	expr := &NestedPrimaryExpression{}
	expr.t = NestedPrimaryExpressionType
	expr.expr = raw
	expr.doExec = expr.doExecute
	return expr
}

func (priExpr *NestedPrimaryExpression) setStack(stack Function) {
	priExpr.stack = stack

	priExpr.expr.setStack(stack)
}

func (priExpr *NestedPrimaryExpression) doExecute() Value {
	return priExpr.expr.execute()
}