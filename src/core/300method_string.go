package core

import (
	"unicode/utf8"
	"strings"
	"fmt"
)

func evalStringMethod(str string, method string, args []interface{}) (res *Value) {
	var rawVal interface{}
	argCount := len(args)
	switch method {
	case "size":
		rawVal = utf8.RuneCountInString(str)

	case "trim":
		rawVal = strings.TrimSpace(str)

	case "replace":
		assert(len(args)<2, "method string.replace must has two parameter.")
		oldVal, ok1 := args[0].(string)
		newVal, ok2 := args[1].(string)
		assert(!ok1 || !ok2, "parameter type error, require replace(string, string)")
		rawVal = strings.Replace(str, oldVal, newVal, -1)

	case "contain":
		assert(len(args)<1, "method string.contain must has one parameter.")
		subStr, ok := args[0].(string)
		assert(!ok, "parameter type error, require contain(string)")
		rawVal = strings.Contains(str, subStr)

	case "lower":
		rawVal = strings.ToLower(str)

	case "upper":
		rawVal = strings.ToUpper(str)

	case "lowerFirst":
		rawVal = strings.ToLower(str[:1]) + str[1:]

	case "upperFirst":
		rawVal = strings.ToUpper(str[:1]) + str[1:]

	case "toTitle":
		rawVal = strings.ToTitle(str)

	case "title":
		rawVal = strings.Title(str)

	case "hasPrefix":
		assert(len(args)<1, "method string.hasPrefix must has one parameter.")
		prefix, ok := args[0].(string)
		assert(!ok, "parameter type error, require hasPrefix(string)")
		rawVal = strings.HasPrefix(str, prefix)

	case "hasSuffix":
		assert(len(args)<1, "method string.hasSuffix must has one parameter.")
		suffix, ok := args[0].(string)
		assert(!ok, "parameter type error, require hasSuffix(string)")
		rawVal = strings.HasSuffix(str, suffix)

	case "sub":
		assert(argCount<1, "method string.sub must has one parameter.")
		startIndex, ok1 := args[0].(int)
		if argCount == 1 {
			assert(!ok1, "parameter type error, require sub(int)")
			rawVal = str[startIndex:]
			break
		}
		endIndex, ok2 := args[1].(int)
		assert(!ok1 || !ok2, "parameter type error, require sub(int, int)")
		assert(startIndex<0 || startIndex>endIndex || endIndex>=len(str), fmt.Sprintf("string out of index, sub(%v, %v)", startIndex, endIndex))
		rawVal = str[startIndex: endIndex]

	case "split":
		assert(len(args)<1, "method string.split must has one parameter.")
		reg, ok := args[0].(string)
		assert(!ok, "parameter type error, require split(string)")
		tmp := toCommonSlice(strings.Split(str, reg))
		return toQKValue(tmp)

	default:
		runtimeExcption(fmt.Sprintf("string.%v is undefined.", method))
	}

	if rawVal == nil {
		return
	}
	return newVal(rawVal)

}