package core

// 捕获function call类型token(运算符 "()")
// funcName()
// funcLiteral() => funcName(){}()
// xxx()[]
// xxx()()
// 运算符 "()"后面可以叠加 "[]", "()"
// 只要前一个表达式返回的是function, 运算符 "()"皆能处理
// -------------------------------------
// 捕获Element类型token(运算符 "[]")
// var[]
// var[:]  => 除了jsonObject或jsonArray, 还可以处理string
// objectLiteral[]
// arrayLiteral[]
// func()[] --> 这种情况归运算符 "()"实现逻辑处理
// xxx[][]。。。  --> 这种情况说明运算符 "[]"可以同类叠加
// 运算符 "[]"后面可以叠加 "[]", "()"
// 只要前一个表达式返回的是jsonObject或jsonArray, 运算符 "[]"皆能处理
type ElemFCallTokenExtractor struct {}


func (extractor *ElemFCallTokenExtractor) check(pre Token, cur Token, next Token, res []Token, raws []Token, curIndex int) bool {
	size := len(res)
	if size < 1 {
		return false
	}
	tailToken := res[size-1]
	if cur.assertSymbol("(") {
		switch tailToken.(type) {
		case *NameToken, *FuncLiteralToken:
			return true
		}
	}
	if cur.assertSymbol("[") {
		switch tailToken.(type) {
		case *NameToken, *ArrLiteralToken, *ObjLiteralToken, *StrToken:
			return true
		}
	}

	return false
}



func (extractor *ElemFCallTokenExtractor) extract(pre Token, cur Token, next Token, res *[]Token, raws []Token, curIndex int) int {
	lastSecond, lastSecondExist := lastSecondToken(*res)
	header := last(*res)
	var subs []Token
	for curIndex < len(raws) {
		var nextToken Token
		if raws[curIndex].assertSymbol("[") {
			endIndex := scopeEndIndex(raws, curIndex, "[", "]")
			exprTokens := raws[curIndex+1 : endIndex]
			exprTokens = parse4ComplexTokens(exprTokens)

			curIndex = endIndex + 1

			if midIndex := nextSymbolIndex(exprTokens, 0, ":"); midIndex > -1 {
				startTokens := exprTokens[:midIndex]
				var endTokens []Token
				if midIndex == len(exprTokens)-1 {
					endTokens = []Token{}
				} else {
					endTokens = exprTokens[midIndex+1:]
				}
				nextToken = newSubListToken(startTokens, endTokens)
			} else {
				nextToken = newElemToken(exprTokens)
			}

		} else if raws[curIndex].assertSymbol("(") {
			endIndex := scopeEndIndex(raws, curIndex, "(", ")")
			argsTokens := parse4ComplexTokens(raws[curIndex+1 : endIndex])

			curIndex = endIndex + 1

			// 排除 if chk() {}； elif chk() {}； for chk() {};
			// forv v : list() {};
			if (!lastSecondExist ||
				!lastSecond.assertKeys("if", "elif", "for") &&
				!lastSecond.assertSymbol(":")) &&
				curIndex < len(raws) && raws[curIndex].assertSymbol("{") {
				// 捕获函数字面值
				endIndex = scopeEndIndex(raws, curIndex, "{", "}")
				bodyTokens := parse4ComplexTokens(raws[curIndex+1 : endIndex])
				nextToken = newFuncLiteralToken("", argsTokens, bodyTokens)
				curIndex = endIndex + 1

			} else {
				nextToken = newFunctionCallToken(argsTokens)
			}
		} else {
			break
		}
		subs = append(subs, nextToken)
	}

	lastTokenSet(*res, newElemFCallToken(header, subs))
	return curIndex
}