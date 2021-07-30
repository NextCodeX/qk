package core

type ElementPrimaryExpression struct {
    name  string            // 变量名或者函数名称
    args  []Expression // 函数调用参数 / 数组索引
    PrimaryExpressionImpl
}

func newElementPrimaryExpression(name string, args []Expression) PrimaryExpression {
    expr := &ElementPrimaryExpression{}
    expr.t = ElementPrimaryExpressionType
    expr.name = name
    expr.args = args
    expr.doExec = expr.doExecute
    return expr
}

func (priExpr *ElementPrimaryExpression) setStack(stack Function) {
    priExpr.stack = stack

    for _, subExpr := range priExpr.args {
        subExpr.setStack(stack)
    }
}

func (priExpr *ElementPrimaryExpression) getName() string {
    return "#element"
}

func (priExpr *ElementPrimaryExpression) doExecute() Value {
    varname := priExpr.name
    varVal := priExpr.getVar(varname)

    argRawVals := priExpr.toGoTypeValues(priExpr.args)
    if varVal.isJsonArray() {
        arr := goArr(varVal)
        index := toIntValue(argRawVals[0])
        return arr.get(index)
    } else if varVal.isJsonObject() {
        obj := goObj(varVal)
        key := toStringValue(argRawVals[0])
        return obj.get(key)
    } else {
        errorf("failed to eval element %v[%v]: %v is not jsonArray or jsonObject", varname, argRawVals[0], varname)
        return nil
    }
}

