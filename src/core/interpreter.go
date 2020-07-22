package core


func Interpret() {
    stack := newVariableStack()

	executeFunctionStatementList(mainFunc.block, stack)
}



