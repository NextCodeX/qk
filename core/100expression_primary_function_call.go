package core


type FunctionCallPrimaryExpression struct {
    name  string            // 变量名或者函数名称
    args  []Expression // 函数调用参数 / 数组索引
    PrimaryExpressionImpl
}

func newFunctionCallPrimaryExpression(name string, args []Expression) PrimaryExpression {
    expr := &FunctionCallPrimaryExpression{}
    expr.t = FunctionCallPrimaryExpressionType
    expr.name = name
    expr.args = args
    expr.doExec = expr.doExecute
    return expr
}

func (priExpr *FunctionCallPrimaryExpression) getName() string {
    return "#functionCall"
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

