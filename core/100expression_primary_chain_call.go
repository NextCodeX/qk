package core

type ChainCallPrimaryExpression struct {
	chain []PrimaryExpression
	head  PrimaryExpression // ChainCall 的头表达式
	PrimaryExpressionImpl
}

func newChainCallPrimaryExpression(headExpr PrimaryExpression, priExprs []PrimaryExpression) PrimaryExpression {
	expr := &ChainCallPrimaryExpression{}
	expr.t = ChainCallPrimaryExpressionType
	expr.head = headExpr
	expr.chain = priExprs
	expr.doExec = expr.doExecute
	return expr
}

func (priExpr *ChainCallPrimaryExpression) setParent(p Function) {
	priExpr.ExpressionAdapter.setParent(p)

	priExpr.head.setParent(p)
	for _, subExpr := range priExpr.chain {
		subExpr.setParent(p)
	}
}

func (priExpr *ChainCallPrimaryExpression) doExecute() Value {
	return priExpr.exprExec(priExpr.chain)
}

func (priExpr *ChainCallPrimaryExpression) exprExec(chainExprs []PrimaryExpression) Value {
	caller := priExpr.head.execute()

	for _, pri := range chainExprs {
		var intermediateResult Value
		if caller.isObject() {
			obj := caller.(Object)
			if pri.isElemFunctionCall() {
				// object.method() / object.arr[]
				nextExpr := pri.(*ElemFunctionCallPrimaryExpression)
				intermediateResult = nextExpr.runWith(obj)

			} else if pri.isVar() {
				// object.attribute
				nextExpr := pri.(*VarPrimaryExpression)
				intermediateResult = nextExpr.getField(obj)

			} else {
				errorf("%v.%v is error", caller.val(), tokensString(pri.tokenList()))
			}
		} else {
			runtimeExcption("invalid chain call expression:", tokensString(priExpr.tokenList()))
		}

		if intermediateResult == nil {
			runtimeExcption("invalid chain call expression:", tokensString(priExpr.tokenList()))
		} else {
			caller = intermediateResult
		}
	}
	return caller
}

func (priExpr *ChainCallPrimaryExpression) beAssigned(res Value) {
	chainExprs := priExpr.chain
	chainLen := len(chainExprs)
	obj := priExpr.exprExec(chainExprs[:chainLen-1])

	tailExpr := chainExprs[chainLen-1]
	if obj.isJsonObject() && tailExpr.isVar() {
		jsonObject := obj.(JSONObject)
		varExpr := tailExpr.(*VarPrimaryExpression)
		varExpr.assign(jsonObject, res)

	} else if obj.isJsonObject() && tailExpr.isElemFunctionCall() {
		clazz := obj.(Object)
		subExpr := tailExpr.(*ElemFunctionCallPrimaryExpression)
		subExpr.beAssignedAfterChainCall(clazz, res)

	} else {
		errorf("(in ChainCall)invalid assign expression: %v = %v", tokensString(priExpr.tokenList()), res.val())
	}
}
