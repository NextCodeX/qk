package core

// 提取匿名函数字面值Token
type AnonymousFuncTokenExtractor struct{}

func (extractor *AnonymousFuncTokenExtractor) check(pre Token, cur Token, next Token, res []Token, raws []Token, curIndex int) bool {
	return cur.assertSymbol("$")
}

func (extractor *AnonymousFuncTokenExtractor) extract(pre Token, cur Token, next Token, res *[]Token, raws []Token, curIndex int) int {
	ts := raws
	i := curIndex

	if endIndex := nextSymbolIndexNotError(ts, i, "->", ";", ")", "]", "}"); endIndex > 0 {
		// $ [args...] -> result
		retTokenRowIndex := ts[endIndex].rowIndex()
		retToken := newKeyToken("return", retTokenRowIndex)
		argsTokens := ts[i+1 : endIndex]
		currentIndex := endIndex

		var bodyTokens []Token
		if ts[currentIndex+1].assertSymbol("{") {
			objLiteral, index := extractObjLiteral(ts, currentIndex+1)
			bodyTokens = []Token{retToken, objLiteral}
			endIndex = index + 1
			goto scopeEnd
		}
		if ts[currentIndex+1].assertSymbol("[") {
			arrLiteral, index := extractArrLiteral(ts, currentIndex+1)
			bodyTokens = []Token{retToken, arrLiteral}
			endIndex = index + 1
			goto scopeEnd
		}
		endIndex = minFuncEndIndex(ts, currentIndex)
		// 这时的endIndex已不是条件判断时的endIndex
		if endIndex < 0 {
			endIndex = len(ts)
		}

		bodyTokens = ts[currentIndex+1 : endIndex]
		bodyTokens = parse4ComplexTokens(bodyTokens)
		bodyTokens = insert(retToken, bodyTokens)

	scopeEnd:
		funcToken := newFuncLiteralToken("", argsTokens, bodyTokens)
		*res = append(*res, funcToken)
		// 不能使endIndex+1, 避免";", ",", ")", "]", "}"这些分隔符被删除
		return endIndex

	} else if endIndex := nextSymbolIndexNotError(ts, i, "{", ";"); endIndex > 0 {
		// $ [args...] {}
		argsTokens := ts[i+1 : endIndex]
		currentIndex := endIndex
		endIndex = scopeEndIndex(ts, currentIndex, "{", "}")
		bodyTokens := ts[currentIndex+1 : endIndex]
		bodyTokens = parse4ComplexTokens(bodyTokens)

		funcToken := newFuncLiteralToken("", argsTokens, bodyTokens)
		*res = append(*res, funcToken)
		return endIndex + 1

	} else if endIndex := minFuncEndIndex(ts, i); endIndex > 0 || endIndex == -1 {
		// $ expression
		var argsTokens []Token

		if endIndex == -1 {
			endIndex = len(ts)
		}
		bodyTokens := ts[i+1 : endIndex]
		bodyTokens = parse4ComplexTokens(bodyTokens)

		funcToken := newFuncLiteralToken("", argsTokens, bodyTokens)
		*res = append(*res, funcToken)
		// 不能使endIndex+1, 避免";", ",", ")", "]", "}"这些分隔符被删除
		return endIndex
	} else {
		runtimeExcption("error function literal:", tokensShow10(ts[i:]))
	}

	return -1
}
