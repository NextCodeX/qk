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
    raw []Token // token列表
}

func (s *Statement) setHeaderInfo(exprs []*Expression) {
    s.preExpr = exprs[0]
    s.condition = exprs[1]
    s.postExpr = exprs[2]
}

func (s *Statement) addExpression(expr *Expression) {
    if len(s.exprs)>0 && s.exprs[len(s.exprs)-1].isListExpression() && !(s.exprs[len(s.exprs)-1].listFinish) {
        lastExpr := s.exprs[len(s.exprs)-1]
        subExprs := &lastExpr.list

        if expr.isListExpression() {
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

func (s *Statement) execute(super, local Variables) *StatementResultType {


    return nil
}

func (s *Statement) String() string {
    var res bytes.Buffer
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

type Function struct {
    super Variables // 父作用域的变量列表
    local Variables // 当前作用域的变量列表
    params []Value // 参数
    res []Value // 返回值
    block []*Statement // 执行语句
    name string
    defToken Token
    raw []Token // token列表
}


func newFunc() *Function {
    return &Function{local:newVariables()}
}

func (f *Function) addStatement(stm *Statement) {
    f.block = append(f.block, stm)
}




