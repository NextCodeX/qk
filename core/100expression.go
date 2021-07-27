package core

import (
    "bytes"
)

type ExpressionType int

const (
    PrimaryExpression ExpressionType = 1 << iota  // 不可再分的原始表达式
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
    t          ExpressionType // 表达式类型
    op         OperationType // 表达式操作类型
    left       *PrimaryExpr
    right      *PrimaryExpr
    list       []*Expression // 多元表达式拆分后的二元表达式列表
    finalExpr  *Expression // 多元表达式中最后执行的表达式
    raw        []Token // 原始Token列表
    res        Value  // 常量折叠缓存值
    receiver   string // 最终赋值变量名
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
    expr.receiver = name
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
