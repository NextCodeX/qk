package core

type ReturnStmtExtractor struct{}

func (se *ReturnStmtExtractor) check(cur Token) bool {
	return cur.assertKey("return")
}
func (se *ReturnStmtExtractor) extract(raws []Token, curIndex int) (Statement, int) {
	stmt := newReturnStatement()
	size := len(raws[curIndex:])
	nextIndex := curIndex + 1
	if size < 2 || (size == 2 && raws[nextIndex].assertSymbol(";")) {
		return stmt, curIndex + 1
	}

	var endIndex int
	size = len(raws)
	for i := nextIndex; i < size; i++ {
		t := raws[i]
		if t.assertSymbols("}", ";") {
			endIndex = i
			break
		}

		stmt.tokenAppend(t)

		if i == size-1 {
			endIndex = i
		}
	}

	return stmt, endIndex
}
