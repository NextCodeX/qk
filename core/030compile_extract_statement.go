package core

type StmtExtractor interface {
	check(cur Token) bool
	extract(raws []Token, curIndex int) (Statement, int)
}

var stmtExtractorList = []StmtExtractor{
	&BreakStmtExtractor{},
	&ContinueStmtExtractor{},
	&ForStmtExtractor{},
	&ForeachStmtExtractor{},
	&ForIndexStmtExtractor{},
	&ForValueStmtExtractor{},
	&IfStmtExtractor{},
	&ReturnStmtExtractor{},
}

func extractStatement(stmt Statement) {
	ts := stmt.tokenList()
	for i := 0; i < len(ts); {
		cur := ts[i]
		var endIndex int
		var subStmt Statement

		for _, extractor := range stmtExtractorList {
			if extractor.check(cur) {
				subStmt, endIndex = extractor.extract(ts, i)
				goto loopEnd
			}
		}

		// 提取 表达式语句
		subStmt, endIndex = extractExpressionStatement(i, ts)

	loopEnd:
		if endIndex > 0 {
			if subStmt != nil {
				stmt.addStmt(subStmt)
			}
			i = endIndex
		}
		i++
	}
}

func extractExpressionStatement(currentIndex int, ts []Token) (Statement, int) {
	size := len(ts)
	if !hasSymbol(ts[currentIndex:], ";") {
		stmt := newExpressionStatement(ts[currentIndex:])
		return stmt, size
	}
	var nextBoundaryIndex = nextSymbolIndex(ts, currentIndex, ";")
	if nextBoundaryIndex > currentIndex {
		stmt := newExpressionStatement(ts[currentIndex:nextBoundaryIndex])
		return stmt, nextBoundaryIndex
	}
	if nextBoundaryIndex > -1 {
		// 处理空语句
		return nil, nextBoundaryIndex
	}
	runtimeExcption("unknown statement: ", tokensShow10(ts[currentIndex:]))
	return nil, -1
}
