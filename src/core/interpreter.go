package core

import "fmt"

func Interpret() {
    stack := newVariableStack()
	stack.push()
    fmt.Println("main stack.list:", len(stack.list))
	executeFunctionStatementList(mainFunc.block, stack)
}



