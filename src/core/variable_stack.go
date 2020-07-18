package core


type VariableStack struct {
	list []Variables
}

func newVariableStack() *VariableStack {
	vars := newVariables()
	varStack := &VariableStack{}
	varStack.list = append(varStack.list, vars)
	return varStack
}

func (stack *VariableStack) push() {
	vars := newVariables()
	stack.list = append(stack.list, vars)
}

func (stack *VariableStack) pop() {
	size := len(stack.list)
	if size<1 {
		return
	}
	stack.list = stack.list[:size-1]
}

func (stack *VariableStack) clone() *VariableStack {
	size := len(stack.list)
	if size<1 {
		return newVariableStack()
	}

	newStack := stack.list[:size-1]
	return &VariableStack{newStack}
}

func (stack *VariableStack) searchVariable(name string) *Variable {
	for i:=len(stack.list)-1; i>=0; i-- {
		vars := stack.list[i]
		res := vars.get(name)
		if res !=nil {
			return res
		}
	}
	runtimeExcption("variable", name, "is undefined")
	return nil
}

func (stack *VariableStack) addLocalVariable(vr *Variable) {
	size := len(stack.list)
	if size<1 {
		runtimeExcption("stack is empty!")
	}

	stack.list[size-1].add(vr)
}
