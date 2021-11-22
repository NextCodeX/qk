package core

type Statement interface {
	addStmt(stmt Statement) // 添加子statement, 并将stack传递给 子statement
	stmts() []Statement
	setStatements(stmts []Statement)

	SourceCode
	tokenAppend(t Token)

	getStack() Function   //函数返回自己，语句返回父函数
	setParent(p Function) // 等同于setStack()
	getParent() Function

	parse() // 解析 各子statement, expression; 并将stack传递给 它们
	execute() StatementResult

	String() string
}
