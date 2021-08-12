package core

import (
    "bytes"
    "fmt"
)

type JSONArray interface {
    parsed() bool
    Size() int
    add(elem Value)
    set(index int, elem Value)
    getElem(index int) Value
    sub(start, end int) Value
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
    ClassObject
}

func rawJSONArray(ts []Token) JSONArray {
    return newJsonArray(nil, ts)
}

func array(v []Value) JSONArray {
    return newJsonArray(v, nil)
}

func newJsonArray(v []Value, ts []Token) JSONArray {
    arr :=  &JSONArrayImpl{valList:v, ts: ts}
    arr.ClassObject.raw = &arr
    arr.ClassObject.name = "JSONArray"
    return arr
}

func (arr *JSONArrayImpl) parsed() bool {
    return arr.valList != nil
}

func (arr *JSONArrayImpl) Size() int {
    return len(arr.valList)
}

func (arr *JSONArrayImpl) add(elem Value) {
    arr.valList = append(arr.valList, elem)
}

func (arr *JSONArrayImpl) set(index int, elem Value) {
    arr.valList[index] = elem
}

func (arr *JSONArrayImpl) getElem(index int) Value {
    return arr.valList[index]
}

func (arr *JSONArrayImpl) sub(start, end int) Value {
    return array(arr.valList[start:end])
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

func (arr *JSONArrayImpl) Add(args []interface{}) {
    for _, arg := range args {
        arr.add(newQKValue(arg))
    }
}
func (arr *JSONArrayImpl) Remove(index int) {
    assert(arr.checkOutofIndex(index), "array out of index")
    newList := make([]Value, 0, arr.Size())
    newList = append(newList, arr.valList[:index]...)
    if index + 1 < arr.Size() {
        newList = append(newList, arr.valList[index+1:]...)
    }
    arr.valList = newList
}
func (arr *JSONArrayImpl) Join(seperator string) string {
    vals := arr.values()
    var res bytes.Buffer
    for i, val := range vals {
        if i > 0 {
            res.WriteString(seperator)
        }
        valStr := fmt.Sprintf("%v", val.val())
        res.WriteString(valStr)
    }
    return res.String()
}



func (arr *JSONArrayImpl) val() interface{} {
    return arr
}
func (arr *JSONArrayImpl) isJsonArray() bool {
    return true
}
