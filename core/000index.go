package core

var mainFunc = newFuncWithoutTokens("main")
// 添加 Quick 系统内部函数
func addInternalFunc(name string, internalFunc func([]interface{})interface{}) {
	mainFunc.addFunc(name, newInternalFunc(name, internalFunc))
}
func addModuleFunc(name string, moduleFunc *FunctionExecutor) {
	mainFunc.addFunc(name, newModuleFunc(name, moduleFunc))
}

func Run(bs []byte) {
	// 词法分析
	ts := ParseTokens(bs)
	printTokensByLine(ts)

	// 语法分析(解析)
	mainFunc.setRaw(ts)
	Compile(mainFunc)

	// 程序执行
	mainFunc.execute()
}

// 指定变量𣏾, 执行qk代码片段.
func evalScript(src string, stack Function) Value {
	ts := ParseTokens([]byte(src))
	tsLen := len(ts)
	if tsLen < 1 {
		return NULL
	}
	if last(ts).assertSymbol(";") {
		ts = ts[:tsLen-1]
	}
	if len(ts) < 1 {
		return NULL
	}
	expr := extractExpression(ts)
	expr.setStack(stack)
	qkValue := expr.execute()
	return qkValue
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

func Compile(stmt Statement) {
	extractStatement(stmt)
	for _, stmt := range stmt.stmts() {
		stmt.parse()
	}
}

