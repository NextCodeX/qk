package core

type ForeachType int

const (
    ForeachIndexValue ForeachType = 1 << iota
    ForeachIndex
    ForeachValue
)

type ForeachStatement struct {
    t ForeachType
    indexName string
    valueName string
    listExpr Expression
    StatementAdapter
}

func newForeachStatement(indexName, valueName string, ts []Token) Statement {
    expr := extractExpression(ts)
    stmt := &ForeachStatement{t:ForeachIndexValue, indexName: indexName, valueName: valueName, listExpr: expr}
    stmt.initStatement(stmt)
    return stmt
}

func newForIndexStatement(indexName string, ts []Token) Statement {
    expr := extractExpression(ts)
    stmt := &ForeachStatement{t:ForeachIndexValue, indexName: indexName, listExpr: expr}
    stmt.initStatement(stmt)
    return stmt
}

func newForValueStatement(valueName string, ts []Token) Statement {
    expr := extractExpression(ts)
    stmt := &ForeachStatement{t:ForeachIndexValue, valueName: valueName, listExpr: expr}
    stmt.initStatement(stmt)
    return stmt
}

func (stmt *ForeachStatement) parse() {
    stmt.listExpr.setStack(stmt.getStack())
    Compile(stmt)
}

func (stmt *ForeachStatement) execute() StatementResult {
    varVal := stmt.listExpr.execute()
    itr, ok := varVal.(Iterator)
    if !ok {
        runtimeExcption(varVal.val(), "is not iterator!")
        return newStatementResult(StatementNormal, NULL)
    }

    var res StatementResult
    indexs := itr.indexs()
    for _, index := range indexs {

        if stmt.t != ForeachValue  {
            i := newQKValue(index)
            stmt.setVar(stmt.indexName, i)
        }
        if stmt.t != ForeachIndex {
            item := itr.getItem(index)
            stmt.setVar(stmt.valueName, item)
        }

        res = stmt.executeStatementList(stmt.block, StmtListTypeFor)

        if res.isBreak() {
            res.setType(StatementNormal)
            return res
        } else if res.isReturn() {
            return res
        }
    }
    // fix foreach: empty loop exception
    if res == nil {
        res = newStatementResult(StatementNormal, NULL)
    }
    return res
}