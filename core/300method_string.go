package core

import (
	"fmt"
	"strings"
)

func (str *StringValue) isObject() bool {
	return true
}
func (str *StringValue) get(key string) Value {
	raw := str.goValue
	switch key {
	case "size":
		return newInternalFunc(key, func(args []interface{})interface{}{
			return str.size()
		})

	case "trim":
		return newInternalFunc(key, func(args []interface{})interface{}{
			return strings.TrimSpace(raw)
		})

	case "replace":
		return newInternalFunc(key, func(args []interface{})interface{}{
			assert(len(args)<2, "method string.replace must has two parameter.")
			oldVal, ok1 := args[0].(string)
			newVal, ok2 := args[1].(string)
			assert(!ok1 || !ok2, "parameter type error, require replace(string, string)")
			return strings.ReplaceAll(raw, oldVal, newVal)
		})

	case "contain":
		return newInternalFunc(key, func(args []interface{})interface{}{
			assert(len(args)<1, "method string.contain must has one parameter.")
			subStr, ok := args[0].(string)
			assert(!ok, "parameter type error, require contain(string)")
			return strings.Contains(raw, subStr)
		})

	case "lower":
		return newInternalFunc(key, func(args []interface{})interface{}{
			return strings.ToLower(raw)
		})

	case "upper":
		return newInternalFunc(key, func(args []interface{})interface{}{
			return strings.ToUpper(raw)
		})

	case "lowerFirst":
		return newInternalFunc(key, func(args []interface{})interface{}{
			return strings.ToLower(str.sub(0, 1)) + str.sub(1, str.size())
		})

	case "upperFirst":
		return newInternalFunc(key, func(args []interface{})interface{}{
			return strings.ToUpper(str.sub(0, 1)) + str.sub(1, str.size())
		})

	case "toTitle":
		return newInternalFunc(key, func(args []interface{})interface{}{
			return strings.ToTitle(raw)
		})

	case "title":
		return newInternalFunc(key, func(args []interface{})interface{}{
			return strings.Title(raw)
		})

	case "hasPrefix":
		return newInternalFunc(key, func(args []interface{})interface{}{
			assert(len(args)<1, "method string.hasPrefix must has one parameter.")
			prefix, ok := args[0].(string)
			assert(!ok, "parameter type error, require hasPrefix(string)")
			return strings.HasPrefix(raw, prefix)
		})

	case "hasSuffix":
		return newInternalFunc(key, func(args []interface{})interface{}{
			assert(len(args)<1, "method string.hasSuffix must has one parameter.")
			suffix, ok := args[0].(string)
			assert(!ok, "parameter type error, require hasSuffix(string)")
			return strings.HasSuffix(raw, suffix)
		})

	case "split":
		return newInternalFunc(key, func(args []interface{})interface{}{
			assert(len(args)<1, "method string.split must has one parameter.")
			reg, ok := args[0].(string)
			assert(!ok, "parameter type error, require split(string)")
			tmp := toCommonSlice(strings.Split(raw, reg))
			return newQKValue(tmp)
		})

	case "eic":
		return newInternalFunc(key, func(args []interface{})interface{}{
			assert(len(args)<1, "method string.eic must has one parameter.")
			strB, ok := args[0].(string)
			assert(!ok, "parameter type error, require eic(string)")
			res := raw == strB || strings.ToLower(raw) == strings.ToLower(strB)
			return newQKValue(res)
		})
	default:
		runtimeExcption(fmt.Sprintf("string.%v() is undefined.", key))
		return NULL
	}
}