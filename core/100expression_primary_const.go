package core


type ConstPrimaryExpression struct {
    val Value
    PrimaryExpressionImpl
}

func newConstPrimaryExpression(val Value) PrimaryExpression {
    expr := &ConstPrimaryExpression{}
    expr.t = ConstPrimaryExpressionType
    expr.val = val
    expr.doExec = expr.doExecute
    return expr
}

func (priExpr *ConstPrimaryExpression) doExecute() Value {
    return priExpr.val
}

