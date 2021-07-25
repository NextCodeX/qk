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

// 用于收集增强for的基本信息
type ForPlusInfo struct {
	indexName string // 索引变量名
	itemName string // 值变量名
	iterator *Expression // 迭代器表达式
}

func newForPlusInfo(indexName, itemName string, iterator *Expression) *ForPlusInfo {
	return &ForPlusInfo{indexName, itemName, iterator }
}
