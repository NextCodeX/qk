package core

type ValueReceiver interface {
	beAssigned(res Value)
}

type ExpressionAdapter struct {
	parent     Function
	localScope bool

	SourceCodeImpl
	ValueStack
}

func (this *ExpressionAdapter) setLocalScope() {
	this.localScope = true
}

func (this *ExpressionAdapter) getVar(name string) Value {
	if isTmpVar(name) {
		tmpVars := this.ValueStack.getVar(tmpVarsKey)
		return goObj(tmpVars).get(name)
	}
	return this.ValueStack.getVar(name)
}

func (this *ExpressionAdapter) setVar(name string, value Value) {
	if isTmpVar(name) {
		tmpVars := this.ValueStack.getVar(tmpVarsKey)
		goObj(tmpVars).put(name, value)
	} else {
		if this.localScope {
			this.ValueStack.setLocalVar(name, value)
			return
		}
		this.ValueStack.setVar(name, value)
	}
}

func (this *ExpressionAdapter) setParent(p Function) {
	this.ValueStack.cur = p
	this.parent = p
}
func (this *ExpressionAdapter) getParant() Function {
	return this.parent
}

// 赋值
func (this *ExpressionAdapter) evalAssign(priExpr Expression, res Value) {
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

func (this *ExpressionAdapter) isPrimaryExpression() bool {
	return false
}
func (this *ExpressionAdapter) isBinaryExpression() bool {
	return false
}
func (this *ExpressionAdapter) isMultiExpression() bool {
	return false
}
