package core



type PrimaryExpressionType int
const (
	VarPrimaryExpressionType PrimaryExpressionType = 1 << iota
	ConstPrimaryExpressionType
	ArrayPrimaryExpressionType
	ObjectPrimaryExpressionType
	ElementPrimaryExpressionType
	AttibutePrimaryExpressionType
	FunctionCallPrimaryExpressionType
	MethodCallPrimaryExpressionType
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

func (this *PrimaryExpr) isArray() bool {
	return (this.t & ArrayPrimaryExpressionType) == ArrayPrimaryExpressionType
}

func (this *PrimaryExpr) isObject() bool {
	return (this.t & ObjectPrimaryExpressionType) == ObjectPrimaryExpressionType
}

func (this *PrimaryExpr) isElement() bool {
	return (this.t & ElementPrimaryExpressionType) == ElementPrimaryExpressionType
}

func (this *PrimaryExpr) isAttibute() bool {
	return (this.t & AttibutePrimaryExpressionType) == AttibutePrimaryExpressionType
}

func (this *PrimaryExpr) isFunctionCall() bool {
	return (this.t & FunctionCallPrimaryExpressionType) == FunctionCallPrimaryExpressionType
}

func (this *PrimaryExpr) isMethodCall() bool {
	return (this.t & MethodCallPrimaryExpressionType) == MethodCallPrimaryExpressionType
}

func (this *PrimaryExpr) isOther() bool {
	return (this.t & OtherPrimaryExpressionType) == OtherPrimaryExpressionType
}


