package core


type Function struct {
	super Variables // 父作用域的变量列表
	local Variables // 当前作用域的变量列表
	params []Value // 参数
	res []Value // 返回值
	block []*Statement // 执行语句
	name string
	defToken Token
	raw []Token // token列表
	compiled bool
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
