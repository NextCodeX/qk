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

		// 捕获function call类型token(运算符 "()")
		// funcName()
		// funcLiteral() => funcName(){}()
		// xxx()[]
		// xxx()()
		// 运算符 "()"后面可以叠加 "[]", "()"
		// 只要前一个表达式返回的是function, 运算符 "()"皆能处理
		if res != nil && checkParenthesesOperator(res, token) {
			tailToken := res[len(res)-1]
			i = extractScopeOperator(tailToken, ts, i)
			goto nextLoop
		}

		// 捕获Element类型token(运算符 "[]")
		// var[]
		// var[:]  => 除了jsonObject或jsonArray, 还可以处理string
		// objectLiteral[]
		// arrayLiteral[]
		// func()[] --> 这种情况归运算符 "()"实现逻辑处理
		// xxx[][]。。。  --> 这种情况说明运算符 "[]"可以同类叠加
		// 运算符 "[]"后面可以叠加 "[]", "()"
		// 只要前一个表达式返回的是jsonObject或jsonArray, 运算符 "[]"皆能处理
		if res != nil && checkBracketsOperator(res, token) {
			tailToken := res[len(res)-1]
			i = extractScopeOperator(tailToken, ts, i)
			goto nextLoop
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
		// 捕获函数字面值
		if token.assertSymbol("$") {
			if endIndex := nextSymbolIndexNotError(ts, i, "{", ";"); endIndex > 0 {
				argsTokens := ts[i+1:endIndex]
				currentIndex := endIndex
				endIndex = scopeEndIndex(ts, currentIndex, "{", "}")
				bodyTokens := ts[currentIndex+1:endIndex]
				bodyTokens = parse4ComplexTokens(bodyTokens)
				funcToken := &TokenImpl{t:FuncLiteral, ts:argsTokens, bodyTokens: bodyTokens}
				res = append(res, funcToken)
				i = endIndex+1
				goto nextLoop

			} else if endIndex := nextSymbolIndexNotError(ts, i, "->", ";", ",", ")", "]", "}"); endIndex > 0 {
				argsTokens := ts[i+1:endIndex]
				currentIndex := endIndex
				endIndex = nextSymbolsIndex(ts, currentIndex, ";", ",", ")", "]", "}")
				bodyTokens := ts[currentIndex+1:endIndex]
				bodyTokens = insert(newToken("return", Identifier), bodyTokens)
				bodyTokens = parse4ComplexTokens(bodyTokens)
				funcToken := &TokenImpl{t:FuncLiteral, ts:argsTokens, bodyTokens: bodyTokens}
				res = append(res, funcToken)
				i = endIndex
				goto nextLoop

			} else {
				runtimeExcption("error function literal:", tokensShow10(ts[i:]))
			}
		}

		// 处理非逻辑运算符 "!"
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
		// obj.attribute
		// obj.method()
		// obj.arr[]
		// "." 运算符只能操作类型为Object的值
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

		// 处理自增, 自减运算符 "++" "--"
		if res != nil {
			if token.assertSymbol("++") {
				last(res).addType(AddSelf)
				goto endCurrentIterate

			} else if token.assertSymbol("--") {
				last(res).addType(SubSelf)
				goto endCurrentIterate

			} else {}
		}

		// token 原样返回
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

func extractScopeOperator(tailToken Token, ts []Token, currentIndex int) int {
	tailToken.addType(ElemFunctionCallMixture)
	for currentIndex < len(ts) {
		var nextToken Token
		if ts[currentIndex].assertSymbol("[") {
			endIndex := scopeEndIndex(ts, currentIndex, "[", "]")
			exprTokens := ts[currentIndex+1:endIndex]
			exprTokens = parse4ComplexTokens(exprTokens)

			currentIndex = endIndex + 1

			if midIndex := nextSymbolIndex(exprTokens, 0, ":"); midIndex>-1 {
				startTokens := exprTokens[:midIndex]
				var endTokens []Token
				if endIndex == len(exprTokens)-1 {
					endTokens = []Token{}
				} else {
					endTokens = exprTokens[midIndex+1:]
				}
				nextToken = &TokenImpl{t:SubList, startExpr: startTokens, endExpr: endTokens}
			} else {
				nextToken = &TokenImpl{t:Element, ts:exprTokens}
			}


		} else if ts[currentIndex].assertSymbol("(") {
			endIndex := scopeEndIndex(ts, currentIndex, "(", ")")
			argsTokens := parse4ComplexTokens(ts[currentIndex+1:endIndex])

			currentIndex = endIndex + 1

			if currentIndex < len(ts) && ts[currentIndex].assertSymbol("{") {
				// 捕获函数字面值
				endIndex = scopeEndIndex(ts, currentIndex, "{", "}")
				bodyTokens := parse4ComplexTokens(ts[currentIndex+1:endIndex])
				nextToken = &TokenImpl{t:FuncLiteral, ts:argsTokens, bodyTokens: bodyTokens}
				currentIndex = endIndex + 1

			} else {
				nextToken = &TokenImpl{t:Fcall, ts:argsTokens}
			}
		} else {
			break
		}
		tailToken.scopeOperatorTokensAppend(nextToken)
	}
	return currentIndex
}

func checkBracketsOperator(res []Token, currentToken Token) bool {
	if currentToken.assertSymbol("[") && len(res) > 0 {
		tailToken := res[len(res)-1]
		if tailToken.isIdentifier() && !match(tailToken.raw(), "for", "foreach", "fori", "forv", "if", "return", "break", "continue", "true", "false") {
			return true
		}
		if tailToken.isArrLiteral() || tailToken.isObjLiteral() || tailToken.isStr() {
			return true
		}
	}
	return false
}

func checkParenthesesOperator(res []Token, currentToken Token) bool {
	if currentToken.assertSymbol("(") && len(res) > 0 {
		tailToken := res[len(res)-1]
		if tailToken.isIdentifier() && !match(tailToken.raw(), "for", "foreach", "fori", "forv", "if", "return", "break", "continue", "true", "false") {
			return true
		}
		if tailToken.isFuncLiteral() {
			return true
		}
	}
	return false
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

	if last, lastExist := lastToken(resTs); lastExist && last.assertSymbols("{", "}", ";", "=") {
		// 为什么symbol ";" 前面不可以有symbol "}"呢？
		// 因为如果symbol ";" 有symbol "}"的话，
		// 这个symbol "}" 只会是各种语句或函数的结束符，
		// 不会是json object字面值的结束符，因为json object字面值已经转成一个复合token
		// 所以就应该消除
		return true
	}

	if next, nextExist := nextToken(currentIndex, ts); nextExist && next.assertSymbols("{") {
		return true
	}
	return false
}

// 判断是否遇到了json数组字面值
func checkJSONArrayLiteral(currentIndex int, ts []Token) bool {
	return ts[currentIndex].assertSymbol("[") && (currentIndex == 0 || (currentIndex > 0 && ((ts[currentIndex-1].isSymbol() && !ts[currentIndex-1].assertSymbol(")")) || ts[currentIndex-1].assertIdentifier("return"))))
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
	funcScope := false
	funcScopeOpenCount := 0
	var elems []Token
	lineIndex := ts[currentIndex].getLineIndex()
	var endLineIndex int
	for i := currentIndex+1; i < size; i++ {
		token := ts[i]
		if token.assertSymbol("{") {
			scopeOpenCount ++
			if funcScope {
				funcScopeOpenCount ++
			}
		}
		if token.assertSymbol("}") {
			scopeOpenCount --
			if scopeOpenCount == 0 {
				nextIndex = i + 1
				endLineIndex = token.getLineIndex()
				break
			}
			if funcScope {
				funcScopeOpenCount --
				if funcScopeOpenCount == 0 {
					funcScope = false
				}
			}
		}
		if token.assertSymbol("$") {
			funcScope = true
		}
		if token.assertSymbol(";") && !funcScope {
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

