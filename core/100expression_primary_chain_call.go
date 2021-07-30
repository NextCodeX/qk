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

func (priExpr *ChainCallPrimaryExpression) getName() string {
    return "#chaincall"
}

func (priExpr *ChainCallPrimaryExpression) doExecute() Value {
    priExpr.head.setStack(priExpr.getStack())
    caller := priExpr.head.execute()

    //fmt.Println("eval chainCall#in: ", priExpr.head.left.isFunctionCall(), priExpr.isNot())
    //fmt.Println("eval chainCall#start: ", caller.val())

    for _, pri := range priExpr.chain {
        var intermediateResult Value
        if caller.isJsonArray() {
            if pri.isFunctionCall() {
                info := pri.(*FunctionCallPrimaryExpression)
                argRawVals := priExpr.toGoTypeValues(info.args)
                intermediateResult = evalJSONArrayMethod(goArr(caller), info.name, argRawVals)
            } else {}
        } else if caller.isJsonObject() {
            if pri.isVar() {
                info := pri.(*VarPrimaryExpression)
                intermediateResult = goObj(caller).get(info.varname)
            } else if pri.isFunctionCall() {
                info := pri.(*FunctionCallPrimaryExpression)
                argRawVals := priExpr.toGoTypeValues(info.args)
                intermediateResult = evalJSONObjectMethod(goObj(caller), info.name, argRawVals)
            } else if pri.isElement() {
                info := pri.(*ElementPrimaryExpression)
                val := goObj(caller).get(info.name)
                argRawVals := priExpr.toGoTypeValues(info.args)
                if val.isJsonArray() {
                    arr := goArr(val)
                    index := toIntValue(argRawVals[0])
                    intermediateResult = arr.get(index)
                } else if val.isJsonObject() {
                    obj := goObj(val)
                    key := toStringValue(obj)
                    intermediateResult = obj.get(key)
                } else { }
            }
        } else if caller.isString() {
            if pri.isFunctionCall() {
                info := pri.(*FunctionCallPrimaryExpression)
                argRawVals := priExpr.toGoTypeValues(info.args)
                intermediateResult = evalStringMethod(goStr(caller), info.name, argRawVals)
            } else {}
        } else if caller.isClass() {
            if pri.isFunctionCall() {
                info := pri.(*FunctionCallPrimaryExpression)
                argRawVals := priExpr.toGoTypeValues(info.args)
                intermediateResult = evalClassMethod(goAny(caller), info.name, argRawVals)
            } else if pri.isVar() {
                info := pri.(*VarPrimaryExpression)
                intermediateResult = evalClassField(goAny(caller), info.varname)
            } else {}
        } else {}

        if intermediateResult == nil {
            runtimeExcption("invalid chain call expression")
        } else {
            caller = intermediateResult
            //fmt.Println("eval chainCall#intermediate: ", intermediateResult.val())
        }
    }
    return caller
}

