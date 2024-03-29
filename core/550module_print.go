package core

import "fmt"

func (this *InternalFunctionSet) Print(args []interface{}) {
	fmt.Print(args...)
}

func (this *InternalFunctionSet) Printf(args []interface{}) {
	argCount := len(args)
	if argCount < 1 {
		runtimeExcption("printf argument is too less")
		return
	}
	format, ok := args[0].(string)
	if !ok {
		runtimeExcption("printf argumant format must be string type.")
		return
	}
	if argCount == 1 {
		fmt.Printf(format)
		return
	}
	fmt.Printf(format, args[1:]...)
}

func (this *InternalFunctionSet) Println(args []interface{}) {
	fmt.Println(args...)
}

func (this *InternalFunctionSet) Echo(args []interface{}) {
	fmt.Println(args...)
}
func (this *InternalFunctionSet) Echof(args []interface{}) {
	argCount := len(args)
	if argCount < 1 {
		runtimeExcption("echof() argument is too less")
		return
	}
	format, ok := args[0].(string)
	if !ok {
		runtimeExcption("echof() argumant format must be string type.")
		return
	}
	format = format + "\n"

	fmt.Printf(format, args[1:]...)
}
