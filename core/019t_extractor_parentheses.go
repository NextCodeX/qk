package core

type ParenthesesTokenExtractor struct{}

func (nt *ParenthesesTokenExtractor) check(pre Token, cur Token, next Token, res []Token, raws []Token, curIndex int) bool {
	return cur.assertSymbol("(")
}
func (nt *ParenthesesTokenExtractor) extract(pre Token, cur Token, next Token, res *[]Token, raws []Token, curIndex int) int {
	endIndex := scopeEndIndex(raws, curIndex, "(", ")")
	tmp := parse4ComplexTokens(raws[curIndex+1 : endIndex])
	*res = append(*res, newExprToken(tmp))
	return endIndex + 1
}
