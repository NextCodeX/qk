package core

type Iterator interface {
	indexs() []interface{}
	getItem(interface{}) *Value
}

func toIterator(v *Value) Iterator {
	if v.isArrayValue() {
		return v.jsonArr
	}
	if v.isObjectValue() {
		return v.jsonObj
	}
	if v.isStringValue() {
		return newStringIterator(v.str)
	}
	return nil
}

type ForPlusInfo struct {
	indexName string // 索引变量名
	itemName string // 值变量名
	iterator string // 迭代器变量名
}

func newForPlusInfo(indexName, itemName, iterator string) *ForPlusInfo {
	return &ForPlusInfo{indexName, itemName, iterator }
}
