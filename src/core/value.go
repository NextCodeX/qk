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
    int_value int
    float_value float64
    bool_value bool
    str_value string
    any_value interface{}
    arr_value *JSONArray
    obj_value *JSONObject
}

func newVal(rawVal interface{}) *Value {
    var val *Value
    switch v := rawVal.(type) {
    case int:
        val = &Value{t: IntValue, int_value: v}
    case float64:
        val = &Value{t: FloatValue, float_value: v}
    case float32:
        val = &Value{t: FloatValue, float_value: float64(v)}
    case bool:
        val = &Value{t: BooleanValue, bool_value: v}
    case string:
        val = &Value{t: StringValue, str_value: v}
    case JSONArray:
        val = &Value{t: ArrayValue, arr_value: &v}
    case JSONObject:
        val = &Value{t: ObjectValue, obj_value: &v}
    default:
        panic(fmt.Sprintln("unknow exception when newVal:", rawVal))
    }
    return val
}


func (v *Value) val() interface{} {
    switch {
    case v.isIntValue(): return v.int_value
    case v.isFloatValue(): return v.float_value
    case v.isBooleanValue(): return v.bool_value
    case v.isStringValue(): return v.str_value
    case v.isArrayValue(): return v.arr_value
    case v.isObjectValue(): return v.obj_value
    case v.isAnyValue(): return v.any_value
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











