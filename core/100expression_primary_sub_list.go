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

func (priExpr *SubListPrimaryExpression) getName() string {
    return "#subList"
}

func (priExpr *SubListPrimaryExpression) doExecute() Value {
    return nil
}

func (priExpr *SubListPrimaryExpression) startEndIndex(defaultEndIndex int) (int, int) {
    var start, end int
    if priExpr.start != nil {
        start = toInt(priExpr.start.execute().val())
    } else {
        start = 0
    }
    if priExpr.end != nil {
        end = toInt(priExpr.end.execute().val())
    } else {
        end = defaultEndIndex
    }
    return start, end
}

func (priExpr *SubListPrimaryExpression) subArr(arr JSONArray) Value {
    startIndex, endIndex := priExpr.startEndIndex(arr.size())
    return arr.sub(startIndex, endIndex)
}

func (priExpr *SubListPrimaryExpression) subStr(str *StringValue) Value {
    startIndex, endIndex := priExpr.startEndIndex(str.size())
    return newQKValue(str.sub(startIndex, endIndex))
}

