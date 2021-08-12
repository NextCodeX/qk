package core

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type JSONObject interface {
    parsed() bool
    init()
    Size() int
	Contain(key string) bool
    Remove(key string)
    put(key string, value Value)
    get(key string) Value
    keys() []string
    values() []Value
	mapVal() map[string]Value
    tokens() []Token
    String() string
	Pretty()
	toJSONObjectString() string
    Iterator
    Value
}

type JSONObjectImpl struct {
    valMap map[string]Value
    ts []Token
	ClassObject
}

func rawJSONObject(ts []Token) JSONObject {
	return newJsonObject(nil, ts)
}

func jsonObject(v map[string]Value) JSONObject {
    return newJsonObject(v, nil)
}

func newJsonObject(v map[string]Value, ts []Token) JSONObject {
	obj := &JSONObjectImpl{valMap:v, ts: ts}
	obj.ClassObject.raw = &obj
	obj.ClassObject.name = "JSONObject"
	return obj
}

func (obj *JSONObjectImpl) init() {
    obj.valMap =  make(map[string]Value)
}

func (obj *JSONObjectImpl) parsed() bool {
    return obj.valMap != nil
}

func (obj *JSONObjectImpl) Size() int {
    return len(obj.valMap)
}

func (obj *JSONObjectImpl) Remove(key string) {
    delete(obj.valMap, key)
}

func (obj *JSONObjectImpl) Contain(key string) bool {
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
	res := obj.ClassObject.get(key)
    if res == nil {
    	return NULL
	}
    return res
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

func (obj *JSONObjectImpl) Pretty() {
	uglyBody := obj.toJSONObjectString()
	var out bytes.Buffer
	err := json.Indent(&out, []byte(uglyBody), "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(out.String())
}

func (obj *JSONObjectImpl) String() string {
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
	return true
}
