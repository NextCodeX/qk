package core

import "fmt"

func evalJSONArrayMethod(arr JSONArray, method string, args []interface{}) (res *Value) {
	if method == "size" {
		return newVal(arr.size())
	}
	if method == "add" {
		for _, arg := range args {
			arr.add(newVal(arg))
		}
		return
	}
	if method == "remove" {
		assert(len(args)<1, "method array.remove must has one parameter.")
		index, ok := args[0].(int)
		assert(!ok, "method array.remove, parameter must be int type")
		arr.remove(index)
		return
	}


	runtimeExcption(fmt.Sprintf("array.%v is undefined!", method))
	return nil
}

