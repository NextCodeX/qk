package core

func extractStatement(stmt Statement) {
	ts := stmt.raw()
	for i := 0; i < len(ts); {
		t := ts[i]
		var endIndex int
		var subStmt Statement

		switch {
		case t.assertIdentifier("if"):
			subStmt, endIndex = extractIfStatement(i, ts)
		case t.assertIdentifier("for"):
			subStmt, endIndex = extractForStatement(i, ts)
		case t.assertIdentifier("foreach"):
			subStmt, endIndex = extractForeachStatement(i, ts)
		case t.assertIdentifier("fori"):
			subStmt, endIndex = extractForIndexStatement(i, ts)
		case t.assertIdentifier("forv"):
			subStmt, endIndex = extractForValueStatement(i, ts)
		case t.assertIdentifier("switch"):
		case t.assertIdentifier("continue"):
			subStmt, endIndex = extractContinueStatement(i, ts)
		case t.assertIdentifier("break"):
			subStmt, endIndex = extractBreakStatement(i, ts)
		case t.assertIdentifier("return"):
			subStmt, endIndex = extractReturnStatement(i, ts)
		default:
			// 提取 表达式语句
			subStmt, endIndex = extractExpressionStatement(i, ts)
		}
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
	runtimeExcption("unknown statement: ", tokensShow10(ts[currentIndex:]))
	return nil, -1
}

func extractContinueStatement(currentIndex int, ts []Token) (Statement, int) {
	stmt := newContinueStatement()
	size := len(ts)
	if currentIndex == size -1 {
		return stmt, size
	}

	endIndex := nextSymbolIndex(ts, currentIndex, ";")
	return stmt, endIndex
}

func extractBreakStatement(currentIndex int, ts []Token) (Statement, int) {
	stmt := newBreakStatement()
	size := len(ts)
	if currentIndex == size -1 {
		return stmt, size
	}


	endIndex := nextSymbolIndex(ts, currentIndex, ";")
	return stmt, endIndex
}

func extractReturnStatement(currentIndex int, ts []Token) (Statement, int) {
	stmt := newReturnStatement()
	size := len(ts[currentIndex:])
	nextIndex := currentIndex + 1
	if size < 2 || (size == 2 && ts[nextIndex].assertSymbol(";")) {
		return stmt, currentIndex + 1
	}

	var endIndex int
	size = len(ts)
	for i:=nextIndex; i<size; i++ {
		t := ts[i]
		if t.assertSymbols("}", ";") {
			endIndex = i
			break
		}

		stmt.rawAppend(t)

		if i==size-1 {
			endIndex = i
		}
	}

	return stmt, endIndex
}


func extractIfStatement(currentIndex int, ts []Token) (Statement, int) {
	var condStmts []Statement
	var defStmt Statement
	size := len(ts)

	nextLoop:
	index := nextSymbolIndex(ts, currentIndex,"{")
	condExprTokens := ts[currentIndex+1:index]
	endIndex := scopeEndIndex(ts, index, "{", "}")
	stmt := newSingleIfStatement(condExprTokens, ts[index+1:endIndex])

	if endIndex+1<size && ts[endIndex+1].assertIdentifier("elif") {
		currentIndex = endIndex+1
		condStmts = append(condStmts, stmt)
		goto nextLoop
	}
	if endIndex+1<size && ts[endIndex+1].assertIdentifier("else") {
		elseEndIndex := scopeEndIndex(ts, endIndex+1, "{", "}")
		if elseEndIndex > 0 {
			defStmt = newMultiStatement(ts[endIndex+3:elseEndIndex])
			endIndex = elseEndIndex
		}
	}

	condStmts = append(condStmts, stmt)

	return newMultiIfStatement(condStmts, defStmt), endIndex
}

func extractForStatement(currentIndex int, ts []Token) (Statement, int) {
	index := nextSymbolIndex(ts, currentIndex,  "{")
	headerTokens := ts[currentIndex+1:index]
	preExprTokens, condExprTokens, postExprTokens := extractForHeaderExpressions(headerTokens)
	stmt := newForStatement(preExprTokens, condExprTokens, postExprTokens)

	endIndex := scopeEndIndex(ts, index, "{", "}")
	stmt.setRaw(ts[index+1:endIndex])

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

func extractForeachStatement(currentIndex int, ts []Token) (Statement, int) {
	headerEndIndex := nextSymbolIndex(ts, currentIndex, "{")
	headerInfo := ts[currentIndex+1: headerEndIndex]
	stmt := newForeachStatement(headerInfo[0].raw(), headerInfo[2].raw(), headerInfo[4:])

	stmtEndIndex := scopeEndIndex(ts, currentIndex, "{", "}")
	stmt.setRaw(ts[headerEndIndex+1: stmtEndIndex])
	return stmt, stmtEndIndex
}

func extractForIndexStatement(currentIndex int, ts []Token) (Statement, int) {
	headerEndIndex := nextSymbolIndex(ts, currentIndex, "{")
	headerInfo := ts[currentIndex+1: headerEndIndex]
	stmt := newForIndexStatement(headerInfo[0].raw(), headerInfo[2:])

	stmtEndIndex := scopeEndIndex(ts, currentIndex, "{", "}")
	stmt.setRaw(ts[headerEndIndex+1: stmtEndIndex])
	return stmt, stmtEndIndex
}

func extractForValueStatement(currentIndex int, ts []Token) (Statement, int) {
	headerEndIndex := nextSymbolIndex(ts, currentIndex, "{")
	headerInfo := ts[currentIndex+1: headerEndIndex]
	stmt := newForValueStatement(headerInfo[0].raw(), headerInfo[2:])

	stmtEndIndex := scopeEndIndex(ts, currentIndex, "{", "}")
	stmt.setRaw(ts[headerEndIndex+1: stmtEndIndex])
	return stmt, stmtEndIndex
}
