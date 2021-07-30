package core


type ExprPrimaryExpression struct {
    expr Expression
    PrimaryExpressionImpl
}

func newExprPrimaryExpression(subExpr Expression) PrimaryExpression {
    expr := &ExprPrimaryExpression{}
    expr.t = ExprPrimaryExpressionType
    expr.expr = subExpr
    expr.doExec = expr.doExecute
    return expr
}

func (priExpr *ExprPrimaryExpression) getName() string {
    return "#expr"
}

func (priExpr *ExprPrimaryExpression) setStack(stack Function) {
    priExpr.stack = stack

    priExpr.expr.setStack(stack)
}

func (priExpr *ExprPrimaryExpression) doExecute() Value {
    priExpr.expr.setStack(priExpr.getStack())
    return priExpr.expr.execute()
}

