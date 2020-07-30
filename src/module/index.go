package module

import (
	"fmt"
	"reflect"
)

var funcs map[string]*FunctionExecutor


func collectFunctionInfo(obj interface{})  {
	fmt.Println("do collectFunctionInfo...")
	v1 := reflect.ValueOf(obj).Elem()
	k1 := v1.Type()
	for i := 0; i < v1.NumMethod(); i++ {
		methodName := k1.Method(i).Name
		methodObject := v1.Method(i)
		if methodName == "GetAge" {
			fmt.Println("age: ", methodObject.Call(nil)[0])
		}

		fmt.Println(methodName, methodObject.Type())
		methodType := methodObject.Type()

		fmt.Println("\nin params:")
		incount := methodType.NumIn()
		for ii := 0; ii < incount; ii++ {
			argType := methodType.In(ii)
			fmt.Println(argType, argType.Kind())
		}

		fmt.Println("\nout params:")
		outcount := methodType.NumOut()
		for ii := 0; ii < outcount; ii++ {
			argType := methodType.Out(ii)
			fmt.Println(argType, argType.Kind())
		}
		fmt.Println("+++++++++++++")
	}
}

func Load() map[string]*FunctionExecutor {
	fmt.Println("register all module...")
	return funcs
}
