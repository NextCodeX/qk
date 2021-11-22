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
		errorf("failed to extract expression from token list: %v[%v]", tokensString(ts), len(ts))
	}

	if hasSymbol(ts, "?") {
		// 解析三目运算符 ?:
		return parseTernaryOperator(ts)
	}

	switch {
	case tlen == 1:
		// 处理一元表达式
		return ts[0].toExpr()

	case tlen == 3:
		// 处理二元表达式
		expr = parseBinaryExpression(ts)

	default:
		// 处理多元表达式
		expr = parseMultivariateExpression(ts)
	}
	if expr == nil {
		runtimeExcption("failed to extract expression from token list: ", len(ts), tokensString(ts))
	} else {
		expr.setTokenList(ts)
	}
	return expr
}

// 解析三目运算符 ?:
func parseTernaryOperator(ts []Token) Expression {
	var receiver PrimaryExpression
	if len(ts) > 2 && ts[1].assertSymbol("=") {
		receiver = ts[0].toExpr()
		ts = clearParentheses(ts[2:])
	}
	var condExpr, ifExpr, elseExpr Expression
	condBoundryIndex := nextSymbolIndex(ts, 0, "?")
	valueBoundryIndex := scopeEndIndex(ts, condBoundryIndex, "?", ":")

	if condBoundryIndex < 1 || valueBoundryIndex <= condBoundryIndex+1 || valueBoundryIndex >= len(ts)-1 {
		errorf("invalid ternary operator expression: %v, %v", len(ts), tokensString(ts))
		return nil
	}

	condExpr = extractExpression(ts[:condBoundryIndex])
	ifExpr = extractExpression(ts[condBoundryIndex+1 : valueBoundryIndex])
	elseExpr = extractExpression(ts[valueBoundryIndex+1:])
	return newTernaryOperatorPrimaryExpression(condExpr, ifExpr, elseExpr, receiver)
}

func parseMultivariateExpression(ts []Token) Expression {
	var expr Expression
	var resVarToken Token
	var multiExprTokens []Token
	if ts[1].assertSymbol("=") {
		resVarToken = ts[0]
		multiExprTokens = clearParentheses(ts[2:])
	} else {
		multiExprTokens = ts
	}

	var exprTokensList [][]Token
	reduceTokensForExpression(resVarToken, multiExprTokens, &exprTokensList)
	//printExprTokens(exprTokensList)

	exprs, finalExpr := generateMulExprFactor(exprTokensList, resVarToken)
	if finalExpr == nil {
		runtimeExcption("failed to get final expression for multiExpression: ", len(ts), tokensString(ts))
	}
	if len(exprs) == 0 {
		// 只有一个表达式时,直接返回一个二元表达式
		return finalExpr
	}

	expr = &MultiExpressionImpl{list: exprs, finalExpr: finalExpr}
	return expr
}

func generateMulExprFactor(exprTokensList [][]Token, resToken Token) ([]BinaryExpression, BinaryExpression) {
	var finalExpr BinaryExpression
	var res []BinaryExpression
	exprsLen := len(exprTokensList)
	for _, tokens := range exprTokensList {
		expr := generateBinaryExpr(tokens)

		// 筛选出最后计算的二元表达式
		// tokens的长度只会是3 或 4; 长度为4时, 表示最后一个为二元表达式的结果接收Token
		if finalExpr == nil {
			// finalExpr 不添加到exprList中
			tokensLen := len(tokens)
			if resToken != nil && tokensLen == 4 && tokens[3].String() == resToken.String() {
				// 当表达式存在 赋值运算"="时, 存在resToken, 且分割出来的二元表达式的结果接收Token等于resToken
				// 说明这个二元表达式是应该最后被执行的
				finalExpr = expr
				continue
			} else if resToken == nil && !expr.isAssign() && tokensLen == 3 {
				// 当表达式不存在 赋值运算"="时, 不存在resToken, 所以分割出来的二元表达式没有结果接收Token时
				// 说明这个二元表达式是应该最后被执行的
				finalExpr = expr
				continue
			} else if exprsLen == 1 {
				finalExpr = expr
				continue
			} else {
			}
		}

		res = append(res, expr)
	}
	return res, finalExpr
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
		expr.setReceiver(ts[3].toExpr())
	} else {
		expr = parseBinaryExpression(ts)
	}

	expr.setTokenList(ts)
	return expr
}

// 分解多元表达式, 并把结果保存至exprTokensList *[][]TokenImpl
func reduceTokensForExpression(res Token, ts []Token, exprTokensList *[][]Token) {
	var exprTokens []Token
	// 处理括号是第一个token的情况
	if ts[0].assertSymbol("(") {
		endIndex := scopeEndIndex(ts, 0, "(", ")")
		tmpvarToken := getTmpVarToken()
		// e.g. (d + (f - c)) * e
		// d + (f - c) => tmpvarToken
		reduceTokensForExpression(tmpvarToken, ts[1:endIndex], exprTokensList)

		// tmpvarToken * e => res
		nextTokens := insert(tmpvarToken, ts[endIndex+1:])
		reduceTokensForExpression(res, nextTokens, exprTokensList)
		return
	}

	size := len(ts)
	if size == 1 {
		if res == nil {
			errorf("single expression(%v) is not result", tokensString(ts))
		}

		exprTokens = append(exprTokens, res)
		exprTokens = append(exprTokens, newSymbolToken("=", ts[0].rowIndex()))
		exprTokens = append(exprTokens, ts[0])
		*exprTokensList = append(*exprTokensList, exprTokens)
		return
	}

	i := 2
	token := ts[2]
	exprTokens = append(exprTokens, ts[:2]...)
	isSymbolParentheses := token.assertSymbol("(")
	if size == 3 {
		// e.g. c + 7
		exprTokens = append(exprTokens, token)
		if res != nil {
			exprTokens = append(exprTokens, res)
		}
		*exprTokensList = append(*exprTokensList, exprTokens)

	} else if !isSymbolParentheses && priorityLE(last(exprTokens), ts[i+1]) {
		// 处理根据运算符优先级, 左向归约的情况
		// e.g. a + 9 + c
		// a + 9 => tmpVarToken
		tmpVarToken := getTmpVarToken()
		exprTokens = append(exprTokens, token)
		exprTokens = append(exprTokens, tmpVarToken)
		*exprTokensList = append(*exprTokensList, exprTokens)

		// tmpVarToken + c => res
		nextTokens := insert(tmpVarToken, ts[i+1:])
		reduceTokensForExpression(res, nextTokens, exprTokensList)

	} else if !isSymbolParentheses && priorityGT(last(exprTokens), ts[i+1]) {
		// 处理根据运算符优先级, 右向归约的情况
		if i+3 == size {
			// e.g. a + b * c
			// a + tmpVarToken => res
			tmpVarToken := getTmpVarToken()
			exprTokens = append(exprTokens, tmpVarToken)
			if res != nil {
				exprTokens = append(exprTokens, res)
			}
			*exprTokensList = append(*exprTokensList, exprTokens)

			// b * c => tmpVarToken
			reduceTokensForExpression(tmpVarToken, ts[i:], exprTokensList)
		} else {
			// e.g. a + b * c - 9
			// b * c => middleResultToken
			middleResultToken := getTmpVarToken()
			middleExprTokens := make([]Token, 3)
			copy(middleExprTokens, ts[i:i+3])
			middleExprTokens = append(middleExprTokens, middleResultToken)
			*exprTokensList = append(*exprTokensList, middleExprTokens)

			// a + middleResultToken - 9 => res
			exprTokens = append(exprTokens, middleResultToken)
			exprTokens = append(exprTokens, ts[i+3:]...)
			reduceTokensForExpression(res, exprTokens, exprTokensList)
		}

	} else if isSymbolParentheses {
		endIndex := scopeEndIndex(ts, i, "(", ")")
		if endIndex == size-1 {
			// e.g. a * (b + c)
			// a * tmpVarToken => res
			tmpVarToken := getTmpVarToken()
			exprTokens = append(exprTokens, tmpVarToken)
			if res != nil {
				exprTokens = append(exprTokens, res)
			}
			*exprTokensList = append(*exprTokensList, exprTokens)

			// b + c => tmpVarToken
			reduceTokensForExpression(tmpVarToken, ts[i+1:endIndex], exprTokensList)

		} else {
			// e.g. a * (b + c / a) - 9
			// b + c / a => middleResultToken
			middleResultToken := getTmpVarToken()
			reduceTokensForExpression(middleResultToken, ts[i+1:endIndex], exprTokensList)

			// a * middleResultToken - 9 => res
			exprTokens = append(exprTokens, middleResultToken)
			exprTokens = append(exprTokens, ts[endIndex+1:]...)
			reduceTokensForExpression(res, exprTokens, exprTokensList)
		}

	} else {
		runtimeExcption("failed to split multiExpression:", len(ts), tokensString(ts))
	}
}

// 运算符优先级比较： 小于等于(用于判断左优先)
func priorityLE(left, right Token) bool {
	lp, rp := operatorPriority(left, right)
	return lp <= rp
}

// 运算符优先级比较： 大于(用于判断右优先)
func priorityGT(left, right Token) bool {
	lp, rp := operatorPriority(left, right)
	return lp > rp
}

func operatorPriority(left Token, right Token) (int, int) {
	l, ok1 := left.(*SymbolToken)
	r, ok2 := right.(*SymbolToken)
	if !ok1 || !ok2 {
		errorf("%v:%v  invalid operator: %v, %v", left.row(), right.row(), left, right)
	}
	lp, rp := l.priority(), r.priority()
	if lp == -1 || rp == -1 {
		errorf("%v:%v  invalid operator: %v, %v", left.row(), right.row(), left, right)
	}
	return lp, rp
}

// 根据token列表获取二元表达式
func parseBinaryExpression(ts []Token) BinaryExpression {
	first, mid, third := ts[0], ts[1], ts[2]
	left, right := first.toExpr(), third.toExpr()

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

	expr := &BinaryExpressionImpl{t: op, left: left, right: right}
	expr.setTokenList(ts)
	return expr
}

func getFuncParamNames(tokens []Token) []string {
	var res []string
	if len(tokens) < 1 {
		return res
	}
	for _, tk := range tokens {
		if nt, ok := tk.(*NameToken); ok {
			res = append(res, nt.name())
		}
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
