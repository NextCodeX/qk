package core

import (
    "fmt"
    "os"
)

type DynamicStrPrimaryExpression struct {
    val Value
    PrimaryExpressionImpl
}

func newDynamicStrPrimaryExpression(val Value) PrimaryExpression {
    expr := &DynamicStrPrimaryExpression{}
    expr.t = DynamicStrPrimaryExpressionType
    expr.val = val
    expr.doExec = expr.doExecute
    return expr
}

func (priExpr *DynamicStrPrimaryExpression) getName() string {
    return "#dynamicString"
}

func (priExpr *DynamicStrPrimaryExpression) doExecute() Value {
    raw := goStr(priExpr.val)
    res := os.Expand(raw, func(key string) string {
        qkValue := evalScript(key, priExpr.getStack())
        return fmt.Sprint(qkValue.val())
    })
    return newQKValue(res)
}

