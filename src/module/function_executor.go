package module

import (
	"reflect"
	"fmt"
)

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
			fmt.Println(methodType.In(ii).Kind())
		}

		fmt.Println("\nout params:")
		outcount := methodType.NumOut()
		for ii := 0; ii < outcount; ii++ {
			fmt.Println(methodType.Out(ii).Kind())
		}
		fmt.Println("+++++++++++++")
	}
}
