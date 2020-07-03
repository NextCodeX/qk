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
)

type Value struct {
    t ValueType
    int_value int
    float_value float64
    bool_value bool
    str_value string
    any_value interface{}
    arr_value []interface{}
    obj_value map[string]interface{}
}

func newVal(rawVal interface{}) Value {
    var val Value
    switch v := rawVal.(type) {
    case int:
        val = Value{t: IntValue, int_value: v}
    case float64:
        val = Value{t: FloatValue, float_value: v}
    case float32:
        val = Value{t: FloatValue, float_value: float64(v)}
    case bool:
        val = Value{t: BooleanValue, bool_value: v}
    case string:
        val = Value{t: StringValue, str_value: v}
    case []interface{}:
        val = Value{t: ArrayValue, arr_value: v}
    case map[string]interface{}:
        val = Value{t: ObjectValue, obj_value: v}
    default:
        panic(fmt.Sprintln("unknow exception when newVal:", rawVal))
    }
    return val
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

type StatementResultType int

const (
    StatementReturn StatementResultType = 1 << iota
    StatementContinue
    StatementBreak
    StatementNormal
)

type StatementResult struct {
    t StatementResultType
    val Value
}

func (this *StatementResult) isStatementReturn() bool {
    return (this.t & StatementReturn) == StatementReturn
}

func (this *StatementResult) isStatementContinue() bool {
    return (this.t & StatementContinue) == StatementContinue
}

func (this *StatementResult) isStatementBreak() bool {
    return (this.t & StatementBreak) == StatementBreak
}

func (this *StatementResult) isStatementNormal() bool {
    return (this.t & StatementNormal) == StatementNormal
}

type Variable struct{
    name string
    val Value
}

func newVar(name string, rawVal interface{}) *Variable {
    res := &Variable{
        name: name,
        val:  newVal(rawVal),
    }

    return res
}

func toVar(name string, rawVal Value) *Variable {
    res := &Variable{
        name: name,
        val:  rawVal,
    }
    return res
}


type Variables map[string]*Variable

func newVariables() Variables {
    return make(map[string]*Variable)
}

func (vs *Variables) add(v *Variable) {
    (*vs)[v.name] = v
}

func (vs *Variables) get(name string) *Variable {
    return (*vs)[name]
}
