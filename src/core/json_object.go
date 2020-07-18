package core


type JSONObject interface {
    size()
    exist(key string)
    put(key string, value interface{})
    get(key string) interface{}
    getValue(key string) *Value
    keys() []string
    rawValues() []interface{}
    Values() []*Value
}
