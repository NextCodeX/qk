package core

type StatementAdapter struct {
	owner  Statement   // 具体Statement的引用
	parent Function    // 父引用
	ts     []Token     // 原始Token列表
	block  []Statement // 执行语句列表
	StatementExecutor
	ValueStack
}

func (stmtAdapter *StatementAdapter) initStatement(stmt Statement) {
	// 为getStack()判断当前的statement 是否为stack引用
	stmtAdapter.owner = stmt

	if fn, ok := stmt.(Function); ok {
		// 初始化main函数的值栈
		stmtAdapter.ValueStack.cur = fn
	}
}

func (stmtAdapter *StatementAdapter) getStack() Function {
	// cur 在Quick中是存放变量的一个地方
	// 对于函数, 它自身就是stack, 它的stack与parent不一致, 它的parent是父函数(上一层stack)
	// 对于非函数的statement, 它们的parent就是stack, parent与stack是一致的, 皆是父函数
	if f, ok := stmtAdapter.owner.(Function); ok {
		return f
	} else {
		return stmtAdapter.parent
	}
}

func (stmtAdapter *StatementAdapter) parse() {
}

func (stmtAdapter *StatementAdapter) tokenList() []Token {
	return stmtAdapter.ts
}

func (stmtAdapter *StatementAdapter) setTokenList(ts []Token) {
	stmtAdapter.ts = ts
}

func (stmtAdapter *StatementAdapter) tokenAppend(t Token) {
	stmtAdapter.ts = append(stmtAdapter.ts, t)
}

func (stmtAdapter *StatementAdapter) stmts() []Statement {
	return stmtAdapter.block
}

func (stmtAdapter *StatementAdapter) addStmt(stmt Statement) {
	// 将stack逐级往下传递
	stack := stmtAdapter.getStack()
	stmt.setParent(stack)

	// 将子statement添加至当前statement列表
	stmtAdapter.block = append(stmtAdapter.block, stmt)
}

func (stmtAdapter *StatementAdapter) getParent() Function {
	return stmtAdapter.parent
}

func (stmtAdapter *StatementAdapter) setParent(p Function) {
	stmtAdapter.parent = p

	// initStatement()[stmt创建] -> addStmt()[父statement传入stack] -> setParent()
	// 为当前statement设置stack, 启用ValueStack
	// 使得foreach这样的statement 可以使用ValueStack
	stmtAdapter.ValueStack.cur = stmtAdapter.getStack()
}

func (stmtAdapter *StatementAdapter) String() string {
	return "statement"
}
