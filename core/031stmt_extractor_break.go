package core

type BreakStmtExtractor struct {}

func (se *BreakStmtExtractor) check(cur Token) bool {
	return cur.assertKey("break")
}
func (se *BreakStmtExtractor) extract(raws []Token, curIndex int) (Statement, int) {
	stmt := newBreakStatement()
	size := len(raws)
	if curIndex == size -1 {
		return stmt, size
	}

	endIndex := nextSymbolIndex(raws, curIndex, ";")
	return stmt, endIndex
}
