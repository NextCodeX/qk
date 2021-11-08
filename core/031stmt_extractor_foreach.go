package core

type ForeachStmtExtractor struct {}

func (se *ForeachStmtExtractor) check(cur Token) bool {
	return cur.assertKey("foreach")
}
func (se *ForeachStmtExtractor) extract(raws []Token, curIndex int) (Statement, int) {
	headerEndIndex := nextSymbolIndex(raws, curIndex, "{")
	headerInfo := raws[curIndex+1: headerEndIndex]
	k, v, exprTokens := parseForeachHeader(headerInfo)
	stmt := newForeachStatement(k, v, exprTokens)

	stmtEndIndex := scopeEndIndex(raws, curIndex, "{", "}")
	stmt.setRaw(raws[headerEndIndex+1: stmtEndIndex])
	return stmt, stmtEndIndex
}

type ForIndexStmtExtractor struct {}

func (se *ForIndexStmtExtractor) check(cur Token) bool {
	return cur.assertKey("fori")
}
func (se *ForIndexStmtExtractor) extract(raws []Token, curIndex int) (Statement, int) {
	headerEndIndex := nextSymbolIndex(raws, curIndex, "{")
	headerInfo := raws[curIndex+1: headerEndIndex]
	_, i, exprTokens := parseForeachHeader(headerInfo)
	stmt := newForIndexStatement(i, exprTokens)

	stmtEndIndex := scopeEndIndex(raws, curIndex, "{", "}")
	stmt.setRaw(raws[headerEndIndex+1: stmtEndIndex])
	return stmt, stmtEndIndex
}

type ForValueStmtExtractor struct {}

func (se *ForValueStmtExtractor) check(cur Token) bool {
	return cur.assertKey("forv")
}
func (se *ForValueStmtExtractor) extract(raws []Token, curIndex int) (Statement, int) {
	headerEndIndex := nextSymbolIndex(raws, curIndex, "{")
	headerInfo := raws[curIndex+1: headerEndIndex]
	_, v, exprTokens := parseForeachHeader(headerInfo)
	stmt := newForValueStatement(v, exprTokens)

	stmtEndIndex := scopeEndIndex(raws, curIndex, "{", "}")
	stmt.setRaw(raws[headerEndIndex+1: stmtEndIndex])
	return stmt, stmtEndIndex
}

func parseForeachHeader(ts []Token) (string, string, []Token) {
	eqSymbolIndex := nextSymbolIndex(ts, 0, ":")
	if eqSymbolIndex < 1 {
		runtimeExcption("invalid foreach: ", tokensString(ts))
	}
	first, end := ts[:eqSymbolIndex], ts[eqSymbolIndex+1:]
	var var1, var2 string
	if len(first) == 1 {
		var1, var2 = "", toName(first[0])
	} else if len(first) == 3 {
		var1, var2 = toName(first[0]), toName(first[2])
	} else {
		runtimeExcption("invalid foreach: ", tokensString(ts))
	}
	return var1, var2, end
}
