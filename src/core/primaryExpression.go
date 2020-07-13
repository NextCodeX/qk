package core



type PrimaryExpressionType int
const (
	VarPrimaryExpressionType PrimaryExpressionType = 1 << iota
	ConstPrimaryExpressionType
	ElementPrimaryExpressionType
	AttibutePrimaryExpressionType
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

func (this *PrimaryExpr) isElement() bool {
	return (this.t & ElementPrimaryExpressionType) == ElementPrimaryExpressionType
}

func (this *PrimaryExpr) isAttibute() bool {
	return (this.t & AttibutePrimaryExpressionType) == AttibutePrimaryExpressionType
}

func (this *PrimaryExpr) isOther() bool {
	return (this.t & OtherPrimaryExpressionType) == OtherPrimaryExpressionType
}


