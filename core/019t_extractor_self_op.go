package core


// 处理自增, 自减运算 "++" "--"
type SelfOpTokenExtractor struct {
	incrFlag bool
}

func (so *SelfOpTokenExtractor) check(pre Token, cur Token, next Token, res []Token, raws []Token, curIndex int) bool {
	if cur.assertSymbol("++") {
		so.incrFlag = true
		return true
	}
	if cur.assertSymbol("--") {
		so.incrFlag = false
		return true
	}
	return false
}

func (so *SelfOpTokenExtractor) extract(pre Token, cur Token, next Token, res *[]Token, raws []Token, curIndex int) int {
	last := last(*res)
	if so.incrFlag {
		lastTokenSet(*res, newSelfIncrToken(last))
	} else {
		lastTokenSet(*res, newSelfDecrToken(last))
	}
	return curIndex + 1
}


