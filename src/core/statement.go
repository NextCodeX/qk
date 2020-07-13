package core

import (
    "bytes"
)

type StatementType int

const (
    ExpressionStatement StatementType = 1 << iota
    IfStatement
    ForStatement
    SwitchStatement
    ReturnStatement
)


type Statement struct {
    t StatementType
    tmpcount int
    exprs []*Expression
    preExpr *Expression
    condition *Expression
    postExpr *Expression
    block []*Statement
    raw []Token // token列表
    compiled bool
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

func (stmt *Statement) setCompiled(flag bool) {
    stmt.compiled = flag
}


func (s *Statement) setHeaderInfo(exprs []*Expression) {
    s.preExpr = exprs[0]
    s.condition = exprs[1]
    s.postExpr = exprs[2]
}

func (s *Statement) addExpression(expr *Expression) {
    if len(s.exprs)>0 && s.exprs[len(s.exprs)-1].isMultiExpression() && !(s.exprs[len(s.exprs)-1].listFinish) {
        lastExpr := s.exprs[len(s.exprs)-1]
        subExprs := &lastExpr.list

        if expr.isMultiExpression() {
            expr = expr.list[0]
        }

        *subExprs = append(*subExprs, expr)

        if !expr.isTmpExpression() {
            lastExpr.listFinish = true
        }
        return
    }
    s.exprs = append(s.exprs, expr)
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

func (s *Statement) isSwitchStatement() bool {
    return (s.t & SwitchStatement) == SwitchStatement
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
        res.WriteString(s.condition.String())
        res.WriteString(" ")
    }
    if (s.t & ForStatement) == ForStatement {
        res.WriteString("header:")
        res.WriteString(s.preExpr.String())
        res.WriteString("; ")
        res.WriteString(s.condition.String())
        res.WriteString("; ")
        res.WriteString(s.postExpr.String())
        res.WriteString(" ")
    }
    return res.String()
}

func newStatement(t StatementType, ts []Token) *Statement {
    return &Statement{
        t:     t,
        exprs: nil,
        raw:   ts,
    }
}








