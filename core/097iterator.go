package core

type Iterator interface {
	indexs() []interface{}
	getItem(interface{}) Value
}

func toIterator(v Value) Iterator {
	if v.isJsonArray() {
		return goArr(v)
	}
	if v.isJsonObject() {
		return goObj(v)
	}
	if v.isString() {
		return newStringIterator(goStr(v))
	}
	return nil
}
