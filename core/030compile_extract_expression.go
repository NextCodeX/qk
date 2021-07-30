package core

func extractExpression(ts []Token) Expression {
	var expr Expression

	// 去括号
	ts = clearParentheses(ts)
	tlen := len(ts)
	if tlen < 1 {
		return expr
	}
	if tlen%2 == 0 {
		runtimeExcption("extractExpression#error expression:", tokensString(ts))
	}

	switch {
	case tlen == 1:
		return parseUnaryExpression(ts)

	case tlen == 3:
		expr = parseBinaryExpression(ts)

	default:
		// 处理多元表达式
		expr = parseMultivariateExpression(ts)
	}
	if expr == nil {
		runtimeExcption("parseExpressionStatement Exception:", tokensString(ts))
	} else {
		expr.setRaw(ts)
	}
	return expr
}

// 获取一元表达式
func parseUnaryExpression(ts []Token) Expression {
	token := ts[0]
	var expr Expression
	if token.isAddSelf() || token.isSubSelf() {
		// 处理自增, 自减
		var op Token
		if token.isSubSelf() {
			op = symbolToken("-")
		} else {
			op = symbolToken("+")
		}
		var tmpTokens []Token
		tmpTokens = append(tmpTokens, token)
		tmpTokens = append(tmpTokens, op)
		tmpTokens = append(tmpTokens, newToken("1", Int))
		tmpTokens = append(tmpTokens, token)
		expr = generateBinaryExpr(tmpTokens)
	} else {
		expr = parsePrimaryExpression(token)
	}
	expr.setRaw(ts)
	return expr
}

func parseMultivariateExpression(ts []Token) Expression {
	var expr Expression
	var resVarToken Token
	var multiExprTokens []Token
	if ts[1].assertSymbol("=") {
		resVarToken = ts[0]
		multiExprTokens = ts[2:]
	} else {
		multiExprTokens = ts
	}
	var exprTokensList [][]Token
	reduceTokensForExpression(resVarToken, multiExprTokens, &exprTokensList)
	//printExprTokens(exprTokensList)

	exprs := generateBinaryExprs(exprTokensList)
	if len(exprs) == 1 {
		// 只有一个表达式时,直接返回一个二元表达式
		return exprs[0]
	}

	finalExpr := getFinalExpr(exprs, resVarToken)

	expr = &MultiExpressionImpl{
		list:      exprs,
		finalExpr: finalExpr,
	}
	return expr
}

func getFinalExpr(exprs []BinaryExpression, resVarToken Token) BinaryExpression {
	var finalExprTokens BinaryExpression
	isAssignExpr := resVarToken != nil
	for _, expr := range exprs {
		if isAssignExpr && expr.getReceiver() == resVarToken.raw() {
			finalExprTokens = expr
			break
		}
		if !isAssignExpr && expr.getReceiver() == "" {
			finalExprTokens = expr
			break
		}
	}
	if finalExprTokens == nil {
		runtimeExcption("failed to getFinalExpr resortExprTokensList")
	}
	return finalExprTokens
}

func generateBinaryExprs(exprTokensList [][]Token) []BinaryExpression {
	var res []BinaryExpression
	for _, tokens := range exprTokensList {
		expr := generateBinaryExpr(tokens)
		expr.setRaw(tokens)
		res = append(res, expr)
	}
	return res
}

// 入参由三个或四个入参组成:
// 因子, 操作符, 因子, 结果变量名(可选)
func generateBinaryExpr(ts []Token) BinaryExpression {
	size := len(ts)
	if size < 3 || size > 4 {
		runtimeExcption("generateBinaryExpr error args:", tokensString(ts))
		return nil
	}
	var expr BinaryExpression
	if size == 4 {
		expr = parseBinaryExpression(ts[:3])
		expr.setReceiver(ts[3].raw())
	} else {
		expr = parseBinaryExpression(ts)
	}
	return expr
}

// 分解多元表达式, 并把结果保存至exprTokensList *[][]TokenImpl
func reduceTokensForExpression(res Token, ts []Token, exprTokensList *[][]Token) {
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
			// e.g. (a + b) / 5
			// 左处理
			exprTokens = append(leftTokens, tmpvarToken)
			*exprTokensList = append(*exprTokensList, exprTokens)

			// 右处理
			nextTokens := insert(tmpvarToken, ts[nextIndex:])
			reduceTokensForExpression(res, nextTokens, exprTokensList)
			return
		} else {
			// e.g. (d + (f - c)) * e
			// 左处理
			reduceTokensForExpression(tmpvarToken, leftTokens, exprTokensList)

			// 右处理
			nextTokens := insert(tmpvarToken, ts[nextIndex:])
			reduceTokensForExpression(res, nextTokens, exprTokensList)
			return
		}
	}

	for i := 0; i < size; i++ {
		tmpSize := len(exprTokens)
		token := ts[i]

		condBoundry := tmpSize == 2
		condFinish := condBoundry && i == size-1
		preCond1 := condBoundry && i < size-1
		// 处理根据运算符优先级, 左向归约的情况
		condShiftLeft1 := preCond1 && last(exprTokens).equal(next(ts, i))
		condShiftLeft2 := preCond1 && last(exprTokens).lower(next(ts, i))
		if condShiftLeft1 || condShiftLeft2 || condFinish {
			// e.g. c + 7
			// a = 23
			exprTokens = append(exprTokens, token)
			if !condFinish {
				// e.g. a + 9 + c
				tmpVarToken := getTmpVarToken()
				nextTokens := insert(tmpVarToken, ts[i+1:])
				exprTokens = append(exprTokens, tmpVarToken)
				reduceTokensForExpression(res, nextTokens, exprTokensList)
			}
			break
		}

		// 处理括号不是第一个token的情况
		condRightParentheses := preCond1 && token.assertSymbol("(")
		if condRightParentheses {
			rightTokens, nextIndex := extractTokensByParentheses(ts[i:])
			nextIndex = i + nextIndex // 转换回切片ts的相应索引
			if nextIndex < size-1 {
				// e.g. a + (9 * d) - 3; rightTokens -> 9 * d; 9 * d => tmpVarToken1
				// 因为括号圈的不是右边整个表达式, 故先求括号值, 再通过运算符优先级求值
				tmpVarToken1 := getTmpVarToken() // 括号内的中间临时值
				reduceTokensForExpression(tmpVarToken1, rightTokens, exprTokensList)

				// 根据运算符优先级的不同, tmpVarToken2可能是左边表达式的或者右边表达式的中间临时值
				tmpVarToken2 := getTmpVarToken()

				nextToken := ts[nextIndex]
				if last(exprTokens).equal(nextToken) || last(exprTokens).lower(nextToken) {
					// e.g. d * tmp + c / 2
					// 左优先
					exprTokens = append(exprTokens, tmpVarToken1) // middle tmp result
					exprTokens = append(exprTokens, tmpVarToken2) // left expr result.

					nextTokens2 := insert(tmpVarToken2, ts[nextIndex:]) //
					reduceTokensForExpression(res, nextTokens2, exprTokensList)
				} else {
					// e.g. e + tmp * f
					// 右优先
					exprTokens = append(exprTokens, tmpVarToken2) // rigtht expr result.

					nextTokens3 := insert(tmpVarToken1, ts[nextIndex:])
					reduceTokensForExpression(tmpVarToken2, nextTokens3, exprTokensList)
				}
				break
			}

			// e.g. a * (b + 3)
			// 因为括号圈的是右边整个表达式时
			tmpVarToken3 := getTmpVarToken()
			exprTokens = append(exprTokens, tmpVarToken3)

			nextTokens := ts[i:]
			reduceTokensForExpression(tmpVarToken3, nextTokens, exprTokensList)
			break
		}

		// e.g.
		// a + b * c
		// a * b - 3
		// 处理根据运算符优先级, 右向归约的情况
		condShiftRight1 := preCond1 && last(exprTokens).upper(next(ts, i))
		if condShiftRight1 {
			tmpVarToken := getTmpVarToken()
			nextTokens := ts[i:]
			exprTokens = append(exprTokens, tmpVarToken)
			reduceTokensForExpression(tmpVarToken, nextTokens, exprTokensList)
			break
		}

		exprTokens = append(exprTokens, token)
	}

	if res != nil && len(exprTokens) == 3 {
		exprTokens = append(exprTokens, res)
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
				nextIndex = i + 1
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
func parseBinaryExpression(ts []Token) BinaryExpression {
	first := ts[0]
	mid := ts[1]
	third := ts[2]
	left := parsePrimaryExpression(first)
	right := parsePrimaryExpression(third)
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
	case mid.assertSymbol("!="):
		op = Opne

	case mid.assertSymbol("="):
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

	case mid.assertSymbol("||"):
		op = Opor
	case mid.assertSymbol("&&"):
		op = Opand

	default:
		runtimeExcption("parseBinaryExpression# invalid expression", tokensString(ts))
	}

	return &BinaryExpressionImpl{
		t:     op,
		left:  left,
		right: right,
	}
}

func parsePrimaryExpression(t Token) PrimaryExpression {
	v := tokenToValue(t)
	var res PrimaryExpression
	if t.isChainCall() {
		var priExprs []PrimaryExpression
		for _, tk := range t.chainTokenList() {
			priExpr := parsePrimaryExpression(tk)
			priExprs = append(priExprs, priExpr)
		}
		var headExpr PrimaryExpression
		if t.isNot() {
			// 避免类型 Not, ChainCall在解析执行时发生冲突
			t.setTyp(^Not & (^ChainCall) & t.typ())
			headExpr = parsePrimaryExpression(t)
			t.addType(Not)
		} else {
			t.setTyp((^ChainCall) & t.typ())
			headExpr = parsePrimaryExpression(t)
		}

		res = newChainCallPrimaryExpression(headExpr, priExprs)
	} else if v != nil {
		if t.isObjLiteral() {
			res = newObjectPrimaryExpression(v)
		} else  if t.isArrLiteral() {
			res = newArrayPrimaryExpression(v)
		} else if t.isDynamicStr() {
			res = newDynamicStrPrimaryExpression(v)
		} else {
			res = newConstPrimaryExpression(v)
		}
	} else if t.isElement() {
		exprs := getArgExprsFromToken(t.tokens())
		res = newElementPrimaryExpression(t.raw(), exprs)

	} else if t.isFcall() {
		exprs := getArgExprsFromToken(t.tokens())
		res = newFunctionCallPrimaryExpression(t.raw(), exprs)

	} else if t.isExpr() {
		expr := extractExpression(t.tokens())
		res = newExprPrimaryExpression(expr)

	} else if t.isIdentifier() {
		res = newVarPrimaryExpression(t.raw())
	} else {
		runtimeExcption("parsePrimaryExpression: unknown token type ->", t.String(), t.TokenTypeName())
	}
	if res != nil && t.isNot() {
		res.addType(NotPrimaryExpressionType)
		res.setNotFlag(t.notFlag())
	}
	return res
}

func getArgExprsFromToken(ts []Token) []Expression {
	var res []Expression
	size := len(ts)
	if size < 1 {
		return res
	}
	if size == 1 || !hasSymbol(ts, ",") {
		expr := extractExpression(ts)
		assert(expr == nil, "failed to parse ExpressionImpl:", tokensString(ts))
		res = append(res, expr)
		return res
	}
	var exprTokens []Token
	var nextIndex int
	for nextIndex >= 0 {
		exprTokens, nextIndex = extractExpressionTokensByComma(nextIndex, ts)
		expr := extractExpression(exprTokens)
		assert(expr == nil, "failed to parse ExpressionImpl:", tokensString(ts))
		res = append(res, expr)
	}
	return res
}

func extractExpressionTokensByComma(currentIndex int, ts []Token) (exprTokens []Token, nextIndex int) {
	size := len(ts)
	for i := currentIndex; i < size; i++ {
		token := ts[i]
		if token.assertSymbol(",") {
			nextIndex = i + 1
			break
		}
		exprTokens = append(exprTokens, token)
		if i == size-1 {
			nextIndex = -1
		}
	}
	return exprTokens, nextIndex
}


