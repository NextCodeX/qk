package core

import (
	"bytes"
	"fmt"
)

func (arr *JSONArrayImpl) isObject() bool {
	return true
}

func (arr *JSONArrayImpl) get(key string) Value {
	switch key {
	case "size":
		return newInternalFunc(key, func(args []interface{})interface{}{
			return newQKValue(arr.size())
		})

	case "add":
		return newInternalFunc(key, func(args []interface{})interface{}{
			for _, arg := range args {
				arr.add(newQKValue(arg))
			}
			return NULL
		})

	case "remove":
		return newInternalFunc(key, func(args []interface{})interface{}{
			assert(len(args)<1, "method array.remove must has one parameter.")
			index := toInt(args[0])
			//assert(!ok, "method array.remove, parameter must be int type")
			arr.remove(index)
			return NULL
		})

	case "join":
		return newInternalFunc(key, func(args []interface{})interface{}{
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
			return newQKValue(res.String())
		})
	default:
		runtimeExcption(fmt.Sprintf("array.%v() is undefined!", key))
		return NULL
	}
}