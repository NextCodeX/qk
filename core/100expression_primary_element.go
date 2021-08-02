package core

type ElementPrimaryExpression struct {
    arg  Expression // 函数调用参数 / 数组索引
    PrimaryExpressionImpl
}

func newElementPrimaryExpression(name string, arg Expression) PrimaryExpression {
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

func (priExpr *ElementPrimaryExpression) getName() string {
    return "#element"
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

func (priExpr *ElementPrimaryExpression) assignToObj(object JSONObject, res Value) {
    index := priExpr.arg.execute()
    object.put(goStr(index), res)
}

func (priExpr *ElementPrimaryExpression) assignToArr(array JSONArray, res Value) {
    index := priExpr.arg.execute()
    //fmt.Println("assignToArr -> arr is ", array.val(), index.val())
    array.set(toInt(index.val()), res)

}

