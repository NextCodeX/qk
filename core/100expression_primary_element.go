package core

type ElementPrimaryExpression struct {
    arg  Expression
    PrimaryExpressionImpl
}

func newElementPrimaryExpression(arg Expression) PrimaryExpression {
    expr := &ElementPrimaryExpression{}
    expr.t = ElementPrimaryExpressionType
    expr.arg = arg
    expr.doExec = expr.doExecute
    return expr
}

func (priExpr *ElementPrimaryExpression) setStack(stack Function) {
    priExpr.stack = stack

    priExpr.arg.setStack(stack)
}

func (priExpr *ElementPrimaryExpression) doExecute() Value {
    runtimeExcption("running ElementPrimaryExpression.doExecute is error")
    return nil
}

func (priExpr *ElementPrimaryExpression) getValue(obj JSONObject) Value {
	return obj.get(goStr(priExpr.arg.execute()))
}

func (priExpr *ElementPrimaryExpression) getArrElem(arr JSONArray) Value {
    return arr.getElem(toInt(priExpr.arg.execute().val()))
}

func (priExpr *ElementPrimaryExpression) getChar(str *StringValue) Value {
    index := priExpr.arg.execute()
    return newQKValue(str.getChar(toInt(index.val())))
}

func (priExpr *ElementPrimaryExpression) assignToObj(object JSONObject, res Value) {
    index := priExpr.arg.execute()
    object.put(goStr(index), res)
}

func (priExpr *ElementPrimaryExpression) assignToArr(array JSONArray, res Value) {
    index := priExpr.arg.execute()
    array.set(toInt(index.val()), res)

}

