package core

type TernaryOperatorPrimaryExpression struct {
	receiver PrimaryExpression
	condExpr Expression
	ifExpr   Expression
	elseExpr Expression
	PrimaryExpressionImpl
}

func newTernaryOperatorPrimaryExpression(condExpr, ifExpr, elseExpr Expression, receiver PrimaryExpression) PrimaryExpression {
	expr := &TernaryOperatorPrimaryExpression{}
	expr.t = TernaryOperatorPrimaryExpressionType
	expr.condExpr = condExpr
	expr.ifExpr = ifExpr
	expr.elseExpr = elseExpr
	expr.receiver = receiver
	expr.doExec = expr.doExecute
	return expr
}

func (priExpr *TernaryOperatorPrimaryExpression) setParent(p Function) {
	priExpr.ExpressionAdapter.setParent(p)

	priExpr.condExpr.setParent(p)
	priExpr.ifExpr.setParent(p)
	priExpr.elseExpr.setParent(p)
	if priExpr.receiver != nil {
		priExpr.receiver.setParent(p)
	}
}

func (priExpr *TernaryOperatorPrimaryExpression) doExecute() Value {
	var res Value

	if toBoolean(priExpr.condExpr.execute()) {
		res = priExpr.ifExpr.execute()
	} else {
		res = priExpr.elseExpr.execute()
	}

	if priExpr.receiver != nil {
		priExpr.evalAssign(priExpr.receiver, res)
	}
	return res
}
