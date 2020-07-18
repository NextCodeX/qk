package core

type JSONArray interface {
    size()
    add(elem interface{})
    set(index int, elem interface{})
    get(index int) interface{}
    getValue(index int) *Value
    rawValueList() []interface{}
    ValueList() []*Value
}
