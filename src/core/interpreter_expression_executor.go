package core

import (
	"fmt"
	"strings"
)

type ExpressionExecutor struct {
	expr *Expression
	stack *VariableStack
	tmpVars *Variables
}

func newExpressionExecutor(expr *Expression, stack *VariableStack, tmpVars *Variables) *ExpressionExecutor {
	executor := &ExpressionExecutor{expr:expr, stack:stack, tmpVars:tmpVars}
	return executor
}

func (executor *ExpressionExecutor) run() (res *Value) {
	expr := executor.expr
	if expr.isConstExpression() || expr.isVarExpression() || expr.isElementExpression() || expr.isAttributeExpression() {
		return executor.leftVal()
	}

	if expr.isBinaryExpression() {
		return executor.executeBinaryExpression()
	}
	if expr.isFunctionCallExpression() {
		return executor.executeFunctionCallExpression()
	}
	if expr.isMultiExpression() {
		return executor.executeMultiExpression()
	}
	runtimeExcption("unknow expression.")
	return nil
}

func (executor *ExpressionExecutor) executeAttributeExpression() (res *Value) {
	expr := executor.expr
	varname := expr.left.caller
	attrname := expr.left.name
	varVal := executor.searchVariable(varname)
	if varVal == nil {
		return NULL
	}
	varRawVal := varVal.val.val()
	obj, ok := varRawVal.(map[string]interface{})
	if !ok {
		return NULL
	}
	return newVal(obj[attrname])
}

func (executor *ExpressionExecutor) executeElementExpression() (res *Value) {
	expr := executor.expr
	varname := expr.left.name
	arrRawVal := executor.searchArrayVariable(varname)

	argVals := executor.toGoTypeValues(expr.left.args)
	argRawVals := executor.getArrayIndexs(len(arrRawVal), argVals)
	return newVal(arrRawVal[argRawVals[0]])
}


func (executor *ExpressionExecutor) executeMultiExpression() (res *Value) {
	expr := executor.expr
	return executor.recursiveEvalMultiExpression(expr.finalExpr, expr.list)
}


func (executor *ExpressionExecutor) recursiveEvalMultiExpression(nextExpr *Expression, exprList []*Expression) *Value {
	left := nextExpr.left
	right := nextExpr.right
	if left.isConst() && right.isConst() {
		return executor.evalNewExpression(nextExpr)
	}

	executor.checkValueExist(left, exprList)
	executor.checkValueExist(right, exprList)

	return executor.evalNewExpression(nextExpr)
}

func (executor *ExpressionExecutor) checkValueExist(primaryExpr *PrimaryExpr, exprList []*Expression) {
	if !primaryExpr.isVar() {
		return 
	}
	varname := primaryExpr.name
	variable := executor.searchVariable(varname)
	if variable != nil {
		return
	}
	nextExpr := executor.getNextExprForMultiExpression(varname, exprList)
	if nextExpr == nil {
		runtimeExcption("executeMultiExpression Exception")
	}
	executor.recursiveEvalMultiExpression(nextExpr, exprList)
}

func (executor *ExpressionExecutor) isTmpVar(name string) bool {
	return strings.HasPrefix(name, "tmp.")
}

func (executor *ExpressionExecutor) getNextExprForMultiExpression(varname string, exprList []*Expression) *Expression {
	for _, subExpr := range exprList {
		if subExpr.tmpname == varname {
			return subExpr
		}
	}
	return nil
}

func (executor *ExpressionExecutor) executeFunctionCallExpression() (res *Value) {
	expr := executor.expr
	functionName := expr.left.name
	args := expr.left.args
	if functionName == "println" {
		argVals := executor.toGoTypeValues(args)
		if len(argVals) < 1 {
			fmt.Println()
		}else{
			fmt.Println(argVals...)
		}
	}
	return nil
}

func (executor *ExpressionExecutor) executeBinaryExpression() (res *Value) {
	expr := executor.expr
	if expr.res != nil && expr.left.isConst() && expr.right.isConst() {
		return expr.res
	}

	switch {
	case expr.isAssign():
		res = executor.evalAssignBinaryExpression()
	case expr.isAssignAfterAdd():
		res = executor.evalAssignAfterAddBinaryExpression()
	case expr.isAssignAfterSub():
		res = executor.evalAssignAfterSubBinaryExpression()
	case expr.isAssignAfterMul():
		res = executor.evalAssignAfterMulBinaryExpression()
	case expr.isAssignAfterDiv():
		res = executor.evalAssignAfterDivBinaryExpression()
	case expr.isAssignAfterMod():
		res = executor.evalAssignAfterModBinaryExpression()

	case expr.isAdd():
		res = executor.evalAddBinaryExpression()
	case expr.isSub():
		res = executor.evalSubBinaryExpression()
	case expr.isMul():
		res = executor.evalMulBinaryExpression()
	case expr.isDiv():
		res = executor.evalDivBinaryExpression()
	case expr.isMod():
		res = executor.evalModBinaryExpression()

	case expr.isEq():
		res = executor.evalEqBinaryExpression()
	case expr.isGt():
		res = executor.evalGtBinaryExpression()
	case expr.isGe():
		res = executor.evalGeBinaryExpression()
	case expr.isLt():
		res = executor.evalLtBinaryExpression()
	case expr.isLe():
		res = executor.evalLeBinaryExpression()

	case expr.isOr():
		res = executor.evalOrBinaryExpression()
	case expr.isAnd():
		res = executor.evalAndBinaryExpression()

	}
	if expr.isTmpExpression() {
		if res == nil {
			res = NULL
		}
		varname := expr.tmpname
		executor.addVar(varname, res)

	}
	// 常量折叠
	if expr.left.isConst() && expr.right.isConst() {
		expr.res = res
	}
	return res
}

func (executor *ExpressionExecutor) evalAndBinaryExpression() (res *Value) {
	expr := executor.expr
	left := executor.leftVal()
	right := executor.rightVal()
	var tmpVal interface{}
	switch {
	case left.isBooleanValue() && right.isBooleanValue():
		tmpVal = left.bool_value && right.bool_value

	default:
		runtimeExcption("evalAndBinaryExpression Exception:", tokensString(expr.raw))
	}
	res = newVal(tmpVal)
	return res
}

func (executor *ExpressionExecutor) evalOrBinaryExpression() (res *Value) {
	expr := executor.expr
	left := executor.leftVal()
	right := executor.rightVal()
	var tmpVal interface{}
	switch {
	case left.isBooleanValue() && right.isBooleanValue():
		tmpVal = left.bool_value || right.bool_value

	default:
		runtimeExcption("evalOrBinaryExpression Exception:", tokensString(expr.raw))
	}
	res = newVal(tmpVal)
	return res
}

func (executor *ExpressionExecutor) evalEqBinaryExpression() (res *Value) {
	left := executor.leftVal()
	right := executor.rightVal()
	var tmpVal interface{}
	switch {
	case left.isBooleanValue() && right.isBooleanValue():
		tmpVal = left.bool_value == left.bool_value
	case left.isIntValue() && right.isIntValue():
		tmpVal = left.int_value == right.int_value
	case left.isFloatValue() && right.isFloatValue():
		tmpVal = left.float_value == right.float_value
	case left.isStringValue() && right.isStringValue():
		tmpVal = left.str_value == right.str_value

	case left.isFloatValue() && right.isIntValue():
		tmpVal = left.float_value == float64(right.int_value)
	case left.isIntValue() && right.isFloatValue():
		tmpVal = float64(left.int_value) == right.float_value

	default:
		tmpVal = false
	}
	res = newVal(tmpVal)
	return res
}

func (executor *ExpressionExecutor) evalGtBinaryExpression() (res *Value) {
	left := executor.leftVal()
	right := executor.rightVal()
	var tmpVal interface{}
	switch {
	case left.isIntValue() && right.isIntValue():
		tmpVal = left.int_value > right.int_value
	case left.isFloatValue() && right.isFloatValue():
		tmpVal = left.float_value > right.float_value
	case left.isStringValue() && right.isStringValue():
		tmpVal = left.str_value > right.str_value

	case left.isFloatValue() && right.isIntValue():
		tmpVal = left.float_value > float64(right.int_value)
	case left.isIntValue() && right.isFloatValue():
		tmpVal = float64(left.int_value) > right.float_value

	default:
		tmpVal = false
	}
	res = newVal(tmpVal)
	return res
}

func (executor *ExpressionExecutor) evalLtBinaryExpression() (res *Value) {
	left := executor.leftVal()
	right := executor.rightVal()
	var tmpVal interface{}
	switch {
	case left.isIntValue() && right.isIntValue():
		tmpVal = left.int_value < right.int_value
	case left.isFloatValue() && right.isFloatValue():
		tmpVal = left.float_value < right.float_value
	case left.isStringValue() && right.isStringValue():
		tmpVal = left.str_value < right.str_value

	case left.isFloatValue() && right.isIntValue():
		tmpVal = left.float_value < float64(right.int_value)
	case left.isIntValue() && right.isFloatValue():
		tmpVal = float64(left.int_value) < right.float_value

	default:
		tmpVal = false
	}
	res = newVal(tmpVal)
	return res
}

func (executor *ExpressionExecutor) evalGeBinaryExpression() (res *Value) {
	tmpVal := executor.evalGtBinaryExpression().bool_value || executor.evalEqBinaryExpression().bool_value
	res = newVal(tmpVal)
	return res
}

func (executor *ExpressionExecutor) evalLeBinaryExpression() (res *Value) {
	tmpVal := executor.evalLtBinaryExpression().bool_value || executor.evalEqBinaryExpression().bool_value
	res = newVal(tmpVal)
	return res
}

func (executor *ExpressionExecutor) evalAssignAfterAddBinaryExpression() (res *Value) {
	expr := executor.expr
	res = executor.evalAddBinaryExpression()

	varname := expr.left.name
	executor.addVar(varname, res)
	return res
}

func (executor *ExpressionExecutor) evalAssignAfterSubBinaryExpression() (res *Value) {
	expr := executor.expr
	res = executor.evalSubBinaryExpression()

	varname := expr.left.name
	executor.addVar(varname, res)
	return res
}

func (executor *ExpressionExecutor) evalAssignAfterMulBinaryExpression() (res *Value) {
	expr := executor.expr
	res = executor.evalMulBinaryExpression()

	varname := expr.left.name
	executor.addVar(varname, res)
	return res
}

func (executor *ExpressionExecutor) evalAssignAfterDivBinaryExpression() (res *Value) {
	expr := executor.expr
	res = executor.evalDivBinaryExpression()

	varname := expr.left.name
	executor.addVar(varname, res)
	return res
}

func (executor *ExpressionExecutor) evalAssignAfterModBinaryExpression() (res *Value) {
	expr := executor.expr
	res = executor.evalModBinaryExpression()

	varname := expr.left.name
	executor.addVar(varname, res)
	return res
}

func (executor *ExpressionExecutor) evalAssignBinaryExpression() (res *Value) {
	expr := executor.expr
	primaryExpr := expr.left
	res = executor.rightVal()

	varname := primaryExpr.name
	if primaryExpr.isElement() {
		arrRawVal := executor.searchArrayVariable(varname)
		argVals := executor.toGoTypeValues(primaryExpr.args)
		argRawVals := executor.getArrayIndexs(len(arrRawVal), argVals)

		arrRawVal[argRawVals[0]] = res.val()
		res = newVal(arrRawVal)
	}else if primaryExpr.isAttibute() {
		varname = primaryExpr.caller
		attrname := primaryExpr.name
		objVal := executor.searchVariable(varname)
		objRawVal := objVal.val.val()
		obj, ok := objRawVal.(map[string]interface{})
		if ok {
			obj[attrname] = res.val()
			res = newVal(obj)
		}
	}


	executor.addVar(varname, res)
	return res
}

func (executor *ExpressionExecutor) evalAddBinaryExpression() (res *Value) {
	left := executor.leftVal()
	right := executor.rightVal()
	var tmpVal interface{}
	switch {
	case left.isIntValue() && right.isIntValue():
		tmpVal = left.int_value + right.int_value

	case left.isFloatValue() && right.isFloatValue():
		tmpVal = left.float_value + right.float_value

	case left.isFloatValue() && right.isIntValue():
		tmpVal = left.float_value + float64(right.int_value)

	case left.isIntValue() && right.isFloatValue():
		tmpVal = float64(left.int_value) + right.float_value

	case left.isIntValue() && right.isStringValue():
		tmpVal = fmt.Sprintf("%v%v", left.int_value, right.str_value)
	case left.isFloatValue() && right.isStringValue():
		tmpVal = fmt.Sprintf("%v%v", left.float_value, right.str_value)
	case left.isBooleanValue() && right.isStringValue():
		tmpVal = fmt.Sprintf("%v%v", left.bool_value, right.str_value)

	case left.isStringValue() && right.isIntValue():
		tmpVal = fmt.Sprintf("%v%v", left.str_value, right.int_value)
	case left.isStringValue() && right.isFloatValue():
		tmpVal = fmt.Sprintf("%v%v", left.str_value, right.float_value)
	case left.isStringValue() && right.isBooleanValue():
		tmpVal = fmt.Sprintf("%v%v", left.str_value, right.bool_value)

	case left.isStringValue() && right.isStringValue():
		tmpVal = fmt.Sprintf("%v%v", left.str_value, right.str_value)

	default:
		runtimeExcption("unknow operation:", left.val(), "+", right.val(), " -> ", executor.expr.RawString())
	}

	res = newVal(tmpVal)
	return res
}

func (executor *ExpressionExecutor) evalSubBinaryExpression() (res *Value) {
	left := executor.leftVal()
	right := executor.rightVal()
	var tmpVal interface{}
	switch {
	case left.isIntValue() && right.isIntValue():
		tmpVal = left.int_value - right.int_value

	case left.isFloatValue() && right.isFloatValue():
		tmpVal = left.float_value - right.float_value

	case left.isFloatValue() && right.isIntValue():
		tmpVal = left.float_value - float64(right.int_value)

	case left.isIntValue() && right.isFloatValue():
		tmpVal = float64(left.int_value) - right.float_value

	default:
		runtimeExcption("unknow operation:", left.val(), "-", right.val())
	}
	res = newVal(tmpVal)
	return res
}

func (executor *ExpressionExecutor) evalMulBinaryExpression() (res *Value) {
	left := executor.leftVal()
	right := executor.rightVal()
	var tmpVal interface{}
	switch {
	case left.isIntValue() && right.isIntValue():
		tmpVal = left.int_value * right.int_value

	case left.isFloatValue() && right.isFloatValue():
		tmpVal = left.float_value * right.float_value

	case left.isFloatValue() && right.isIntValue():
		tmpVal = left.float_value * float64(right.int_value)

	case left.isIntValue() && right.isFloatValue():
		tmpVal = float64(left.int_value) * right.float_value

	default:
		runtimeExcption("unknow operation:", left.val(), "*", right.val())
	}
	res = newVal(tmpVal)
	return res
}

func (executor *ExpressionExecutor) evalDivBinaryExpression() (res *Value) {
	left := executor.leftVal()
	right := executor.rightVal()

	executor.checkDivZeroOperation(right)

	var tmpVal interface{}
	switch {
	case left.isIntValue() && right.isIntValue():
		tmpVal = left.int_value / right.int_value

	case left.isFloatValue() && right.isFloatValue():
		tmpVal = left.float_value / right.float_value

	case left.isFloatValue() && right.isIntValue():
		tmpVal = left.float_value / float64(right.int_value)

	case left.isIntValue() && right.isFloatValue():
		tmpVal = float64(left.int_value) / right.float_value

	default:
		runtimeExcption("unknow operation:", left.val(), "/", right.val())
	}
	res = newVal(tmpVal)
	return res
}

func (executor *ExpressionExecutor) evalModBinaryExpression() (res *Value) {
	left := executor.leftVal()
	right := executor.rightVal()

	executor.checkDivZeroOperation(right)

	var tmpVal interface{}
	switch {
	case left.isIntValue() && right.isIntValue():
		tmpVal = right.int_value % left.int_value

	default:
		runtimeExcption("unknow operation:", left.val(), "%", right.val())
	}
	res = newVal(tmpVal)
	return res
}

func (executor *ExpressionExecutor) checkDivZeroOperation(val *Value) {
	var flag bool
	if val.isIntValue() {
		flag = val.int_value == 0
	}
	if val.isFloatValue() {
		flag = val.float_value == 0
	}
	if flag {
		runtimeExcption("Invalid Operation: divide zero")
	}
}

func (executor *ExpressionExecutor) getArrayIndexs(arrSize int, objs []interface{}) []int {
	if objs == nil {
		runtimeExcption("array indexs is null:", fmt.Sprintln(objs...))
		return nil
	}
	var res []int
	for _, obj := range objs {
		var argRawVal int
		argRawVal, ok := obj.(int)
		if !ok {
			runtimeExcption("index type error:", obj)
			return nil
		}
		if argRawVal >= arrSize || argRawVal<0 {
			runtimeExcption("array index out of bounds:", argRawVal)
			return nil
		}
		res = append(res, argRawVal)
	}
	return res
}


func (executor *ExpressionExecutor) toGoTypeValues(exprs []*Expression) []interface{} {
	var res []interface{}
	for _, expr := range exprs {
		if expr == nil {
			continue
		}
		goValue := executor.evalNewExpression(expr)
		v := goValue.val()
		res = append(res, v)
	}
	return res
}

func (executor *ExpressionExecutor) evalNewExpression(nextExpr *Expression) *Value {
	exprExecutor := newExpressionExecutor(nextExpr, executor.stack, executor.tmpVars)
	return exprExecutor.run()
}

func (executor *ExpressionExecutor) leftVal() *Value {
	return executor.evalPrimaryExpr(executor.expr.left)
}

func (executor *ExpressionExecutor) rightVal() *Value {
	return executor.evalPrimaryExpr(executor.expr.right)
}

func (executor *ExpressionExecutor) evalPrimaryExpr(primaryExpr *PrimaryExpr) *Value {
	if primaryExpr == nil {
		return NULL
	}
	if primaryExpr.isConst() {
		return primaryExpr.res
	}
	if primaryExpr.isVar() {
		varname := primaryExpr.name
		variable := executor.searchVariable(varname)
		if variable == nil {
			return NULL
		}
		return variable.val
	}
	if primaryExpr.isElement() {
		return executor.executeElementExpression()
	}
	if primaryExpr.isAttibute() {
		return executor.executeAttributeExpression()
	}
	return NULL
}


func (executor *ExpressionExecutor) searchVariable(name string) *Variable {
	if executor.isTmpVar(name) {
		return executor.searchTmpVariable(name)
	}

	return executor.stack.searchVariable(name)
}

func (executor *ExpressionExecutor) searchTmpVariable(name string) *Variable {
	return executor.tmpVars.get(name)
}

func (executor *ExpressionExecutor) searchArrayVariable(varname string) []interface{} {
	varObj := executor.searchVariable(varname)
	varVal := varObj.val.val()

	var ok bool
	arrVal, ok := varVal.([]interface{})
	if !ok {
		varValType := fmt.Sprintf("%T", varVal)
		runtimeExcption("error operator: ", varname, "is not a array", varValType, varVal)
		return nil
	}
	return arrVal
}

//func (executor *ExpressionExecutor) addVariable(vr *Variable)  {
//	executor.stack.addLocalVariable(vr)
//}

func (executor *ExpressionExecutor) addVar(name string, val *Value)  {
	if executor.isTmpVar(name) {
		executor.addTmpVar(name, val)
		return
	}

	variable := toVar(name,  val)
	executor.stack.addLocalVariable(variable)
}

func (executor *ExpressionExecutor) addTmpVar(name string, val *Value)  {
	variable := toVar(name,  val)
	executor.tmpVars.add(variable)
}




