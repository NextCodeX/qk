package core

type ExpressionAdapter struct {
    ts []Token
    ValueStack
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

func toGoTypeValues(exprs []Expression) []interface{} {
    var res []interface{}
    for _, expr := range exprs {
        if expr == nil {
            continue
        }
        val := expr.execute()
        rawVal := val.val()
        res = append(res, rawVal)
    }
    return res
}

func evalValues(exprs []Expression) []Value {
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