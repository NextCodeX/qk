package core

type ObjLiteralTokenExtractor struct{}

func (ol *ObjLiteralTokenExtractor) check(pre Token, cur Token, next Token, res []Token, raws []Token, curIndex int) bool {
	// 判断是否遇到了json对象字面值
	if !cur.assertSymbol("{") {
		return false
	}
	if curIndex == 0 {
		return true
	}
	if curIndex > 0 {
		last := last(res)
		if last.assertKey("return") {
			return true
		}
		if last.isSymbol() && !last.assertSymbol(")") {
			return true
		}
	}
	return false
}

func (ol *ObjLiteralTokenExtractor) extract(pre Token, cur Token, next Token, res *[]Token, raws []Token, curIndex int) int {
	t, endIndex := extractObjLiteral(raws, curIndex)
	*res = append(*res, t)
	return endIndex + 1
}

func extractObjLiteral(raws []Token, curIndex int) (Token, int) {
	endIndex := scopeEndIndex(raws, curIndex, "{", "}")
	elems := raws[curIndex+1 : endIndex]
	elemsLen := len(elems)
	if elemsLen > 0 && last(elems).assertSymbol(",") {
		elems = elems[:elemsLen-1]
	}

	// 合并对象字面值内的复合token
	elems = parse4ComplexTokens(elems)
	t := newObjLiteralToken(elems)
	return t, endIndex
}
