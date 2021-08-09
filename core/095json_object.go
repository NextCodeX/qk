package core

import (
	"bytes"
	"fmt"
)

type JSONObject interface {
	getName() string
    parsed() bool
    init()
    size() int
    exist(key string) bool
    remove(key string)
    put(key string, value Value)
    get(key string) Value
    keys() []string
    values() []Value
	mapVal() map[string]Value
    tokens() []Token
    String() string
	toJSONObjectString() string
    Iterator
    Value
}

type JSONObjectImpl struct {
	name string
    valMap map[string]Value
    ts []Token
    jsonObjectFlag bool
    ValueAdapter
}

func newJSONObject(ts []Token) JSONObject {
    return &JSONObjectImpl{ts:ts, jsonObjectFlag: true}
}

func toJSONObject(v map[string]Value) JSONObject {
    return &JSONObjectImpl{valMap:v, jsonObjectFlag: true}
}

func newClass(name string, v map[string]Value) JSONObject {
	return &JSONObjectImpl{name:name, valMap:v, jsonObjectFlag: false}
}

func (obj *JSONObjectImpl) getName() string {
	if obj.name == "" {
		return "json object"
	}
    return obj.name
}

func (obj *JSONObjectImpl) init() {
    obj.valMap =  make(map[string]Value)
}

func (obj *JSONObjectImpl) parsed() bool {
    return obj.valMap != nil
}

func (obj *JSONObjectImpl) size() int {
    return len(obj.valMap)
}

func (obj *JSONObjectImpl) remove(key string) {
    delete(obj.valMap, key)
}

func (obj *JSONObjectImpl) exist(key string) bool {
    _, ok := obj.valMap[key]
    return ok
}

func (obj *JSONObjectImpl) put(key string, value Value) {
    obj.valMap[key] = value
}


func (obj *JSONObjectImpl) get(key string) Value {
    v, ok := obj.valMap[key]
    if ok {
        return v
    }
    if obj.jsonObjectFlag {
    	return obj.returnFakeMethod(key)
	}
    return NULL
}

func (obj *JSONObjectImpl) keys() []string {
    var keys []string
    for key := range obj.valMap {
        keys = append(keys, key)
    }
    return keys
}

func (obj *JSONObjectImpl) values() []Value {
    var vals []Value
    for _, v := range obj.valMap {
        vals = append(vals, v)
    }
    return vals
}

func (obj *JSONObjectImpl) mapVal() map[string]Value {
    return obj.valMap
}

func (obj *JSONObjectImpl) tokens() []Token {
    return obj.ts
}

func (obj *JSONObjectImpl) String() string {
	if !obj.isJsonObject() {
		return "class " + obj.getName()
	}
	return obj.toJSONObjectString()
}

func (obj *JSONObjectImpl) toJSONObjectString() string {
	var res bytes.Buffer
	res.WriteString("{")
	var i int
	for k, v := range obj.valMap {
		kstr := fmt.Sprintf(`"%v"`, k)
		var rawVal interface{}
		if v.isString() {
			rawVal = fmt.Sprintf(`"%v"`, goStr(v))
		} else if v.isJsonObject() {
			rawVal = goObj(v).toJSONObjectString()
		} else if v.isJsonArray() {
			rawVal = goArr(v).toJSONArrayString()
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
    for key := range obj.valMap {
        res = append(res, key)
    }
    return res
}

func (obj *JSONObjectImpl) getItem(index interface{}) Value {
    key := index.(string)
    return obj.valMap[key]
}

func (obj *JSONObjectImpl) val() interface{} {
	return obj
}
func (obj *JSONObjectImpl) isJsonObject() bool {
	return obj.jsonObjectFlag
}
