package core

import (
	"fmt"
	"bytes"
)

func evalJSONArrayMethod(arr JSONArray, method string, args []interface{}) (res *Value) {

	switch method {
	case "size":
		return newQkValue(arr.size())

	case "add":
		for _, arg := range args {
			arr.add(newQkValue(arg))
		}

	case "remove":
		assert(len(args)<1, "method array.remove must has one parameter.")
		index, ok := args[0].(int)
		assert(!ok, "method array.remove, parameter must be int type")
		arr.remove(index)

	case "join":
		assert(len(args)<1, "method array.join must has one parameter.")
		seperator, ok := args[0].(string)
		assert(!ok, "method array.join, parameter must be string type")
		vals := arr.values()
		var res bytes.Buffer
		for i, val := range vals {
			if i > 0 {
				res.WriteString(seperator)
			}
			valStr := fmt.Sprintf("%v", val.val())
			res.WriteString(valStr)
		}
		return newQkValue(res.String())

	case "sub":
		// todo
		//assert(len(args)<1, "method array.remove must has one parameter.")
		//startIndex, ok := args[0].(int)
		//assert(!ok, "method array.remove, parameter must be int type")
		//if len(args) < 2 {
		//	return
		//}

	default:
		runtimeExcption(fmt.Sprintf("array.%v is undefined!", method))
	}

	return nil
}

