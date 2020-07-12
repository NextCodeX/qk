package core

import "fmt"


func Interpret() {
    local := newVariables()
    executeStatementList(local)
}

func executeStatementList(local Variables) {
    for _, stmt := range mainFunc.block {
        executeStatement(stmt, newVarScope(nil, &local))
    }
}

func executeStatement(stmt *Statement, vars *VarScope) *StatementResultType {
    if stmt.isExpressionStatement() {
        for _, expr := range stmt.exprs {
            expr.vars = vars
            executeExpression(expr)
        }
    }

    return nil
}

func executeExpression(expr *Expression) (res *Value) {
    if expr.isConstExpression() || expr.isVarExpression() || expr.isElementExpression() || expr.isAttributeExpression() {
        return expr.leftVal()
    }

    if expr.isBinaryExpression() {
        return executeBinaryExpression(expr)
    }
    if expr.isFunctionCallExpression() {
        return executeFunctionCallExpression(expr)
    }
    if expr.isMultiExpression() {
    	return executeMultiExpression(expr)
	}

    return
}

func executeAttributeExpression(expr *Expression) (res *Value) {
	varname := expr.left.caller
	attrname := expr.left.name
	varVal := expr.searchVariable(varname)
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

func executeElementExpression(expr *Expression) (res *Value) {
	varname := expr.left.name
	arrRawVal := expr.searchArrayVariable(varname)

	argVals := toGoTypeValues(expr.left.args, expr.vars)
	argRawVals := getArrayIndexs(len(arrRawVal), argVals)
	return newVal(arrRawVal[argRawVals[0]])
}


func executeMultiExpression(expr *Expression) (res *Value) {
	finalExpr := expr.finalExpr
	finalExpr.vars = expr.vars

	return delayCalculate(finalExpr, expr.list)
}

func delayCalculate(expr *Expression, exprList []*Expression) *Value {
	left := expr.left
	right := expr.right
	if left.isConst() && right.isConst() {
		return executeBinaryExpression(expr)
	}

	checkValueExist(left, expr, exprList)
	checkValueExist(right, expr, exprList)

	return executeBinaryExpression(expr)
}

func checkValueExist(primaryExpr *PrimaryExpr, expr *Expression, exprList []*Expression) {
	if primaryExpr.isVar() {
		variable := expr.searchVariable(primaryExpr.name)
		if variable == nil {
			subExpr := getSubExprForMultiExpression(primaryExpr.name, exprList)
			if subExpr == nil {
				runtimeExcption("executeMultiExpression Exception")
			}
			subExpr.vars = expr.vars
			executeBinaryExpression(subExpr)
		}
	}
}

func getSubExprForMultiExpression(varname string, exprList []*Expression) *Expression {
	for _, subExpr := range exprList {
		if subExpr.tmpname == varname {
			return subExpr
		}
	}
	return nil
}

func executeFunctionCallExpression(expr *Expression) (res *Value) {
	//fmt.Println("executeFunctionCallExpression", tokensString(expr.raw))
    functionName := expr.left.name
    args := expr.left.args
    if functionName == "println" {
        argVals := toGoTypeValues(args, expr.vars)
        if len(argVals) < 1 {
            fmt.Println()
        }else{
            fmt.Println(argVals...)
        }
    }
    return nil
}

func executeBinaryExpression(expr *Expression) (res *Value) {
	if expr.res != nil && expr.left.isConst() && expr.right.isConst() {
		return expr.res
	}

    switch {
    case expr.isAssign():
        res = evalAssignBinaryExpression(expr)
	case expr.isAssignAfterAdd():
		res = evalAssignAfterAddBinaryExpression(expr)
	case expr.isAssignAfterSub():
		res = evalAssignAfterSubBinaryExpression(expr)
	case expr.isAssignAfterMul():
		res = evalAssignAfterMulBinaryExpression(expr)
	case expr.isAssignAfterDiv():
		res = evalAssignAfterDivBinaryExpression(expr)
	case expr.isAssignAfterMod():
		res = evalAssignAfterModBinaryExpression(expr)

    case expr.isAdd():
        res = evalAddBinaryExpression(expr)
	case expr.isSub():
		res = evalSubBinaryExpression(expr)
	case expr.isMul():
		res = evalMulBinaryExpression(expr)
	case expr.isDiv():
		res = evalDivBinaryExpression(expr)
	case expr.isMod():
		res = evalModBinaryExpression(expr)

	case expr.isEq():
		res = evalEqBinaryExpression(expr)
	case expr.isGt():
		res = evalGtBinaryExpression(expr)
	case expr.isGe():
		res = evalGeBinaryExpression(expr)
	case expr.isLt():
		res = evalLtBinaryExpression(expr)
	case expr.isLe():
		res = evalLeBinaryExpression(expr)

	case expr.isOr():
		res = evalOrBinaryExpression(expr)
	case expr.isAnd():
		res = evalAndBinaryExpression(expr)

    }
    if expr.isTmpExpression() {
        if res == nil {
            res = NULL
        }
        expr.addVar(expr.tmpname, res)
    }
    // 常量折叠
    if expr.left.isConst() && expr.right.isConst() {
    	expr.res = res
	}
    return res
}

func evalAndBinaryExpression(expr *Expression) (res *Value) {
	left := expr.leftVal()
	right := expr.rightVal()
	var tmpVal interface{}
	switch {
	case left.isBooleanValue() && right.isBooleanValue():
		tmpVal = left.bool_value && left.bool_value

	default:
		runtimeExcption("evalAndBinaryExpression Exception:", tokensString(expr.raw))
	}
	res = newVal(tmpVal)
	return res
}

func evalOrBinaryExpression(expr *Expression) (res *Value) {
	left := expr.leftVal()
	right := expr.rightVal()
	var tmpVal interface{}
	switch {
	case left.isBooleanValue() && right.isBooleanValue():
		tmpVal = left.bool_value || left.bool_value

	default:
		runtimeExcption("evalOrBinaryExpression Exception:", tokensString(expr.raw))
	}
	res = newVal(tmpVal)
	return res
}

func evalEqBinaryExpression(expr *Expression) (res *Value) {
	left := expr.leftVal()
	right := expr.rightVal()
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

func evalGtBinaryExpression(expr *Expression) (res *Value) {
	left := expr.leftVal()
	right := expr.rightVal()
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

func evalLtBinaryExpression(expr *Expression) (res *Value) {
	left := expr.leftVal()
	right := expr.rightVal()
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

func evalGeBinaryExpression(expr *Expression) (res *Value) {
	tmpVal := evalGtBinaryExpression(expr).bool_value || evalEqBinaryExpression(expr).bool_value
	res = newVal(tmpVal)
	return res
}

func evalLeBinaryExpression(expr *Expression) (res *Value) {
	tmpVal := evalLtBinaryExpression(expr).bool_value || evalEqBinaryExpression(expr).bool_value
	res = newVal(tmpVal)
	return res
}

func evalAssignAfterAddBinaryExpression(expr *Expression) (res *Value) {
	res = evalAddBinaryExpression(expr)

	varname := expr.left.name
	expr.addVar(varname, res)
	return res
}

func evalAssignAfterSubBinaryExpression(expr *Expression) (res *Value) {
	res = evalSubBinaryExpression(expr)

	varname := expr.left.name
	expr.addVar(varname, res)
	return res
}

func evalAssignAfterMulBinaryExpression(expr *Expression) (res *Value) {
	res = evalMulBinaryExpression(expr)

	varname := expr.left.name
	expr.addVar(varname, res)
	return res
}

func evalAssignAfterDivBinaryExpression(expr *Expression) (res *Value) {
	res = evalDivBinaryExpression(expr)

	varname := expr.left.name
	expr.addVar(varname, res)
	return res
}

func evalAssignAfterModBinaryExpression(expr *Expression) (res *Value) {
	res = evalModBinaryExpression(expr)

	varname := expr.left.name
	expr.addVar(varname, res)
	return res
}

func evalAssignBinaryExpression(expr *Expression) (res *Value) {
	primaryExpr := expr.left
	res = expr.rightVal()

	varname := primaryExpr.name
	if primaryExpr.isElement() {
		arrRawVal := expr.searchArrayVariable(varname)
		argVals := toGoTypeValues(primaryExpr.args, expr.vars)
		argRawVals := getArrayIndexs(len(arrRawVal), argVals)

		arrRawVal[argRawVals[0]] = res.val()
		res = newVal(arrRawVal)
	}else if primaryExpr.isAttibute() {
		varname = primaryExpr.caller
		attrname := primaryExpr.name
		objVal := expr.searchVariable(varname)
		objRawVal := objVal.val.val()
		obj, ok := objRawVal.(map[string]interface{})
		if ok {
			obj[attrname] = res.val()
			res = newVal(obj)
		}
	}


	expr.addVar(varname, res)
	return res
}

func evalAddBinaryExpression(expr *Expression) (res *Value) {
    left := expr.leftVal()
    right := expr.rightVal()
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
		runtimeExcption("unknow operation:", left.val(), "+", right.val())
    }

	res = newVal(tmpVal)
    return res
}

func evalSubBinaryExpression(expr *Expression) (res *Value) {
	left := expr.leftVal()
	right := expr.rightVal()
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

func evalMulBinaryExpression(expr *Expression) (res *Value) {
	left := expr.leftVal()
	right := expr.rightVal()
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

func evalDivBinaryExpression(expr *Expression) (res *Value) {
	left := expr.leftVal()
	right := expr.rightVal()

	checkDivZeroOperation(right)

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

func evalModBinaryExpression(expr *Expression) (res *Value) {
	left := expr.leftVal()
	right := expr.rightVal()

	checkDivZeroOperation(right)

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

func checkDivZeroOperation(val *Value) {
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

func getArrayIndexs(arrSize int, objs []interface{}) []int {
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


func toGoTypeValues(exprs []*Expression, vars *VarScope) []interface{} {
    var res []interface{}
    for _, expr := range exprs {
    	//fmt.Println("expr:",expr==nil, expr,", vars:", vars, "< toGoTypeValues")
        if expr == nil {
        	continue
		}
    	expr.vars = vars
        goValue := executeExpression(expr)
        v := goValue.val()
        res = append(res, v)
    }
    return res
}

