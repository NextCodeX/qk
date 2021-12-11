package core

type MultiAssignedPrimaryExpression struct {
	receivers []Expression
	valueList []Expression
	PrimaryExpressionImpl
}

func newMultiAssignedPrimaryExpression(receivers []Expression, valueList []Expression) PrimaryExpression {
	expr := &MultiAssignedPrimaryExpression{}
	expr.t = MultiAssignedPrimaryExpressionType
	expr.receivers = receivers
	expr.valueList = valueList
	expr.doExec = expr.doExecute
	return expr
}

func (priExpr *MultiAssignedPrimaryExpression) setParent(p Function) {
	priExpr.ExpressionAdapter.setParent(p)

	for _, item := range priExpr.receivers {
		if priExpr.localScope {
			item.setLocalScope()
		}
		item.setParent(p)
	}
	for _, item := range priExpr.valueList {
		item.setParent(p)
	}
}

func (priExpr *MultiAssignedPrimaryExpression) doExecute() Value {
	if priExpr.valueList == nil {
		for _, item := range priExpr.receivers {
			priExpr.evalAssign(item, NULL)
		}
		return NULL
	}

	for i, item := range priExpr.receivers {
		priExpr.evalAssign(item, priExpr.valueList[i].execute())
	}

	return NULL
}
