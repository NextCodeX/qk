package core

import (
    "bytes"
    "fmt"
)

type JSONArray interface {
    parsed() bool
    size() int
    add(elem Value)
    remove(index int)
    set(index int, elem Value)
    get(index int) Value
    checkOutofIndex(index int) bool
    values() []Value
    tokens() []Token
    String() string
    toJSONArrayString() string
    Iterator
    Value
}

type JSONArrayImpl struct {
    valList []Value
    ts []Token
}

func newJSONArray(ts []Token) JSONArray {
    return &JSONArrayImpl{ts:ts}
}

func toJSONArray(v []Value) JSONArray {
    return &JSONArrayImpl{valList:v}
}

func (arr *JSONArrayImpl) parsed() bool {
    return arr.valList != nil
}

func (arr *JSONArrayImpl) size() int {
    return len(arr.valList)
}

func (arr *JSONArrayImpl) add(elem Value) {
    arr.valList = append(arr.valList, elem)
}

func (arr *JSONArrayImpl)  remove(index int) {
    assert(arr.checkOutofIndex(index), "array out of index")
    newList := make([]Value, 0, arr.size())
    newList = append(newList, arr.valList[:index]...)
    if index + 1 < arr.size() {
        newList = append(newList, arr.valList[index+1:]...)
    }
    arr.valList = newList
}

func (arr *JSONArrayImpl) set(index int, elem Value) {
    arr.valList[index] = elem
}

func (arr *JSONArrayImpl) get(index int) Value {
    return arr.valList[index]
}

func (arr *JSONArrayImpl) checkOutofIndex(index int) bool {
    return index < 0 || index >= len(arr.valList)
}

func (arr *JSONArrayImpl) values() []Value {
    return arr.valList
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
    for i, item := range arr.valList {
        var rawVal interface{}
        if item.isString() {
            rawVal = fmt.Sprintf(`"%v"`, goStr(item))
        } else if item.isJsonObject() {
            rawVal = goObj(item).toJSONObjectString()
        } else if item.isJsonArray() {
            rawVal = goArr(item).toJSONArrayString()
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
    for i := range arr.valList {
        res = append(res, i)
    }
    return res
}

func (arr *JSONArrayImpl) getItem(index interface{}) Value {
    i := index.(int)
    return arr.valList[i]
}

func (arr *JSONArrayImpl) val() interface{} {
    return arr
}
func (arr *JSONArrayImpl) isNULL() bool {
    return false
}
func (arr *JSONArrayImpl) isInt() bool {
    return false
}
func (arr *JSONArrayImpl) isFloat() bool {
    return false
}
func (arr *JSONArrayImpl) isBoolean() bool {
    return false
}
func (arr *JSONArrayImpl) isString() bool {
    return false
}
func (arr *JSONArrayImpl) isAny() bool {
    return false
}
func (arr *JSONArrayImpl) isClass() bool {
    return false
}
func (arr *JSONArrayImpl) isJsonArray() bool {
    return true
}
func (arr *JSONArrayImpl) isJsonObject() bool {
    return false
}
