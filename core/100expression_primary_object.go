package core

type ObjectPrimaryExpression struct {
	objTokens []Token
	PrimaryExpressionImpl
}

func newObjectPrimaryExpression(ts []Token) PrimaryExpression {
	expr := &ObjectPrimaryExpression{}
	expr.t = ObjectPrimaryExpressionType
	expr.objTokens = ts
	expr.doExec = expr.doExecute
	return expr
}

func (priExpr *ObjectPrimaryExpression) doExecute() Value {
	object := emptyJsonObject()

	ts := clearBraces(priExpr.objTokens)
	size := len(ts)

	if size < 1 {
		return object
	}

	for i := 0; i < size; i++ {
		var nextCommaIndex int
		var exprTokens []Token
		if i+2 >= size {
			runtimeExcption("error jsonobject literal:", tokensString(ts))
		} else if ts[i+2].assertSymbol("[") {
			complexToken, endIndex := extractArrLiteral(ts, i+2)
			nextCommaIndex = endIndex + 1
			exprTokens = append(exprTokens, complexToken)
		} else if ts[i+2].assertSymbol("{") {
			complexToken, endIndex := extractObjLiteral(ts, i+2)
			nextCommaIndex = endIndex + 1
			exprTokens = append(exprTokens, complexToken)
		} else {
			nextCommaIndex = nextSymbolIndex(ts, i, ",")
			if nextCommaIndex < 0 {
				nextCommaIndex = size
			}
			exprTokens = ts[i+2 : nextCommaIndex]
		}

		token := ts[i]
		keyname := token.String()

		expr := extractExpression(exprTokens)
		expr.setStack(priExpr.getStack())
		val := expr.execute()
		object.put(keyname, val)

		i = nextCommaIndex
	}
	return object
}
