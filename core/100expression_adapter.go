package core

type ValueReceiver interface {
    beAssigned(res Value)
}

type ExpressionAdapter struct {
    ts []Token
    ValueStack
}

func (ea *ExpressionAdapter) getVar(name string) Value {
    if isTmpVar(name) {
        tmpVars := ea.ValueStack.getVar(tmpVarsKey)
        return goObj(tmpVars).get(name)
    }
    return ea.ValueStack.getVar(name)
}

func (ea *ExpressionAdapter) setVar(name string, value Value) {
    if isTmpVar(name) {
        tmpVars := ea.ValueStack.getVar(tmpVarsKey)
        goObj(tmpVars).put(name, value)
    } else {
        ea.ValueStack.setVar(name, value)
    }
}

func (ea *ExpressionAdapter) raw() []Token {
    return ea.ts
}
func (ea *ExpressionAdapter) setRaw(ts []Token) {
    ea.ts = ts
}

func (ea *ExpressionAdapter) setStack(stack Function) {
    ea.ValueStack.stack = stack
}
func (ea *ExpressionAdapter) getStack() Function {
    return ea.ValueStack.stack
}

// 赋值
func (ea *ExpressionAdapter) evalAssign(priExpr Expression, res Value) {
    if expr, ok := priExpr.(ValueReceiver); ok {
        expr.beAssigned(res)
    } else {
        errorf("invalid assign expression: %v[%v] = %v", priExpr, tokensString(priExpr.raw()), res.val())
    }
}

func evalGoValues(exprs []Expression) []interface{} {
    var res []interface{}
    for _, expr := range exprs {
        if expr == nil {
            continue
        }
        val := expr.execute()
        res = append(res, val.val())
    }
    return res
}

func evalQKValues(exprs []Expression) []Value {
    var res []Value
    for _, expr := range exprs {
        if expr == nil {
            continue
        }
        val := expr.execute()
        res = append(res, val)
    }
    return res
}

func (ea *ExpressionAdapter) isPrimaryExpression() bool {
    return false
}
func (ea *ExpressionAdapter) isBinaryExpression() bool {
    return false
}
func (ea *ExpressionAdapter) isMultiExpression() bool {
    return false
}