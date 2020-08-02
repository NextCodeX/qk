package core


type Function struct {
	super Variables // 父作用域的变量列表
	local Variables // 当前作用域的变量列表
	params []Value // 参数
	res []Value // 返回值
	block []*Statement // 执行语句
	name string // 函数名
	paramNames []string // 形参名列表
	defToken Token  // 函数定义token, 包含函数名,参数名
	raw []Token // token列表
	compiled bool // 是否已编译
}

//func newFunc(name string)

func newFunc(name string) *Function {
	return &Function{local:newVariables(), name:name}
}

func (f *Function) addStatement(stm *Statement) {
	f.block = append(f.block, stm)
}

func (f *Function) stmts() []*Statement {
	return f.block
}

func (f *Function) getRaw() []Token {
	return f.raw
}

func (f *Function) setRaw(ts []Token) {
	f.raw = ts
}

func (f *Function) isCompiled() bool {
	return f.compiled
}

func (f *Function) setCompiled() {
	f.compiled = true
}
