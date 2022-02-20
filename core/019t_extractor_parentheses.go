package core

type ParenthesesTokenExtractor struct{}

func (nt *ParenthesesTokenExtractor) check(pre Token, cur Token, next Token, res []Token, raws []Token, curIndex int) bool {
	return cur.assertSymbol("(")
}
func (nt *ParenthesesTokenExtractor) extract(pre Token, cur Token, next Token, res *[]Token, raws []Token, curIndex int) int {
	endIndex := scopeEndIndex(raws, curIndex, "(", ")")
	*res = append(*res, newExprToken(raws[curIndex+1:endIndex]))
	return endIndex + 1
}
