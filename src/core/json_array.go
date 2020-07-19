package core

type JSONArray interface {
    setParsed()
    parsed() bool
    size() int
    add(elem *Value)
    set(index int, elem *Value)
    get(index int) *Value
    checkOutofIndex(index int) bool
    values() []*Value
    tokens() []Token
}

type JSONArrayImpl struct {
    val []*Value
    ts []Token
    parsedFlag bool
}

func newJSONArray(ts []Token) JSONArray {
    return &JSONArrayImpl{ts:ts}
}

func (obj *JSONArrayImpl) setParsed() {
    obj.parsedFlag = true
}

func (obj *JSONArrayImpl) parsed() bool {
    return obj.val != nil
}

func (arr *JSONArrayImpl) size() int {
    return len(arr.val)
}

func (arr *JSONArrayImpl) add(elem *Value) {
    arr.val = append(arr.val, elem)
}

func (arr *JSONArrayImpl) set(index int, elem *Value) {
    arr.val[index] = elem
}

func (arr *JSONArrayImpl) get(index int) *Value {
    return arr.val[index]
}

func (arr *JSONArrayImpl) checkOutofIndex(index int) bool {
    if index<0 || index >= len(arr.val) {
        return true
    }
    return false
}

func (arr *JSONArrayImpl) values() []*Value {
    return arr.val
}

func (obj *JSONArrayImpl) tokens() []Token {
    return obj.ts
}
