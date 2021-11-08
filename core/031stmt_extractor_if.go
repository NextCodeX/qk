package core

type IfStmtExtractor struct {}

func (se *IfStmtExtractor) check(cur Token) bool {
	return cur.assertKey("if")
}
func (se *IfStmtExtractor) extract(raws []Token, curIndex int) (Statement, int) {
	var condStmts []Statement
	var defStmt Statement
	size := len(raws)

nextLoop:
	index := nextSymbolIndex(raws, curIndex,"{")
	condExprTokens := raws[curIndex+1:index]
	endIndex := scopeEndIndex(raws, index, "{", "}")
	stmt := newSingleIfStatement(condExprTokens, raws[index+1:endIndex])

	if endIndex+1<size && raws[endIndex+1].assertKey("elif") {
		curIndex = endIndex+1
		condStmts = append(condStmts, stmt)
		goto nextLoop
	}
	if endIndex+1<size && raws[endIndex+1].assertKey("else") {
		elseEndIndex := scopeEndIndex(raws, endIndex+1, "{", "}")
		if elseEndIndex > 0 {
			defStmt = newMultiStatement(raws[endIndex+3:elseEndIndex])
			endIndex = elseEndIndex
		}
	}

	condStmts = append(condStmts, stmt)

	return newMultiIfStatement(condStmts, defStmt), endIndex
}