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

// 用于收集增强for的基本信息
type ForPlusInfo struct {
	indexName string // 索引变量名
	itemName string // 值变量名
	iterator *Expression // 迭代器表达式
}

func newForPlusInfo(indexName, itemName string, iterator *Expression) *ForPlusInfo {
	return &ForPlusInfo{indexName, itemName, iterator }
}
