package core

import (
    "bytes"
)

type StatementType int

const (
    ExpressionStatement StatementType = 1 << iota // 表达式语句
    IfStatement // if语句
    ForStatement // for循环语句
    ForeachStatement // 遍历语句
    ForIndexStatement // 索引遍历语句
    ForItemStatement // 值遍历语句
    SwitchStatement // switch 语句
    MultiStatement // 多重语句
    ContinueStatement // continue 语句
    BreakStatement // break 语句
    ReturnStatement // 返回语句
)


type Statement struct {
    t StatementType
    expr *Expression // 用于expressionStatement.
    preExprTokens []Token // 用于普通for语句
    condExprTokens []Token // 用于if,for语句
    postExprTokens []Token // 用于普通for语句
    preExpr *Expression  // 用于for语句
    condExpr *Expression // 用于if,for语句
    postExpr *Expression // 用于for语句
    condStmts []*Statement // 用于if, switch语句
    defStmt *Statement // 用于if, switch语句
    block []*Statement // stmt核心组成(编译后的信息)
    raw []Token // stmt核心组成, token列表(编译前的信息)
    compiled bool // 该语句是否已编译
    fpi *ForPlusInfo // 增强for, 相关的信息
}


func (stmt *Statement) addStatement(stm *Statement) {
    stmt.block = append(stmt.block, stm)
}

func (stmt *Statement) stmts() []*Statement {
    return stmt.block
}

func (stmt *Statement) getRaw() []Token {
    return stmt.raw
}

func (stmt *Statement) setRaw(ts []Token) {
    stmt.raw = ts
}

func (stmt *Statement) isCompiled() bool {
    return stmt.compiled
}

func (stmt *Statement) setCompiled() {
    stmt.compiled = true
}


func (stmt *Statement) isExpressionStatement() bool {
    return (stmt.t & ExpressionStatement) == ExpressionStatement
}

func (stmt *Statement) isIfStatement() bool {
    return (stmt.t & IfStatement) == IfStatement
}

func (stmt *Statement) isForStatement() bool {
    return (stmt.t & ForStatement) == ForStatement
}

func (stmt *Statement) isForeachStatement() bool {
    return (stmt.t & ForeachStatement) == ForeachStatement
}

func (stmt *Statement) isForIndexStatement() bool {
    return (stmt.t & ForIndexStatement) == ForIndexStatement
}

func (stmt *Statement) isForItemStatement() bool {
    return (stmt.t & ForItemStatement) == ForItemStatement
}

func (stmt *Statement) isSwitchStatement() bool {
    return (stmt.t & SwitchStatement) == SwitchStatement
}

func (stmt *Statement) isMultiStatement() bool {
    return (stmt.t & MultiStatement) == MultiStatement
}

func (stmt *Statement) isContinueStatement() bool {
    return (stmt.t & ContinueStatement) == ContinueStatement
}

func (stmt *Statement) isBreakStatement() bool {
    return (stmt.t & BreakStatement) == BreakStatement
}

func (stmt *Statement) isReturnStatement() bool {
    return (stmt.t & ReturnStatement) == ReturnStatement
}


func (stmt *Statement) String() string {
    var res bytes.Buffer
    if stmt.isReturnStatement() {
        res.WriteString(" return: ")
    }
    for _, t := range stmt.raw {
        res.WriteString(t.String())
        res.WriteString(" ")
    }
    if (stmt.t & IfStatement) == IfStatement {
        res.WriteString("condition:")
        res.WriteString(tokensString(stmt.condExprTokens))
        res.WriteString(" ")
    }
    if (stmt.t & ForStatement) == ForStatement {
        res.WriteString("header:")
        res.WriteString(stmt.preExpr.String())
        res.WriteString("; ")
        res.WriteString(tokensString(stmt.condExprTokens))
        res.WriteString("; ")
        res.WriteString(stmt.postExpr.String())
        res.WriteString(" ")
    }
    return res.String()
}

func newStatement(t StatementType, ts []Token) *Statement {
    return &Statement{
        t:     t,
        raw:   ts,
    }
}








