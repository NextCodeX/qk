package core



// 该函数用于： 去掉无用的';', 合并token生成函数调用token(Fcall), 方法调用token(Mtcall)等复合token
func parse4ComplexTokens(ts []Token) []Token {
	var res []Token
	size := len(ts)
	for i:=0; i<size; {
		token := ts[i]
		pre, preExist := preToken(i, ts)
		next, nextExist := nextToken(i, ts)
		var t Token
		var nextIndex int

		// 处理无用分号
		if token.str == ";" && ((preExist && pre.assertSymbols("{","}")) || (nextExist && next.assertSymbols("{", ";"))) {
			goto end_current_iterate
		}

		// 捕获数组的字面值Token
		if token.assertSymbol("[") && preExist && pre.assertSymbols("=", "(") {
			t, nextIndex = extractArrayLiteral(i+1, ts)
			if nextIndex > i {
				res = append(res, t)
				i = nextIndex
				goto next_loop
			}
		}
		// 捕获对象的字面值Token
		if token.assertSymbol("{") && preExist && pre.assertSymbols("=", "(") {
			t, nextIndex = extractObjectLiteral(i+1, ts)
			if nextIndex > i {
				res = append(res, t)
				res = append(res, symbolToken(";"))
				i = nextIndex
				goto next_loop
			}
		}

		if !token.isIdentifier() || !nextExist {
			goto token_collect
		}

		// 捕获Attribute类型token
		t, nextIndex = extractAttribute(i, ts)
		if nextIndex > i {
			// 捕获Mtcall类型token
			if nextIndex < size && ts[nextIndex].assertSymbol("(") {
				t, nextIndex = extractMethodCall(i, ts)
			}
			res = append(res, t)
			i = nextIndex
			goto next_loop
		}

		// 捕获Fcall类型token
		t, nextIndex = extractFunctionCall(i, ts)
		if nextIndex > i {
			// 标记Fdef类型token
			if nextIndex < size && ts[nextIndex].assertSymbol("{") {
				t.t = Fdef | t.t
			}
			res = append(res, t)
			i = nextIndex
			goto next_loop
		}

		// 捕获Element类型token
		t, nextIndex = extractElement(i, ts)
		if nextIndex > i {
			res = append(res, t)
			i = nextIndex
			goto next_loop
		}

		// token 原样返回
	token_collect:
		res = append(res, token)

	end_current_iterate:
		i++
	next_loop:
	}
	return res
}

func extractArrayLiteral(currentIndex int, ts []Token) (t Token, nextIndex int) {
	size := len(ts)
	scopeOpenCount := 1
	var elems []Token
	for i := currentIndex; i < size; i++ {
		token := ts[i]
		if token.assertSymbol("]") {
			scopeOpenCount --
			nextIndex = i + 1
			break
		}
		if token.isSymbol() && !match(token.str, ",") {
			msg := printCurrentPositionTokens(ts, i)
			panic("extract ArrayLiteral Exception, illegal character:" + msg)
		}
		elems = append(elems, token)
	}
	if scopeOpenCount > 0 {
		panic("extract ArrayLiteral Exception: no match final character \"]\"")
	}
	t = Token{
		str:    "[]",
		t:      ArrLiteral | Complex,
		ts:     elems,
	}
	return t, nextIndex
}

func extractObjectLiteral(currentIndex int, ts []Token) (t Token, nextIndex int) {
	size := len(ts)
	scopeOpenCount := 1
	var elems []Token
	for i := currentIndex; i < size; i++ {
		token := ts[i]
		if token.assertSymbol("{") {
			scopeOpenCount ++
		}
		if token.assertSymbol("}") {
			scopeOpenCount --
			if scopeOpenCount == 0 {
				nextIndex = i + 1
				break
			}
		}
		if token.isSymbol() && !match(token.str,",", ":", "[", "]", "{", "}") {
			msg := printCurrentPositionTokens(ts, i)
			panic("extract element ObjectLiteral, illegal character: " + msg + " -type " + token.TokenTypeName())
		}
		elems = append(elems, token)
	}
	if scopeOpenCount > 0 {
		panic("extract element ObjectLiteral: no match final character \"}\"")
	}
	t = Token{
		str:    "{}",
		t:      ObjLiteral | Complex,
		ts:     elems,
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
			panic("extract element index Exception, illegal character:"+token.str)
		}
		*indexs = append(*indexs, token)
	}
	if scopeOpenCount > 0 {
		panic("extract element index Exception: no match final character \"]\"")
	}
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
			panic("extract call args Exception, illegal character:"+msg)
		}
		args = append(args, token)
	}
	if scopeOpenCount > 0 {
		panic("extract call args Exception: no match final character \")\"")
	}
	return args, nextIndex
}

func extractAttribute(currentIndex int, ts []Token) (t Token, nextIndex int) {
	size := len(ts)
	if size - currentIndex < 3 {
		return
	}
	if !ts[currentIndex+1].assertSymbol(".")  || !ts[currentIndex+2].isIdentifier() {
		return
	}

	token := Token{
		str:    ts[currentIndex+2].str,
		t:      Attribute | Complex,
		caller: ts[currentIndex].str,
	}
	return token, currentIndex+3
}
