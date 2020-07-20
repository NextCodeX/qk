package core

func extractStatement(stmts StatementList, ts []Token) {
	stmts.setRaw(ts)
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
	stmt := &Statement{t:IfStatement}

	index := nextSymbolIndex(ts, currentIndex,"{")
	stmt.condition = &Expression{
		raw:   ts[currentIndex+1:index],
	}

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

func extractForStatement(currentIndex int, ts []Token) (*Statement, int) {
	stmt := &Statement{t:IfStatement}
	index := nextSymbolIndex(ts, currentIndex,  "{")
	exprs := extractForHeaderExpressions(ts[currentIndex+1:index])
	stmt.setHeaderInfo(exprs)
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

func extractForHeaderExpressions(ts []Token) []*Expression {
	res := make([]*Expression, 3)
	size := len(ts)
	// for语句没用";"分隔，表达式即为condition expression
	if !hasSymbol(ts, ";") {
		res[1] = &Expression{
			t:     BinaryExpression,
			raw:   ts,
		}
		return res
	}

	//for 语句存在";"分隔符时
	// extract initialize expression
	index := nextSymbolIndex( ts, 0,";")
	if index > 2 {
		res[0] = &Expression{
			t:     BinaryExpression,
			raw:   ts[:index],
		}
	}

	// extract condition expression
	preIndex := index+1
	index = nextSymbolIndex(ts, preIndex, ";")
	if index - preIndex > 2 {
		res[1] = &Expression{
			t:     BinaryExpression,
			raw:   ts[preIndex:index],
		}
	} else if preIndex<size && !hasSymbol(ts[preIndex:], ";") {
		res[1] = &Expression{
			t:     BinaryExpression,
			raw:   ts[preIndex:],
		}
		index = size
	}

	// extract post expression
	preIndex = index+1
	if preIndex < size {
		res[2] = &Expression{
			t:     BinaryExpression,
			raw:   ts[preIndex:],
		}
	}
	return res
}
