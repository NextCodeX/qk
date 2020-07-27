package core

import (
    "bytes"
)

type StatementType int

const (
    ExpressionStatement StatementType = 1 << iota
    IfStatement
    ForStatement
    ForeachStatement
    ForIndexStatement
    ForItemStatement
    SwitchStatement
    MultiStatement
    ContinueStatement
    BreakStatement
    ReturnStatement
)


type Statement struct {
    t StatementType
    expr *Expression // 用于expressionStatement.
    preExprTokens []Token // 用于普通for语句
    condExprTokens []Token // 用于if,for语句
    postExprTokens []Token // 用于普通for语句
    preExpr *Expression
    condExpr *Expression
    postExpr *Expression
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


func (s *Statement) isExpressionStatement() bool {
    return (s.t & ExpressionStatement) == ExpressionStatement
}

func (s *Statement) isIfStatement() bool {
    return (s.t & IfStatement) == IfStatement
}

func (s *Statement) isForStatement() bool {
    return (s.t & ForStatement) == ForStatement
}

func (s *Statement) isForeachStatement() bool {
    return (s.t & ForeachStatement) == ForeachStatement
}

func (s *Statement) isForIndexStatement() bool {
    return (s.t & ForIndexStatement) == ForIndexStatement
}

func (s *Statement) isForItemStatement() bool {
    return (s.t & ForItemStatement) == ForItemStatement
}

func (s *Statement) isSwitchStatement() bool {
    return (s.t & SwitchStatement) == SwitchStatement
}

func (s *Statement) isMultiStatement() bool {
    return (s.t & MultiStatement) == MultiStatement
}

func (s *Statement) isContinueStatement() bool {
    return (s.t & ContinueStatement) == ContinueStatement
}

func (s *Statement) isBreakStatement() bool {
    return (s.t & BreakStatement) == BreakStatement
}

func (s *Statement) isReturnStatement() bool {
    return (s.t & ReturnStatement) == ReturnStatement
}


func (s *Statement) String() string {
    var res bytes.Buffer
    if s.isReturnStatement() {
        res.WriteString(" return: ")
    }
    for _, t := range s.raw {
        res.WriteString(t.String())
        res.WriteString(" ")
    }
    if (s.t & IfStatement) == IfStatement {
        res.WriteString("condition:")
        res.WriteString(tokensString(s.condExprTokens))
        res.WriteString(" ")
    }
    if (s.t & ForStatement) == ForStatement {
        res.WriteString("header:")
        res.WriteString(s.preExpr.String())
        res.WriteString("; ")
        res.WriteString(tokensString(s.condExprTokens))
        res.WriteString("; ")
        res.WriteString(s.postExpr.String())
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








