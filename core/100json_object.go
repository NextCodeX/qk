package core

import (
	"bytes"
	"fmt"
)

type JSONObject interface {
    parsed() bool
    init()
    size() int
    exist(key string) bool
    remove(key string)
    put(key string, value *Value)
    get(key string) *Value
    keys() []string
    values() []*Value
    tokens() []Token
    String() string
	toJSONObjectString() string
    Iterator
}

type JSONObjectImpl struct {
    val map[string]*Value
    ts []Token
    parsedFlag bool
}

func newJSONObject(ts []Token) JSONObject {
    return &JSONObjectImpl{ts:ts}
}

func toJSONObject(v map[string]*Value) JSONObject {
    return &JSONObjectImpl{val:v, parsedFlag:true}
}

func (obj *JSONObjectImpl) init() {
    obj.parsedFlag = true
    obj.val =  make(map[string]*Value)
}

func (obj *JSONObjectImpl) parsed() bool {
    return obj.val != nil
}

func (obj *JSONObjectImpl) size() int {
    return len(obj.val)
}

func (obj *JSONObjectImpl) remove(key string) {
    delete(obj.val, key)
}

func (obj *JSONObjectImpl) exist(key string) bool {
    _, ok := obj.val[key]
    return ok
}

func (obj *JSONObjectImpl) put(key string, value *Value) {
    obj.val[key] = value
}


func (obj *JSONObjectImpl) get(key string) *Value {
    v, ok := obj.val[key]
    if ok {
        return v
    }
    return NULL
}

func (obj *JSONObjectImpl) keys() []string {
    var keys []string
    for key := range obj.val {
        keys = append(keys, key)
    }
    return keys
}

func (obj *JSONObjectImpl) values() []*Value {
    var vals []*Value
    for _, v := range obj.val {
        vals = append(vals, v)
    }
    return vals
}

func (obj *JSONObjectImpl) tokens() []Token {
    return obj.ts
}

func (obj *JSONObjectImpl) String() string {
	return obj.toJSONObjectString()
}

func (obj *JSONObjectImpl) toJSONObjectString() string {
	var res bytes.Buffer
	res.WriteString("{")
	var i int
	for k, v := range obj.val {
		kstr := fmt.Sprintf(`"%v"`, k)
		var rawVal interface{}
		if v.isStringValue() {
			rawVal = fmt.Sprintf(`"%v"`, v.str)
		} else if v.isObjectValue() {
			rawVal = v.jsonObj.toJSONObjectString()
		} else if v.isArrayValue() {
			rawVal = v.jsonArr.toJSONArrayString()
		} else {
			rawVal = v.val()
		}
		if i < 1 {
			i++
			res.WriteString(fmt.Sprintf("%v:%v", kstr, rawVal))
		} else {
			res.WriteString(fmt.Sprintf(", %v:%v", kstr, rawVal))
		}
	}
	res.WriteString("}")

	return res.String()
}

func (obj *JSONObjectImpl) indexs() []interface{} {
    var res []interface{}
    for key := range obj.val {
        res = append(res, key)
    }
    return res
}

func (obj *JSONObjectImpl) getItem(index interface{}) *Value {
    key := index.(string)
    return obj.val[key]
}


