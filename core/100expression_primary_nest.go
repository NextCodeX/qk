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

func (priExpr *NestedPrimaryExpression) setParent(p Function) {
	priExpr.ExpressionAdapter.setParent(p)

	priExpr.expr.setParent(p)
}

func (priExpr *NestedPrimaryExpression) doExecute() Value {
	return priExpr.expr.execute()
}
