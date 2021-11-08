package core


type NotPrimaryExpression struct {
	toBool bool
	NestedPrimaryExpression
	PrimaryExpressionImpl
}

func newNotPrimaryExpression(toBool bool, raw Expression) PrimaryExpression {
	expr := &NotPrimaryExpression{}
	expr.t = NotPrimaryExpressionType
	expr.toBool = toBool
	expr.expr = raw
	expr.doExec = expr.doExecute
	return expr
}

func (priExpr *NotPrimaryExpression) doExecute() Value {
	flag := toBoolean(priExpr.expr.execute())
	if !priExpr.toBool {
		flag = !flag
	}
	return newQKValue(flag)
}