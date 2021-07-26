package core



type PrimaryExpressionType int
const (
	VarPrimaryExpressionType PrimaryExpressionType = 1 << iota
	ConstPrimaryExpressionType
	DynamicStrPrimaryExpressionType
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
	args []*Expression // 函数调用参数 / 数组索引
	res Value  // 常量值
}

func (priExpr *PrimaryExpr) isVar() bool {
	return (priExpr.t & VarPrimaryExpressionType) == VarPrimaryExpressionType
}

func (priExpr *PrimaryExpr) isConst() bool {
	return (priExpr.t & ConstPrimaryExpressionType) == ConstPrimaryExpressionType
}

func (priExpr *PrimaryExpr) isDynamicStr() bool {
	return (priExpr.t & DynamicStrPrimaryExpressionType) == DynamicStrPrimaryExpressionType
}

func (priExpr *PrimaryExpr) isArray() bool {
	return (priExpr.t & ArrayPrimaryExpressionType) == ArrayPrimaryExpressionType
}

func (priExpr *PrimaryExpr) isObject() bool {
	return (priExpr.t & ObjectPrimaryExpressionType) == ObjectPrimaryExpressionType
}

func (priExpr *PrimaryExpr) isElement() bool {
	return (priExpr.t & ElementPrimaryExpressionType) == ElementPrimaryExpressionType
}

func (priExpr *PrimaryExpr) isAttibute() bool {
	return (priExpr.t & AttibutePrimaryExpressionType) == AttibutePrimaryExpressionType
}

func (priExpr *PrimaryExpr) isFunctionCall() bool {
	return (priExpr.t & FunctionCallPrimaryExpressionType) == FunctionCallPrimaryExpressionType
}

func (priExpr *PrimaryExpr) isMethodCall() bool {
	return (priExpr.t & MethodCallPrimaryExpressionType) == MethodCallPrimaryExpressionType
}

func (priExpr *PrimaryExpr) isOther() bool {
	return (priExpr.t & OtherPrimaryExpressionType) == OtherPrimaryExpressionType
}


