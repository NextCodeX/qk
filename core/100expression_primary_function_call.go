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
    functionName := priExpr.name
    args := priExpr.args

    customFunc, ok := funcList[functionName]
    if ok {
        argVals := priExpr.evalValues(args)
        customFunc.setArgs(argVals)
        res :=  customFunc.execute()
        if res != nil {
            return res.value()
        }
    } else if isPrint(functionName) {
        argRawVals := priExpr.toGoTypeValues(args)
        executePrintFunc(functionName, argRawVals)
    } else if isModuleFunc(functionName) {
        argRawVals := priExpr.toGoTypeValues(args)
        return executeModuleFunc(functionName, argRawVals)
    } else {
        errorf("function %v() is not defined!", functionName)
    }
    return nil
}

