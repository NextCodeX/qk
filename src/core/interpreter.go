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
    if expr.isConstExpression() || expr.isVarExpression() {
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

func executeMultiExpression(expr *Expression) (res *Value) {
	for _, subExpr := range expr.list {
		subExpr.vars = expr.vars
		res = executeBinaryExpression(subExpr)
		fmt.Println("executeMultiExpression: ", res.int_value)
	}
	return res
}

func executeFunctionCallExpression(expr *Expression) (res *Value) {
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
        varname := expr.left.name
        res = expr.right.res
        expr.addVar(varname, res)

    case expr.isAdd():
        res = evalAddBinaryExpression(expr)
	case expr.isSub():
		res = evalSubBinaryExpression(expr)
	case expr.isMul():
		res = evalMulBinaryExpression(expr)
	case expr.isDiv():
		res = evalDivBinaryExpression(expr)

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


func toGoTypeValues(exprs []*Expression, vars *VarScope) []interface{} {
    var res []interface{}
    for _, expr := range exprs {
        expr.vars = vars
        goValue := executeExpression(expr)
        v := goValue.val()
        res = append(res, v)
    }
    return res
}

