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
    JSONObjectExpression
    JSONArrayExpression
    VarExpression
    AttributeExpression
    ElementExpression
    FunctionCallExpression
    MethodCallExpression
    PrimaryExpression  // 不可再分的原始表达式

    BinaryExpression // 二元表达式
    MultiExpression // 多元表达式
    AssignExpression // 用于表示赋值的二元表达式
)

type OperationType int
const (

    // 逻辑运算
    Opeq OperationType = 1 << iota //等于 equal to
    Opgt // 大于 greater than
    Oplt // 小于 less than
    Opge // 大于等于 greater than or equal to
    Ople // 小于等于 less than or equal to

    Opor // 逻辑或
    Opand // 逻辑与

    // 算术运算
    Opadd // 相加
    Opsub // 相减
    Opmul // 相乘
    Opdiv // 相除
    Opmod // 求余(余数 remainder/残余数 modulo)

    // 赋值运算
    Opassign // 赋值
    OpassignAfterAdd // 相加后赋值
    OpassignAfterSub // 相减后赋值
    OpassignAfterMul // 相乘后赋值
    OpassignAfterDiv // 相除后赋值
    OpassignAfterMod // 求余后赋值

)

type Expression struct {
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
func (this *Expression) isAssignAfterAdd() bool {
    return (this.op & OpassignAfterAdd) == OpassignAfterAdd
}
func (this *Expression) isAssignAfterSub() bool {
    return (this.op & OpassignAfterSub) ==OpassignAfterSub
}
func (this *Expression) isAssignAfterMul() bool {
    return (this.op & OpassignAfterMul) ==OpassignAfterMul
}
func (this *Expression) isAssignAfterDiv() bool {
    return (this.op & OpassignAfterDiv) ==OpassignAfterDiv
}
func (this *Expression) isAssignAfterMod() bool {
    return (this.op & OpassignAfterMod) ==OpassignAfterMod
}



func (this *Expression) isEq() bool {
    return (this.op & Opeq) ==Opeq
}
func (this *Expression) isGt() bool {
    return (this.op & Opgt) ==Opgt
}
func (this *Expression) isLt() bool {
    return (this.op & Oplt) ==Oplt
}
func (this *Expression) isGe() bool {
    return (this.op & Opge) ==Opge
}
func (this *Expression) isLe() bool {
    return (this.op & Ople) ==Ople
}

func (this *Expression) isOr() bool {
    return (this.op & Opor) ==Opor
}
func (this *Expression) isAnd() bool {
    return (this.op & Opand) ==Opand
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
func (this *Expression) isMod() bool {
    return (this.op & Opmod) ==Opmod
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

func (this *Expression) isJSONObjectExpression() bool {
    return (this.t & JSONObjectExpression) == JSONObjectExpression
}

func (this *Expression) isJSONArrayExpression() bool {
    return (this.t & JSONArrayExpression) == JSONArrayExpression
}

func (this *Expression) isVarExpression() bool {
    return (this.t & VarExpression) == VarExpression
}

func (this *Expression) isAttributeExpression() bool {
	return (this.t & AttributeExpression) == AttributeExpression
}

func (this *Expression) isElementExpression() bool {
	return (this.t & ElementExpression) == ElementExpression
}

func (this *Expression) isFunctionCallExpression() bool {
    return (this.t & FunctionCallExpression) == FunctionCallExpression
}

func (this *Expression) isMethodCallExpression() bool {
    return (this.t & MethodCallExpression) == MethodCallExpression
}

func (this *Expression) isPrimaryExpression() bool {
    return (this.t & PrimaryExpression) == PrimaryExpression
}

func (this *Expression) isBinaryExpression() bool {
    return (this.t & BinaryExpression) == BinaryExpression
}

func (this *Expression) isMultiExpression() bool {
    return (this.t & MultiExpression) == MultiExpression
}

func (this *Expression) isAssignExpression() bool {
    return (this.t & AssignExpression) == AssignExpression
}

func (expr *Expression) setTmpname(name string) {
    expr.t = expr.t | AssignExpression
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
    if expr.isPrimaryExpression() {
        res.WriteString("primary expression, ")
    }
    if expr.isVarExpression() {
        res.WriteString("var expression, ")
    }
	if expr.isAttributeExpression() {
		res.WriteString("attribute expression, ")
	}
	if expr.isElementExpression() {
		res.WriteString("element expression, ")
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
    if expr.isAssignExpression() {
        res.WriteString("assign expression, ")
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


func (expr *Expression) RawString() string {
	if len(expr.raw) < 1 {
		return "unrecod"
	}
	return tokensString(expr.raw)
}
