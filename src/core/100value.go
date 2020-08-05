package core

import "fmt"

type ValueType int

const (
    IntValue ValueType = 1 << iota
    FloatValue
    BooleanValue
    StringValue
    AnyValue
    ArrayValue
    ObjectValue
    NULLValue
)

// 空值
var NULL = &Value{
    t:           NULLValue,
}

type Value struct {
    t ValueType
    int int
    float float64
    bool bool
    str string
    any interface{}
    jsonArr JSONArray
    jsonObj JSONObject
}

func newQkValue(rawVal interface{}) *Value {
    if rawVal == nil {
        return NULL
    }
    var val *Value
    switch v := rawVal.(type) {
    case int:
        val = &Value{t: IntValue, int: v}
    case float64:
        val = &Value{t: FloatValue, float: v}
    case float32:
        val = &Value{t: FloatValue, float: float64(v)}
    case bool:
        val = &Value{t: BooleanValue, bool: v}
    case string:
        val = &Value{t: StringValue, str: v}
    case JSONArray:
        val = &Value{t: ArrayValue, jsonArr: v}
    case JSONObject:
        val = &Value{t: ObjectValue, jsonObj: v}
    default:
        panic(fmt.Sprintln("unknow exception when newVal:", rawVal))
    }
    return val
}


func (v *Value) val() interface{} {
    switch {
    case v.isIntValue(): return v.int
    case v.isFloatValue(): return v.float
    case v.isBooleanValue(): return v.bool
    case v.isStringValue(): return v.str
    case v.isArrayValue(): return v.jsonArr
    case v.isObjectValue(): return v.jsonObj
    case v.isAnyValue(): return v.any
    }
    return nil
}

func (v *Value) isNULL() bool {
    return (v.t & NULLValue) == NULLValue
}

func (v *Value) isIntValue() bool {
    return (v.t & IntValue) == IntValue
}

func (v *Value) isFloatValue() bool {
    return (v.t & FloatValue) == FloatValue
}

func (v *Value) isBooleanValue() bool {
    return (v.t & BooleanValue) == BooleanValue
}

func (v *Value) isStringValue() bool {
    return (v.t & StringValue) == StringValue
}

func (v *Value) isAnyValue() bool {
    return (v.t & AnyValue) == AnyValue
}

func (v *Value) isArrayValue() bool {
    return (v.t & ArrayValue) == ArrayValue
}

func (v *Value) isObjectValue() bool {
    return (v.t & ObjectValue) == ObjectValue
}











