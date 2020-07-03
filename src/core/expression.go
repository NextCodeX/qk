package core

import (
    "bytes"
)

type ExpressionType int

const (
    IntExpression ExpressionType = 1 << iota
    FloatExpression
    BooleanExpression
    StringExpression
    ConstExpression
    IdentifierExpression
    AddExpression
    SubExpression
    MulExpression
    DivExpression
    AssignExpression
    FunctionCallExpression
    MethodCallExpression
    BinaryExpression
    ListExpression
    TmpExpression
)

type OperationType int

const (
    Opeq OperationType = 1 << iota
    Opgt
    Oplt
    Opge
    Ople
    Opassign
)


type Expression struct {
    t ExpressionType
    raw []Token
    list []*Expression
    listFinish bool
    left PrimaryExpr
    op OperationType
    right PrimaryExpr
    result *Value
    tmpname string
}

func newListExpression() *Expression {
    return &Expression{
        t:          ListExpression,
        listFinish: false,
    }
}

func (this *Expression) isIntExpression() bool {
    return (this.t & IntExpression) == IntExpression
}

func (this *Expression) isFloatExpression() bool {
    return (this.t & FloatExpression) == FloatExpression
}

func (this *Expression) isBooleanExpression() bool {
    return (this.t & BooleanExpression) == BooleanExpression
}

func (this *Expression) isStringExpression() bool {
    return (this.t & StringExpression) == StringExpression
}

func (this *Expression) isConstExpression() bool {
    return (this.t & ConstExpression) == ConstExpression
}

func (this *Expression) isIdentifierExpression() bool {
    return (this.t & IdentifierExpression) == IdentifierExpression
}

func (this *Expression) isAssignExpression() bool {
    return (this.t & AssignExpression) == AssignExpression
}

func (this *Expression) isFunctionCallExpression() bool {
    return (this.t & FunctionCallExpression) == FunctionCallExpression
}

func (this *Expression) isMethodCallExpression() bool {
    return (this.t & MethodCallExpression) == MethodCallExpression
}

func (this *Expression) isBinaryExpression() bool {
    return (this.t & BinaryExpression) == BinaryExpression
}

func (this *Expression) isListExpression() bool {
    return (this.t & ListExpression) == ListExpression
}

func (this *Expression) isTmpExpression() bool {
    return (this.t & TmpExpression) == TmpExpression
}

func (s *Expression) execute(super, local Variables) *StatementResultType {


    return nil
}

func (expr *Expression) String() string {
    if expr == nil {
        return ""
    }
    var res bytes.Buffer
    for _, t := range expr.raw {
        res.WriteString(t.String())
        res.WriteString(" ")
    }
    return res.String()
}

type PrimaryExpressionType int

const (
    Varname PrimaryExpressionType = 1 << iota
    Expr
    Const
    Fill
)

type PrimaryExpr struct {
    t PrimaryExpressionType
    name string
    args []Expression
    result *Value
}

func (this *PrimaryExpr) isVarname() bool {
    return (this.t & Varname) == Varname
}

func (this *PrimaryExpr) isExpr() bool {
    return (this.t & Expr) == Expr
}

func (this *PrimaryExpr) isConst() bool {
    return (this.t & Const) == Const
}

