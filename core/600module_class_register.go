package core

import (
	"reflect"
)


// 成员变量相关信息
type FieldInfo struct {
	name    string // field name
	t reflect.Type // field go type
	v reflect.Value // field go internal value
}

func (f *FieldInfo) set(val interface{}) {
	f.v = reflect.ValueOf(val)
}

func (f *FieldInfo) get() interface{} {
	return f.v.Interface()
}

func collectFieldInfo(objPtr interface{}) (res map[string]*FieldInfo) {
	res = make(map[string]*FieldInfo)
	v := reflect.ValueOf(objPtr).Elem()
	k := v.Type()
	for i := 0; i < v.NumField(); i++ {
		key := k.Field(i)
		val := v.Field(i)
		if !val.CanInterface() { //CanInterface(): 判断该成员变量是否能被获取值
			continue
		}
		fieldName := formatName(key.Name)
		res[fieldName] = &FieldInfo{name: fieldName, t:val.Type(), v:val}
	}
	return res
}



