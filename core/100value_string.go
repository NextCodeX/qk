package core


type StringValue struct {
   goValue string
}

func newStringValue(raw string) *StringValue {
    return &StringValue{raw}
}

func (str *StringValue) val() interface{} {
    return str.goValue
}
func (str *StringValue) isNULL() bool {
    return false
}
func (str *StringValue) isInt() bool {
    return false
}
func (str *StringValue) isFloat() bool {
    return false
}
func (str *StringValue) isBoolean() bool {
    return false
}
func (str *StringValue) isString() bool {
    return true
}
func (str *StringValue) isAny() bool {
    return false
}
func (str *StringValue) isClass() bool {
    return false
}
func (str *StringValue) isJsonArray() bool {
    return false
}
func (str *StringValue) isJsonObject() bool {
    return false
}

