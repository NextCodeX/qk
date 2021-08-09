package core


type ObjectPrimaryExpression struct {
    obj JSONObject
    PrimaryExpressionImpl
}

func newObjectPrimaryExpression(val Value) PrimaryExpression {
    expr := &ObjectPrimaryExpression{}
    expr.t = ObjectPrimaryExpressionType
    expr.obj = goObj(val)
    expr.doExec = expr.doExecute
    return expr
}

func (priExpr *ObjectPrimaryExpression) getName() string {
    return "#jsonObject"
}

func (priExpr *ObjectPrimaryExpression) doExecute() Value {
    object := priExpr.obj
    if object.parsed() {
        return object
    }
    object.init()
    ts := clearBraces(object.tokens())
    size := len(ts)

    if size < 1 {
        return object
    }

    for i:=0; i<size; i++ {
        var nextCommaIndex int
        var exprTokens []Token
        if i+2 >= size {
            runtimeExcption("error jsonobject literal:", tokensString(ts))
        } else if ts[i+2].assertSymbol("[") {
            complexToken, endIndex := extractArrayLiteral(i+2, ts)
            nextCommaIndex = endIndex+1
            exprTokens = append(exprTokens, complexToken)
        } else if ts[i+2].assertSymbol("{") {
            complexToken, endIndex := extractObjectLiteral(i+2, ts)
            nextCommaIndex = endIndex+1
            exprTokens = append(exprTokens, complexToken)
        } else {
            nextCommaIndex = nextSymbolIndex(ts, i, ",")
            if nextCommaIndex < 0 {
                nextCommaIndex = size
            }
            exprTokens = ts[i+2:nextCommaIndex]
        }

        token := ts[i]
        keyname := token.raw()

        expr := extractExpression(exprTokens)
        expr.setStack(priExpr.getStack())
        val := expr.execute()
        if fn, ok := val.(Function); ok {
            fn.setPreVar("this", object)
        }
        object.put(keyname, val)
        i = nextCommaIndex
    }
    return object
}

