package core

import (
    "fmt"
    "reflect"
)

type ValueType int

const (
    IntValue ValueType = 1 << iota // 整型
    FloatValue // 浮点型
    BooleanValue // 布尔类型
    StringValue // 字符串
    AnyValue // 任意值
    ArrayValue // json 数组
    ObjectValue // json 对象
    NULLValue // 空值
)

// 空值
var NULL = &Value{
    t:           NULLValue,
}

type Value struct {
    t ValueType
    integer int
    decimal float64
    boolean bool
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
        val = &Value{t: IntValue, integer: v}
    case int64:
        val = &Value{t: IntValue, integer: int(v)}
    case int32:
        val = &Value{t: IntValue, integer: int(v)}
    case float64:
        val = &Value{t: FloatValue, decimal: v}
    case float32:
        val = &Value{t: FloatValue, decimal: float64(v)}
    case bool:
        val = &Value{t: BooleanValue, boolean: v}
    case string:
        val = &Value{t: StringValue, str: v}
    case JSONArray:
        val = &Value{t: ArrayValue, jsonArr: v}
    case JSONObject:
        val = &Value{t: ObjectValue, jsonObj: v}
    default:
        val = &Value{t: AnyValue, any: v}
    }
    return val
}


func (v *Value) val() interface{} {
    switch {
        case v.isIntValue(): return v.integer
        case v.isFloatValue(): return v.decimal
        case v.isBooleanValue(): return v.boolean
        case v.isStringValue(): return v.str
        case v.isArrayValue(): return v.jsonArr
        case v.isObjectValue(): return v.jsonObj
        case v.isAnyValue(): {
            si, ok := v.any.(fmt.Stringer)
            if ok {
                return si.String()
            }
            return v.any
        }
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

func (v *Value) isClass() bool {
    return v.isAnyValue() && reflect.TypeOf(v.any).AssignableTo(ClassType)
}

func (v *Value) isArrayValue() bool {
    return (v.t & ArrayValue) == ArrayValue
}

func (v *Value) isObjectValue() bool {
    return (v.t & ObjectValue) == ObjectValue
}








