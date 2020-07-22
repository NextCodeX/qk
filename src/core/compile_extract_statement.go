package core


func extractStatement(stmts StatementList) {
	ts := stmts.getRaw()
	for i := 0; i < len(ts); {
		t := ts[i]
		var endIndex int
		var stmt *Statement

		if !t.isIdentifier() && !t.isComplex() {
			goto next_loop
		}
		switch t.str {
		case "if":
			stmt, endIndex = extractIfStatement(i, ts)
		case "for":
			stmt, endIndex = extractForStatement(i, ts)
		case "switch":
		case "return":
			stmt, endIndex = extractReturnStatement(i, ts)
		default:
			if t.isFdef() {
				f, endIndex1 := extractFunction(i, ts)
				funcList[f.name] = f
				i = endIndex1
				goto next_loop
			}
			stmt, endIndex = extractExpressionStatement(i, ts)

		}
		if endIndex > 0 {
			stmts.addStatement(stmt)
			i = endIndex
		}

	next_loop:
		i++
	}
}

func extractExpressionStatement(currentIndex int, ts []Token) (*Statement, int) {
	size := len(ts)
	if !hasSymbol(ts[currentIndex:], ";") && currentIndex<size {
		stmt := newStatement(ExpressionStatement, ts[currentIndex:])
		return stmt, size
	}
	var nextBoundaryIndex = nextSymbolIndex(ts, currentIndex, ";")
	if nextBoundaryIndex > currentIndex {
		stmt := newStatement(ExpressionStatement, ts[currentIndex:nextBoundaryIndex])
		return stmt, nextBoundaryIndex
	}
	return nil, -1
}

func extractReturnStatement(currentIndex int, ts []Token) (*Statement, int) {
	stmt := &Statement{t:ReturnStatement}
	size := len(ts)
	var endIndex int
	for i:=currentIndex+1; i<size; i++ {
		t := ts[i]
		if t.assertSymbols("}", ";") {
			endIndex = i
			break
		}

		stmt.raw = append(stmt.raw, t)

		if i==size-1 {
			endIndex = i
		}
	}
	return stmt, endIndex
}

func extractFunction(currentIndex int, ts []Token) (*Function, int) {
	var nextIndex int
	f := newFunc(ts[currentIndex].str)
	f.defToken = ts[currentIndex]
	var blockTokens []Token
	size := len(ts)
	scopeOpenCount := 1
	for i:=currentIndex+2; i<size; i++ {
		token := ts[i]
		if token.assertSymbol("{") {
			scopeOpenCount ++
		}
		if token.assertSymbol("}") {
			scopeOpenCount --
			if scopeOpenCount == 0 {
				nextIndex = i
				break
			}
		}
		blockTokens = append(blockTokens, token)
	}
	if scopeOpenCount > 0 {
		panic("parse function statement exception!")
	}
	f.setRaw(blockTokens)
	return f, nextIndex
}

func extractIfStatement(currentIndex int, ts []Token) (*Statement, int) {
	var condStmts []*Statement
	var defStmt *Statement

	nextLoop:
	stmt := &Statement{t:IfStatement}
	index := nextSymbolIndex(ts, currentIndex,"{")
	stmt.condExprTokens = ts[currentIndex+1:index]

	scopeOpenCount := 1
	var endIndex int
	size := len(ts)
	for i:=index+1; i<size; i++ {
		t := ts[i]
		if t.assertSymbol("{") {
			scopeOpenCount++
		}
		if t.assertSymbol("}") {
			scopeOpenCount--
			if scopeOpenCount == 0 {
				endIndex = i
				break
			}
		}
		stmt.raw = append(stmt.raw, t)
	}

	if endIndex+1<size && ts[endIndex+1].assertIdentifier("elif") {
		currentIndex = endIndex+1
		condStmts = append(condStmts, stmt)
		goto nextLoop
	}
	if endIndex+1<size && ts[endIndex+1].assertIdentifier("else") {
		elseEndIndex := scopeEndIndex(ts, endIndex+1, "{", "}")
		if elseEndIndex > 0 {
			defStmt = newStatement(MultiStatement, ts[endIndex+3:elseEndIndex])
			endIndex = elseEndIndex
		}
	}

	condStmts = append(condStmts, stmt)

	ifStmt := newStatement(IfStatement, ts[currentIndex:endIndex+1])
	ifStmt.condStmts = condStmts
	ifStmt.defStmt = defStmt
	return ifStmt, endIndex
}

func extractForStatement(currentIndex int, ts []Token) (*Statement, int) {
	stmt := &Statement{t:IfStatement}
	index := nextSymbolIndex(ts, currentIndex,  "{")
	headerTokens := ts[currentIndex+1:index]
	stmt.preExprTokens, stmt.condExprTokens, stmt.postExprTokens = extractForHeaderExpressions(headerTokens)

	scopeOpenCount := 1
	var endIndex int
	for i:=index+1; i<len(ts); i++ {
		t := ts[i]
		if t.assertSymbol("{") {
			scopeOpenCount++
		}
		if t.assertSymbol("}") {
			scopeOpenCount--
			if scopeOpenCount == 0 {
				endIndex = i
				break
			}
		}
		stmt.raw = append(stmt.raw, t)
	}
	return stmt, endIndex
}

func extractForHeaderExpressions(ts []Token) (preTokens, condTokens, postTokens []Token) {
	size := len(ts)
	// for语句没用";"分隔，表达式即为condition expression
	if !hasSymbol(ts, ";") {
		condTokens = ts
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
