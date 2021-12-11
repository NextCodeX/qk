package core

import "fmt"

type OperationType int

const (
	// 逻辑运算
	Opeq OperationType = 1 << iota //等于 equal to
	Opne                           // not equal to
	Opgt                           // 大于 greater than
	Oplt                           // 小于 less than
	Opge                           // 大于等于 greater than or equal to
	Ople                           // 小于等于 less than or equal to

	Opor  // 逻辑或
	Opand // 逻辑与

	// 算术运算
	Opadd // 相加
	Opsub // 相减
	Opmul // 相乘
	Opdiv // 相除
	Opmod // 求余(余数 remainder/残余数 modulo)

	// 赋值运算
	Opassign         // 赋值
	OpassignAfterAdd // 相加后赋值
	OpassignAfterSub // 相减后赋值
	OpassignAfterMul // 相乘后赋值
	OpassignAfterDiv // 相除后赋值
	OpassignAfterMod // 求余后赋值
)

type BinaryExpression interface {
	execute() Value        // 执行二元表达式
	resultCache(res Value) // 表达式结果赋值于receiver

	evalAndBinaryExpression() Value
	evalOrBinaryExpression() Value

	evalEqBinaryExpression() (res Value)
	evalNeBinaryExpression() (res Value)
	evalGtBinaryExpression() (res Value)
	evalLtBinaryExpression() (res Value)
	evalGeBinaryExpression() (res Value)
	evalLeBinaryExpression() (res Value)

	evalAssignAfterAddBinaryExpression() (res Value)
	evalAssignAfterSubBinaryExpression() (res Value)
	evalAssignAfterMulBinaryExpression() (res Value)
	evalAssignAfterDivBinaryExpression() (res Value)
	evalAssignAfterModBinaryExpression() (res Value)
	evalAssignBinaryExpression() Value

	evalAddBinaryExpression() (res Value)
	evalSubBinaryExpression() (res Value)
	evalMulBinaryExpression() (res Value)
	evalDivBinaryExpression() (res Value)
	evalModBinaryExpression() (res Value)

	isAssign() bool
	isAssignAfterAdd() bool
	isAssignAfterSub() bool
	isAssignAfterMul() bool
	isAssignAfterDiv() bool
	isAssignAfterMod() bool
	isEq() bool
	isNe() bool
	isGt() bool
	isLt() bool
	isGe() bool
	isLe() bool
	isOr() bool
	isAnd() bool
	isAdd() bool
	isSub() bool
	isMul() bool
	isDiv() bool
	isMod() bool

	rightVal() Value
	leftVal() Value
	rightExpr() PrimaryExpression
	leftExpr() PrimaryExpression
	setReceiver(name PrimaryExpression)
	getReceiver() PrimaryExpression
	Expression
}

type BinaryExpressionImpl struct {
	t        OperationType
	left     PrimaryExpression
	right    PrimaryExpression
	receiver PrimaryExpression // 结果接收表达式
	res      Value             // 用于常量折叠
	ExpressionAdapter
}

func newBinaryExpression() BinaryExpression {
	return &BinaryExpressionImpl{}
}

func (binExpr *BinaryExpressionImpl) setParent(p Function) {
	binExpr.ExpressionAdapter.setParent(p)

	binExpr.left.setParent(p)
	if binExpr.localScope {
		binExpr.left.setLocalScope()
	}

	binExpr.right.setParent(p)
	if binExpr.receiver != nil {
		if binExpr.localScope {
			binExpr.receiver.setLocalScope()
		}
		binExpr.receiver.setParent(p)
	}
}

func (binExpr *BinaryExpressionImpl) execute() Value {
	left := binExpr.left
	right := binExpr.right

	if binExpr.res != nil && left.isConst() && right.isConst() {
		return binExpr.res
	}

	var res Value
	switch {
	case binExpr.isAssign():
		res = binExpr.evalAssignBinaryExpression()
	case binExpr.isAssignAfterAdd():
		res = binExpr.evalAssignAfterAddBinaryExpression()
	case binExpr.isAssignAfterSub():
		res = binExpr.evalAssignAfterSubBinaryExpression()
	case binExpr.isAssignAfterMul():
		res = binExpr.evalAssignAfterMulBinaryExpression()
	case binExpr.isAssignAfterDiv():
		res = binExpr.evalAssignAfterDivBinaryExpression()
	case binExpr.isAssignAfterMod():
		res = binExpr.evalAssignAfterModBinaryExpression()

	case binExpr.isAdd():
		res = binExpr.evalAddBinaryExpression()
	case binExpr.isSub():
		res = binExpr.evalSubBinaryExpression()
	case binExpr.isMul():
		res = binExpr.evalMulBinaryExpression()
	case binExpr.isDiv():
		res = binExpr.evalDivBinaryExpression()
	case binExpr.isMod():
		res = binExpr.evalModBinaryExpression()

	case binExpr.isEq():
		res = binExpr.evalEqBinaryExpression()
	case binExpr.isNe():
		res = binExpr.evalNeBinaryExpression()
	case binExpr.isGt():
		res = binExpr.evalGtBinaryExpression()
	case binExpr.isGe():
		res = binExpr.evalGeBinaryExpression()
	case binExpr.isLt():
		res = binExpr.evalLtBinaryExpression()
	case binExpr.isLe():
		res = binExpr.evalLeBinaryExpression()

	case binExpr.isOr():
		res = binExpr.evalOrBinaryExpression()
	case binExpr.isAnd():
		res = binExpr.evalAndBinaryExpression()
	}
	if res == nil {
		res = NULL
	}
	binExpr.resultCache(res)
	// 常量折叠
	if left.isConst() && right.isConst() {
		binExpr.res = res
	}
	return res
}

// 表达式结果赋值于receiver
func (binExpr *BinaryExpressionImpl) resultCache(res Value) {
	if binExpr.receiver == nil {
		return
	}
	binExpr.evalAssign(binExpr.receiver, res)
}

func (binExpr *BinaryExpressionImpl) setReceiver(name PrimaryExpression) {
	binExpr.receiver = name
}
func (binExpr *BinaryExpressionImpl) getReceiver() PrimaryExpression {
	return binExpr.receiver
}

func (binExpr *BinaryExpressionImpl) rightVal() Value {
	return binExpr.right.execute()
}
func (binExpr *BinaryExpressionImpl) leftVal() Value {
	return binExpr.left.execute()
}

func (binExpr *BinaryExpressionImpl) rightExpr() PrimaryExpression {
	return binExpr.right
}
func (binExpr *BinaryExpressionImpl) leftExpr() PrimaryExpression {
	return binExpr.left
}

func (binExpr *BinaryExpressionImpl) evalAndBinaryExpression() Value {
	left := binExpr.leftVal()
	if !toBoolean(left) {
		return left
	} else {
		return binExpr.rightVal()
	}
}

func (binExpr *BinaryExpressionImpl) evalOrBinaryExpression() Value {
	left := binExpr.leftVal()
	if toBoolean(left) {
		return left
	} else {
		return binExpr.rightVal()
	}
}

func (binExpr *BinaryExpressionImpl) evalEqBinaryExpression() (res Value) {
	left := binExpr.leftVal()
	right := binExpr.rightVal()
	var tmpVal interface{}
	switch {
	case left.isBoolean() && right.isBoolean():
		tmpVal = goBool(left) == goBool(right)
	case left.isInt() && right.isInt():
		tmpVal = goInt(left) == goInt(right)
	case left.isFloat() && right.isFloat():
		tmpVal = goFloat(left) == goFloat(right)
	case left.isString() && right.isString():
		tmpVal = goStr(left) == goStr(right)

	case left.isFloat() && right.isInt():
		tmpVal = goFloat(left) == float64(goInt(right))
	case left.isInt() && right.isFloat():
		tmpVal = float64(goInt(left)) == goFloat(right)

	default:
		tmpVal = left.val() == right.val()
		//errorf("invalid expression: %v == %v", left.val(), right.val())
	}
	res = newQKValue(tmpVal)
	return res
}

func (binExpr *BinaryExpressionImpl) evalNeBinaryExpression() (res Value) {
	left := binExpr.leftVal()
	right := binExpr.rightVal()
	var tmpVal interface{}
	switch {
	case left.isBoolean() && right.isBoolean():
		tmpVal = goBool(left) != goBool(right)
	case left.isInt() && right.isInt():
		tmpVal = goInt(left) != goInt(right)
	case left.isFloat() && right.isFloat():
		tmpVal = goFloat(left) != goFloat(right)
	case left.isString() && right.isString():
		tmpVal = goStr(left) != goStr(right)

	case left.isFloat() && right.isInt():
		tmpVal = goFloat(left) != float64(goInt(right))
	case left.isInt() && right.isFloat():
		tmpVal = float64(goInt(left)) != goFloat(right)

	default:
		errorf("invalid expression: %v != %v", left.val(), right.val())
	}
	res = newQKValue(tmpVal)
	return res
}

func (binExpr *BinaryExpressionImpl) evalGtBinaryExpression() (res Value) {
	left := binExpr.leftVal()
	right := binExpr.rightVal()
	var tmpVal interface{}
	switch {
	case left.isInt() && right.isInt():
		tmpVal = goInt(left) > goInt(right)
	case left.isFloat() && right.isFloat():
		tmpVal = goFloat(left) > goFloat(right)
	case left.isString() && right.isString():
		tmpVal = goStr(left) > goStr(right)

	case left.isFloat() && right.isInt():
		tmpVal = goFloat(left) > float64(goInt(right))
	case left.isInt() && right.isFloat():
		tmpVal = float64(goInt(left)) > goFloat(right)

	default:
		errorf("invalid expression: %v > %v", left.val(), right.val())
	}
	res = newQKValue(tmpVal)
	return res
}

func (binExpr *BinaryExpressionImpl) evalLtBinaryExpression() (res Value) {
	left := binExpr.leftVal()
	right := binExpr.rightVal()
	var tmpVal interface{}
	switch {
	case left.isInt() && right.isInt():
		tmpVal = goInt(left) < goInt(right)
	case left.isFloat() && right.isFloat():
		tmpVal = goFloat(left) < goFloat(right)
	case left.isString() && right.isString():
		tmpVal = goStr(left) < goStr(right)

	case left.isFloat() && right.isInt():
		tmpVal = goFloat(left) < float64(goInt(right))
	case left.isInt() && right.isFloat():
		tmpVal = float64(goInt(left)) < goFloat(right)

	default:
		errorf("invalid expression: %v < %v", left.val(), right.val())
	}
	res = newQKValue(tmpVal)
	return res
}

func (binExpr *BinaryExpressionImpl) evalGeBinaryExpression() (res Value) {
	left := binExpr.leftVal()
	right := binExpr.rightVal()
	var tmpVal interface{}
	switch {
	case left.isInt() && right.isInt():
		tmpVal = goInt(left) >= goInt(right)
	case left.isFloat() && right.isFloat():
		tmpVal = goFloat(left) >= goFloat(right)
	case left.isString() && right.isString():
		tmpVal = goStr(left) >= goStr(right)

	case left.isFloat() && right.isInt():
		tmpVal = goFloat(left) >= float64(goInt(right))
	case left.isInt() && right.isFloat():
		tmpVal = float64(goInt(left)) >= goFloat(right)

	default:
		errorf("invalid expression: %v >= %v", left.val(), right.val())
	}
	res = newQKValue(tmpVal)
	return res
}

func (binExpr *BinaryExpressionImpl) evalLeBinaryExpression() (res Value) {
	left := binExpr.leftVal()
	right := binExpr.rightVal()
	var tmpVal interface{}
	switch {
	case left.isInt() && right.isInt():
		tmpVal = goInt(left) <= goInt(right)
	case left.isFloat() && right.isFloat():
		tmpVal = goFloat(left) <= goFloat(right)
	case left.isString() && right.isString():
		tmpVal = goStr(left) <= goStr(right)

	case left.isFloat() && right.isInt():
		tmpVal = goFloat(left) <= float64(goInt(right))
	case left.isInt() && right.isFloat():
		tmpVal = float64(goInt(left)) <= goFloat(right)

	default:
		errorf("invalid expression: %v <= %v", left.val(), right.val())
	}
	res = newQKValue(tmpVal)
	return res
}

func (binExpr *BinaryExpressionImpl) evalAssignAfterAddBinaryExpression() (res Value) {
	res = binExpr.evalAddBinaryExpression()
	binExpr.evalAssign(binExpr.left, res)
	return res
}

func (binExpr *BinaryExpressionImpl) evalAssignAfterSubBinaryExpression() (res Value) {
	res = binExpr.evalSubBinaryExpression()
	binExpr.evalAssign(binExpr.left, res)
	return res
}

func (binExpr *BinaryExpressionImpl) evalAssignAfterMulBinaryExpression() (res Value) {
	res = binExpr.evalMulBinaryExpression()
	binExpr.evalAssign(binExpr.left, res)
	return res
}

func (binExpr *BinaryExpressionImpl) evalAssignAfterDivBinaryExpression() (res Value) {
	res = binExpr.evalDivBinaryExpression()
	binExpr.evalAssign(binExpr.left, res)
	return res
}

func (binExpr *BinaryExpressionImpl) evalAssignAfterModBinaryExpression() (res Value) {
	res = binExpr.evalModBinaryExpression()
	binExpr.evalAssign(binExpr.left, res)
	return res
}

func (binExpr *BinaryExpressionImpl) evalAssignBinaryExpression() Value {
	res := binExpr.rightVal()
	binExpr.evalAssign(binExpr.left, res)
	return res
}

func (binExpr *BinaryExpressionImpl) evalAddBinaryExpression() (res Value) {
	left := binExpr.leftVal()
	right := binExpr.rightVal()
	var tmpVal interface{}
	switch {
	case left.isInt() && right.isInt():
		tmpVal = goInt(left) + goInt(right)

	case left.isFloat() && right.isFloat():
		tmpVal = goFloat(left) + goFloat(right)

	case left.isFloat() && right.isInt():
		tmpVal = goFloat(left) + float64(goInt(right))

	case left.isInt() && right.isFloat():
		tmpVal = float64(goInt(left)) + goFloat(right)

	case left.isString() || right.isString():
		tmpVal = fmt.Sprintf("%v%v", left.val(), right.val())

	default:
		runtimeExcption("invalid binary expression:", left.val(), "+", right.val(), " -> ", tokensString(binExpr.tokenList()))
	}

	res = newQKValue(tmpVal)
	return res
}

func (binExpr *BinaryExpressionImpl) evalSubBinaryExpression() (res Value) {
	left := binExpr.leftVal()
	right := binExpr.rightVal()
	var tmpVal interface{}
	switch {
	case left.isInt() && right.isInt():
		tmpVal = goInt(left) - goInt(right)

	case left.isFloat() && right.isFloat():
		tmpVal = goFloat(left) - goFloat(right)

	case left.isFloat() && right.isInt():
		tmpVal = goFloat(left) - float64(goInt(right))

	case left.isInt() && right.isFloat():
		tmpVal = float64(goInt(left)) - goFloat(right)

	default:
		runtimeExcption("unknow operation:", left.val(), "-", right.val())
	}
	res = newQKValue(tmpVal)
	return res
}

func (binExpr *BinaryExpressionImpl) evalMulBinaryExpression() (res Value) {
	left := binExpr.leftVal()
	right := binExpr.rightVal()
	var tmpVal interface{}
	switch {
	case left.isInt() && right.isInt():
		tmpVal = goInt(left) * goInt(right)

	case left.isFloat() && right.isFloat():
		tmpVal = goFloat(left) * goFloat(right)

	case left.isFloat() && right.isInt():
		tmpVal = goFloat(left) * float64(goInt(right))

	case left.isInt() && right.isFloat():
		tmpVal = float64(goInt(left)) * goFloat(right)

	default:
		runtimeExcption("unknow operation:", left.val(), "*", right.val())
	}
	res = newQKValue(tmpVal)
	return res
}

func (binExpr *BinaryExpressionImpl) evalDivBinaryExpression() (res Value) {
	left := binExpr.leftVal()
	right := binExpr.rightVal()

	if (right.isInt() && goInt(right) == 0) || (right.isFloat() && goFloat(right) == 0) {
		runtimeExcption("Invalid Operation: divide zero")
	}

	var tmpVal interface{}
	switch {
	case left.isInt() && right.isInt():
		tmpVal = goInt(left) / goInt(right)

	case left.isFloat() && right.isFloat():
		tmpVal = goFloat(left) / goFloat(right)

	case left.isFloat() && right.isInt():
		tmpVal = goFloat(left) / float64(goInt(right))

	case left.isInt() && right.isFloat():
		tmpVal = float64(goInt(left)) / goFloat(right)

	default:
		runtimeExcption("unknow operation:", left.val(), "/", right.val())
	}
	res = newQKValue(tmpVal)
	return res
}

func (binExpr *BinaryExpressionImpl) evalModBinaryExpression() (res Value) {
	left := binExpr.leftVal()
	right := binExpr.rightVal()

	if (right.isInt() && goInt(right) == 0) || (right.isFloat() && goFloat(right) == 0) {
		runtimeExcption("Invalid Operation: divide zero")
	}

	var tmpVal interface{}
	switch {
	case left.isInt() && right.isInt():
		tmpVal = goInt(left) % goInt(right)

	default:
		errorf("invalid expression: %v %v %v", left.val(), "%", right.val())
	}
	res = newQKValue(tmpVal)
	return res
}

func (binExpr *BinaryExpressionImpl) isAssign() bool {
	return (binExpr.t & Opassign) == Opassign
}
func (binExpr *BinaryExpressionImpl) isAssignAfterAdd() bool {
	return (binExpr.t & OpassignAfterAdd) == OpassignAfterAdd
}
func (binExpr *BinaryExpressionImpl) isAssignAfterSub() bool {
	return (binExpr.t & OpassignAfterSub) == OpassignAfterSub
}
func (binExpr *BinaryExpressionImpl) isAssignAfterMul() bool {
	return (binExpr.t & OpassignAfterMul) == OpassignAfterMul
}
func (binExpr *BinaryExpressionImpl) isAssignAfterDiv() bool {
	return (binExpr.t & OpassignAfterDiv) == OpassignAfterDiv
}
func (binExpr *BinaryExpressionImpl) isAssignAfterMod() bool {
	return (binExpr.t & OpassignAfterMod) == OpassignAfterMod
}

func (binExpr *BinaryExpressionImpl) isEq() bool {
	return (binExpr.t & Opeq) == Opeq
}
func (binExpr *BinaryExpressionImpl) isNe() bool {
	return (binExpr.t & Opne) == Opne
}
func (binExpr *BinaryExpressionImpl) isGt() bool {
	return (binExpr.t & Opgt) == Opgt
}
func (binExpr *BinaryExpressionImpl) isLt() bool {
	return (binExpr.t & Oplt) == Oplt
}
func (binExpr *BinaryExpressionImpl) isGe() bool {
	return (binExpr.t & Opge) == Opge
}
func (binExpr *BinaryExpressionImpl) isLe() bool {
	return (binExpr.t & Ople) == Ople
}

func (binExpr *BinaryExpressionImpl) isOr() bool {
	return (binExpr.t & Opor) == Opor
}
func (binExpr *BinaryExpressionImpl) isAnd() bool {
	return (binExpr.t & Opand) == Opand
}

func (binExpr *BinaryExpressionImpl) isAdd() bool {
	return (binExpr.t & Opadd) == Opadd
}

func (binExpr *BinaryExpressionImpl) isSub() bool {
	return (binExpr.t & Opsub) == Opsub
}

func (binExpr *BinaryExpressionImpl) isMul() bool {
	return (binExpr.t & Opmul) == Opmul
}

func (binExpr *BinaryExpressionImpl) isDiv() bool {
	return (binExpr.t & Opdiv) == Opdiv
}
func (binExpr *BinaryExpressionImpl) isMod() bool {
	return (binExpr.t & Opmod) == Opmod
}
