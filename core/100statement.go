package core

type Statement interface {
	addStmt(stmt Statement) // 添加子statement, 并将stack传递给 子statement
	stmts() []Statement
	tokenList() []Token
	setTokenList(ts []Token)
	tokenAppend(t Token)
	setParent(p Function) // 等同于setStack()
	getParent() Function

	setVar(name string, value Value)

	parse() // 解析 各子statement, expression; 并将stack传递给 它们
	execute() StatementResult

	String() string
}
