package core

type FunctionPrimaryExpression struct {
    val Function
    PrimaryExpressionImpl
}

func newFunctionPrimaryExpression(val Function) PrimaryExpression {
    expr := &FunctionPrimaryExpression{}
    expr.t = FunctionPrimaryExpressionType
    expr.val = val
    expr.doExec = expr.doExecute
    return expr
}

func (priExpr *FunctionPrimaryExpression) setStack(stack Function) {
    priExpr.stack = stack

    priExpr.val.setParent(stack)
}

func (priExpr *FunctionPrimaryExpression) getName() string {
    return "#function"
}

func (priExpr *FunctionPrimaryExpression) doExecute() Value {
    fn := priExpr.val
    funcName := fn.getName()
    if funcName != "" {
        priExpr.setVar(funcName, fn)
    }

    return fn
}

