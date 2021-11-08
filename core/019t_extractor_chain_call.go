package core

// obj.attribute
// obj.method()
// obj.arr[]
// "." 运算符只能操作类型为Object的值
type ChainCallTokenExtractor struct {}

func (cc *ChainCallTokenExtractor) check(pre Token, cur Token, next Token, res []Token, raws []Token, curIndex int) bool {
	return cur.assertSymbol(".")
}

func (cc *ChainCallTokenExtractor) extract(pre Token, cur Token, next Token, res *[]Token, raws []Token, curIndex int) int {
	subs, nextIndex := extractChainCall(raws, cur, curIndex)

	header := last(*res)
	lastTokenSet(*res, newChainCallToken(header, subs))

	return nextIndex
}

func extractChainCall(raws []Token, cur Token, curIndex int) ([]Token, int) {
	rawsLen := len(raws)
	var subs []Token
	var tmps []Token
	for !(cur.isSymbol() && !cur.assertSymbols(".", "[", "(")) {
		if cur.assertSymbol(".") {
			curIndex ++
			cur = raws[curIndex]

			var sub Token
			if len(tmps) < 1 {
				continue
			} else if len(tmps) > 1 {
				sub = newElemFCallToken(tmps[0], tmps[1:])
			} else {
				sub = tmps[0]
			}
			subs = append(subs, sub)
			tmps = nil

			continue
		} else if cur.assertSymbol("[") {
			endIndex := scopeEndIndex(raws, curIndex, "[", "]")
			exprTokens := raws[curIndex+1 : endIndex]
			exprTokens = parse4ComplexTokens(exprTokens)

			if midIndex := nextSymbolIndex(exprTokens, 0, ":"); midIndex > -1 {
				startTokens := exprTokens[:midIndex]
				var endTokens []Token
				if midIndex == len(exprTokens)-1 {
					endTokens = []Token{}
				} else {
					endTokens = exprTokens[midIndex+1:]
				}
				tmps = append(tmps, newSubListToken(startTokens, endTokens))
			} else {
				tmps = append(tmps, newElemToken(exprTokens))
			}

			curIndex = endIndex + 1
			if curIndex >= rawsLen {break}
			cur = raws[curIndex]

		} else if cur.assertSymbol("(") {
			endIndex := scopeEndIndex(raws, curIndex, "(", ")")
			argsTokens := parse4ComplexTokens(raws[curIndex+1 : endIndex])
			tmps = append(tmps, newFunctionCallToken(argsTokens))

			curIndex = endIndex + 1
			if curIndex >= rawsLen {break}
			cur = raws[curIndex]

		} else {
			if _,ok := cur.(*NameToken); !ok {
				runtimeExcption(cur.rowIndex(), "failed to parse chain call:", tokensShow10(raws[curIndex:]))
			}
			tmps = append(tmps, cur)

			curIndex++
			if curIndex >= rawsLen {break}
			cur = raws[curIndex]
		}
	}

	if len(tmps) > 1 {
		subs = append(subs, newElemFCallToken(tmps[0], tmps[1:]))
	} else if len(tmps) == 1 {
		subs = append(subs, tmps[0])
	}

	return subs, curIndex
}
