package core

type SelfIncrPrimaryExpression struct {
	NestedPrimaryExpression
	PrimaryExpressionImpl
}

func newSelfIncrPrimaryExpression(raw Expression) PrimaryExpression {
	expr := &SelfIncrPrimaryExpression{}
	expr.t = SelfIncrPrimaryExpressionType
	expr.expr = raw
	expr.doExec = expr.doExecute
	return expr
}

func (priExpr *SelfIncrPrimaryExpression) doExecute() Value {
	origin := priExpr.expr.execute()
	var tmp interface{}
	if i, ok := origin.(*IntValue); ok {
		tmp = i.goValue + 1
	} else if f, ok := origin.(*FloatValue); ok {
		tmp = f.goValue + 1
	} else {
		runtimeExcption("invalid SelfIncr operation: ", tokensString(priExpr.tokenList()), origin)
	}
	res := newQKValue(tmp)
	priExpr.evalAssign(priExpr.expr, res)
	return res
}
