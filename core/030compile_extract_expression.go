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
		runtimeExcption("failed to extract expression from token list: ", len(ts), tokensString(ts))
	}

	switch {
	case tlen == 1:
		// 处理一元表达式
		return parseUnaryExpression(ts)

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
			token.setTyp(^SubSelf & token.typ())
			op = symbolToken("-")
		} else {
			token.setTyp(^AddSelf & token.typ())
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
	return expr
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
	printExprTokens(exprTokensList)

	exprs, finalExpr := generateMulExprFactor(exprTokensList, resVarToken)
	if finalExpr == nil {
		runtimeExcption("failed to get final expression for multiExpression: ", len(ts), tokensString(ts))
	}
	if len(exprs) == 0 {
		// 只有一个表达式时,直接返回一个二元表达式
		return finalExpr
	}

	expr = &MultiExpressionImpl{list:exprs, finalExpr: finalExpr}
	expr.setRaw(ts)
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
			if resToken != nil && tokensLen == 4 && tokens[3].raw() == resToken.raw() {
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
			} else {}
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
		expr.setReceiver(parsePrimaryExpression(ts[3]))
	} else {
		expr = parseBinaryExpression(ts)
	}

	expr.setRaw(ts)
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
		exprTokens = append(exprTokens, symbolToken("="))
		exprTokens = append(exprTokens, ts[0])
		*exprTokensList = append(*exprTokensList, exprTokens)
		return
	}
	for i := 0; i < size; i++ {
		token := ts[i]
		if len(exprTokens) < 2 {
			exprTokens = append(exprTokens, token)
			continue
		}

		isSymbolParentheses := token.assertSymbol("(")
		if i == size - 1 {
			// e.g. c + 7
			exprTokens = append(exprTokens, token)
			if res != nil {
				exprTokens = append(exprTokens, res)
			}
			*exprTokensList = append(*exprTokensList, exprTokens)

		} else if  !isSymbolParentheses && (last(exprTokens).equal(ts[i+1]) || last(exprTokens).lower(ts[i+1])) {
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

		} else if !isSymbolParentheses && last(exprTokens).upper(ts[i+1]) {
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
			if endIndex == size - 1 {
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
		return
	}
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

	expr := &BinaryExpressionImpl{t:op, left:left, right:right}
	expr.setRaw(ts)
	return expr
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
		tmp := t.typ()
		// 排除类型 Not: 避免类型 Not, ChainCall在解析执行时发生冲突
		// 排除类型 ChainCall: 避免无限递归
		t.setTyp(^Not & (^ChainCall) & t.typ())
		headExpr = parsePrimaryExpression(t)
		t.setTyp(tmp)

		res = newChainCallPrimaryExpression(headExpr, priExprs)

	} else if t.isElemFunctionCallMixture() {
		subTokens := t.getScopeOperatorTokens()
		if t.isIdentifier() && len(subTokens) == 1 && subTokens[0].isFuncLiteral() {
			funcToken := subTokens[0]
			funcToken.setRaw(t.raw())
			return parsePrimaryExpression(funcToken)
		}

		var priExprs []PrimaryExpression
		for _, tk := range subTokens {
			priExpr := parsePrimaryExpression(tk)
			priExprs = append(priExprs, priExpr)
		}

		var headExpr PrimaryExpression
		tmp := t.typ()
		// 排除类型 Not: 避免类型 Not, ElemFunctionCallMixture在解析执行时发生冲突
		// 排除类型 ElemFunctionCallMixture: 避免无限递归
		t.setTyp(^Not & (^ElemFunctionCallMixture) & t.typ())
		headExpr = parsePrimaryExpression(t)
		t.setTyp(tmp)

		res = newElemFunctionCallPrimaryExpression(headExpr, priExprs)

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
		expr := extractExpression(t.tokens())
		res = newElementPrimaryExpression(t.raw(), expr)

	} else if t.isFcall() {
		exprs := getArgExprsFromToken(t.tokens())
		res = newFunctionCallPrimaryExpression(t.raw(), exprs)

	} else if t.isSubList() {
		start := extractExpression(t.startExprTokens())
		end := extractExpression(t.endExprTokens())
		res = newSubListPrimaryExpression(start, end)

	} else if t.isFuncLiteral() {
		paramNames := getFuncParamNames(t.tokens())
		f := newFunc(t.raw(), t.getBodyTokens(), paramNames)
		Compile(f)
		res = newFunctionPrimaryExpression(f)

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

	res.setRaw(tokenArray(t))
	return res
}

func getFuncParamNames(tokens []Token) []string {
	var res []string
	if len(tokens) < 1 {
		return res
	}
	for _, tk := range tokens {
		if tk.isIdentifier() {
			res = append(res, tk.raw())
		} else {
			runtimeExcption("invalid function parameter name: ", tokensString(tokens))
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


