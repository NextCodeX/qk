package core

import (
	"strconv"
	"strings"
)


func toInt(any interface{}) int {
	switch v := any.(type) {
	case int32:
		return int(v)
	case int64:
		return int(v)
	case int:
		return v
	case float32:
		return int(v)
	case float64:
		return int(v)
	case string:
		i, err := strconv.Atoi(v)
		assert(err!=nil, err, "Value:", any)
		return i
	case Value:
		return toInt(v.val())
	default:
		runtimeExcption("failed to int value", any)
	}
	return -1
}

func tokenToValue(t Token)  Value {
	if t.isArrLiteral() {
		v := rawJSONArray(t.tokens())
		return newQKValue(v)

	} else if t.isObjLiteral() {
		v := rawJSONObject(t.tokens())
		return newQKValue(v)

	} else if t.isFloat() {
		f, err := strconv.ParseFloat(t.raw(), 64)
		assert(err != nil, "failed to parse float", t.String(), "line:", t.getLineIndex())
		return newQKValue(f)

	} else if t.isInt() {
		i, err := strconv.Atoi(t.raw())
		assert(err != nil, "failed to parse int", t.String(), "line:", t.getLineIndex())
		return newQKValue(i)

	} else if t.isDynamicStr() {
		return newQKValue(t.raw())

	} else if t.isStr() {
		str := strings.Replace(t.raw(), "\\\\", "\\", -1)
		str = strings.Replace(str, "\\r", "\r", -1) // 对 \r 进行转义
		str = strings.Replace(str, "\\n", "\n", -1) // 对 \n 进行转义
		str = strings.Replace(str, "\\t", "\t", -1) // 对 \t 进行转义
		return newQKValue(str)

	} else if t.assertIdentifier("true") || t.assertIdentifier("false") {
		b, err := strconv.ParseBool(t.raw())
		assert(err != nil, t.String(), "line:", t.getLineIndex())
		return newQKValue(b)

	} else if t.assertIdentifier("null") {
		return NULL

	} else {
		return nil
	}
}


// QK Value 转 go 类型bool
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
	} else if raw.isAny() || raw.isObject() {
		return raw.val() != nil
	} else {
		runtimeExcption("toBoolean: unknown value type: ", raw)
		return false
	}
}