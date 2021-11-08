package core


type NotTokenExtractor struct {}


func (nt *NotTokenExtractor) check(pre Token, cur Token, next Token, res []Token, raws []Token, curIndex int) bool {
	return cur.assertSymbol("!")
}
func (nt *NotTokenExtractor) extract(pre Token, cur Token, next Token, res *[]Token, raws []Token, curIndex int) int {
	count := 0
	for cur.assertSymbol("!") {
		count++
		curIndex ++
		if curIndex >= len(raws) {
			runtimeExcption(cur.rowIndex(), " failed to extract not token")
		}
		cur = raws[curIndex]
	}
	flag := count % 2 == 0
	if cur.assertSymbol("(") {
		endIndex := scopeEndIndex(raws, curIndex, "(", ")")
		ts := parse4ComplexTokens(raws[curIndex+1 : endIndex])
		*res = append(*res, newNotToken(flag, ts...))
		return endIndex + 1
	}

	chain, nextIndex := extractChainCall(raws, cur, curIndex)
	if len(chain) == 1 {
		*res = append(*res, newNotToken(flag, chain[0]))
	} else if len(chain) > 1 {
		*res = append(*res, newNotToken(flag, newChainCallToken(chain[0], chain[1:])))
	} else {}
	return nextIndex
}
