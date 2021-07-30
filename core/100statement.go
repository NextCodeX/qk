package core

type Statement interface {
    addStmt(stmt Statement)
    stmts() []Statement
    raw() []Token
    setRaw(ts []Token)
    rawAppend(t Token)
    setParent(p Function)
    getParent() Function
    parse()
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

