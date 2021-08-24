package core


type SubListPrimaryExpression struct {
    start Expression
    end Expression
    PrimaryExpressionImpl
}

func newSubListPrimaryExpression(start, end Expression) PrimaryExpression {
    expr := &SubListPrimaryExpression{}
    expr.t = SubListPrimaryExpressionType
    expr.start = start
    expr.end = end
    expr.doExec = expr.doExecute
    return expr
}

func (priExpr *SubListPrimaryExpression) setStack(stack Function) {
    priExpr.stack = stack

    if priExpr.start != nil {
        priExpr.start.setStack(stack)
    }
    if priExpr.end != nil {
        priExpr.end.setStack(stack)
    }
}

func (priExpr *SubListPrimaryExpression) doExecute() Value {
    return nil
}

func (priExpr *SubListPrimaryExpression) startEndIndex(defaultEndIndex int) (int, int) {
    var start, end int
    if priExpr.start != nil {
        start = toInt(priExpr.start.execute().val())
        if start < 0 {
            // compatible
           start = 0
        }
    } else {
        start = 0
    }
    if priExpr.end != nil {
        end = toInt(priExpr.end.execute().val())
        if end > defaultEndIndex {
            // compatible
            end = defaultEndIndex
        }
    } else {
        end = defaultEndIndex
    }
    return start, end
}

func (priExpr *SubListPrimaryExpression) subArr(arr JSONArray) Value {
    startIndex, endIndex := priExpr.startEndIndex(arr.Size())
    return arr.sub(startIndex, endIndex)
}

func (priExpr *SubListPrimaryExpression) subByteArray(arr *ByteArrayValue) Value {
    startIndex, endIndex := priExpr.startEndIndex(arr.Size())
    return newQKValue(arr.sub(startIndex, endIndex))
}

func (priExpr *SubListPrimaryExpression) subStr(str *StringValue) Value {
    startIndex, endIndex := priExpr.startEndIndex(str.Size())
    return newQKValue(str.sub(startIndex, endIndex))
}

