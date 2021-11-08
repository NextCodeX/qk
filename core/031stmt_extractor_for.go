package core


type ForStmtExtractor struct {}

func (se *ForStmtExtractor) check(cur Token) bool {
	return cur.assertKey("for")
}
func (se *ForStmtExtractor) extract(raws []Token, curIndex int) (Statement, int) {
	index := nextSymbolIndex(raws, curIndex,  "{")
	headerTokens := raws[curIndex+1:index]
	preExprTokens, condExprTokens, postExprTokens := extractForHeaderExpressions(headerTokens)
	stmt := newForStatement(preExprTokens, condExprTokens, postExprTokens)

	endIndex := scopeEndIndex(raws, index, "{", "}")
	stmt.setRaw(raws[index+1:endIndex])

	return stmt, endIndex
}

func extractForHeaderExpressions(ts []Token) (preTokens, condTokens, postTokens []Token) {
	size := len(ts)
	// for语句没用";"分隔，表达式即为condition expression
	if !hasSymbol(ts, ";") {
		if len(ts) > 0 {
			condTokens = ts
		}
		return
	}

	//for 语句存在";"分隔符时
	// extract initialize expression
	currentIndex := 0
	boundaryIndex := nextSymbolIndex(ts, currentIndex,";")
	if boundaryIndex > 0 {
		preTokens = ts[:boundaryIndex]
	}

	// extract condition expression
	currentIndex = boundaryIndex+1
	boundaryIndex = nextSymbolIndex(ts, currentIndex, ";")
	if boundaryIndex > currentIndex {
		condTokens = ts[currentIndex:boundaryIndex]
	} else if currentIndex<size && !hasSymbol(ts[currentIndex:], ";") {
		condTokens = ts[currentIndex:]
		return
	}

	// extract post expression
	currentIndex = boundaryIndex+1
	if currentIndex < size {
		postTokens = ts[currentIndex:]
	}
	return
}