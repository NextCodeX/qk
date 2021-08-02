package core

type Statement interface {
    addStmt(stmt Statement)  // 添加子statement, 并将stack传递给 子statement
    stmts() []Statement
    raw() []Token
    setRaw(ts []Token)
    rawAppend(t Token)
    setParent(p Function) // 等同于setStack()
    getParent() Function

    getStack() Function

    parse() // 解析 各子statement, expression; 并将stack传递给 它们
    execute() StatementResult

    isExpressionStatement() bool
    isIfStatement() bool
    isForStatement() bool
    isForeachStatement() bool
    isForIndexStatement() bool
    isForItemStatement() bool
    isSwitchStatement() bool
    isMultiStatement() bool
    isContinueStatement() bool
    isBreakStatement() bool
    isReturnStatement() bool
    String() string
}

