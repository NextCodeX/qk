package core

import (
    "bytes"
    "strings"
)

type ExpressionType int

const (
    IntExpression ExpressionType = 1 << iota
    FloatExpression
    BooleanExpression
    StringExpression
    ConstExpression
    VarExpression
    FunctionCallExpression
    MethodCallExpression
    BinaryExpression
    MultiExpression
    TmpExpression
)

type OperationType int
const (
    Opeq OperationType = 1 << iota
    Opgt
    Oplt
    Opge
    Ople
    Opadd
    Opsub
    Opmul
    Opdiv
    Opassign
)

type Expression struct {
    vars *VarScope
    t ExpressionType
    op OperationType
    left *PrimaryExpr
    right *PrimaryExpr
    list []*Expression
    finalExpr *Expression
    listFinish bool
    raw []Token
    res *Value
    tmpname string
}

func (this *Expression) isAssign() bool {
    return (this.op & Opassign) == Opassign
}

func (this *Expression) isAdd() bool {
    return (this.op & Opadd) == Opadd
}

func (this *Expression) isSub() bool {
    return (this.op & Opsub) == Opsub
}

func (this *Expression) isMul() bool {
    return (this.op & Opmul) == Opmul
}

func (this *Expression) isDiv() bool {
    return (this.op & Opdiv) == Opdiv
}

func newMultiExpression() *Expression {
    return &Expression{
        t:          MultiExpression,
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

func (this *Expression) isVarExpression() bool {
    return (this.t & VarExpression) == VarExpression
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

func (this *Expression) isMultiExpression() bool {
    return (this.t & MultiExpression) == MultiExpression
}

func (this *Expression) isTmpExpression() bool {
    return (this.t & TmpExpression) == TmpExpression
}

func (expr *Expression) searchVariable(name string) *Variable {
    res := expr.vars.local.get(name)
    if res != nil {
        return res
    }
    if expr.vars.super == nil {
        return nil
    }
    res = expr.vars.super.get(name)
    if res != nil {
        return res
    }
    return nil
}

func (expr *Expression) addVariable(vr *Variable)  {
    expr.vars.local.add(vr)
}

func (expr *Expression) addVar(name string, val *Value)  {
    variable := toVar(name,  val)
    expr.vars.local.add(variable)
}

func (expr *Expression) leftVal() *Value {
    if expr.left == nil {
        return NULL
    }
    if expr.left.isConst() {
        return expr.left.res
    }
    if expr.left.isVar() {
        varname := expr.left.name
        variable := expr.searchVariable(varname)
        if variable == nil {
            return NULL
        }
        return variable.val
    }
    return NULL
}

func (expr *Expression) rightVal() *Value {
    if expr.right == nil {
        return NULL
    }
    if expr.right.isConst() {
        return expr.right.res
    }
    if expr.right.isVar() {
        varname := expr.right.name
        variable := expr.searchVariable(varname)
        if variable == nil {
            return NULL
        }
        return variable.val
    }
    return NULL
}

func (expr *Expression) setTmpname(name string) {
    expr.t = expr.t | TmpExpression
    expr.tmpname = name
}


func (expr *Expression) TypeString() string {
    if expr == nil {
        return ""
    }
    var res bytes.Buffer
    if expr.isIntExpression() {
        res.WriteString("int expression, ")
    }
    if expr.isBooleanExpression() {
        res.WriteString("bool expression, ")
    }
    if expr.isStringExpression() {
        res.WriteString("string expression, ")
    }
    if expr.isConstExpression() {
        res.WriteString("const expression, ")
    }
    if expr.isVarExpression() {
        res.WriteString("var expression, ")
    }
    if expr.isFunctionCallExpression() {
        res.WriteString("function call expression, ")
    }
    if expr.isMethodCallExpression() {
        res.WriteString("method call expression, ")
    }
    if expr.isBinaryExpression() {
        res.WriteString("binary expression, ")
    }
    if expr.isMultiExpression() {
        res.WriteString("multi expression, ")
    }
    if expr.isTmpExpression() {
        res.WriteString("tmp expression, ")
    }
    return strings.Trim(strings.TrimSpace(res.String()), ",")
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
    VarPrimaryExpressionType PrimaryExpressionType = 1 << iota
    ConstPrimaryExpressionType
    OtherPrimaryExpressionType PrimaryExpressionType = 0
)

type PrimaryExpr struct {
    t PrimaryExpressionType
    caller string // 调用者名称
    name string  // 变量名或者函数名称
    args []*Expression // 参数变量名
    res *Value  // 常量值
}

func (this *PrimaryExpr) isVar() bool {
    return (this.t & VarPrimaryExpressionType) == VarPrimaryExpressionType
}

func (this *PrimaryExpr) isConst() bool {
    return (this.t & ConstPrimaryExpressionType) == ConstPrimaryExpressionType
}

func (this *PrimaryExpr) isOther() bool {
    return (this.t & OtherPrimaryExpressionType) == OtherPrimaryExpressionType
}

