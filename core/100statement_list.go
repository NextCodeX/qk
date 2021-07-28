package core

type StatementList interface {
	stmts() []*Statement     // 语句列表
	addStatement(*Statement) // 添加语句至语句列表
	getRaw() []Token     // 获取StatementList的Token列表
	setRaw([]Token)      // 设置StatementList的Token列表
	isCompiled() bool        // 该StatementList是否完成语法分析
	setCompiled()            // 标记为已完成语法分析
}
