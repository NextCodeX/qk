package core

import (
    "bytes"
    "fmt"
)

type JSONArray interface {
    parsed() bool
    size() int
    add(elem *Value)
    remove(index int)
    set(index int, elem *Value)
    get(index int) *Value
    checkOutofIndex(index int) bool
    values() []*Value
    tokens() []Token
    String() string
    toJSONArrayString() string
    Iterator
}

type JSONArrayImpl struct {
    val []*Value
    ts []Token
}

func newJSONArray(ts []Token) JSONArray {
    return &JSONArrayImpl{ts:ts}
}

func toJSONArray(v []*Value) JSONArray {
    return &JSONArrayImpl{val:v}
}

func (arr *JSONArrayImpl) parsed() bool {
    return arr.val != nil
}

func (arr *JSONArrayImpl) size() int {
    return len(arr.val)
}

func (arr *JSONArrayImpl) add(elem *Value) {
    arr.val = append(arr.val, elem)
}

func (arr *JSONArrayImpl)  remove(index int) {
    assert(arr.checkOutofIndex(index), "array out of index")
    newList := make([]*Value, 0, arr.size())
    newList = append(newList, arr.val[:index]...)
    if index + 1 < arr.size() {
        newList = append(newList, arr.val[index+1:]...)
    }
    arr.val = newList
}

func (arr *JSONArrayImpl) set(index int, elem *Value) {
    arr.val[index] = elem
}

func (arr *JSONArrayImpl) get(index int) *Value {
    return arr.val[index]
}

func (arr *JSONArrayImpl) checkOutofIndex(index int) bool {
    return index < 0 || index >= len(arr.val)
}

func (arr *JSONArrayImpl) values() []*Value {
    return arr.val
}

func (arr *JSONArrayImpl) tokens() []Token {
    return arr.ts
}

func (arr *JSONArrayImpl) String() string {
    return arr.toJSONArrayString()
}

func (arr *JSONArrayImpl) toJSONArrayString() string {
    var res bytes.Buffer
    res.WriteString("[")
    for i, item := range arr.val {
        var rawVal interface{}
        if item.isStringValue() {
            rawVal = fmt.Sprintf(`"%v"`, item.str)
        } else if item.isObjectValue() {
            rawVal = item.jsonObj.toJSONObjectString()
        } else if item.isArrayValue() {
            rawVal = item.jsonArr.toJSONArrayString()
        } else {
            rawVal = item.val()
        }
        if i < 1 {
            res.WriteString(fmt.Sprintf("%v", rawVal))
        } else {
            res.WriteString(fmt.Sprintf(", %v", rawVal))
        }
    }
    res.WriteString("]")
    return res.String()
}

func (arr *JSONArrayImpl) indexs() []interface{} {
    var res []interface{}
    for i := range arr.val {
        res = append(res, i)
    }
    return res
}

func (arr *JSONArrayImpl) getItem(index interface{}) *Value {
    i := index.(int)
    return arr.val[i]
}
