package core


type TernaryOperatorPrimaryExpression struct {
    receiver PrimaryExpression
    condExpr Expression
    ifExpr Expression
    elseExpr Expression
    PrimaryExpressionImpl
}

func newTernaryOperatorPrimaryExpression(condExpr, ifExpr, elseExpr Expression, receiver PrimaryExpression) PrimaryExpression {
    expr := &TernaryOperatorPrimaryExpression{}
    expr.t = TernaryOperatorPrimaryExpressionType
    expr.condExpr = condExpr
    expr.ifExpr = ifExpr
    expr.elseExpr = elseExpr
    expr.receiver = receiver
    expr.doExec = expr.doExecute
    return expr
}

func (priExpr *TernaryOperatorPrimaryExpression) setStack(stack Function) {
    priExpr.stack = stack

    priExpr.condExpr.setStack(stack)
    priExpr.ifExpr.setStack(stack)
    priExpr.elseExpr.setStack(stack)
    if priExpr.receiver != nil {
        priExpr.receiver.setStack(stack)
    }
}

func (priExpr *TernaryOperatorPrimaryExpression) doExecute() Value {
    var res Value

    if toBoolean(priExpr.condExpr.execute()) {
        res = priExpr.ifExpr.execute()
    } else {
        res = priExpr.elseExpr.execute()
    }

    if priExpr.receiver != nil {
        priExpr.evalAssign(priExpr.receiver, res)
    }
    return res
}

