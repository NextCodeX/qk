package core

import (
	"fmt"
	"os"
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
	}else if expr.isBinaryExpression() {
		return executor.executeBinaryExpression()
	}else if expr.isMultiExpression() {
		return executor.executeMultiExpression()
	} else {
		runtimeExcption("expression is not supported:", expr.RawString())
	}
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
		arr := varVal.jsonArr
		index := toIntValue(attrname)
		return arr.get(index)
	}

	if varVal.isObjectValue() {
		obj := varVal.jsonObj
		return obj.get(attrname)
	}

	if varVal.isClass() {
		return evalClassField(varVal.any, attrname)
	}

	runtimeExcption("eval attribute exception:", executor.expr.RawString())
	return nil
}

func (executor *ExpressionExecutor) executeElementExpression(primaryExpr *PrimaryExpr) (res *Value) {
	varname := primaryExpr.name
	varVal := executor.searchVariable(varname)

	argRawVals := executor.toGoTypeValues(primaryExpr.args)
	if varVal.isArrayValue() {
		arr := varVal.jsonArr
		index := toIntValue(argRawVals[0])
		return arr.get(index)
	} else if varVal.isObjectValue() {
		obj := varVal.jsonObj
		key := toStringValue(argRawVals[0])
		return obj.get(key)
	} else {
		runtimeExcption(fmt.Sprintf("failed to eval element %v[%v]: %v is not jsonArray or jsonObject", varname, argRawVals[0], varname))
		return nil
	}
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

	executor.calculateIfNotExist(left, exprList)
	executor.calculateIfNotExist(right, exprList)

	return executor.evalNewExpression(nextExpr)
}

// 检查相应的变量是否已计算，若未计算则进行计算
func (executor *ExpressionExecutor) calculateIfNotExist(primaryExpr *PrimaryExpr, exprList []*Expression) {
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

	if isPrint(functionName) {
		argRawVals := executor.toGoTypeValues(args)
		executePrintFunc(functionName, argRawVals)
	}

	if isModuleFunc(functionName) {
		argRawVals := executor.toGoTypeValues(args)
		return executeModuleFunc(functionName, argRawVals)
	}
	return nil
}

func (executor *ExpressionExecutor) executeCustomFunction(f *Function, args []*Value) (res *Value) {
	executor.stack.push()
	for i, paramName := range f.paramNames {
		arg := args[i]
		executor.addVar(paramName, arg)
	}
	res =  executeFunctionStatementList(f.block, executor.stack)

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
	case expr.isNe():
		res = executor.evalNeBinaryExpression()
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
	if expr.isAssignExpression() {
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
		tmpVal = left.boolean && right.boolean

	default:
		runtimeExcption("evalAndBinaryExpression Exception:", tokensString(expr.raw))
	}
	res = newQKValue(tmpVal)
	return res
}

func (executor *ExpressionExecutor) evalOrBinaryExpression() (res *Value) {
	expr := executor.expr
	left := executor.leftVal()
	right := executor.rightVal()
	var tmpVal interface{}
	switch {
	case left.isBooleanValue() && right.isBooleanValue():
		tmpVal = left.boolean || right.boolean

	default:
		runtimeExcption("evalOrBinaryExpression Exception:", tokensString(expr.raw))
	}
	res = newQKValue(tmpVal)
	return res
}

func (executor *ExpressionExecutor) evalEqBinaryExpression() (res *Value) {
	left := executor.leftVal()
	right := executor.rightVal()
	var tmpVal interface{}
	switch {
	case left.isBooleanValue() && right.isBooleanValue():
		tmpVal = left.boolean == left.boolean
	case left.isIntValue() && right.isIntValue():
		tmpVal = left.integer == right.integer
	case left.isFloatValue() && right.isFloatValue():
		tmpVal = left.decimal == right.decimal
	case left.isStringValue() && right.isStringValue():
		tmpVal = left.str == right.str

	case left.isFloatValue() && right.isIntValue():
		tmpVal = left.decimal == float64(right.integer)
	case left.isIntValue() && right.isFloatValue():
		tmpVal = float64(left.integer) == right.decimal

	default:
		tmpVal = false
	}
	res = newQKValue(tmpVal)
	return res
}

func (executor *ExpressionExecutor) evalNeBinaryExpression() (res *Value) {
	left := executor.leftVal()
	right := executor.rightVal()
	var tmpVal interface{}
	switch {
	case left.isBooleanValue() && right.isBooleanValue():
		tmpVal = left.boolean != left.boolean
	case left.isIntValue() && right.isIntValue():
		tmpVal = left.integer != right.integer
	case left.isFloatValue() && right.isFloatValue():
		tmpVal = left.decimal != right.decimal
	case left.isStringValue() && right.isStringValue():
		tmpVal = left.str != right.str

	case left.isFloatValue() && right.isIntValue():
		tmpVal = left.decimal != float64(right.integer)
	case left.isIntValue() && right.isFloatValue():
		tmpVal = float64(left.integer) != right.decimal

	default:
		tmpVal = false
	}
	res = newQKValue(tmpVal)
	return res
}

func (executor *ExpressionExecutor) evalGtBinaryExpression() (res *Value) {
	left := executor.leftVal()
	right := executor.rightVal()
	var tmpVal interface{}
	switch {
	case left.isIntValue() && right.isIntValue():
		tmpVal = left.integer > right.integer
	case left.isFloatValue() && right.isFloatValue():
		tmpVal = left.decimal > right.decimal
	case left.isStringValue() && right.isStringValue():
		tmpVal = left.str > right.str

	case left.isFloatValue() && right.isIntValue():
		tmpVal = left.decimal > float64(right.integer)
	case left.isIntValue() && right.isFloatValue():
		tmpVal = float64(left.integer) > right.decimal

	default:
		tmpVal = false
	}
	res = newQKValue(tmpVal)
	return res
}

func (executor *ExpressionExecutor) evalLtBinaryExpression() (res *Value) {
	left := executor.leftVal()
	right := executor.rightVal()
	var tmpVal interface{}
	switch {
	case left.isIntValue() && right.isIntValue():
		tmpVal = left.integer < right.integer
	case left.isFloatValue() && right.isFloatValue():
		tmpVal = left.decimal < right.decimal
	case left.isStringValue() && right.isStringValue():
		tmpVal = left.str < right.str

	case left.isFloatValue() && right.isIntValue():
		tmpVal = left.decimal < float64(right.integer)
	case left.isIntValue() && right.isFloatValue():
		tmpVal = float64(left.integer) < right.decimal

	default:
		tmpVal = false
	}
	res = newQKValue(tmpVal)
	return res
}

func (executor *ExpressionExecutor) evalGeBinaryExpression() (res *Value) {
	tmpVal := executor.evalGtBinaryExpression().boolean || executor.evalEqBinaryExpression().boolean
	res = newQKValue(tmpVal)
	return res
}

func (executor *ExpressionExecutor) evalLeBinaryExpression() (res *Value) {
	tmpVal := executor.evalLtBinaryExpression().boolean || executor.evalEqBinaryExpression().boolean
	res = newQKValue(tmpVal)
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
			arr := varVal.jsonArr
			arr.set(index, res)
			return
		}
		if varVal.isObjectValue() {
			key := toStringValue(argRawVals[0])
			obj := varVal.jsonObj
			obj.put(key, res)
			return
		}

	} else if primaryExpr.isAttibute() {
		varname = primaryExpr.caller
		attrname := primaryExpr.name
		varVal := executor.searchVariable(varname)
		if varVal.isObjectValue() {
			obj := varVal.jsonObj
			obj.put(attrname, res)
			return
		}
		if varVal.isArrayValue() {
			index := toIntValue(attrname)
			arr := varVal.jsonArr
			arr.set(index, res)
			return
		}
		if varVal.isClass() {
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
		tmpVal = left.integer + right.integer

	case left.isFloatValue() && right.isFloatValue():
		tmpVal = left.decimal + right.decimal

	case left.isFloatValue() && right.isIntValue():
		tmpVal = left.decimal + float64(right.integer)

	case left.isIntValue() && right.isFloatValue():
		tmpVal = float64(left.integer) + right.decimal

	case left.isStringValue() || right.isStringValue():
		tmpVal = fmt.Sprintf("%v%v", left.val(), right.val())

	default:
		runtimeExcption("unknow operation:", left.val(), "+", right.val(), " -> ", executor.expr.RawString())
	}

	res = newQKValue(tmpVal)
	return res
}

func (executor *ExpressionExecutor) evalSubBinaryExpression() (res *Value) {
	left := executor.leftVal()
	right := executor.rightVal()
	var tmpVal interface{}
	switch {
	case left.isIntValue() && right.isIntValue():
		tmpVal = left.integer - right.integer

	case left.isFloatValue() && right.isFloatValue():
		tmpVal = left.decimal - right.decimal

	case left.isFloatValue() && right.isIntValue():
		tmpVal = left.decimal - float64(right.integer)

	case left.isIntValue() && right.isFloatValue():
		tmpVal = float64(left.integer) - right.decimal

	default:
		runtimeExcption("unknow operation:", left.val(), "-", right.val())
	}
	res = newQKValue(tmpVal)
	return res
}

func (executor *ExpressionExecutor) evalMulBinaryExpression() (res *Value) {
	left := executor.leftVal()
	right := executor.rightVal()
	var tmpVal interface{}
	switch {
	case left.isIntValue() && right.isIntValue():
		tmpVal = left.integer * right.integer

	case left.isFloatValue() && right.isFloatValue():
		tmpVal = left.decimal * right.decimal

	case left.isFloatValue() && right.isIntValue():
		tmpVal = left.decimal * float64(right.integer)

	case left.isIntValue() && right.isFloatValue():
		tmpVal = float64(left.integer) * right.decimal

	default:
		runtimeExcption("unknow operation:", left.val(), "*", right.val())
	}
	res = newQKValue(tmpVal)
	return res
}

func (executor *ExpressionExecutor) evalDivBinaryExpression() (res *Value) {
	left := executor.leftVal()
	right := executor.rightVal()

	executor.checkDivZeroOperation(right)

	var tmpVal interface{}
	switch {
	case left.isIntValue() && right.isIntValue():
		tmpVal = left.integer / right.integer

	case left.isFloatValue() && right.isFloatValue():
		tmpVal = left.decimal / right.decimal

	case left.isFloatValue() && right.isIntValue():
		tmpVal = left.decimal / float64(right.integer)

	case left.isIntValue() && right.isFloatValue():
		tmpVal = float64(left.integer) / right.decimal

	default:
		runtimeExcption("unknow operation:", left.val(), "/", right.val())
	}
	res = newQKValue(tmpVal)
	return res
}

func (executor *ExpressionExecutor) evalModBinaryExpression() (res *Value) {
	left := executor.leftVal()
	right := executor.rightVal()

	executor.checkDivZeroOperation(right)

	var tmpVal interface{}
	switch {
	case left.isIntValue() && right.isIntValue():
		tmpVal = left.integer % right.integer

	default:
		runtimeExcption("unknow operation:", left.val(), "%", right.val())
	}
	res = newQKValue(tmpVal)
	return res
}

func (executor *ExpressionExecutor) checkDivZeroOperation(val *Value) {
	var flag bool
	if val.isIntValue() {
		flag = val.integer == 0
	}
	if val.isFloatValue() {
		flag = val.decimal == 0
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
			executor.parseJSONObject(v.jsonObj)
		} else if primaryExpr.isArray() {
			executor.parseJSONArray(v.jsonArr)
		} else if primaryExpr.isDynamicStr() {
			v = executor.parseDynamicStr(v.str)
		} else {}
		return v
	} else if primaryExpr.isVar() {
		varname := primaryExpr.name
		varVal := executor.searchVariable(varname)
		if varVal == nil {
			return NULL
		}
		return varVal
	} else if primaryExpr.isElement() {
		return executor.executeElementExpression(primaryExpr)
	} else if primaryExpr.isAttibute() {
		return executor.executeAttributeExpression(primaryExpr)
	} else if primaryExpr.isFunctionCall() {
		return executor.executeFunctionCallExpression(primaryExpr)
	} else if primaryExpr.isMethodCall() {
		return executor.executeMethodCallExpression(primaryExpr)
	} else {
		return NULL
	}
}

func (executor *ExpressionExecutor) executeMethodCallExpression(primaryExpr *PrimaryExpr) (res *Value) {
	caller := primaryExpr.caller
	methodName := primaryExpr.name
	args := primaryExpr.args

	variable := executor.stack.searchVariable(caller)
	argRawVals := executor.toGoTypeValues(args)
	if variable.isArrayValue() {
		return evalJSONArrayMethod(variable.jsonArr, methodName, argRawVals)
	}
	if variable.isObjectValue() {
		return evalJSONObjectMethod(variable.jsonObj, methodName, argRawVals)
	}
	if variable.isStringValue() {
		return evalStringMethod(variable.str, methodName, argRawVals)
	}
	if variable.isClass() {
		return evalClassMethod(variable.any, methodName, argRawVals)
	}

	return nil
}

func (executor *ExpressionExecutor) parseDynamicStr(raw string) *Value {
	res := os.Expand(raw, func(key string) string {
		qkValue := evalScript(key, executor.stack)
		return fmt.Sprint(qkValue.val())
	})
	return newQKValue(res)
}

func (executor *ExpressionExecutor) parseJSONObject(object JSONObject) {
	if object.parsed() {
		return
	}
	object.init()
	ts := clearBraces(object.tokens())
	size := len(ts)

	if size < 1 {
		return
	}

	for i:=0; i<size; i++ {

		var nextCommaIndex int
		var exprTokens []Token
		if i+2 >= size {
			runtimeExcption("error jsonobject literal:", tokensString(ts))
		} else if ts[i+2].assertSymbol("[") {
			complexToken, endIndex := extractArrayLiteral(i+2, ts)
			nextCommaIndex = endIndex+1
			exprTokens = append(exprTokens, complexToken)
		} else if ts[i+2].assertSymbol("{") {
			complexToken, endIndex := extractObjectLiteral(i+2, ts)
			nextCommaIndex = endIndex+1
			exprTokens = append(exprTokens, complexToken)
		} else {
			nextCommaIndex = nextSymbolIndex(ts, i, ",")
			if nextCommaIndex < 0 {
				nextCommaIndex = size
			}
			exprTokens = ts[i+2:nextCommaIndex]
		}

		token := ts[i]
		keyname := token.str

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

	if size < 1 {
		return
	}

	for i:=0; i<size; i++ {
		var nextCommaIndex int
		var exprTokens []Token
		if ts[i].assertSymbol("[") {
			complexToken, endIndex := extractArrayLiteral(i, ts)
			nextCommaIndex = endIndex+1
			exprTokens = append(exprTokens, complexToken)
		} else if ts[i].assertSymbol("{") {
			complexToken, endIndex := extractObjectLiteral(i, ts)
			nextCommaIndex = endIndex+1
			exprTokens = append(exprTokens, complexToken)
		} else {
			nextCommaIndex = nextSymbolIndex(ts, i, ",")
			if nextCommaIndex < 0 {
				nextCommaIndex = size
			}
			exprTokens = ts[i:nextCommaIndex]
		}

		expr := extractExpression(exprTokens)
		val := executor.evalNewExpression(expr)
		array.add(val)
		i = nextCommaIndex
	}
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




