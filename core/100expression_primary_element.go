package core

type ElementPrimaryExpression struct {
	NestedPrimaryExpression
	PrimaryExpressionImpl
}

func newElementPrimaryExpression(arg Expression) PrimaryExpression {
	expr := &ElementPrimaryExpression{}
	expr.t = ElementPrimaryExpressionType
	expr.expr = arg
	expr.doExec = expr.doExecute
	return expr
}

func (priExpr *ElementPrimaryExpression) doExecute() Value {
	runtimeExcption("running ElementPrimaryExpression.doExecute is error")
	return nil
}

func (priExpr *ElementPrimaryExpression) getValue(obj JSONObject) Value {
	return obj.get(goStr(priExpr.expr.execute()))
}

func (priExpr *ElementPrimaryExpression) getArrElem(arr JSONArray) Value {
	return arr.getElem(toInt(priExpr.expr.execute().val()))
}

func (priExpr *ElementPrimaryExpression) getChar(str *StringValue) Value {
	index := priExpr.expr.execute()
	return newQKValue(str.getChar(toInt(index.val())))
}

func (priExpr *ElementPrimaryExpression) assignToObj(object JSONObject, res Value) {
	index := priExpr.expr.execute()
	object.put(goStr(index), res)
}

func (priExpr *ElementPrimaryExpression) assignToArr(array JSONArray, res Value) {
	index := priExpr.expr.execute()
	array.set(toInt(index.val()), res)

}
