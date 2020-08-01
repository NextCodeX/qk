package core


type JSONObject interface {
    parsed() bool
    init()
    size() int
    exist(key string) bool
    put(key string, value *Value)
    get(key string) *Value
    keys() []string
    values() []*Value
    tokens() []Token
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


