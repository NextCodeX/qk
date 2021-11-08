package core


type ContinueStmtExtractor struct {}

func (se *ContinueStmtExtractor) check(cur Token) bool {
	return cur.assertKey("continue")
}
func (se *ContinueStmtExtractor) extract(raws []Token, curIndex int) (Statement, int) {
	stmt := newContinueStatement()
	size := len(raws)
	if curIndex == size -1 {
		return stmt, size
	}

	endIndex := nextSymbolIndex(raws, curIndex, ";")
	return stmt, endIndex
}