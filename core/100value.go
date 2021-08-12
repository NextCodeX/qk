package core


// 空值
var NULL = newNULLValue()

type Value interface {
   String() string
   val() interface{}
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

func newQKValue(rawVal interface{}) Value {
    if rawVal == nil {
        return NULL
    }
    var val Value
    switch v := rawVal.(type) {
    case []byte:
        val = newByteArrayValue(v)
    case int:
        val = newIntValue(int64(v))
    case int64:
        val = newIntValue(v)
    case int32:
        val = newIntValue(int64(v))
    case float64:
        val = newFloatValue(v)
    case float32:
        val = newFloatValue(float64(v))
    case bool:
        val = newBooleanValue(v)
    case string:
        val = newStringValue(v)
    case JSONArray:
        val = v
    case JSONObject:
        val = v
    case Function:
        val = v
    case Value:
        val = v
    case []Value:
        val = array(v)
    case map[string]Value:
        val = jsonObject(v)
    case map[string]string:
        mapRes := make(map[string]Value)
        for key, value := range v {
            mapRes[key] = newQKValue(value)
        }
        tmp := jsonObject(mapRes)
        return tmp

    case map[string]interface{}:
        mapRes := make(map[string]Value)
        for key, value := range v {
            mapRes[key] = newQKValue(value)
        }
        tmp := jsonObject(mapRes)
        return tmp

    case []string:
        var arrRes []Value
        for _, item := range v {
            arrRes = append(arrRes, newQKValue(item))
        }
        tmp := array(arrRes)
        return tmp

    case []interface{}:
        var arrRes []Value
        for _, item := range v {
            arrRes = append(arrRes, newQKValue(item))
        }
        tmp := array(arrRes)
        return tmp

    default:
        val = newAnyValue(v)
    }
    return val
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
        runtimeExcption("value is not string")
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




