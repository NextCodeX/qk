package core

type ChainCallPrimaryExpression struct {
    chain []PrimaryExpression
    head  PrimaryExpression // ChainCall 的头表达式
    PrimaryExpressionImpl
}

func newChainCallPrimaryExpression(headExpr PrimaryExpression,  priExprs []PrimaryExpression) PrimaryExpression {
    expr := &ChainCallPrimaryExpression{}
    expr.t = ChainCallPrimaryExpressionType
    expr.head = headExpr
    expr.chain = priExprs
    expr.doExec = expr.doExecute
    return expr
}

func (priExpr *ChainCallPrimaryExpression) setStack(stack Function) {
    priExpr.stack = stack

    priExpr.head.setStack(stack)
    for _, subExpr := range priExpr.chain {
        subExpr.setStack(stack)
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
                errorf("%v.%v is error", caller.val(), tokensString(pri.raw()))
            }
        } else {
            runtimeExcption("invalid chain call expression:", tokensString(priExpr.raw()))
        }

        if intermediateResult == nil {
            runtimeExcption("invalid chain call expression:", tokensString(priExpr.raw()))
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
        errorf("(in ChainCall)invalid assign expression: %v = %v", tokensString(priExpr.raw()), res.val())
    }
}

