package core

import "reflect"

func toCommonMap(any interface{}) map[string]interface{} {
	res := make(map[string]interface{})
	switch v := any.(type) {
	case map[string]string:
		for key, val := range v {
			res[key] = val
		}
	default:
		runtimeExcption("toCommonMap# unknown type:", reflect.TypeOf(any))
	}
	return res
}

func toCommonSlice(any interface{}) []interface{} {
	var res []interface{}
	switch v := any.(type) {
	case []string:
		for _, item := range v {
			res = append(res, item)
		}
	default:
		runtimeExcption("toCommonSlice# unknown type:", reflect.TypeOf(any))
	}
	return res
}


func toQKValue(v interface{}) *Value {
	typ := reflect.TypeOf(v)
	kind := typ.Kind()
	switch kind {
	case reflect.Map:
		mapRes := make(map[string]*Value)
		m := v.(map[string]interface{})
		for key, value := range m {
			var qkVal *Value
			if isDecomposable(value) {
				qkVal = toQKValue(value)
			} else {
				qkVal = newVal(value)
			}
			mapRes[key] = qkVal
		}
		tmp := toJSONObject(mapRes)
		return newVal(tmp)

	case reflect.Slice:
		var arrRes []*Value
		list := v.([]interface{})
		for _, item := range list {
			var qkVal *Value
			if isDecomposable(item) {
				qkVal = toQKValue(item)
			} else {
				qkVal = newVal(item)
			}
			arrRes = append(arrRes, qkVal)
		}
		tmp := toJSONArray(arrRes)
		return newVal(tmp)

	default:
		return newVal(v)
	}
}

func isDecomposable(v interface{}) bool {
	kind := reflect.TypeOf(v).Kind()
	return kind == reflect.Map || kind == reflect.Slice
}