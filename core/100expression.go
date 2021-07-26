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
    Opne // not equal to
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
    res Value
    tmpname string
}

func (expr *Expression) isAssign() bool {
    return (expr.op & Opassign) == Opassign
}
func (expr *Expression) isAssignAfterAdd() bool {
    return (expr.op & OpassignAfterAdd) == OpassignAfterAdd
}
func (expr *Expression) isAssignAfterSub() bool {
    return (expr.op & OpassignAfterSub) ==OpassignAfterSub
}
func (expr *Expression) isAssignAfterMul() bool {
    return (expr.op & OpassignAfterMul) ==OpassignAfterMul
}
func (expr *Expression) isAssignAfterDiv() bool {
    return (expr.op & OpassignAfterDiv) ==OpassignAfterDiv
}
func (expr *Expression) isAssignAfterMod() bool {
    return (expr.op & OpassignAfterMod) ==OpassignAfterMod
}



func (expr *Expression) isEq() bool {
    return (expr.op & Opeq) == Opeq
}
func (expr *Expression) isNe() bool {
    return (expr.op & Opne) == Opne
}
func (expr *Expression) isGt() bool {
    return (expr.op & Opgt) == Opgt
}
func (expr *Expression) isLt() bool {
    return (expr.op & Oplt) == Oplt
}
func (expr *Expression) isGe() bool {
    return (expr.op & Opge) == Opge
}
func (expr *Expression) isLe() bool {
    return (expr.op & Ople) == Ople
}

func (expr *Expression) isOr() bool {
    return (expr.op & Opor) == Opor
}
func (expr *Expression) isAnd() bool {
    return (expr.op & Opand) == Opand
}



func (expr *Expression) isAdd() bool {
    return (expr.op & Opadd) == Opadd
}

func (expr *Expression) isSub() bool {
    return (expr.op & Opsub) == Opsub
}

func (expr *Expression) isMul() bool {
    return (expr.op & Opmul) == Opmul
}

func (expr *Expression) isDiv() bool {
    return (expr.op & Opdiv) == Opdiv
}
func (expr *Expression) isMod() bool {
    return (expr.op & Opmod) ==Opmod
}

func (expr *Expression) isIntExpression() bool {
    return (expr.t & IntExpression) == IntExpression
}

func (expr *Expression) isFloatExpression() bool {
    return (expr.t & FloatExpression) == FloatExpression
}

func (expr *Expression) isBooleanExpression() bool {
    return (expr.t & BooleanExpression) == BooleanExpression
}

func (expr *Expression) isStringExpression() bool {
    return (expr.t & StringExpression) == StringExpression
}

func (expr *Expression) isJSONObjectExpression() bool {
    return (expr.t & JSONObjectExpression) == JSONObjectExpression
}

func (expr *Expression) isJSONArrayExpression() bool {
    return (expr.t & JSONArrayExpression) == JSONArrayExpression
}

func (expr *Expression) isVarExpression() bool {
    return (expr.t & VarExpression) == VarExpression
}

func (expr *Expression) isAttributeExpression() bool {
	return (expr.t & AttributeExpression) == AttributeExpression
}

func (expr *Expression) isElementExpression() bool {
	return (expr.t & ElementExpression) == ElementExpression
}

func (expr *Expression) isFunctionCallExpression() bool {
    return (expr.t & FunctionCallExpression) == FunctionCallExpression
}

func (expr *Expression) isMethodCallExpression() bool {
    return (expr.t & MethodCallExpression) == MethodCallExpression
}

func (expr *Expression) isPrimaryExpression() bool {
    return (expr.t & PrimaryExpression) == PrimaryExpression
}

func (expr *Expression) isBinaryExpression() bool {
    return (expr.t & BinaryExpression) == BinaryExpression
}

func (expr *Expression) isMultiExpression() bool {
    return (expr.t & MultiExpression) == MultiExpression
}

func (expr *Expression) isAssignExpression() bool {
    return (expr.t & AssignExpression) == AssignExpression
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
