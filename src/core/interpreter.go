package core

func Interpret() {
    stack := newVariableStack()
	stack.push()
	executeFunctionStatementList(mainFunc.block, stack)
}



