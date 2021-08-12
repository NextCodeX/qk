package core


type FunctionCallPrimaryExpression struct {
    args  []Expression // 函数调用参数
    PrimaryExpressionImpl
}

func newFunctionCallPrimaryExpression(args []Expression) PrimaryExpression {
    expr := &FunctionCallPrimaryExpression{}
    expr.t = FunctionCallPrimaryExpressionType
    expr.args = args
    expr.doExec = expr.doExecute
    return expr
}

func (priExpr *FunctionCallPrimaryExpression) setStack(stack Function) {
    priExpr.stack = stack

    for _, subExpr := range priExpr.args {
        subExpr.setStack(stack)
    }
}

func (priExpr *FunctionCallPrimaryExpression) doExecute() Value {
    runtimeExcption("running FunctionCallPrimaryExpression.doExecute is error")
    return nil
}

func (priExpr *FunctionCallPrimaryExpression) callFunc(fn Function) Value {
    if fn.isInternalFunc() {
        rawArgs := evalGoValues(priExpr.args)
        fn.setRawArgs(rawArgs)
    } else {
        args := evalQKValues(priExpr.args)
        fn.setArgs(args)
    }

    r := fn.execute()
    if r != nil {
       return r.value()
    } else {
        return NULL
    }
}

