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
}

type JSONObjectImpl struct {
    val map[string]*Value
    ts []Token
    parsedFlag bool
}

func newJSONObject(ts []Token) JSONObject {
    return &JSONObjectImpl{ts:ts}
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


