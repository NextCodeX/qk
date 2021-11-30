package core

type Expression interface {
	SourceCode
	setLocalScope()

	setParent(p Function) // 设置表达式的stack
	getParant() Function  // 获取表达式的stack

	execute() Value // 执行表达式
}
