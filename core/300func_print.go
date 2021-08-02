package core

import "fmt"

func init()  {
	addInternalFunc("print", func(args []interface{}) interface{} {
		fmt.Print(args...)
		return nil
	})

	addInternalFunc("printf", func(args []interface{}) (res interface{}) {
		argCount := len(args)
		if argCount < 1 {
			runtimeExcption("printf argument is too less")
			return res
		}
		format, ok := args[0].(string)
		if !ok {
			runtimeExcption("printf argumant format must be string type.")
			return res
		}
		if argCount == 1 {
			fmt.Printf(format)
			return res
		}
		fmt.Printf(format, args[1:]...)
		return res
	})

	addInternalFunc("println", func(args []interface{}) interface{} {
		fmt.Println(args...)
		return nil
	})
}



