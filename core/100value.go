package core

import (
	"fmt"
)

var (
	NULL  = newNULLValue() // 空值
	TRUE  = newQKValue(true)
	FALSE = newQKValue(false)
)

type Value interface {
	String() string
	val() interface{}
	typeName() string
	isNULL() bool
	isByteArray() bool
	isInt() bool
	isFloat() bool
	isBoolean() bool
	isString() bool
	isJsonArray() bool
	isJsonObject() bool
	isFunction() bool
	isObject() bool
	isAny() bool
}

func newQKValue(rawVal any) Value {
	if rawVal == nil || rawVal == NULL {
		return NULL
	}
	switch v := rawVal.(type) {
	case Value: return v
	case JSONArray: return v
	case JSONObject: return v
	case Function: return v

	case []byte: return newByteArrayValue(v)
	case int: return newIntValue(int64(v))
	case int64: return newIntValue(v)
	case int32: return newIntValue(int64(v))
	case float64: return newFloatValue(v)
	case float32: return newFloatValue(float64(v))
	case bool: return newBooleanValue(v)
	case string: return newStringValue(v)

	case map[string]Value: return jsonObject(v)
	case map[string]string: return jsonObject(v)
	case map[string]interface{}: return jsonObject(v)

	case []Value: return array(v)
	case []string: return array(v)
	case [][]byte: return array(v)
	case []interface{}: return array(v)

	default: return newAnyValue(v)
	}
}

func goBytes(val Value) []byte {
	if v, ok := val.(*ByteArrayValue); ok {
		return v.goValue
	} else {
		runtimeExcption("value is not ByteArray")
		return nil
	}
}

func goInt(val Value) int64 {
	if v, ok := val.(*IntValue); ok {
		return v.goValue
	} else {
		runtimeExcption("value is not int")
		return -1
	}
}

func goFloat(val Value) float64 {
	if v, ok := val.(*FloatValue); ok {
		return v.goValue
	} else {
		runtimeExcption("value is not float")
		return -1
	}
}

func goBool(val Value) bool {
	if v, ok := val.(*BooleanValue); ok {
		return v.goValue
	} else {
		runtimeExcption("value is not boolean")
		return false
	}
}

func goStr(val Value) string {
	if v, ok := val.(*StringValue); ok {
		return v.goValue
	} else {
		fmt.Println("value is not string:", val, val.typeName())
		return ""
	}
}

func goAny(val Value) interface{} {
	if v, ok := val.(*AnyValue); ok {
		return v.goValue
	} else {
		runtimeExcption("value is not any type")
		return nil
	}
}

func goArr(val Value) JSONArray {
	if v, ok := val.(JSONArray); ok {
		return v
	} else {
		runtimeExcption("value is not json array")
		return nil
	}
}

func goObj(val Value) JSONObject {
	if v, ok := val.(JSONObject); ok {
		return v
	} else {
		runtimeExcption("value is not json object")
		return nil
	}
}

func goQKObj(val Value) Object {
	if v, ok := val.(Object); ok {
		return v
	} else {
		runtimeExcption("value is not object")
		return nil
	}
}
