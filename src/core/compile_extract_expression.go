package core

import (
	"strconv"
	"fmt"
)

// 用于统计临时变量名,或计算得到临时变量名
var tmpcount int

func getTmpVarToken() Token {
	tmpname := getTmpname()
	return varToken(tmpname)
}

func getTmpname() string {
	name := fmt.Sprintf("tmp.%v", tmpcount)
	tmpcount++
	return name
}

func extractExpression(ts []Token) *Expression {
	var expr *Expression

	// 去括号
	ts = clearParentheses(ts)

	tlen := len(ts)

	if tlen % 2 == 0 {
		runtimeExcption("error expression:", tokensString(ts))
	}

	if tlen < 1 {
		return expr
	}
	switch {
	case tlen == 1:
		return parseUnaryExpression(ts)

	case tlen == 3:
		expr = parseBinaryExpression(ts)

	default:
		// 处理多维表达式
		expr = parseMultivariateExpression(ts)
	}
	if expr == nil {
		runtimeExcption("parseExpressionStatement Exception:", tokensString(ts))
	}
	expr.raw = ts
	return expr
}

// 获取一元表达式
func parseUnaryExpression(ts []Token) *Expression {
	token := ts[0]

	// 处理自增, 自减
	if token.isAddSelf() || token.isSubSelf() {
		op := symbolToken("+")
		if token.isSubSelf() {
			op = symbolToken("-")
		}
		var tmpTokens []Token
		tmpTokens = append(tmpTokens, token)
		tmpTokens = append(tmpTokens, op)
		tmpTokens = append(tmpTokens, newToken("1", Int))
		tmpTokens = append(tmpTokens, token)
		return generateBinaryExpr(tmpTokens)
	}

	expr := &Expression{}
	primaryExpr := parsePrimaryExpression(&token)
	expr.left = primaryExpr
	if token.isObjLiteral() {
		expr.t = JSONObjectExpression | ConstExpression
		return expr
	}
	if token.isArrLiteral() {
		expr.t = JSONArrayExpression | ConstExpression
		return expr
	}
	if token.isStr() {
		expr.t = StringExpression | ConstExpression
		return expr
	}
	if token.isInt() {
		expr.t = IntExpression | ConstExpression
		return expr
	}
	if token.isFloat() {
		expr.t = FloatExpression | ConstExpression
		return expr
	}
	if token.isIdentifier() {
		if primaryExpr.isVar() {
			expr.t = VarExpression
		} else {
			expr.t = BooleanExpression | ConstExpression
		}
		return expr
	}
	if token.isElement() {
		expr.t = ElementExpression
		return expr
	}
	if token.isAttribute() {
		expr.t = AttributeExpression
		return expr
	}
	if token.isFcall() {
		expr.t = FunctionCallExpression
		return expr
	}

	return nil
}

func parseMultivariateExpression(ts []Token) (expr *Expression) {
	var resVarToken *Token
	var multiExprTokens []Token
	if ts[1].assertSymbol("=") {
		resVarToken = &(ts[0])
		multiExprTokens = ts[2:]
	} else {
		multiExprTokens = ts
	}
	var exprTokensList [][]Token
	reduceTokensForExpression(resVarToken, multiExprTokens, &exprTokensList)
	printExprTokens(exprTokensList)

	exprs := generateBinaryExprs(exprTokensList)
	if len(exprs) == 1 {
		// 只有一个表达式时,直接返回一个二元表达式
		return exprs[0]
	}

	finalExpr := getFinalExpr(exprs, resVarToken)

	expr = &Expression{
		t:MultiExpression,
		list:exprs,
		finalExpr:finalExpr,
	}
	return expr
}

func getFinalExpr(exprs []*Expression, resVarToken *Token) *Expression {
	var finalExprTokens *Expression
	isAssignExpr := resVarToken != nil
	for _, expr := range exprs {
		if isAssignExpr && expr.tmpname == resVarToken.str {
			finalExprTokens = expr
			break
		}
		if !isAssignExpr && expr.tmpname == "" {
			finalExprTokens = expr
			break
		}
	}
	if finalExprTokens == nil {
		runtimeExcption("failed to getFinalExpr resortExprTokensList")
	}
	return finalExprTokens
}

func generateBinaryExprs(exprTokensList [][]Token) []*Expression {
	var res []*Expression
	for _, tokens := range exprTokensList {
		expr := generateBinaryExpr(tokens)
		expr.raw = tokens
		res = append(res, expr)
	}
	return res
}

// 入参由三个或四个入参组成:
// 因子, 操作符, 因子, 结果变量名(可选)
func generateBinaryExpr(ts []Token) *Expression {
	size := len(ts)
	if size < 3 || size > 4 {
		runtimeExcption("generateBinaryExpr error args:", tokensString(ts))
		return nil
	}
	var expr *Expression
	if size == 4 {
		expr = parseBinaryExpression(ts[:3])
		expr.setTmpname(ts[3].str)
	} else {
		expr = parseBinaryExpression(ts)
	}
	return expr
}

// 分解多元表达式, 并把结果保存至exprTokensList *[][]Token
func reduceTokensForExpression(res *Token, ts []Token, exprTokensList *[][]Token) {
	var exprTokens []Token

	ts = clearParentheses(ts)

	size := len(ts)
	if size < 3 {
		runtimeExcption("reduceTokensForExpression Exception:", tokensString(ts))
	}

	// 处理括号是第一个token的情况
	if ts[0].assertSymbol("(") {
		leftTokens, nextIndex := extractTokensByParentheses(ts)
		tmpvarToken := getTmpVarToken()
		if !hasSymbol(leftTokens, "(") && len(leftTokens) == 3 {
			// 左处理
			exprTokens = append(leftTokens, tmpvarToken)
			*exprTokensList = append(*exprTokensList, exprTokens)

			// 右处理
			nextTokens := insert(tmpvarToken, ts[nextIndex:])
			reduceTokensForExpression(res, nextTokens, exprTokensList)
			return
		} else {
			// 左处理
			reduceTokensForExpression(&tmpvarToken, leftTokens, exprTokensList)

			// 右处理
			nextTokens := insert(tmpvarToken, ts[nextIndex:])
			reduceTokensForExpression(res, nextTokens, exprTokensList)
			return
		}
	}

	for i:=0; i<size; i++ {
		tmpSize := len(exprTokens)
		token := ts[i]

		condBoundry := tmpSize == 2
		condFinish := condBoundry && i==size-1
		preCond1 := condBoundry && i<size-1
		//if preCond1 {
		//	fmt.Printf("pre.priority:  %v:%v; next.priority: %v:%v \n", last(exprTokens).str, last(exprTokens).priority(), next(ts, i).str, next(ts, i).priority() )
		//}
		// 处理根据运算符优先级, 左向归约的情况
		condShiftLeft1 := preCond1 && last(exprTokens).equal(next(ts,i))
		condShiftLeft2 := preCond1 && last(exprTokens).lower(next(ts,i))
		//fmt.Printf("condFinish: %v, preCond1: %v, condShiftLeft1: %v, condShiftLeft2: %v\n", condFinish, preCond1, condShiftLeft1, condShiftLeft2)
		if condShiftLeft1 || condShiftLeft2 || condFinish {
			exprTokens = append(exprTokens, token)
			if !condFinish {
				tmpVarToken := getTmpVarToken()
				nextTokens := insert(tmpVarToken, ts[i+1:])
				exprTokens = append(exprTokens, tmpVarToken)
				reduceTokensForExpression(res, nextTokens, exprTokensList)
			}
			break
		}


		// 处理括号不是第一个token的情况
		condRightParentheses := preCond1 && token.assertSymbol("(")
		//fmt.Printf("condRightParentheses: %v; \n", condRightParentheses)
		if condRightParentheses {
			rightTokens, nextIndex := extractTokensByParentheses(ts[i:])
			nextIndex  = i+nextIndex // 转换回切片ts的相应索引
			if nextIndex < size - 1 {
				// 因为括号圈的不是右边整个表达式, 故先求括号值, 再通过运算符优先级求值
				tmpVarToken1 := getTmpVarToken() // 括号内的中间临时值
				reduceTokensForExpression(&tmpVarToken1, rightTokens, exprTokensList)

				// 根据运算符优先级的不同, tmpVarToken2可能是左边表达式的或者右边表达式的中间临时值
				tmpVarToken2 := getTmpVarToken()

				nextToken := ts[nextIndex]
				if last(exprTokens).equal(&nextToken) || last(exprTokens).lower(&nextToken) {
					// 左优先
					exprTokens = append(exprTokens, tmpVarToken1) // interval tmp result
					exprTokens = append(exprTokens, tmpVarToken2) // left expr result.

					nextTokens2 := insert(tmpVarToken2, ts[nextIndex:]) //
					reduceTokensForExpression(res, nextTokens2, exprTokensList)
				}else {
					// 右优先
					exprTokens = append(exprTokens, tmpVarToken2) // rigtht expr result.


					nextTokens3 := insert(tmpVarToken1, ts[nextIndex:])
					reduceTokensForExpression(&tmpVarToken2, nextTokens3, exprTokensList)
				}
				break
			}

			// 因为括号圈的是右边整个表达式时
			tmpVarToken3 := getTmpVarToken()
			exprTokens = append(exprTokens, tmpVarToken3)

			nextTokens := ts[i:]
			reduceTokensForExpression(&tmpVarToken3, nextTokens, exprTokensList)
			break
		}

		// 处理根据运算符优先级, 右向归约的情况
		condShiftRight1 := preCond1 && last(exprTokens).upper(next(ts,i))
		//fmt.Printf("condShiftRight1: %v; \n", condShiftRight1)
		if  condShiftRight1 {
			tmpVarToken := getTmpVarToken()
			nextTokens := ts[i:]
			exprTokens = append(exprTokens, tmpVarToken)
			reduceTokensForExpression(&tmpVarToken, nextTokens, exprTokensList)
			break
		}

		exprTokens = append(exprTokens, token)
	}

	if res != nil && len(exprTokens) == 3 {
		exprTokens = append(exprTokens, *res)
	}

	*exprTokensList = append(*exprTokensList, exprTokens)
}

// 获取括号内的表达式token列表
func extractTokensByParentheses(ts []Token) (res []Token, nextIndex int) {
	scopeOpen := 0
	for i, token := range ts {
		if token.assertSymbol("(") {
			scopeOpen ++
		}
		if token.assertSymbol(")") {
			scopeOpen --
			if scopeOpen == 0 {
				res = ts[1:i]
				nextIndex = i+1
				break
			}
		}
	}
	if scopeOpen != 0 {
		runtimeExcption("extractTokensByParentheses Exception:", tokensString(ts))
	}
	return res, nextIndex
}





// 根据token列表获取二元表达式
func parseBinaryExpression(ts []Token) *Expression {
	first := ts[0]
	mid := ts[1]
	third := ts[2]
	left := parsePrimaryExpression(&first)
	right := parsePrimaryExpression(&third)
	var op OperationType
	switch {

	case mid.assertSymbol("+"):
		op = Opadd
	case mid.assertSymbol("-"):
		op = Opsub
	case mid.assertSymbol("*"):
		op = Opmul
	case mid.assertSymbol("/"):
		op = Opdiv
	case mid.assertSymbol("%"):
		op = Opmod

	case mid.assertSymbol(">"):
		op = Opgt
	case mid.assertSymbol(">="):
		op = Opge
	case mid.assertSymbol("<"):
		op = Oplt
	case mid.assertSymbol("<="):
		op = Ople
	case mid.assertSymbol("=="):
		op = Opeq

	case  mid.assertSymbol("="):
		op = Opassign
	case mid.assertSymbol("+="):
		op = OpassignAfterAdd
	case mid.assertSymbol("-="):
		op = OpassignAfterSub
	case mid.assertSymbol("*="):
		op = OpassignAfterMul
	case mid.assertSymbol("/="):
		op = OpassignAfterDiv
	case mid.assertSymbol("%="):
		op = OpassignAfterMod

	default:
		runtimeExcption("parseBinaryExpression Exception:", tokensString(ts))
	}

	expr := &Expression{
		t:     BinaryExpression,
		op:    op,
		left:  left,
		right: right,
	}
	return expr
}


func parsePrimaryExpression(t *Token) *PrimaryExpr {
	v := tokenToValue(t)
	var res *PrimaryExpr
	if v != nil {
		primaryExprType := ConstPrimaryExpressionType
		if t.isObjLiteral() {
			primaryExprType = primaryExprType | ObjectPrimaryExpressionType
		}
		if t.isArrLiteral() {
			primaryExprType = primaryExprType | ArrayPrimaryExpressionType
		}
		res = &PrimaryExpr{res:v, t:primaryExprType}

	} else if t.isElement() {
		exprs := getArgExprsFromToken(t.ts)
		res = &PrimaryExpr{name:t.str, args: exprs, t:ElementPrimaryExpressionType}

	} else if t.isAttribute() {
		res = &PrimaryExpr{name:t.str, caller:t.caller, t:AttibutePrimaryExpressionType}

	} else if t.isFcall() {
		exprs := getArgExprsFromToken(t.ts)
		res = &PrimaryExpr{name:t.str, args: exprs, t:OtherPrimaryExpressionType}

	} else {
		res = &PrimaryExpr{name:t.str, t:VarPrimaryExpressionType}
	}
	return res
}

func getArgExprsFromToken(ts []Token) []*Expression {
	var res []*Expression
	size := len(ts)
	if size < 1 {
		return res
	}
	if size == 1 || !hasSymbol(ts, ",") {
		expr := extractExpression(ts)
		assert(expr==nil, "failed to parse Expression:", tokensString(ts))
		res = append(res, expr)
		return res
	}
	var exprTokens []Token
	var nextIndex int
	for nextIndex >= 0 {
		exprTokens, nextIndex = extractExpressionTokensByComma(nextIndex, ts)
		expr := extractExpression(exprTokens)
		assert(expr==nil, "failed to parse Expression:", tokensString(ts))
		res = append(res, expr)
	}
	return res
}

func extractExpressionTokensByComma(currentIndex int, ts []Token) (exprTokens []Token, nextIndex int) {
	size := len(ts)
	for i:=currentIndex; i<size; i++ {
		token := ts[i]
		if token.assertSymbol(",") {
			nextIndex = i+1
			break
		}
		exprTokens = append(exprTokens, token)
		if i == size - 1 {
			nextIndex = -1
		}
	}
	return exprTokens, nextIndex
}

func tokenToValue(t *Token) (v *Value) {
	if t.isArrLiteral() {
		v := newJSONArray(t.ts)
		return newVal(v)
	}
	if t.isObjLiteral() {
		v := newJSONObject(t.ts)
		return newVal(v)
	}
	if t.isFloat() {
		f, err := strconv.ParseFloat(t.str, 64)
		assert(err!=nil, t.String(), "line:", t.lineIndex)
		v = newVal(f)
		return
	}
	if t.isInt() {
		i, err := strconv.Atoi(t.str)
		assert(err!=nil, t.String(), "line:", t.lineIndex)
		v = newVal(i)
		return
	}
	if t.isStr() {
		v = newVal(fmt.Sprintf("%v", t.str))
		return
	}
	if t.isIdentifier() && (t.str == "true" || t.str == "false") {
		b, err := strconv.ParseBool(t.str)
		assert(err!=nil, t.String(), "line:", t.lineIndex)
		v = newVal(b)
		return
	}
	return nil
}

