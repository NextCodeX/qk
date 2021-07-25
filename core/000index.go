package core

var (
	funcList = make(map[string]*Function)
	mainFunc = newFunc("main")
)

//const DEBUG_MODE = true
const DEBUG_MODE = false

func Run(bs []byte) {
	// 词法分析
	ts := ParseTokens(bs)
	//printTokensByLine(ts)

	// 语法分析
	mainFunc.setRaw(ts)
	Compile(mainFunc)
	//printFunc()

	// 解析并执行
	Interpret()
}

func ParseTokens(bs []byte) []Token {
	// 提取原始token列表
	ts := parse4PrimaryTokens(bs)

	// 语法预处理
	// 提取'++', '--'等运算符以及负数表达式
	ts = parse4OperatorTokens(ts)
	// 去掉无用的';', 合并token生成函数调用token(Fcall), 方法调用token(Mtcall)等复合token
	ts = parse4ComplexTokens(ts)
	return ts
}

func Compile(stmts StatementList) {
	if stmts == nil {
		return
	}
	if stmts.isCompiled() {
		return
	} else {
		stmts.setCompiled()
	}
	extractStatement(stmts)
	parseStatementList(stmts.stmts())

	for _, customFunc := range funcList {
		Compile(customFunc)
	}
}

func Interpret() {
	stack := newVariableStack()
	stack.push() // 执行方法前，向变量栈(list)添加的一个变量池(map)
	executeFunctionStatementList(mainFunc.block, stack)
}


