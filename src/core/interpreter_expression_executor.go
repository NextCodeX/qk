package core

import (
	"fmt"
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
	if expr.isPrimaryExpression() {
		return executor.leftVal()
	}

	if expr.isBinaryExpression() {
		return executor.executeBinaryExpression()
	}
	if expr.isMultiExpression() {
		return executor.executeMultiExpression()
	}
	runtimeExcption("expression is not supported:", expr.RawString())
	return nil
}

func (executor *ExpressionExecutor) executeAttributeExpression(primaryExpr *PrimaryExpr) (res *Value) {
	varname := primaryExpr.caller
	attrname := primaryExpr.name
	varVal := executor.searchVariable(varname)
	if varVal == nil {
		return NULL
	}

	if varVal.isArrayValue() {
		arr := varVal.arr_value
		index := toIntValue(attrname)
		return arr.get(index)
	}

	if varVal.isObjectValue() {
		obj := varVal.obj_value
		return obj.get(attrname)
	}

	runtimeExcption("eval attribute exception:", executor.expr.RawString())
	return nil
}

func (executor *ExpressionExecutor) executeElementExpression(primaryExpr *PrimaryExpr) (res *Value) {
	varname := primaryExpr.name
	varVal := executor.searchVariable(varname)

	argRawVals := executor.toGoTypeValues(primaryExpr.args)
	if varVal.isArrayValue() {
		arr := varVal.arr_value
		index := toIntValue(argRawVals[0])
		return arr.get(index)
	}

	if varVal.isObjectValue() {
		obj := varVal.obj_value
		key := toStringValue(argRawVals[0])
		return obj.get(key)
	}

	runtimeExcption("eval element exception:", executor.expr.RawString())
	return nil
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

func (executor *ExpressionExecutor) getNextExprForMultiExpression(varname string, exprList []*Expression) *Expression {
	for _, subExpr := range exprList {
		if subExpr.tmpname == varname {
			return subExpr
		}
	}
	return nil
}

func (executor *ExpressionExecutor) executeFunctionCallExpression(primaryExpr *PrimaryExpr) (res *Value) {
	functionName := primaryExpr.name
	args := primaryExpr.args

	customFunc, ok := funcList[functionName]
	if ok {
		argVals := executor.evalValues(args)
		return executor.executeCustomFunction(customFunc, argVals)
	}

	if functionName == "println" {
		argRawVals := executor.toGoTypeValues(args)
		if len(argRawVals) < 1 {
			fmt.Println()
		}else{
			fmt.Println(argRawVals...)
		}
	}
	return nil
}

func (executor *ExpressionExecutor) executeCustomFunction(f *Function, args []*Value) (res *Value) {
	executor.stack.push()
	defer executor.stack.pop()
	for i, paramName := range f.paramNames {
		arg := args[i]
		executor.addVar(paramName, arg)
	}
	executeStatementList(f.block, executor.stack)
	//executor.stack.printVars()
	res =  executor.searchVariable(funcResultName)

	return res
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
		varVal := executor.searchVariable(varname)
		argRawVals := executor.toGoTypeValues(primaryExpr.args)
		if varVal.isArrayValue() {
			index := toIntValue(argRawVals[0])
			arr := varVal.arr_value
			arr.set(index, res)
			return
		}
		if varVal.isObjectValue() {
			key := toStringValue(argRawVals[0])
			obj := varVal.obj_value
			obj.put(key, res)
			return
		}

	}else if primaryExpr.isAttibute() {
		varname = primaryExpr.caller
		attrname := primaryExpr.name
		varVal := executor.searchVariable(varname)
		if varVal.isObjectValue() {
			obj := varVal.obj_value
			obj.put(attrname, res)
			return
		}
		if varVal.isArrayValue() {
			index := toIntValue(attrname)
			arr := varVal.arr_value
			arr.set(index, res)
			return
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
		val := executor.evalNewExpression(expr)
		rawVal := val.val()
		res = append(res, rawVal)
	}
	return res
}

func (executor *ExpressionExecutor) evalValues(exprs []*Expression) []*Value {
	var res []*Value
	for _, expr := range exprs {
		if expr == nil {
			continue
		}
		val := executor.evalNewExpression(expr)
		res = append(res, val)
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
		v := primaryExpr.res
		if primaryExpr.isObject() {
			executor.parseJSONObject(v.obj_value)
		}
		if primaryExpr.isArray() {
			executor.parseJSONArray(v.arr_value)
		}
		return v
	}
	if primaryExpr.isVar() {
		varname := primaryExpr.name
		varVal := executor.searchVariable(varname)
		if varVal == nil {
			return NULL
		}
		return varVal
	}
	if primaryExpr.isElement() {
		return executor.executeElementExpression(primaryExpr)
	}
	if primaryExpr.isAttibute() {
		return executor.executeAttributeExpression(primaryExpr)
	}
	if primaryExpr.isFunctionCall() {
		return executor.executeFunctionCallExpression(primaryExpr)
	}

	return NULL
}

func (executor *ExpressionExecutor) parseJSONObject(object JSONObject) {
	if object.parsed() {
		return
	}
	object.init()
	ts := clearBraces(object.tokens())
	size := len(ts)
	for i:=0; i<size; i++ {
		token := ts[i]
		nextCommaIndex := nextSymbolIndex(ts, i, ",")
		if nextCommaIndex < 0 {
			nextCommaIndex = size
		}
		keyname := token.str
		exprTokens := ts[i+2:nextCommaIndex]
		expr := extractExpression(exprTokens)
		val := executor.evalNewExpression(expr)
		object.put(keyname, val)
		i = nextCommaIndex
	}
}

func (executor *ExpressionExecutor) parseJSONArray(array JSONArray) {
	if array.parsed() {
		return
	}
	ts := clearBrackets(array.tokens())
	size := len(ts)
	for i:=0; i<size; i++ {
		nextCommaIndex := nextSymbolIndex(ts, i, ",")
		if nextCommaIndex < 0 {
			nextCommaIndex = size
		}
		exprTokens := ts[i:nextCommaIndex]
		expr := extractExpression(exprTokens)
		val := executor.evalNewExpression(expr)
		array.add(val)
		i = nextCommaIndex
	}
	array.setParsed()
}

func (executor *ExpressionExecutor) searchVariable(name string) *Value {
	if isTmpVar(name) {
		return executor.searchTmpVariable(name)
	}

	return executor.stack.searchVariable(name)
}

func (executor *ExpressionExecutor) searchTmpVariable(name string) *Value {
	return executor.tmpVars.get(name)
}

func (executor *ExpressionExecutor) addVar(name string, val *Value)  {
	if isTmpVar(name) {
		executor.addTmpVar(name, val)
		return
	}

	executor.stack.addLocalVariable(name, val)
}

func (executor *ExpressionExecutor) addTmpVar(name string, val *Value)  {
	executor.tmpVars.add(name,  val)
}




