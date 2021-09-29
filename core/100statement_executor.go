package core

// 语句列表类型， 配合StatementResult实现return, continue, break
type StmtListType int

const (
	StmtListTypeFunc StmtListType = 1 << iota
	StmtListTypeIf
	StmtListTypeFor
	StmtListTypeNormal
)

// StatementExecutor 用于执行statement列表
type StatementExecutor struct{}

func (stmtExec *StatementExecutor) executeStatementList(stmts []Statement, t StmtListType) StatementResult {
	var res StatementResult
	for _, stmt := range stmts {

		res = stmt.execute()

		if res == nil {
			println("executeStatement return error")
			break
		}

		if res.isContinue() {
			if t == StmtListTypeFor {
				res.setType(StatementNormal)
			}
			break
		} else if res.isReturn() || res.isBreak() {
			break
		}
	}
	// 修复空语句异常
	if res == nil {
		res = newStatementResult(StatementNormal, NULL)
	}
	return res
}
