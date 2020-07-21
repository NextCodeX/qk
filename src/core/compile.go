package core



func Compile(stmts StatementList) {
    if stmts == nil {
        return
    }
    if stmts.isCompiled() {
        return
    }else {
        stmts.setCompiled()
    }
    extractStatement(stmts)
    parseStatementList(stmts.stmts())

    for _, customFunc := range funcList {
       Compile(customFunc)
    }
}




