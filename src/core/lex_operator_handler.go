package core

import "fmt"

// 提取多符号运算符(>=, <=...)
func parse4OperatorTokens(ts []Token) []Token {
	var res []Token
	for _, token := range ts {
		last, lastExist := lastToken(res)

		currentIsEqual := token.assertSymbol("=")
		condEqualMerge := lastExist && last.assertSymbols("=", ">", "<", "+", "-", "*", "/", "%")
		condEqual := currentIsEqual && condEqualMerge

		currentIsOr := token.assertSymbol("|")
		condOrMerge := lastExist && last.assertSymbols("|")
		condOr := currentIsOr && condOrMerge

		currentIsAnd := token.assertSymbol("&")
		condAndMerge := lastExist && last.assertSymbols("&")
		condAnd := currentIsAnd && condAndMerge

		currentIsAdd := token.assertSymbol("+")
		condAddMerge := lastExist && last.assertSymbols("+")
		condAdd := currentIsAdd && condAddMerge

		currentIsSub := token.assertSymbol("-")
		condSubMerge := lastExist && last.assertSymbols("-")
		condSub := currentIsSub && condSubMerge

		if condEqual || condAnd || condOr || condAdd || condSub {
			res = tailTokenMerge(res, token)
			if newTokens, ok := extractAddSubSelfToken(condAdd, condSub, res); ok {
				res = newTokens
			}
			continue
		}

		res = append(res, token)
	}
	return res
}

// 提取自增，自减token
func extractAddSubSelfToken(condAdd bool, condSub bool, ts []Token) (res []Token, ok bool) {
	size := len(ts)
	if (!condAdd && !condSub) || size<2 || !ts[size-2].isIdentifier() {
		// 判断自增，自减运算符是否能前一个token合并为自增，自减token
		return
	}

	var newTokenType TokenType
	if condAdd {
		newTokenType = AddSelf
	} else {
		newTokenType = SubSelf
	}
	tailIndex := size - 2
	tail := ts[tailIndex]
	tail.t = tail.t | newTokenType

	res = ts[:size-1]
	res[tailIndex] = tail

	return res, true
}

// 当前token与前一个token合并。
// 常用于提取"++", "--"这样的运算符。
func tailTokenMerge(ts []Token, t Token) []Token {
	size := len(ts)
	tail := ts[size-1]
	tail.str = fmt.Sprintf("%v%v", tail.str, t.str)
	ts[size - 1] = tail
	return ts
}
