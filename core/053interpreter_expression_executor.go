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

func (executor *ExpressionExecutor) run() Value {
	expr := executor.expr
	var res Value
	if expr.isPrimaryExpression() {
		res = executor.leftVal()
	} else if expr.isBinaryExpression() {
		res = executor.executeBinaryExpression()
	} else if expr.isMultiExpression() {
		res = executor.executeMultiExpression()
	} else {
		runtimeExcption("expression is not supported:", expr.RawString())
	}
	if res == nil {
		res = NULL
	}
	return res

}

func (executor *ExpressionExecutor) executeAttributeExpression(primaryExpr *PrimaryExpr) (res Value) {
	varname := primaryExpr.caller
	attrname := primaryExpr.name
	varVal := executor.searchVariable(varname)
	if varVal == nil {
		return NULL
	}

	if varVal.isJsonArray() {
		arr := goArr(varVal)
		index := toIntValue(attrname)
		return arr.get(index)
	}

	if varVal.isJsonObject() {
		obj := goObj(varVal)
		return obj.get(attrname)
	}

	if varVal.isClass() {
		return evalClassField(goAny(varVal), attrname)
	}

	runtimeExcption("eval attribute exception:", executor.expr.RawString())
	return nil
}

func (executor *ExpressionExecutor) executeElementExpression(primaryExpr *PrimaryExpr) (res Value) {
	varname := primaryExpr.name
	varVal := executor.searchVariable(varname)

	argRawVals := executor.toGoTypeValues(primaryExpr.args)
	if varVal.isJsonArray() {
		arr := goArr(varVal)
		index := toIntValue(argRawVals[0])
		return arr.get(index)
	} else if varVal.isJsonObject() {
		obj := goObj(varVal)
		key := toStringValue(argRawVals[0])
		return obj.get(key)
	} else {
		runtimeExcption(fmt.Sprintf("failed to eval element %v[%v]: %v is not jsonArray or jsonObject", varname, argRawVals[0], varname))
		return nil
	}
}


func (executor *ExpressionExecutor) executeMultiExpression() (res Value) {
	expr := executor.expr
	return executor.recursiveEvalMultiExpression(expr.finalExpr, expr.list)
}


func (executor *ExpressionExecutor) recursiveEvalMultiExpression(nextExpr *Expression, exprList []*Expression) Value {
	left := nextExpr.left
	right := nextExpr.right
	if left.isConst() && right.isConst() {
		return executor.evalSubExpression(nextExpr)
	}

	executor.calculateIfNotExist(left, exprList)
	executor.calculateIfNotExist(right, exprList)

	return executor.evalSubExpression(nextExpr)
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
		if subExpr.receiver == varname {
			return subExpr
		}
	}
	return nil
}

func (executor *ExpressionExecutor) executeFunctionCallExpression(primaryExpr *PrimaryExpr) (res Value) {
	functionName := primaryExpr.name
	args := primaryExpr.args

	customFunc, ok := funcList[functionName]
	if ok {
		argVals := executor.evalValues(args)
		return executor.executeCustomFunction(customFunc, argVals)
	} else if isPrint(functionName) {
		argRawVals := executor.toGoTypeValues(args)
		executePrintFunc(functionName, argRawVals)
	} else if isModuleFunc(functionName) {
		argRawVals := executor.toGoTypeValues(args)
		return executeModuleFunc(functionName, argRawVals)
	} else {
		errorf("function %v() is not defined!", functionName)
	}
	return nil
}

func (executor *ExpressionExecutor) executeCustomFunction(f *Function, args []Value) (res Value) {
	executor.stack.push()
	for i, paramName := range f.paramNames {
		arg := args[i]
		executor.addVar(paramName, arg)
	}
	res =  executeFunctionStatementList(f.block, executor.stack)

	return res
}

func (executor *ExpressionExecutor) executeBinaryExpression() (res Value) {
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
	if res == nil {
		res = NULL
	}
	if expr.receiver != "" {
		varname := expr.receiver
		executor.addVar(varname, res)
	}
	// 常量折叠
	if expr.left.isConst() && expr.right.isConst() {
		expr.res = res
	}
	return res
}

func (executor *ExpressionExecutor) evalAndBinaryExpression() Value {
	left := executor.leftVal()
	if left.isBoolean() && !goBool(left) {
		return newQKValue(false)
	}
	if !left.isBoolean() && !toBoolean(left) {
		return newQKValue(false)
	}
	right := executor.rightVal()
	var res bool
	switch {
	case left.isBoolean() && right.isBoolean():
		res = goBool(left) && goBool(right)
	case !left.isBoolean() && right.isBoolean():
		res = toBoolean(left) && goBool(right)
	case left.isBoolean() && !right.isBoolean():
		res = goBool(left) && toBoolean(right)

	default:
		errorf("invalid expression: %v && %v", left.val(), right.val())
	}
	return newQKValue(res)
}

func (executor *ExpressionExecutor) evalOrBinaryExpression() Value {
	left := executor.leftVal()
	if left.isBoolean() && goBool(left) {
		return newQKValue(true)
	}
	if !left.isBoolean() && toBoolean(left) {
		return newQKValue(true)
	}
	right := executor.rightVal()
	var res bool
	switch {
	case left.isBoolean() && right.isBoolean():
		res = goBool(left) || goBool(right)
	case !left.isBoolean() && right.isBoolean():
		res = toBoolean(left) || goBool(right)
	case left.isBoolean() && !right.isBoolean():
		res = goBool(left) || toBoolean(right)

	default:
		errorf("invalid expression: %v || %v", left.val(), right.val())
	}
	return newQKValue(res)
}

func (executor *ExpressionExecutor) evalEqBinaryExpression() (res Value) {
	left := executor.leftVal()
	right := executor.rightVal()
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
		errorf("invalid expression: %v == %v", left.val(), right.val())
	}
	res = newQKValue(tmpVal)
	return res
}

func (executor *ExpressionExecutor) evalNeBinaryExpression() (res Value) {
	left := executor.leftVal()
	right := executor.rightVal()
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

func (executor *ExpressionExecutor) evalGtBinaryExpression() (res Value) {
	left := executor.leftVal()
	right := executor.rightVal()
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

func (executor *ExpressionExecutor) evalLtBinaryExpression() (res Value) {
	left := executor.leftVal()
	right := executor.rightVal()
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

func (executor *ExpressionExecutor) evalGeBinaryExpression() (res Value) {
	left := executor.leftVal()
	right := executor.rightVal()
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

func (executor *ExpressionExecutor) evalLeBinaryExpression() (res Value) {
	left := executor.leftVal()
	right := executor.rightVal()
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

func (executor *ExpressionExecutor) evalAssignAfterAddBinaryExpression() (res Value) {
	res = executor.evalAddBinaryExpression()
	executor.evalAssign(executor.expr.left, res)
	return res
}

func (executor *ExpressionExecutor) evalAssignAfterSubBinaryExpression() (res Value) {
	res = executor.evalSubBinaryExpression()
	executor.evalAssign(executor.expr.left, res)
	return res
}

func (executor *ExpressionExecutor) evalAssignAfterMulBinaryExpression() (res Value) {
	res = executor.evalMulBinaryExpression()
	executor.evalAssign(executor.expr.left, res)
	return res
}

func (executor *ExpressionExecutor) evalAssignAfterDivBinaryExpression() (res Value) {
	res = executor.evalDivBinaryExpression()
	executor.evalAssign(executor.expr.left, res)
	return res
}

func (executor *ExpressionExecutor) evalAssignAfterModBinaryExpression() (res Value) {
	res = executor.evalModBinaryExpression()
	executor.evalAssign(executor.expr.left, res)
	return res
}

func (executor *ExpressionExecutor) evalAssignBinaryExpression() Value {
	res := executor.rightVal()
	executor.evalAssign(executor.expr.left, res)
	return res
}

func (executor *ExpressionExecutor) evalAssign(priExpr *PrimaryExpr, res Value) {
	varname := priExpr.name
	if priExpr.isElement() {
		varVal := executor.searchVariable(varname)
		argRawVals := executor.toGoTypeValues(priExpr.args)
		if varVal.isJsonArray() {
			index := toIntValue(argRawVals[0])
			arr := goArr(varVal)
			arr.set(index, res)
			return
		}
		if varVal.isJsonObject() {
			key := toStringValue(argRawVals[0])
			obj := goObj(varVal)
			obj.put(key, res)
			return
		}

	} else if priExpr.isAttibute() {
		varname = priExpr.caller
		attrname := priExpr.name
		varVal := executor.searchVariable(varname)
		if varVal.isJsonObject() {
			obj := goObj(varVal)
			obj.put(attrname, res)
			return
		}
		if varVal.isJsonArray() {
			index := toIntValue(attrname)
			arr := goArr(varVal)
			arr.set(index, res)
			return
		}
		if varVal.isClass() {
			return
		}
	} else if priExpr.isVar() {
		executor.addVar(varname, res)
	} else {
		errorf("invalid assign expression")
	}
}

func (executor *ExpressionExecutor) evalAddBinaryExpression() (res Value) {
	left := executor.leftVal()
	right := executor.rightVal()
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
		runtimeExcption("unknow operation:", left.val(), "+", right.val(), " -> ", executor.expr.RawString())
	}

	res = newQKValue(tmpVal)
	return res
}

func (executor *ExpressionExecutor) evalSubBinaryExpression() (res Value) {
	left := executor.leftVal()
	right := executor.rightVal()
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

func (executor *ExpressionExecutor) evalMulBinaryExpression() (res Value) {
	left := executor.leftVal()
	right := executor.rightVal()
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

func (executor *ExpressionExecutor) evalDivBinaryExpression() (res Value) {
	left := executor.leftVal()
	right := executor.rightVal()

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

func (executor *ExpressionExecutor) evalModBinaryExpression() (res Value) {
	left := executor.leftVal()
	right := executor.rightVal()

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
		val := executor.evalSubExpression(expr)
		rawVal := val.val()
		res = append(res, rawVal)
	}
	return res
}

func (executor *ExpressionExecutor) evalValues(exprs []*Expression) []Value {
	var res []Value
	for _, expr := range exprs {
		if expr == nil {
			continue
		}
		val := executor.evalSubExpression(expr)
		res = append(res, val)
	}
	return res
}

func (executor *ExpressionExecutor) evalSubExpression(nextExpr *Expression) Value {
	exprExecutor := newExpressionExecutor(nextExpr, executor.stack, executor.tmpVars)
	return exprExecutor.run()
}

func (executor *ExpressionExecutor) leftVal() Value {
	return executor.evalPrimaryExpr(executor.expr.left)
}

func (executor *ExpressionExecutor) rightVal() Value {
	return executor.evalPrimaryExpr(executor.expr.right)
}

func (executor *ExpressionExecutor) evalPrimaryExpr(primaryExpr *PrimaryExpr) Value {
	if primaryExpr == nil {
		return NULL
	}
	var res Value

	if primaryExpr.isConst() {
		v := primaryExpr.res
		if primaryExpr.isObject() {
			executor.parseJSONObject(goObj(v))
		} else if primaryExpr.isArray() {
			executor.parseJSONArray(goArr(v))
		} else if primaryExpr.isDynamicStr() {
			v = executor.parseDynamicStr(goStr(v))
		} else {}
		res = v
	} else if primaryExpr.isExpr() {
		res = executeExpression(primaryExpr.head, executor.stack)
	} else if primaryExpr.isChainCall() {
		res = executor.evalChainCall(primaryExpr)
	} else if primaryExpr.isVar() {
		varname := primaryExpr.name
		res = executor.searchVariable(varname)
	} else if primaryExpr.isElement() {
		res =  executor.executeElementExpression(primaryExpr)
	} else if primaryExpr.isAttibute() {
		res =  executor.executeAttributeExpression(primaryExpr)
	} else if primaryExpr.isFunctionCall() {
		res =  executor.executeFunctionCallExpression(primaryExpr)
	} else if primaryExpr.isMethodCall() {
		res =  executor.executeMethodCallExpression(primaryExpr)
	} else {
		runtimeExcption("ExpressionExecutor#evalPrimaryExpr: unknown primary expression type")
	}

	if primaryExpr.isNot() { // 是否为非类型表达式
		flag := toBoolean(res)
		if primaryExpr.not { // 是否需要对表达式进行非逻辑运算
			return newQKValue(!flag)
		} else {
			return newQKValue(flag)
		}
	} else {
		if res == nil {
			res = NULL
		}
		return res
	}
}

func (executor *ExpressionExecutor) executeMethodCallExpression(primaryExpr *PrimaryExpr) (res Value) {
	caller := primaryExpr.caller
	methodName := primaryExpr.name
	args := primaryExpr.args

	variable := executor.stack.searchVariable(caller)
	argRawVals := executor.toGoTypeValues(args)
	if variable.isJsonArray() {
		return evalJSONArrayMethod(goArr(variable), methodName, argRawVals)
	}
	if variable.isJsonObject() {
		return evalJSONObjectMethod(goObj(variable), methodName, argRawVals)
	}
	if variable.isString() {
		return evalStringMethod(goStr(variable), methodName, argRawVals)
	}
	if variable.isClass() {
		return evalClassMethod(goAny(variable), methodName, argRawVals)
	}

	return nil
}

func (executor *ExpressionExecutor) parseDynamicStr(raw string) Value {
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
		keyname := token.raw()

		expr := extractExpression(exprTokens)
		val := executor.evalSubExpression(expr)
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
		val := executor.evalSubExpression(expr)
		array.add(val)
		i = nextCommaIndex
	}
}

func (executor *ExpressionExecutor) searchVariable(name string) Value {
	if isTmpVar(name) {
		return executor.searchTmpVariable(name)
	}

	return executor.stack.searchVariable(name)
}

func (executor *ExpressionExecutor) searchTmpVariable(name string) Value {
	return executor.tmpVars.get(name)
}

func (executor *ExpressionExecutor) addVar(name string, val Value)  {
	if isTmpVar(name) {
		executor.addTmpVar(name, val)
		return
	}

	executor.stack.addLocalVariable(name, val)
}

func (executor *ExpressionExecutor) addTmpVar(name string, val Value)  {
	executor.tmpVars.add(name,  val)
}
func (executor *ExpressionExecutor) evalChainCall(priExpr *PrimaryExpr) Value {
	caller := executeExpression(priExpr.head, executor.stack)

	//fmt.Println("eval chainCall#in: ", priExpr.head.left.isFunctionCall(), priExpr.isNot())
	//fmt.Println("eval chainCall#start: ", caller.val())

	for _, pri := range priExpr.chain {
		var intermediateResult Value
		if caller.isJsonArray() {
			if pri.isFunctionCall() {
				argRawVals := executor.toGoTypeValues(pri.args)
				intermediateResult = evalJSONArrayMethod(goArr(caller), pri.name, argRawVals)
			} else {}
		} else if caller.isJsonObject() {
			if pri.isVar() {
				intermediateResult = goObj(caller).get(pri.name)
			} else if pri.isFunctionCall() {
				argRawVals := executor.toGoTypeValues(pri.args)
				intermediateResult = evalJSONObjectMethod(goObj(caller), pri.name, argRawVals)
			} else if pri.isElement() {
				val := goObj(caller).get(pri.name)
				argRawVals := executor.toGoTypeValues(pri.args)
				if val.isJsonArray() {
					arr := goArr(val)
					index := toIntValue(argRawVals[0])
					intermediateResult = arr.get(index)
				} else if val.isJsonObject() {
					obj := goObj(val)
					key := toStringValue(obj)
					intermediateResult = obj.get(key)
				} else { }
			}
		} else if caller.isString() {
			if pri.isFunctionCall() {
				argRawVals := executor.toGoTypeValues(pri.args)
				intermediateResult = evalStringMethod(goStr(caller), pri.name, argRawVals)
			} else {}
		} else if caller.isClass() {
			if pri.isFunctionCall() {
				argRawVals := executor.toGoTypeValues(pri.args)
				intermediateResult = evalClassMethod(goAny(caller), pri.name, argRawVals)
			} else if pri.isVar() {
				intermediateResult = evalClassField(goAny(caller), pri.name)
			} else {}
		} else {}

		if intermediateResult == nil {
			runtimeExcption("invalid chain call expression")
		} else {
			caller = intermediateResult
			//fmt.Println("eval chainCall#intermediate: ", intermediateResult.val())
		}
	}
	return caller
}




