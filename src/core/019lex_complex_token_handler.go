package core


// 该函数用于： 去掉无用的';', 合并token生成函数调用token(Fcall), 方法调用token(Mtcall)等复合token
func parse4ComplexTokens(ts []Token) []Token {
	var res []Token
	size := len(ts)
	for i:=0; i<size; {
		token := ts[i]

		var complexToken Token
		var nextIndex int

		// 处理无用分号
		if checkUselessBoundary(i, ts, res) {
			goto endCurrentIterate
		}

		// 捕获数组的字面值Token
		if checkJSONArrayLiteral(i, ts) {
			complexToken, nextIndex = extractArrayLiteral(i, ts)
			if nextIndex > i {
				res = append(res, complexToken)
				i = nextIndex
				goto nextLoop
			}
		}
		// 捕获对象的字面值Token
		if checkJSONObjectLiteral(i, ts) {
			complexToken, nextIndex = extractObjectLiteral(i, ts)
			if nextIndex > i {
				res = append(res, complexToken)
				i = nextIndex
				goto nextLoop
			}
		}

		if !token.isIdentifier() || i == size-1 {
			goto tokenCollect
		}

		// 捕获Attribute类型token
		complexToken, nextIndex = extractAttribute(i, ts)
		if nextIndex > i {
			// 捕获Mtcall类型token
			if nextIndex < size && ts[nextIndex].assertSymbol("(") {
				complexToken, nextIndex = extractMethodCall(i, ts)
			}
			res = append(res, complexToken)
			i = nextIndex
			goto nextLoop
		}

		// 捕获Fcall类型token
		complexToken, nextIndex = extractFunctionCall(i, ts)
		if nextIndex > i {
			// 标记Fdef类型token
			if nextIndex < size && ts[nextIndex].assertSymbol("{") {
				complexToken.t = Fdef | complexToken.t
			}
			res = append(res, complexToken)
			i = nextIndex
			goto nextLoop
		}

		// 捕获Element类型token
		complexToken, nextIndex = extractElement(i, ts)
		if nextIndex > i {
			res = append(res, complexToken)
			i = nextIndex
			goto nextLoop
		}

		// token 原样返回
	tokenCollect:
		res = append(res, token)

	endCurrentIterate:
		i++
	nextLoop:
	}
	return res
}

// 判断当前token是否为无用";"
func checkUselessBoundary(currentIndex int, ts []Token, resTs []Token) bool  {
	if !ts[currentIndex].assertSymbol(";") {
		return false
	}
	last, lastExist := lastToken(resTs)
	if lastExist && last.assertSymbols("{", "}", ";", "=") {
		return true
	}
	next, nextExist := nextToken(currentIndex, ts)
	if nextExist && next.assertSymbols("{") {
		return true
	}
	return false
}

// 判断是否遇到了json数组字面值
func checkJSONArrayLiteral(currentIndex int, ts []Token) bool {
	token := ts[currentIndex]
	if !token.assertSymbol("[") {
		return false
	}
	if currentIndex == 0 && last(ts).assertSymbol("]") {
		return true
	}
	pre, preExist := preToken(currentIndex, ts)
	if preExist && pre.assertSymbols("=", "(") {
		return true
	}
	if preExist && pre.assertIdentifier("return") {
		return true
	}
	return false
}

// 判断是否遇到了json对象字面值
func checkJSONObjectLiteral(currentIndex int, ts []Token) bool {
	token := ts[currentIndex]
	if !token.assertSymbol("{") {
		return false
	}
	if currentIndex == 0 && last(ts).assertSymbol("}") {
		return true
	}
	pre, preExist := preToken(currentIndex, ts)
	if preExist && pre.assertSymbols("=", "(") {
		return true
	}
	if preExist && pre.assertIdentifier("return") {
		return true
	}
	return false
}

func extractArrayLiteral(currentIndex int, ts []Token) (t Token, nextIndex int) {
	size := len(ts)
	scopeOpenCount := 1
	var elems []Token
	lineIndex := ts[currentIndex].lineIndex
	var endLineIndex int
	for i := currentIndex+1; i < size; i++ {
		token := ts[i]
		if token.assertSymbol("]") {
			scopeOpenCount --
			nextIndex = i + 1
			endLineIndex = token.lineIndex
			break
		}
		if token.assertSymbol(";") {
			continue
		}
		elems = append(elems, token)
	}
	if scopeOpenCount > 0 {
		msg := printCurrentPositionTokens(ts, currentIndex)
		runtimeExcption("extract ArrayLiteral Exception: no match final character \"]\"", msg)
	}

	// 合并数组字面值内的复合token
	elems = parse4ComplexTokens(elems)

	t = Token{
		str:    "[]",
		t:      ArrLiteral | Complex,
		ts:     elems,
		lineIndex:lineIndex,
		endLineIndex:endLineIndex,
	}
	return t, nextIndex
}

func extractObjectLiteral(currentIndex int, ts []Token) (t Token, nextIndex int) {
	size := len(ts)
	scopeOpenCount := 1
	var elems []Token
	lineIndex := ts[currentIndex].lineIndex
	var endLineIndex int
	for i := currentIndex+1; i < size; i++ {
		token := ts[i]
		if token.assertSymbol("{") {
			scopeOpenCount ++
		}
		if token.assertSymbol("}") {
			scopeOpenCount --
			if scopeOpenCount == 0 {
				nextIndex = i + 1
				endLineIndex = token.lineIndex
				break
			}
		}
		if token.assertSymbol(";") {
			continue
		}
		elems = append(elems, token)
	}
	if scopeOpenCount > 0 {
		msg := printCurrentPositionTokens(ts, currentIndex)
		runtimeExcption("extract element ObjectLiteral: no match final character \"}\"", msg)
	}

	// 合并对象字面值内的复合token
	elems = parse4ComplexTokens(elems)

	t = Token{
		str:    "{}",
		t:      ObjLiteral | Complex,
		ts:     elems,
		lineIndex:lineIndex,
		endLineIndex:endLineIndex,
	}
	return t, nextIndex
}

func extractElement(currentIndex int, ts []Token) (t Token, nextIndex int) {
	size := len(ts)
	// 检测不符合元素定义直接返回
	if size - currentIndex < 3 || !ts[currentIndex+1].assertSymbol("[") {
		return
	}
	var indexs []Token
	extractElementIndexTokens(currentIndex+2, ts, &nextIndex, &indexs)

	t = Token {
		str:    ts[currentIndex].str,
		t:      Element | Complex,
		ts:     indexs,
		lineIndex:ts[currentIndex].lineIndex,
	}
	return t, nextIndex
}

func extractElementIndexTokens(currentIndex int, ts []Token, nextIndex *int, indexs *[]Token) {
	size := len(ts)
	scopeOpenCount := 1
	for i := currentIndex; i < size; i++ {
		token := ts[i]
		if token.assertSymbol("]") {
			scopeOpenCount --
			*nextIndex = i + 1
			break
		}
		if token.isSymbol() && !match(token.str, "{", "}", ",", ";", "[", "=") {
			runtimeExcption("extract element index Exception, illegal character:"+token.str)
		}
		*indexs = append(*indexs, token)
	}
	if scopeOpenCount > 0 {
		runtimeExcption("extract element index Exception: no match final character \"]\"")
	}

	// 合并索引内的复合token
	*indexs = parse4ComplexTokens(*indexs)

	if *nextIndex < size && ts[*nextIndex].assertSymbol("[") {
		*indexs = append(*indexs, symbolToken(","))
		extractElementIndexTokens(*nextIndex+1, ts, nextIndex, indexs)
	}
}

func extractFunctionCall(currentIndex int, ts []Token) (t Token, nextIndex int) {
	size := len(ts)
	// 检测不符合函数调用定义直接返回
	if size - currentIndex < 3 || !ts[currentIndex+1].assertSymbol("(") {
		return
	}

	args, nextIndex := getCallArgsTokens(currentIndex + 2, ts)

	t = Token{
		str:    ts[currentIndex].str,
		t:      Fcall | Complex,
		ts:     args,
		lineIndex:ts[currentIndex].lineIndex,
	}
	return t, nextIndex
}

func extractMethodCall(currentIndex int, ts []Token) (t Token, nextIndex int) {
	args, nextIndex := getCallArgsTokens(currentIndex + 4, ts)

	t = Token{
		str:    ts[currentIndex+2].str,
		t:      Mtcall | Complex,
		caller: ts[currentIndex].str,
		ts:     args,
		lineIndex:ts[currentIndex].lineIndex,
	}
	return t, nextIndex
}

func getCallArgsTokens(currentIndex int, ts []Token) (args []Token, nextIndex int) {
	size := len(ts)
	scopeOpenCount := 1
	for i := currentIndex; i < size; i++ {
		token := ts[i]
		if token.assertSymbol("(") {
			scopeOpenCount ++
		}
		if token.assertSymbol(")") {
			scopeOpenCount --
			if scopeOpenCount == 0 {
				nextIndex = i + 1
				break
			}
		}
		if token.assertSymbols("{", "}", ";", "=") {
			msg := printCurrentPositionTokens(ts, i)
			runtimeExcption("extract call args Exception, illegal character:"+msg)
		}
		args = append(args, token)
	}
	if scopeOpenCount > 0 {
		runtimeExcption("extract call args Exception: no match final character \")\"", tokensString(args))
	}
	// 合并参数内的复合token
	args = parse4ComplexTokens(args)

	return args, nextIndex
}

func extractAttribute(currentIndex int, ts []Token) (t Token, nextIndex int) {
	size := len(ts)
	if  currentIndex + 2 >= size {
		return
	}
	third := ts[currentIndex+2]
	if !ts[currentIndex+1].assertSymbol(".")  || (!third.isIdentifier() && !third.isInt()) {
		return
	}

	token := Token{
		str:    ts[currentIndex+2].str,
		t:      Attribute | Complex,
		caller: ts[currentIndex].str,
		lineIndex:ts[currentIndex].lineIndex,
	}
	return token, currentIndex+3
}
