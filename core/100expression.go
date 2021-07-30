package core

type Expression interface {
    raw() []Token
    setRaw(ts []Token)
    setStack(stack Function)
    getStack() Function

    execute() Value // 执行表达式

    isPrimaryExpression() bool
    isBinaryExpression() bool
    isMultiExpression() bool
}