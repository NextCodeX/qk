package core

type ArrLiteralTokenExtractor struct{}

func (al *ArrLiteralTokenExtractor) check(pre Token, cur Token, next Token, res []Token, raws []Token, curIndex int) bool {
	// 判断是否遇到了json数组字面值
	if !cur.assertSymbol("[") {
		return false
	}
	if curIndex == 0 {
		return true
	}
	if curIndex > 0 {
		last := last(res)
		if last.isSymbol() && !last.assertSymbol(")") {
			return true
		}
		if last.assertKey("return") {
			return true
		}
	}
	return false
}
func (al *ArrLiteralTokenExtractor) extract(pre Token, cur Token, next Token, res *[]Token, raws []Token, curIndex int) int {
	t, endIndex := extractArrLiteral(raws, curIndex)
	*res = append(*res, t)
	return endIndex + 1
}

func extractArrLiteral(raws []Token, curIndex int) (Token, int) {
	var elems []Token
	endIndex := scopeEndIndex(raws, curIndex, "[", "]")
	elems = raws[curIndex+1 : endIndex]

	elemsLen := len(elems)
	if elemsLen > 0 && last(elems).assertSymbol(",") {
		elems = elems[:elemsLen-1]
	}

	// 合并数组字面值内的复合token
	elems = parse4ComplexTokens(elems)
	t := newArrLiteralToken(elems)
	return t, endIndex
}
