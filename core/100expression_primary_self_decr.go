package core

type SelfDecrPrimaryExpression struct {
	NestedPrimaryExpression
	PrimaryExpressionImpl
}

func newSelfDecrPrimaryExpression(raw Expression) PrimaryExpression {
	expr := &SelfDecrPrimaryExpression{}
	expr.t = SelfDecrPrimaryExpressionType
	expr.expr = raw
	expr.doExec = expr.doExecute
	return expr
}

func (priExpr *SelfDecrPrimaryExpression) doExecute() Value {
	origin := priExpr.expr.execute()
	var tmp interface{}
	if i, ok := origin.(*IntValue); ok {
		tmp = i.goValue - 1
	} else if f, ok := origin.(*FloatValue); ok {
		tmp = f.goValue - 1
	} else {
		runtimeExcption("invalid SelfDecr Operation: ", tokensString(priExpr.raw()), origin)
	}
	res := newQKValue(tmp)
	priExpr.evalAssign(priExpr.expr, res)
	return res
}
