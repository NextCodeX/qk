package core


type ArrayPrimaryExpression struct {
    arrTokens []Token
    PrimaryExpressionImpl
}

func newArrayPrimaryExpression(ts []Token) PrimaryExpression {
    expr := &ArrayPrimaryExpression{}
    expr.t = ArrayPrimaryExpressionType
    expr.arrTokens = ts
    expr.doExec = expr.doExecute
    return expr
}

func (priExpr *ArrayPrimaryExpression) doExecute() Value {
    array := emptyArray()

    ts := clearBrackets(priExpr.arrTokens)
    size := len(ts)

    if size < 1 {
        return array
    }

    for i:=0; i<size; i++ {
        var nextCommaIndex int
        var exprTokens []Token
        if ts[i].assertSymbol("[") {
            complexToken, endIndex := extractArrLiteral(ts, i)
            nextCommaIndex = endIndex+1
            exprTokens = append(exprTokens, complexToken)
        } else if ts[i].assertSymbol("{") {
            complexToken, endIndex := extractObjLiteral(ts, i)
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

