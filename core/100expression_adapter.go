package core

type ValueReceiver interface {
	beAssigned(res Value)
}

type ExpressionAdapter struct {
	parent Function
	SourceCodeImpl
	ValueStack
}

func (ea *ExpressionAdapter) getVar(name string) Value {
	if isTmpVar(name) {
		tmpVars := ea.ValueStack.getVar(tmpVarsKey)
		return goObj(tmpVars).get(name)
	}
	return ea.ValueStack.getVar(name)
}

func (ea *ExpressionAdapter) setVar(name string, value Value) {
	if isTmpVar(name) {
		tmpVars := ea.ValueStack.getVar(tmpVarsKey)
		goObj(tmpVars).put(name, value)
	} else {
		ea.ValueStack.setVar(name, value)
	}
}

func (ea *ExpressionAdapter) setParent(p Function) {
	ea.ValueStack.cur = p
	ea.parent = p
}
func (ea *ExpressionAdapter) getParant() Function {
	return ea.parent
}

// 赋值
func (ea *ExpressionAdapter) evalAssign(priExpr Expression, res Value) {
	if expr, ok := priExpr.(ValueReceiver); ok {
		expr.beAssigned(res)
	} else {
		errorf("invalid assign expression: %v[%v] = %v", priExpr, tokensString(priExpr.tokenList()), res.val())
	}
}

func evalGoValues(exprs []Expression) []interface{} {
	var res []interface{}
	for _, expr := range exprs {
		if expr == nil {
			continue
		}
		val := expr.execute()
		res = append(res, val.val())
	}
	return res
}

func evalQKValues(exprs []Expression) []Value {
	var res []Value
	for _, expr := range exprs {
		if expr == nil {
			continue
		}
		val := expr.execute()
		res = append(res, val)
	}
	return res
}

func (ea *ExpressionAdapter) isPrimaryExpression() bool {
	return false
}
func (ea *ExpressionAdapter) isBinaryExpression() bool {
	return false
}
func (ea *ExpressionAdapter) isMultiExpression() bool {
	return false
}
