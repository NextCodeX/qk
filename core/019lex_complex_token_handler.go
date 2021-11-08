package core

type TokenExtractor interface {
	check(pre Token, cur Token, next Token, res []Token, raws []Token, curIndex int) bool
	extract(pre Token, cur Token, next Token, res *[]Token, raws []Token, curIndex int) int
}

var tokenExtractorList = []TokenExtractor {
	&AnonymousFuncTokenExtractor{},
	&ArrLiteralTokenExtractor{},
	&ObjLiteralTokenExtractor{},
	&ChainCallTokenExtractor{},
	&ElemFCallTokenExtractor{},
	&NotTokenExtractor{},
	&SelfOpTokenExtractor{},
}

// 该函数用于： 去掉无用的';', 合并token生成函数调用token(Fcall), 方法调用token(Mtcall)等复合token
func parse4ComplexTokens(ts []Token) []Token {
	var res []Token
	size := len(ts)
	for i := 0; i < size; {
		var pre, cur, next Token
		if i > 0 {
			pre = ts[i-1]
		}
		if i + 1 < size {
			next = ts[i+1]
		}
		cur = ts[i]
		for _, extractor := range tokenExtractorList {
			if extractor.check(pre, cur, next, res, ts, i) {
				i = extractor.extract(pre, cur, next, &res, ts, i)
				goto nextLoop
			}
		}

		// token 原样返回
		res = append(res, cur)
		i++
		nextLoop:
	}

	return res
}

