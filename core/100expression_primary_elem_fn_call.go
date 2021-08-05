package core

type ElemFunctionCallPrimaryExpression struct {
    chain []PrimaryExpression
    head  PrimaryExpression
    PrimaryExpressionImpl
}

func newElemFunctionCallPrimaryExpression(headExpr PrimaryExpression,  priExprs []PrimaryExpression) PrimaryExpression {
    expr := &ElemFunctionCallPrimaryExpression{}
    expr.t = ElemFunctionCallPrimaryExpressionType
    expr.head = headExpr
    expr.chain = priExprs
    expr.doExec = expr.doExecute
    return expr
}

func (priExpr *ElemFunctionCallPrimaryExpression) setStack(stack Function) {
    priExpr.stack = stack

    priExpr.head.setStack(stack)
    for _, subExpr := range priExpr.chain {
        subExpr.setStack(stack)
    }
}

func (priExpr *ElemFunctionCallPrimaryExpression) getName() string {
    return "#elementFunctionCallMixture"
}

func (priExpr *ElemFunctionCallPrimaryExpression) doExecute() Value {
    return priExpr.exprExec(priExpr.chain)
}

func (priExpr *ElemFunctionCallPrimaryExpression) exprExec(chainExprs []PrimaryExpression) Value {
    currentObj := priExpr.head.execute()
    chain := chainExprs
    if priExpr.head.isVar() {
        chainLen := len(chainExprs)
        if  chainLen > 0 && chainExprs[0].isFunction() {
            subExpr := chainExprs[0].(*FunctionPrimaryExpression)
            fn := subExpr.execute().(Function)

            headExpr := priExpr.head.(*VarPrimaryExpression)
            headExpr.beAssigned(fn)

            if chainLen == 1 {
                return nil
            } else {
                currentObj = fn
                chain = chainExprs[1:]
            }
        }
    }

    return priExpr.runScopeChain(currentObj, chain)
}



func (priExpr *ElemFunctionCallPrimaryExpression) runScopeChain(currentObj Value, chain []PrimaryExpression) Value {
    for _, subExpr := range chain {

        if currentObj == nil {
            errorf("invalid expression: null%v", tokensString(subExpr.raw()))
            break
        }

        var intermediateVal Value
        if currentObj.isJsonArray() && subExpr.isSubList() {
            // arr[:]
            arr := goArr(currentObj)
            nextExpr := subExpr.(*SubListPrimaryExpression)
            intermediateVal = nextExpr.subArr(arr)

        } else if currentObj.isString() && subExpr.isSubList() {
            // str[:]
            str := currentObj.(*StringValue)
            nextExpr := subExpr.(*SubListPrimaryExpression)
            intermediateVal = nextExpr.subStr(str)

        } else if currentObj.isString() && subExpr.isElement() {
            // str[]
            str := currentObj.(*StringValue)
            nextExpr := subExpr.(*ElementPrimaryExpression)
            intermediateVal = nextExpr.getChar(str)

        } else if currentObj.isJsonArray() && subExpr.isElement() {
            // arr[]
            arr := goArr(currentObj)
            nextExpr := subExpr.(*ElementPrimaryExpression)
            intermediateVal =  nextExpr.getArrElem(arr)

        } else if currentObj.isJsonObject() && subExpr.isElement() {
            // object[key]
            obj := goObj(currentObj)
            nextExpr := subExpr.(*ElementPrimaryExpression)
            intermediateVal =  nextExpr.getValue(obj)

        } else if currentObj.isFunction() && subExpr.isFunctionCall() {
            // func()
            fn := currentObj.(Function)
            nextExpr := subExpr.(*FunctionCallPrimaryExpression)
            intermediateVal = nextExpr.callFunc(fn)

        } else {
            errorf("invalid mixture expression: %v{%v} chainLen: %v", currentObj.val(), tokensString(subExpr.raw()), len(chain))
        }

        if intermediateVal != nil  {
          currentObj = intermediateVal
        } else {
          runtimeExcption("failed to run Element and Function Call Mixture: %v%v", currentObj.val(), tokensString(subExpr.raw()))
        }
    }
    //fmt.Println("mix#runScopeChain result -> ", currentObj.val())
    //fmt.Println("*****************************")
    return currentObj
}

// 与Chain Call表达式结合, 连调
func (priExpr *ElemFunctionCallPrimaryExpression) runWith(obj Object) Value {
	varExpr, ok := priExpr.head.(*VarPrimaryExpression)
	if !ok {
	    errorf("object.%v is error", tokensString(priExpr.raw()))
    }
	currentObj := varExpr.getAttribute(obj)
	return priExpr.runScopeChain(currentObj, priExpr.chain)
}
// 与Chain Call表达式结合, 被赋值
func (priExpr *ElemFunctionCallPrimaryExpression) beAssignedAfterChainCall(obj Object, res Value) {
	varExpr, ok := priExpr.head.(*VarPrimaryExpression)
	if !ok {
	    errorf("object.%v is error", tokensString(priExpr.raw()))
    }
	currentObj := varExpr.getAttribute(obj)
    chainExprs := priExpr.chain
    chainLen := len(chainExprs)
    finalVal := priExpr.runScopeChain(currentObj, chainExprs[:chainLen-1])
    tailExpr := chainExprs[chainLen-1]
    priExpr.set(finalVal, tailExpr, res)
}

// 被赋值
func (priExpr *ElemFunctionCallPrimaryExpression) beAssigned(res Value) {
	chainExprs := priExpr.chain
	chainLen := len(chainExprs)
    obj := priExpr.exprExec(chainExprs[:chainLen-1])
    tailExpr := chainExprs[chainLen-1]
    priExpr.set(obj, tailExpr, res)
}

func (priExpr *ElemFunctionCallPrimaryExpression) set(obj Value, tailExpr PrimaryExpression, res Value) {
    //fmt.Println("mix#beAssigned#typeAsssert:", obj.val(), obj.isJsonObject(), obj.isJsonArray(), " | ", tailExpr.isVar(), tailExpr.isElement())
    if obj.isJsonObject() && tailExpr.isElement() {
        jsonObject := obj.(JSONObject)
        elemExpr := tailExpr.(*ElementPrimaryExpression)
        elemExpr.assignToObj(jsonObject, res)

    } else if obj.isJsonArray() && tailExpr.isElement() {
        jsonArray := obj.(JSONArray)
        elemExpr := tailExpr.(*ElementPrimaryExpression)
        elemExpr.assignToArr(jsonArray, res)

    } else {
        errorf("(in elem&fn call v1)invalid assign expression: %v = %v", tokensString(priExpr.raw()), res.val())
    }
}

