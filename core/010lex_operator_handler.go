package core

import "fmt"

// 提取多符号运算符(>=, <=...)
func parse4OperatorTokens(ts []Token) []Token {
	var res []Token
	for _, token := range ts {
		last, lastExist := lastToken(res)

		condArrow := func() bool {
			// ->
			return token.assertSymbol(">") && lastExist && last.assertSymbol("-")
		}

		condEqual := func() bool {
			// ==, !=, >=, <=, -=, +=, *=, /=
			return token.assertSymbol("=") && lastExist && last.assertSymbols("=", ">", "<", "+", "-", "*", "/", "%", "!")
		}

		condOr := func() bool {
			// ||
			return token.assertSymbol("|") && lastExist && last.assertSymbols("|")
		}

		condAnd := func() bool {
			// &&
			return token.assertSymbol("&") && lastExist && last.assertSymbols("&")
		}

		condAdd := func() bool {
			// ++
			return token.assertSymbol("+") && lastExist && last.assertSymbol("+")
		}

		condSub := func() bool {
			// --
			return token.assertSymbol("-") && lastExist && last.assertSymbol("-")
		}

		condNegative := func() bool {
			// -number(负数处理条件判断)
			lastSecond, lastSecondExist := lastSecondToken(res)
			return (token.isInt() || token.isFloat()) && (lastExist && last.assertSymbol("-")) && (lastSecondExist && (lastSecond.assertIdentifier("return") || lastSecond.assertSymbols("+", "-", "*", "/", "=", ",", "(", ":", "[", "->")))
		}

		if condArrow() || condEqual() || condAnd() || condOr() || condAdd() || condSub() || condNegative() {
			res = tailTokenMerge(res, token)
			continue
		}

		res = append(res, token)
	}
	return res
}

// 当前token与前一个token合并。
// 常用于提取"++", "--"这样的运算符。
// token类型取自后一个token
func tailTokenMerge(ts []Token, t Token) []Token {
	size := len(ts)
	tail := ts[size-1]

	tail.setRaw(fmt.Sprintf("%v%v", tail.raw(), t.raw()))
	tail.setTyp(t.typ())

	ts[size - 1] = tail
	return ts
}
