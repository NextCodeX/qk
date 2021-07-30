package core

import "fmt"

type Function interface {
	getName() string
	getParent() Function
	getLocalVars() Variables
	setArgs(args []Value)
	Statement
	Value
}

type FunctionImpl struct {
	local      Variables        // 当前作用域的变量列表
	name       string           // 函数名
	paramNames []string         // 形参名列表
	args []Value
	StatementAdapter
	ValueAdapter
}

func newFunc(name string, ts []Token, paramNames []string) Function {
	f := &FunctionImpl{name:name}
	f.StatementAdapter.ts = ts
	f.paramNames = paramNames
	f.initStatement(f)
	return f
}

func (f *FunctionImpl) getName() string {
	return f.name
}
func (f *FunctionImpl) getLocalVars() Variables {
	return f.local
}

func (f *FunctionImpl) parse() {
	fmt.Println("parse function!!!")
}

func (f *FunctionImpl) setArgs(args []Value) {
	f.args = args
}

func (f *FunctionImpl) execute() StatementResult {
	f.local = newVariables()

	for i, paramName := range f.paramNames {
		var arg Value
		if i >= len(f.args) {
			arg = NULL
		} else {
			arg = f.args[i]
		}
		f.local.add(paramName, arg)
	}

	return f.executeStatementList(f.block, StmtListTypeFunc)
}

func (f *FunctionImpl) setParent(p Function) {
	f.parent = p
}

func (f *FunctionImpl) val() interface{} {
	return f.name + "()"
}
func (f *FunctionImpl) isFunction() bool {
	return true
}

func (f *FunctionImpl) String() string {
	return f.name + "()"
}