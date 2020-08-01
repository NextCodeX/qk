package core

import "fmt"

func isPrint(funcName string) bool {
	return match(funcName, "println", "printf", "print")
}

func executePrintFunc(funcName string, args []interface{}) (res *Value) {
	argCount := len(args)
	if funcName == "println" {
		fmt.Println(args...)
		return
	}

	if funcName == "printf" {
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
		return
	}

	if funcName == "print" {
		fmt.Print(args...)
	}
	return
}

