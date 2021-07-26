package core

import (
	"reflect"
	"strconv"
	"fmt"
	"strings"
)


func toIntValue(any interface{}) int {
	switch v := any.(type) {
	case int:
		return v
	case float32:
		return int(v)
	case float64:
		return int(v)
	case string:
		i, err := strconv.Atoi(v)
		assert(err!=nil, err, "failed to int value:", any)
		return i
	default:
		runtimeExcption("failed to int value", any)
	}
	return -1
}

func toStringValue(any interface{}) string {
	return fmt.Sprintf("%v", any)
}

func toCommonMap(any interface{}) map[string]interface{} {
	res := make(map[string]interface{})
	switch v := any.(type) {
	case map[string]string:
		for key, val := range v {
			res[key] = val
		}
	default:
		runtimeExcption("toCommonMap# unknown type:", reflect.TypeOf(any))
	}
	return res
}

func toCommonSlice(any interface{}) []interface{} {
	var res []interface{}
	switch v := any.(type) {
	case []string:
		for _, item := range v {
			res = append(res, item)
		}
	default:
		runtimeExcption("toCommonSlice# unknown type:", reflect.TypeOf(any))
	}
	return res
}


func tokenToValue(t *Token)  Value {
	if t.isArrLiteral() {
		v := newJSONArray(t.ts)
		return newQKValue(v)
	} else if t.isObjLiteral() {
		v := newJSONObject(t.ts)
		return newQKValue(v)
	} else if t.isFloat() {
		f, err := strconv.ParseFloat(t.str, 64)
		assert(err != nil, "failed to parse float", t.String(), "line:", t.lineIndex)
		return newQKValue(f)
	} else if t.isInt() {
		i, err := strconv.Atoi(t.str)
		assert(err != nil, "failed to parse int", t.String(), "line:", t.lineIndex)
		return newQKValue(i)
	} else if t.isDynamicStr() {
		return newQKValue(t.str)
	} else if t.isStr() {
		str := strings.Replace(t.str, "\\\\", "\\", -1)
		str = strings.Replace(str, "\\r", "\r", -1) // 对 \r 进行转义
		str = strings.Replace(str, "\\n", "\n", -1) // 对 \n 进行转义
		str = strings.Replace(str, "\\t", "\t", -1) // 对 \t 进行转义
		return newQKValue(str)
	} else if t.isIdentifier() && (t.str == "true" || t.str == "false") {
		b, err := strconv.ParseBool(t.str)
		assert(err != nil, t.String(), "line:", t.lineIndex)
		return newQKValue(b)
	} else {
		return nil
	}
}

func toQKValue(v interface{}) Value {
	if v == nil {
		return NULL
	}
	typ := reflect.TypeOf(v)
	kind := typ.Kind()
	switch kind {
	case reflect.Map:
		mapRes := make(map[string]Value)
		m := v.(map[string]interface{})
		for key, value := range m {
			var qkVal Value
			if isDecomposable(value) {
				qkVal = toQKValue(value)
			} else {
				qkVal = newQKValue(value)
			}
			mapRes[key] = qkVal
		}
		tmp := toJSONObject(mapRes)
		return newQKValue(tmp)

	case reflect.Slice:
		var arrRes []Value
		list := v.([]interface{})
		for _, item := range list {
			var qkVal Value
			if isDecomposable(item) {
				qkVal = toQKValue(item)
			} else {
				qkVal = newQKValue(item)
			}
			arrRes = append(arrRes, qkVal)
		}
		tmp := toJSONArray(arrRes)
		return newQKValue(tmp)

	default:
		return newQKValue(v)
	}
}

// 是否为Map或Slice类型
func isDecomposable(v interface{}) bool {
	if v == nil {
		return false
	}
	kind := reflect.TypeOf(v).Kind()
	return kind == reflect.Map || kind == reflect.Slice
}

func toBoolean(raw Value) bool {
	if raw == nil || raw.isNULL() {
		return false
	}
	if raw.isInt() {
		return raw.val().(int64) != 0
	} else if raw.isFloat() {
		return raw.val().(float64) != 0
	} else if raw.isBoolean() {
		return raw.val().(bool)
	} else if raw.isString() {
		return raw.val().(string) != ""
	} else if raw.isJsonArray() {
		return raw.val() != nil
	} else if raw.isJsonObject() {
		return raw.val() != nil
	} else if raw.isAny() {
		return raw.val() != nil
	} else {
		runtimeExcption("toBoolean: unknown value type")
		return false
	}
}