package core

type Hover struct {
	ready bool
	index int
	lastLen int
}

func initHover() *Hover {
	return &Hover{false, 0, -1}
}

// 该函数用于： 去掉无用的';', 合并token生成函数调用token(Fcall), 方法调用token(Mtcall)等复合token
func parse4ComplexTokens(ts []Token) []Token {
	var res []Token
	size := len(ts)
	var chainToken Token
	dotHover := initHover()
	notOperatorHover := initHover()
	for i:=0; i<size; {
		token := ts[i]

		var complexToken Token
		var nextIndex int

		// 处理无用分号
		if checkUselessBoundary(i, ts, res) {
			goto endCurrentIterate
		}

		// 处理非逻辑运算符
		if token.assertSymbol("!") {
			notOperatorHover.ready = true
			notOperatorHover.index ++
			if notOperatorHover.lastLen == -1 {
				notOperatorHover.lastLen = len(res)
			}
			goto endCurrentIterate
		} else {
			if notOperatorHover.ready && checkConflict(dotHover, i, ts) {
				lastIndex := len(res) - 1
				if len(res) == notOperatorHover.lastLen && token.assertSymbol("(")  {
					endIndex := scopeEndIndex(ts, i, "(", ")")
					exprTokens := ts[i+1:endIndex]
					exprTokens = parse4ComplexTokens(exprTokens)
					res = append(res, &TokenImpl{
						t: Expr | Not,
						ts: exprTokens,
						not: notOperatorHover.index % 2 == 1,
					})
					i = endIndex + 1
					notOperatorHover = initHover()
					goto nextLoop
				} else if res != nil && len(res) > notOperatorHover.lastLen {
					res[lastIndex].setNotFlag(notOperatorHover.index % 2 == 1)
					res[lastIndex].addType(Not)
					notOperatorHover = initHover()
				} else {}
			}
		}

		// 处理 "." 运算符
		if token.assertSymbol(".") && !dotHover.ready {
			if chainToken == nil {
				chainToken = last(res)
				chainToken.addType(ChainCall)
			}

			dotHover.ready = true
			dotHover.lastLen = len(res)
			goto endCurrentIterate
		}
		if dotHover.ready && dotHover.lastLen < len(res) {
			last, _ := lastToken(res)
			if chainToken != nil {
				chainToken.chainTokenListAppend(last)
			} else {
				runtimeExcption("parse4ComplexTokens: chainToken is not initialized")
			}
			res = res[:len(res)-1]
			if !token.assertSymbol(".") {
				chainToken = nil
				dotHover = initHover()
			} else {
				goto endCurrentIterate
			}
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
		//complexToken, nextIndex = extractAttribute(i, ts)
		//if nextIndex > i {
		//	// 捕获Mtcall类型token
		//	if nextIndex < size && ts[nextIndex].assertSymbol("(") {
		//		complexToken, nextIndex = extractMethodCall(i, ts)
		//	}
		//	res = append(res, complexToken)
		//	i = nextIndex
		//	goto nextLoop
		//}

		// 捕获Fcall类型token
		complexToken, nextIndex = extractFunctionCall(i, ts)
		if nextIndex > i {
			// 标记Fdef类型token
			if nextIndex < size && ts[nextIndex].assertSymbol("{") {
				complexToken.setTyp(Fdef | Complex)
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


	if dotHover.ready && dotHover.lastLen < len(res) {
		// 𥈱悬停监听有延迟， 需要在循环结束后进行数据进行收尾
		last, _ := lastToken(res)
		if chainToken != nil {
			chainToken.chainTokenListAppend(last)
		} else {
			runtimeExcption("parse4ComplexTokens: chainToken is not initialized")
		}
		res = res[:len(res)-1]
	}
	if notOperatorHover.ready && !dotHover.ready && res != nil && len(res) > notOperatorHover.lastLen {
		// 𥈱悬停监听有延迟， 需要在循环结束后进行数据进行收尾
		lastIndex := len(res) - 1
		res[lastIndex].setNotFlag(notOperatorHover.index % 2 == 1)
		res[lastIndex].addType(Not)
	}

	return res
}

func checkConflict(dotHover *Hover, currentIndex int, ts []Token) bool {
	if !dotHover.ready {
		return true
	}
	if currentIndex + 1 < len(ts) && next(ts, currentIndex).assertSymbol(".") {
		return false
	} else {
		return true
	}
}

// 判断当前token是否为无用";"
func checkUselessBoundary(currentIndex int, ts []Token, resTs []Token) bool  {
	if !ts[currentIndex].assertSymbol(";") {
		return false
	}
	last, lastExist := lastToken(resTs)
	if lastExist && last.assertSymbols("{", "}", ";", "=") {
		// 为什么symbol ";" 前面不可以有symbol "}"呢？
		// 因为如果symbol ";" 有symbol "}"的话，
		// 这个symbol "}" 只会是各种语句或函数的结束符，
		// 不会是json object字面值的结束符，因为json object字面值已经转成一个复合token
		// 所以就应该消除
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
	return ts[currentIndex].assertSymbol("[") && (currentIndex == 0 || (currentIndex > 0 && (!ts[currentIndex-1].isIdentifier() || ts[currentIndex-1].assertIdentifier("return"))))
}

// 判断是否遇到了json对象字面值
func checkJSONObjectLiteral(currentIndex int, ts []Token) bool {
	return ts[currentIndex].assertSymbol("{") && (currentIndex == 0 || (currentIndex > 0 && ((ts[currentIndex-1].isSymbol() && !ts[currentIndex-1].assertSymbol(")")) || ts[currentIndex-1].assertIdentifier("return"))))
}

func extractArrayLiteral(currentIndex int, ts []Token) (t Token, nextIndex int) {
	size := len(ts)
	scopeOpenCount := 1
	var elems []Token
	lineIndex := ts[currentIndex].getLineIndex()
	var endLineIndex int
	for i := currentIndex+1; i < size; i++ {
		token := ts[i]
		if token.assertSymbol("[") {
			scopeOpenCount ++
		}
		if token.assertSymbol("]") {
			scopeOpenCount --
		}
		if scopeOpenCount == 0 {
			nextIndex = i + 1
			endLineIndex = token.getLineIndex()
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
	elemsLen := len(elems)
	if elemsLen > 0 && last(elems).assertSymbol(",") {
		elems = elems[:elemsLen-1]
	}

	// 合并数组字面值内的复合token
	elems = parse4ComplexTokens(elems)
	t = &TokenImpl{
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
	lineIndex := ts[currentIndex].getLineIndex()
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
				endLineIndex = token.getLineIndex()
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
	elemsLen := len(elems)
	if elemsLen > 0 && last(elems).assertSymbol(",") {
		elems = elems[:elemsLen-1]
	}

	// 合并对象字面值内的复合token
	elems = parse4ComplexTokens(elems)
	t = &TokenImpl{
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

	t = &TokenImpl{
		str:    ts[currentIndex].raw(),
		t:      Element | Complex,
		ts:     indexs,
		lineIndex:ts[currentIndex].getLineIndex(),
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
		if token.isSymbol() && !match(token.raw(), "{", "}", ",", ";", "[", "=") {
			runtimeExcption("extract element index Exception, illegal character:"+token.raw())
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

	t = &TokenImpl{
		str:    ts[currentIndex].raw(),
		t:      Fcall | Complex,
		ts:     args,
		lineIndex:ts[currentIndex].getLineIndex(),
	}
	return t, nextIndex
}

func extractMethodCall(currentIndex int, ts []Token) (t Token, nextIndex int) {
	args, nextIndex := getCallArgsTokens(currentIndex + 4, ts)

	t = &TokenImpl{
		str:    ts[currentIndex+2].raw(),
		t:      Mtcall | Complex,
		//caller: ts[currentIndex].str,
		ts:     args,
		lineIndex:ts[currentIndex].getLineIndex(),
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
		if token.assertSymbols(";", "=") {
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

	token := &TokenImpl{
		str:    ts[currentIndex+2].raw(),
		t:      Attribute | Complex,
		//caller: ts[currentIndex].raw(),
		lineIndex:ts[currentIndex].getLineIndex(),
	}
	return token, currentIndex+3
}
