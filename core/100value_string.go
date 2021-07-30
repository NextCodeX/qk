package core


type StringValue struct {
   goValue string
   ValueAdapter
}

func newStringValue(raw string) Value {
    return &StringValue{goValue: raw}
}

func (str *StringValue) val() interface{} {
    return str.goValue
}

func (str *StringValue) isString() bool {
    return true
}
