package core


type ArrayPrimaryExpression struct {
    arr JSONArray
    PrimaryExpressionImpl
}

func newArrayPrimaryExpression(val Value) PrimaryExpression {
    expr := &ArrayPrimaryExpression{}
    expr.t = ArrayPrimaryExpressionType
    expr.arr = goArr(val)
    expr.doExec = expr.doExecute
    return expr
}

func (priExpr *ArrayPrimaryExpression) doExecute() Value {
    array := priExpr.arr
    if array.parsed() {
        return array
    }
    ts := clearBrackets(array.tokens())
    size := len(ts)

    if size < 1 {
        return array
    }

    for i:=0; i<size; i++ {
        var nextCommaIndex int
        var exprTokens []Token
        if ts[i].assertSymbol("[") {
            complexToken, endIndex := extractArrayLiteral(i, ts)
            nextCommaIndex = endIndex+1
            exprTokens = append(exprTokens, complexToken)
        } else if ts[i].assertSymbol("{") {
            complexToken, endIndex := extractObjectLiteral(i, ts)
            nextCommaIndex = endIndex+1
            exprTokens = append(exprTokens, complexToken)
        } else {
            nextCommaIndex = nextSymbolIndex(ts, i, ",")
            if nextCommaIndex < 0 {
                nextCommaIndex = size
            }
            exprTokens = ts[i:nextCommaIndex]
        }

        expr := extractExpression(exprTokens)
        expr.setStack(priExpr.getStack())
        val := expr.execute()
        array.add(val)
        i = nextCommaIndex
    }
    return array
}

