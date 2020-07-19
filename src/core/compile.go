package core

var (
    funcList = make(map[string]*Function)
    mainFunc = newFunc("main")
)

func Compile(stmts StatementList, ts []Token) {
    if stmts.isCompiled() {
        return
    }else {
        stmts.setCompiled(true)
    }
    extractStatement(stmts, ts)
    parseStatementList(stmts.stmts())

    for _, customFunc := range funcList {
       Compile(customFunc, customFunc.getRaw())
    }
}


func parseStatementList(stmts []*Statement) {
    for _, stmt := range stmts {
        parseStatement(stmt)
    }
}

func parseStatement(stmt *Statement) {
    ts := stmt.raw
    switch {
    case stmt.isExpressionStatement():
        expr := extractExpression(ts)
        stmt.addExpression(expr)

    case stmt.isIfStatement():
    case stmt.isForStatement():
    case stmt.isSwitchStatement():
    case stmt.isReturnStatement():
    }
}

