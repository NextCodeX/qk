package core

import (
    "fmt"
    "os"
)

type DynamicStrPrimaryExpression struct {
    tmpl string
    PrimaryExpressionImpl
}

func newDynamicStrPrimaryExpression(tmpl string) PrimaryExpression {
    expr := &DynamicStrPrimaryExpression{}
    expr.t = DynamicStrPrimaryExpressionType
    expr.tmpl = tmpl
    expr.doExec = expr.doExecute
    return expr
}

func (priExpr *DynamicStrPrimaryExpression) doExecute() Value {
    res := os.Expand(priExpr.tmpl, func(key string) string {
        qkValue := evalScript(key, priExpr.getStack())
        return fmt.Sprint(qkValue.val())
    })
    return newQKValue(res)
}

