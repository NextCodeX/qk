package core

type StatementList interface {
	stmts() []*Statement
	addStatement(*Statement)
	getRaw() []Token
	setRaw([]Token)
	isCompiled() bool
	setCompiled()
}
